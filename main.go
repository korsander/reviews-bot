package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/korsander/reviews-bot/rtm"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	api := slack.New(
		slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	rtm.StartMessagesHandle(api)
}
