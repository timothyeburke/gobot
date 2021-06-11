package bot

import (
	"math/rand"
	"regexp"
)

func GCP(b *Bot, m *Message, u *User) {
    if !regexp.MustCompile(`(?i)Is GCP down`).MatchString(m.Message) {
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
		"All signs point to yes.",
		"As I see it, yes.",
		"Better not to tell you now. There's a P1.",
		"If you have to ask...",
		"It is certain.",
		"It is decidedly so.",
		"Most likely.",
		"This is not looking good...",
		"Can't reach status.cloud.google.com, which isn't a good sign.",
		"Signs point to \"We're experiencing an issue affecting multiple...\".",
		"Without a doubt.",
		"Yes.",
		"Yes – definitely.",
		"Absolutely.",
	}
	m.Message = answers[rand.Intn(len(answers))]
	b.Send(m)

	m.Message = "But check status.cloud.google.com to be sure."
	b.Send(m)
}
