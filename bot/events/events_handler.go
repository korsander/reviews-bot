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
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error parse: %s\n", err)
			return
		}
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
				_, _, err = api.PostMessage(
					ev.Channel,
					slack.MsgOptionText("Would you like to build app?", false),
					getSelectPlatformMsgOptions(),
				)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("error send: %s\n", err)
				}
			}
		}
	})

	http.HandleFunc("/slack/interactive-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error parse: %s\n", err)
			return
		}
		body := buf.String()
		log.Printf("incoming request %s\n", body)
	})

	log.Println("\n[INFO] Server listening")

	err := http.ListenAndServeTLS(eventsAddr, fullChain, privKey, nil)
	if err != nil {
		log.Printf("error listen: %s", err)
	}
}

func NewDivider() slack.Block {
	return slack.DividerBlock{
		Type: "divider",
	}
}

func NewSection(text string) slack.Block {
	return slack.SectionBlock{
		Type: "section",
		Text: &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: text,
		},
	}
}

func getSelectPlatformMsgOptions() slack.MsgOption {
	hello := NewSection("Welcome to Jenkins build bot")
	divider := NewDivider()
	prompt := NewSection("Which app would you like to build?")

	actions := slack.NewActionBlock(
		"test",
		slack.NewButtonBlockElement(
			"test1",
			"android",
			slack.NewTextBlockObject("plain_text", "Android", false, false),
		),
		slack.NewButtonBlockElement(
			"test2",
			"ios",
			slack.NewTextBlockObject("plain_text", "iOS", false, false),
		),
	)

	return slack.MsgOptionBlocks(hello, divider, prompt, actions)
}
