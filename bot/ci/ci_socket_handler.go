package ci

import (
	"github.com/gorilla/websocket"
	"github.com/korsander/reviews-bot/lib"
	"log"
	"net/http"
	"os"
)

func HandleCISocket() {
	http.HandleFunc(lib.CISocketPath, handleCiSocket)
	ciAddr := os.Getenv("CI_ADDR")
	fullChain := os.Getenv("CERT_FULL_CHAIN")
	privKey := os.Getenv("CERT_PRIVATE_KEY")
	log.Fatal(http.ListenAndServeTLS(ciAddr, fullChain, privKey, nil))
}

func SendCommandToCI(command string) {
	err := connection.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("error sending:", err)
	}
}

var upgrader = websocket.Upgrader{}
var connection *websocket.Conn

func handleCiSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("\n[INFO] Listen for ci")
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
