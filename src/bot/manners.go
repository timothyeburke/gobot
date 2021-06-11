package bot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ConfigureManners(b *Bot, m *Message, u *User) {
	if regexp.MustCompile(fmt.Sprintf("^<@%s>\\pZ+DISABLE MANNERS", b.ID)).MatchString(m.Message) {
		var settings map[string]string
		var result = b.Brain.Get("manners")
		_ = json.Unmarshal(result, &settings)

		if settings == nil {
			return
		}

		var channel, _ = b.GetChannel(m.Channel)
		delete(settings, channel.Name)

		j, _ := json.Marshal(settings)
		b.Brain.Save("manners", j)

		m.Message = "Disabled manners. `@here`/`@channel` all you want."
		b.Send(m)
		return
	}

	var regex = regexp.MustCompile(fmt.Sprintf("^<@%s>\\pZ+ENABLE MANNERS (.+)", b.ID))
	if !regex.MatchString(m.Message) {
		return
	}

	var match = strings.TrimSpace(regex.FindStringSubmatch(m.Message)[1])
	var channel, _ = b.GetChannel(m.Channel)

	var split = strings.Split(match, " ")
	var group = split[0]
	var handleRegex = regexp.MustCompile("@(.+)>")
	var handleResult = handleRegex.FindStringSubmatch(group)

	if len(handleResult) < 2 {
		m.Message = "Slack Team not found. Cannot enforce manners."
		b.Send(m)
		return
	}

	var handle = handleResult[1]
	groups, _ := b.GetUserGroups()
	var userGroup UserGroup
	for _, g := range groups {
		if g.Handle == handle {
			userGroup = g
			break
		}
	}

	if userGroup.ID == "" {
		m.Message = "Slack Team not found. Cannot enforce manners."
		b.Send(m)
		return
	}

	var settings map[string]string
	var result = b.Brain.Get("manners")
	_ = json.Unmarshal(result, &settings)

	if settings == nil {
		settings = make(map[string]string)
	}

	settings[channel.Name] = handle
	j, _ := json.Marshal(settings)
	b.Brain.Save("manners", j)

	m.Message = fmt.Sprintf("Enforcing manners, only members of `@%s` (%s) may `@here` or `@channel`", handle, userGroup.Name)
	b.Send(m)
}

func Manners(b *Bot, m *Message, u *User) {
	if !regexp.MustCompile(`(<!channel>|<!here>)`).MatchString(m.Message) {
		return
	}

	var Channel, err = b.GetChannel(m.Channel)
	if err != nil {
		return
	}

	var settings map[string]string
	var result = b.Brain.Get("manners")
	if result != nil {
		err = json.Unmarshal(result, &settings)
		if err != nil {
			return
		}
	} else {
		j, _ := json.Marshal(settings)
		b.Brain.Save("manners", j)
		return
	}

	handle, ok := settings[Channel.Name]
	if !ok {
		return
	}

	userGroups, _ := b.GetUserGroups()
	var userGroup UserGroup
	for _, group := range userGroups {
		if group.Handle == handle {
			userGroup = group
			break
		}
	}

	users, _ := b.GetGroupUsers(userGroup.ID)
	for _, user := range users {
		if u.ID == user {
			return
		}
	}

	var template = "Please do not use `@channel`/`@here` in this channel. Use `@%s` to ask for assistance from %s."
	m.Message = fmt.Sprintf(template, handle, userGroup.Name)
	b.Send(m)
}
