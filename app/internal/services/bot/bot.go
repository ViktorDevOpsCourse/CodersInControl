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
	ctx                  context.Context
	client               *Client
	actionProcessorQueue chan<- *BotAction
	auth                 *Auth
}

func NewSlackBot(ctx context.Context, options SlackOptions) *SlackBot {
	c := NewClient(ctx, options.ClientOptions)
	return &SlackBot{
		ctx:                  ctx,
		client:               c,
		actionProcessorQueue: options.BotOptions.ActionProcessorQueue,
		auth:                 NewAuth(c.api, options.AuthOptions),
	}
}

func (s *SlackBot) Run() {
	go func() {
		// connect to slack bot via socket
		err := s.client.event.Run()
		if err != nil {
			panic(err)
		}
	}()

	go s.listenEvents()
}

func (s *SlackBot) listenEvents() {
	for event := range s.client.event.Events {
		s.handleEvents(event)
	}
}

func (s *SlackBot) handleEvents(event socketmode.Event) {
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
		s.handleEventTypeEventsAPI(event)
	default:
		log.Errorf("Unexpected event type received: %s", event.Type)
	}
}

func (s *SlackBot) handleEventTypeEventsAPI(event socketmode.Event) {
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return
	}

	s.client.event.Ack(*event.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:

		switch innerEventData := eventsAPIEvent.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// replay to thread
			if innerEventData.ThreadTimeStamp == "" {
				innerEventData.ThreadTimeStamp = innerEventData.TimeStamp
			}
			s.appMentionEventHandler(innerEventData)
		default:
			s.client.event.Debugf("unsupported Events API event received: %v", eventsAPIEvent.Type)
		}
	}
}

func (s *SlackBot) appMentionEventHandler(event *slackevents.AppMentionEvent) {
	log := logger.FromContext(s.ctx)

	isAllow, err := s.auth.hasPermissions(event.User)
	if err != nil {
		log.Error(err)
		s.callBackMessage(event.Channel, "something went wrong", event.ThreadTimeStamp)
		return
	}

	if !isAllow {
		log.Infof("user `%s` do not have permissions", event.User)
		s.callBackMessage(event.Channel, "nice try :stuck_out_tongue_closed_eyes:, you do not have permissions :no_entry::police_car:", event.ThreadTimeStamp)
		return
	}

	act, err := NewAction(event, s.callBackMessage)
	if err != nil {
		log.Errorf("Failed create action `%s`", err)
		return
	}

	s.actionProcessorQueue <- act
}

func (s *SlackBot) callBackMessage(channel, message, messageTimestamp string) {
	log := logger.FromDefaultContext()
	_, _, err := s.client.api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionTS(messageTimestamp))
	if err != nil {
		log.Errorf("Failed response on appMentionEventHandler. Error `%s`", err)
	}
}
