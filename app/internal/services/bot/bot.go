package bot

import (
	"context"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

type Bot interface {
}

type SlackBot struct {
	ctx    context.Context
	client *Client
}

func NewSlackBot(ctx context.Context, options SlackOptions) *SlackBot {
	return &SlackBot{
		ctx:    ctx,
		client: NewClient(ctx, options.ClientOptions),
	}
}

func (s *SlackBot) Run() {
	go func() {
		err := s.client.event.Run()
		if err != nil {
			panic(err)
		}
	}()

	go s.listenEvents()
}

func (s *SlackBot) listenEvents() {
	for event := range s.client.event.Events {
		s.handleEvent(event)
	}
}

func (s *SlackBot) handleEvent(event socketmode.Event) {
	log := logger.FromContext(s.ctx)

	switch event.Type {
	case socketmode.EventTypeConnecting:
		log.Info("Connecting to Slack with Socket Mode...")
	case socketmode.EventTypeConnectionError:
		log.Info("Connection failed. Retrying later...")
	case socketmode.EventTypeConnected:
		log.Info("Connected to Slack with Socket Mode.")
	case socketmode.EventTypeHello:
		log.Info("The client has successfully connected to the server.")
	case socketmode.EventTypeEventsAPI:
		eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
		if !ok {
			return
		}

		s.client.event.Ack(*event.Request)

		switch eventsAPIEvent.Type {
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch innerEventData := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				if innerEventData.ThreadTimeStamp == "" {
					innerEventData.ThreadTimeStamp = innerEventData.TimeStamp
				}
				ts := innerEventData.ThreadTimeStamp

				s.client.api.PostMessage(innerEventData.Channel, slack.MsgOptionText(":bow: I don't know how to help you", false), slack.MsgOptionTS(ts))
			default:
				s.client.event.Debugf("unsupported Events API event received: %v", eventsAPIEvent.Type)
			}
		}
	default:
		log.Errorf("Unexpected event type received: %s", event.Type)
	}
}
