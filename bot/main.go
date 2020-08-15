package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/korsander/reviews-bot/bot/events"
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
	signingSecret := os.Getenv("SIGNING_SECRET")
	api := slack.New(
		slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	go events.StartEventsHandle(api, signingSecret)
}

func handleCISocket() {
	http.HandleFunc(lib.CISocketPath, handleCiSocket)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

var upgrader = websocket.Upgrader{}
var connection *websocket.Conn

func handleCiSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	connection, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
}

func SendCommandToCI(command string) {
	err := connection.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("error sending:", err)
	}
}
