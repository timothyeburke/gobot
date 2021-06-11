package bot

import (
	"math/rand"
	"regexp"
)

func Goose(b *Bot, m *Message, u *User) {
	if !regexp.MustCompile(`(?i)(goose|geese)`).MatchString(m.Message) {
		return
	}

	var allowedRooms = map[string]bool {
		"just-yelling-training": true,
	}
	var Channel, _ = b.GetChannel(m.Channel)
	if !allowedRooms[Channel.Name] {
		return
	}

	var answers = []string{
		"HONK",
		"honk honk",
		"HOOOOOOOOOOOOONNNNNNKKKKKKK",
		":goose:",
		":honk:",
		"I think I will cause problems on purpose.",
	}
	m.Message = answers[rand.Intn(len(answers))]
	b.Send(m)
}
