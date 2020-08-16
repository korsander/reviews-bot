package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/korsander/reviews-bot/bot/cfg"
	"github.com/korsander/reviews-bot/bot/events"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

type ActionType string

const (
	SelectPlatform ActionType = "select_platform"
	SelectApp      ActionType = "select_app"
	BuildApp       ActionType = "build_app"
)

type InteractivityHandler func(payload slack.InteractionCallback, action *slack.BlockAction) error

type HandlerMap map[ActionType]InteractivityHandler

type SlackService struct {
	cfg        cfg.Config
	client     *slack.Client
	handlerMap HandlerMap
}

func NewSlackService(
	cfg cfg.Config,
	client *slack.Client,
	opts ...func(HandlerMap),
) *SlackService {
	handlerMap := make(HandlerMap)

	for _, opt := range opts {
		opt(handlerMap)
	}
	return &SlackService{
		cfg:        cfg,
		client:     client,
		handlerMap: handlerMap,
	}
}

func (s *SlackService) Mount(router *mux.Router) {
	router.HandleFunc("/slack/events", s.Forward)
	router.HandleFunc("/slack/interactive-endpoint", s.HandleInteractivity)
}

func WithActionHandler(actionType ActionType, handler InteractivityHandler) func(handlerMap HandlerMap) {
	return func(handlerMap HandlerMap) {
		handlerMap[actionType] = handler
	}
}

func (s *SlackService) HandleInteractivity(rw http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	for _, action := range payload.ActionCallback.BlockActions {
		err := s.handlerMap[ActionType(action.ActionID)](payload, action)
		if err != nil {
			fmt.Printf("No handler found received action: %v", err)
		}
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *SlackService) Forward(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, err := s.parseMessage(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.handleEvent(eventsAPIEvent, body, w)
}

func (s *SlackService) parseMessage(body string) (event slackevents.EventsAPIEvent, err error) {
	return slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: s.cfg.VerificationToken}),
	)
}

func (s *SlackService) handleEvent(event slackevents.EventsAPIEvent, body string, w http.ResponseWriter) {
	if event.Type == slackevents.URLVerification {
		s.handleVerification(body, w)
	}
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := events.HandleCallbackEvent(w, s.client, ev)
			if err != nil {
				fmt.Printf("Print error %+v", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *SlackService) handleVerification(body string, w http.ResponseWriter) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text")
	_, _ = w.Write([]byte(r.Challenge))
}

func (s *SlackService) sendMessage(channelId string, message string) error {
	_, _, err := s.client.PostMessage(channelId, slack.MsgOptionText(message, false))
	return err
}
