package bot

import (
	"fmt"
	"regexp"
	"strings"
)

func React(b *Bot, m *Message, u *User) {
	var allowedRooms = map[string]bool {
		"just-yelling": true,
	}
	var Channel, _ = b.GetChannel(m.Channel)
	if !allowedRooms[Channel.Name] {
		return
	}

	var emoji = map[string]string{
		"new jersey": "new-jersey",
		"california": "california-flag",
		"colorado": "flag-us-co",
		"denver": "denver",
		"india": "flag-in",
		"ithaca": "ithaca",
		"seattle": "seattle",
		"awesome": "awesome",
	}

	keys := make([]string, 0, len(emoji))
	for key := range emoji {
		keys = append(keys, key)
	}

	var r = fmt.Sprintf("(?i)(%s)", strings.Join(keys, "|"))
	var regex = regexp.MustCompile(r)
	if !regex.MatchString(m.Message) {
		return
	}

	var matches = regex.FindAllString(m.Message, -1)
	for i:= 0; i < len(matches); i++ {
		var match = strings.TrimSpace(matches[i])

		b.React(m, emoji[match])
	}
}
