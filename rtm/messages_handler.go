package rtm

import (
	"fmt"
	"github.com/slack-go/slack"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func StartMessagesHandle(api *slack.Client) {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			handleMessage(rtm, ev)

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func handleMessage(rtm *slack.RTM, e *slack.MessageEvent) {
	if mustSkip(e.Msg.SubType) {
		return
	}

	ch := e.Msg.Channel

	fmt.Printf("Message: %v\n", e)
	if containsAppeal(e.Msg.Text) {
		rtm.SendMessage(rtm.NewOutgoingMessage("I am!", ch))
	}
}

func mustSkip(subtype string) bool {
	switch subtype {
	case "bot_message":
		return true
	case "channel_join":
		return true
	default:
		return false
	}
}

func containsAppeal(msg string) bool {
	for _, appeal := range appeals {
		start, _ := russianSearch.IndexString(msg, appeal)
		if start >= 0 {
			return true
		}
	}

	return false
}

var appeals = [...]string{"бот", "робот", "мастер"}
var russianSearch = search.New(language.Russian, search.IgnoreCase)
