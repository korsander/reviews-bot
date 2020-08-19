package handlers

import (
	"github.com/korsander/reviews-bot/bot/service"
	"github.com/slack-go/slack"
	"log"
)

func BuildActionHandler(client *slack.Client) service.InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		log.Print("handle build")

		return nil
	}
}
