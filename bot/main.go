package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/korsander/reviews-bot/bot/ci"
	"github.com/korsander/reviews-bot/bot/events"
	"github.com/slack-go/slack"
	"os"
)

func main() {
	loadEnvironmentVars()
	startListenSlackMessages()
	ci.HandleCISocket()
}

func loadEnvironmentVars() {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}
}

func startListenSlackMessages() {
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)

	go events.StartEventsHandle(api)
}
