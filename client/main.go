package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/korsander/reviews-bot/lib"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var u url.URL
var connection *websocket.Conn
var done chan struct{}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	serverHost := os.Getenv("BOT_HOST")

	u = url.URL{Scheme: "ws", Host: serverHost, Path: lib.CISocketPath}
	log.Printf("connecting to %s", u.String())

	tryConnect(u)
	defer connection.Close()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			connection = nil
		case <-ticker.C:
			if connection == nil {
				tryConnect(u)
				continue
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}

func tryConnect(u url.URL) {
	log.Println("try to connect")
	var err error
	connection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		return
	}

	done = make(chan struct{})
	log.Printf("connected to %s", u.Host)

	go func() {
		defer close(done)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}
