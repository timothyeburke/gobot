package main

import (
	"log"
	"os"

	"github.com/timothyeburke/gobot/bot"
)

func main() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	apiToken := os.Getenv("SLACK_API_TOKEN")

	slack, err := bot.New(botToken, apiToken)
	if err != nil {
		log.Fatal(err)
	}

	slack.Listen()
}
