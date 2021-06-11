package bot

import (
	"regexp"
)

func Unflip(b *Bot, m *Message, u *User) {
    if !regexp.MustCompile(`(?i)[┻|┸]{1}[━-]*[┻|┸]{1}`).MatchString(m.Message) {
		return
	}

	m.Message = "┬─┬ノ( º _ ºノ)"
	b.Send(m)
}
