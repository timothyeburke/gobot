package bot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConfigureAfterHours(b *Bot, m *Message, u *User) {
	if regexp.MustCompile(fmt.Sprintf("^<@%s>\\pZ+DISABLE AFTER HOURS", b.ID)).MatchString(m.Message) {
		var settings map[string][]string
		var result = b.Brain.Get("after-hours")
		_ = json.Unmarshal(result, &settings)

		if settings == nil {
			return
		}

		var channel, _ = b.GetChannel(m.Channel)
		delete(settings, channel.Name)

		j, _ := json.Marshal(settings)
		b.Brain.Save("after-hours", j)

		m.Message = "Disabled after-hours message."
		b.Send(m)
		return
	}

	var regex = regexp.MustCompile(fmt.Sprintf("^<@%s>\\pZ+ENABLE AFTER HOURS (\\d+)-(\\d+)", b.ID))
	if !regex.MatchString(m.Message) {
		return
	}

	var channel, _ = b.GetChannel(m.Channel)
	var matches = regex.FindStringSubmatch(m.Message)
	var dayStart = matches[1]
	var dayEnd = matches[2]

	var settings map[string][]string
	var result = b.Brain.Get("after-hours")
	_ = json.Unmarshal(result, &settings)

	if settings == nil {
		settings = make(map[string][]string)
	}

	settings[channel.Name] = []string{ dayStart, dayEnd }
	j, _ := json.Marshal(settings)
	b.Brain.Save("after-hours", j)

	m.Message = fmt.Sprintf("Outside of %s to %s will reply with off hours P1 message.", dayStart, dayEnd)
	b.Send(m)
}

func AfterHours(b *Bot, m *Message, u *User) {
	if m.Thread != "" || regexp.MustCompile(fmt.Sprintf("^<@%s>", b.ID)).MatchString(m.Message) {
		return
	}

	var settings map[string][]string
	var result = b.Brain.Get("after-hours")
	err := json.Unmarshal(result, &settings)
	if err != nil {
		return
	}

	var channel, _ = b.GetChannel(m.Channel)
	if settings[channel.Name] == nil {
		return
	}

	var start, _ = strconv.Atoi(settings[channel.Name][0])
	var end, _ = strconv.Atoi(settings[channel.Name][1])
	var message = strings.Join([]string{
		":wave: Thanks for reaching out! It's off-hours for most of us at the ",
		"moment :sleeping:, but someone will get back to you during business hours.",
		"\n*If this is a business critical issue, please ",
		"open a P1 to page our on-call engineer.*",
	}, "")

	ny, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(ny)

	startTime := time.Date(now.Year(), now.Month(), now.Day(), start, 0, 0, 0, ny)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), end, 0, 0, 0, ny)

	var dow = now.Weekday()
	if now.Before(startTime) || now.After(endTime) || dow == time.Saturday || dow == time.Sunday {
		m.Message = message
		m.Thread = m.TS
		b.Send(m)
	}
}
