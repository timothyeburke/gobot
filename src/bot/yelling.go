package bot

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"log"

	"github.com/mb-14/gomarkov"
)

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Yelling(b *Bot, m *Message, u *User) {
	if !regexp.MustCompile(`(^[A-Z][^a-z]+$)`).MatchString(m.Message) {
		return
	}

	var allowedRooms = map[string]bool {
		"just-yelling-training": true,
		"just-yelling": true,
	}
	var Channel, _ = b.GetChannel(m.Channel)
	if !allowedRooms[Channel.Name] {
		return
	}

	var yells []string
	var result = b.Brain.Get("yells")
	json.Unmarshal(result, &yells)

	_, found := find(yells, m.Message)
	if !found {
		log.Println("NEW YELL FOUND!")
		log.Println(m.Message)
		j, _ := json.Marshal(append(yells, m.Message))
		b.Brain.Save("yells", j)
	}

	rand.Seed(time.Now().Unix())

	const alwaysMarkov = true

	var hasPunctuation = regexp.MustCompile(`[!?.]\s*$`)
	if (alwaysMarkov || hasPunctuation.MatchString(m.Message)) {
		chain := gomarkov.NewChain(1)

		rand.Shuffle(len(yells), func(i, j int) { yells[i], yells[j] = yells[j], yells[i] })
		for _, phrase := range yells {
			chain.Add(strings.Split(phrase, " "))
		}

		tokens := []string{gomarkov.StartToken}
		for tokens[len(tokens)-1] != gomarkov.EndToken {
			next, _ := chain.Generate(tokens[(len(tokens) - 1):])
			tokens = append(tokens, next)
		}
		m.Message = strings.Join(tokens[1:len(tokens)-1], " ")
	} else {
		m.Message = yells[rand.Intn(len(yells))]
	}

	b.Send(m)
}
