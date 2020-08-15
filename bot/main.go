package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	rtm "github.com/korsander/reviews-bot/bot/rtm"
	"github.com/korsander/reviews-bot/lib"
	"github.com/slack-go/slack"
	"log"
	"net/http"
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
	http.HandleFunc(lib.CISocketPath, echo)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
