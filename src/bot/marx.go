package bot

import (
	"math/rand"
	"regexp"
)

func Marx(b *Bot, m *Message, u *User) {
    if !regexp.MustCompile(`(?i)(marx)`).MatchString(m.Message) {
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
		// Groucho
		"Outside of a dog, a book is man's best friend. Inside of a dog it's too dark to read.",
		"Those are my principles, and if you don't like them... well, I have others.",
		"I sent the club a wire stating, PLEASE ACCEPT MY RESIGNATION. I DON'T WANT TO BELONG TO ANY CLUB THAT WILL ACCEPT ME AS A MEMBER.",
		"Politics is the art of looking for trouble, finding it everywhere, diagnosing it incorrectly and applying the wrong remedies.",
		"Time flies like an arrow; fruit flies like a banana",
		"I never forget a face, but in your case I'll be glad to make an exception.",
		"Either this man is dead or my watch has stopped.",
		"I find television very educating. Every time somebody turns on the set, I go into the other room and read a book.",
		"I've had a perfectly wonderful evening. But this wasn't it.",
		"One morning I shot an elephant in my pajamas. How he got into my pajamas I'll never know.",
		// Karl
		"The need of a constantly expanding market for its products chases the bourgeoisie over the whole surface of the globe. It must nestle everywhere, settle everywhere, establish connexions everywhere.",
		"The road to Hell is paved with good intentions.",
		"Under the ideal measure of values there lurks the hard cash.",
		"Scientific truth is always paradox, if judged by everyday experience, which catches only the delusive appearance of things.",
	}
	m.Message = answers[rand.Intn(len(answers))]
	b.Send(m)
}
