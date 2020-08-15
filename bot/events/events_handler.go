package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"log"
	"net/http"
	"os"
)

var verificationToken = os.Getenv("VERIFICATION_TOKEN")
var fullChain = os.Getenv("CERT_FULL_CHAIN")
var privKey = os.Getenv("CERT_PRIVATE_KEY")
var eventsAddr = os.Getenv("EVENTS_ADDR")

func StartEventsHandle(api *slack.Client) {
	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming request")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		eventsAPIEvent, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error parse: %s\n", err)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			fmt.Println("[INFO] challenge event")
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			_, err = w.Write([]byte(r.Challenge))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}
	})
	log.Println("\n[INFO] Server listening")

	err := http.ListenAndServeTLS(eventsAddr, fullChain, privKey, nil)
	if err != nil {
		log.Printf("error listen: %s", err)
	}
}
