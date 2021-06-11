package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
)

// Plugins Go Here
var plugins = []func(*Bot, *Message, *User) {
	// AfterHours,
	// ConfigureAfterHours,
	// ConfigureManners,
	GCP,
	Goose,
	// Manners,
	Marx,
	// React,
	// Socialize,
	// Unflip,
	Yelling,
}

func (b *Bot) process(message Message) {
	if message.Type != "message" {
		return
	}

	user, _ := b.GetUser(message.UserID)
	if user.IsBot {
		return
	}

	for _, plugin := range plugins {
		var m Message
		j, _ := json.Marshal(message)
		json.Unmarshal(j, &m)

		plugin(b, &m, user)
	}
}

func New(botToken string, apiToken string) (*Bot, error) {
	bot := Bot{
		Token: botToken,
		APIToken: apiToken,
		Brain: Brain{},
		Cache: cache.New(15*time.Minute, 30*time.Minute),
	}

	bot.Brain.Bot = &bot

	return bot.handshake()
}

func (b *Bot) React(message *Message, emoji string) {
	parts := []string {
		"https://slack.com/api/reactions.add?",
		fmt.Sprintf("token=%s", b.Token),
		fmt.Sprintf("&channel=%s", message.Channel),
		fmt.Sprintf("&name=%s", url.QueryEscape(emoji)),
		fmt.Sprintf("&timestamp=%s", message.TS),
	}

	uri := strings.Join(parts, "")
	_, err := http.Post(uri, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) Send(message *Message) {
	b.Socket.WriteJSON(message)
}

func (b *Bot) SendAttachments(message *Message, attachments []MessageAttachment) {
	j, _ := json.Marshal(attachments)

	parts := []string{
		"https://slack.com/api/chat.postMessage?",
		fmt.Sprintf("token=%s", b.Token),
		"&as_user=true",
		fmt.Sprintf("&channel=%s", message.Channel),
		fmt.Sprintf("&text=%s", url.QueryEscape(message.Message)),
		fmt.Sprintf("&attachments=%s", url.QueryEscape(string(j))),
	}

	uri := strings.Join(parts, "")
	_, err := http.Post(uri, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) GetChannel(id string) (*Channel, error) {
	key := fmt.Sprintf("channel/%s", id)
	cached, found := b.Cache.Get(key)
	if found {
		return cached.(*Channel), nil
	}

	res, _ := http.Get(fmt.Sprintf("https://slack.com/api/conversations.info?token=%s&channel=%s", b.Token, id))

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var response channelResponse
	_ = json.Unmarshal(body, &response)

	b.Cache.Set(key, &response.Channel, cache.DefaultExpiration)

	return &response.Channel, nil
}

func (b *Bot) GetUserGroups() ([]UserGroup, error) {
	key := "usergroups"
	cached, found := b.Cache.Get(key)
	if found {
		return cached.([]UserGroup), nil
	}

	var url = "https://slack.com/api/usergroups.list?token=%s"
	res, _ := http.Get(fmt.Sprintf(url, b.APIToken))

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var response userGroupResponse
	_ = json.Unmarshal(body, &response)

	b.Cache.Set(key, response.Usergroups, cache.DefaultExpiration)

	return response.Usergroups, nil
}

func (b *Bot) GetGroupUsers(id string) ([]string, error) {
	key := fmt.Sprintf("usergroup.users/%s", id)
	cached, found := b.Cache.Get(key)
	if found {
		return cached.([]string), nil
	}

	var url = "https://slack.com/api/usergroups.users.list?token=%s&usergroup=%s"
	res, _ := http.Get(fmt.Sprintf(url, b.APIToken, id))

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var response userGroupUsersResponse
	_ = json.Unmarshal(body, &response)

	b.Cache.Set(key, response.Users, cache.DefaultExpiration)

	return response.Users, nil
}

func (b *Bot) GetUser(id string) (*User, error) {
	key := fmt.Sprintf("user/%s", id)
	cached, found := b.Cache.Get(key)
	if found {
		return cached.(*User), nil
	}

	res, err := http.Get(fmt.Sprintf("https://slack.com/api/users.info?token=%s&user=%s", b.Token, id))
	if err != nil {
		return nil, errors.New("Failed to connect to Slack RTM API")
	}

	// Check for HTTP status code
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed with HTTP Code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read body from response")
	}

	var response userResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %s", body)
	}

	b.Cache.Set(key, &response.User, cache.DefaultExpiration)

	return &response.User, nil
}

func (b *Bot) handshake() (*Bot, error) {
	res, err := http.Get(fmt.Sprintf("https://slack.com/api/rtm.connect?token=%s", b.Token))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Slack RTM API: %s", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed with HTTP Code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read body from response")
	}

	var response handshakeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %s", body)
	}

	if !response.Ok {
		return nil, fmt.Errorf(response.Error)
	}

	b.ID = response.Self.ID
	b.Name = response.Self.Name
	b.Team = response.Team.Name
	b.Socket, _, err = websocket.DefaultDialer.Dial(response.URL, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Websocket: %s", err)
	}

	log.Println(fmt.Sprintf("Logged into %s workspace as %s.", b.Team, b.Name))

	return b, nil
}

func (b *Bot) Listen() {
	var msg Message

	log.Println("Listening.")

	for {
		_, message, err := b.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Socket error. Exiting.")
				log.Println(err)
			}
			break
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			// This happens for message types we don't process or care about
			continue
		}

		go b.process(msg)
		msg = Message{}
	}
}
