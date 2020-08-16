package events

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

func HandleCallbackEvent(w http.ResponseWriter, api *slack.Client, event *slackevents.AppMentionEvent) error {
	var err error
	_, _, err = api.PostMessage(
		event.Channel,
		slack.MsgOptionText("Would you like to build app?", false),
		getSelectPlatformMsgOptions(),
	)
	return err
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
