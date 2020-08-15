package main

import (
	"fmt"
	"github.com/joho/godotenv"
	rtm "github.com/korsander/reviews-bot/bot/rtm"
	"github.com/korsander/reviews-bot/lib"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {
	loadEnvironmentVars()
	startListenSlackMessages()
	handleCISocket()
}

func loadEnvironmentVars() {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}
}

func startListenSlackMessages() {
	slackToken := os.Getenv("SLACK_TOKEN")
	api := slack.New(
		slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	go func() {
		rtm.StartMessagesHandle(api)
	}()
}

func handleCISocket() {
	println(lib.CISocketPath)
}
