package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/korsander/reviews-bot/bot/cfg"
	"github.com/korsander/reviews-bot/bot/ci"
	"github.com/korsander/reviews-bot/bot/service"
	"github.com/slack-go/slack"
	"net/http"
)

func main() {
	config := cfg.LoadConfig()
	api := slack.New(config.SlackToken)

	startSlackService(config, api)

	ci.HandleCISocket()
}

func startSlackService(config cfg.Config, api *slack.Client) {
	router := mux.NewRouter()
	slackService := service.NewSlackService(
		config,
		api,
	)

	slackService.Mount(router)

	go func() {
		if err := http.ListenAndServeTLS(
			config.EventsAddr,
			config.CertChain,
			config.CertPrivate,
			router,
		); err != nil {
			fmt.Printf("Server stopped immediately: %v", err)
		}
	}()
}
