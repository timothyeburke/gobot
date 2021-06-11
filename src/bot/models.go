package bot

import (
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
)

type Bot struct {
	ID       string
	Brain    Brain
	Cache    *cache.Cache
	Name     string
	Socket   *websocket.Conn
	Team     string
	Token    string
	APIToken string
}

type enterprise struct {
	ID   string `json:"id"`
	Name string `json:"enterprise_name"`
}

type channelResponse struct {
	ID      string  `json:"id"`
	Channel Channel `json:"channel"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type handshakeResponseTeam struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type handshakeResponseSelf struct {
	ID         string     `json:"id"`
	Enterprise enterprise `json:"enterprise_user"`
	Name       string     `json:"name"`
}

type handshakeResponse struct {
	Ok    bool                  `json:"ok"`
	Error string                `json:"error"`
	Self  handshakeResponseSelf `json:"self"`
	Team  handshakeResponseTeam `json:"team"`
	URL   string                `json:"url"`
}

type Message struct {
	ID      uint64 `json:"id"`
	Bot     *Bot
	Channel string `json:"channel"`
	Message string `json:"text"`
	Subtype string `json:"subtype"`
	Thread  string `json:"thread_ts"`
	TS      string `json:"ts"`
	Type    string `json:"type"`
	UserID  string `json:"user"`
}

type MessageAttachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Fallback  string `json:"fallback"`
	ImageURL  string `json:"image_url"`
	TS        string `json:"ts"`
	Actions   []MessageAttachmentAction `json:"actions"`
}

type MessageAttachmentAction struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Style string `json:"style"`
	Url   string `json:"url"`
}

type userResponse struct {
	Ok   bool `json:"ok"`
	User User `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
}

type userGroupUsersResponse struct {
	OK    bool     `json:"ok"`
	Users []string `json:"users"`
}

type UserGroup struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Name   string `json:"name"`
}

type userGroupResponse struct {
	OK         bool        `json:"ok"`
	Usergroups []UserGroup `json:"usergroups"`
}
