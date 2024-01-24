package bot

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"regexp"
)

var re = regexp.MustCompile(`<@(\S+)> (promote|list|diff|rollback) (\S+)@(\S+) to (stage|qa|prod)`)

type Bot interface {
}

type SlackBot struct {
	ctx         context.Context
	client      *Client
	actionQueue chan<- actions.Action
}

func NewSlackBot(ctx context.Context, options SlackOptions) *SlackBot {
	return &SlackBot{
		ctx:         ctx,
		client:      NewClient(ctx, options.ClientOptions),
		actionQueue: options.BotOptions.ActionProcessorQueue,
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
			s.appMentionEventHandler(innerEventData)
		default:
			s.client.event.Debugf("unsupported Events API event received: %v", eventsAPIEvent.Type)
		}
	}
}

func (s *SlackBot) appMentionEventHandler(event *slackevents.AppMentionEvent) {
	// TODO validate access for user
	log := logger.FromContext(s.ctx)

	if event.ThreadTimeStamp == "" {
		event.ThreadTimeStamp = event.TimeStamp
	}
	ts := event.ThreadTimeStamp

	command := parseEventText(event.Text)
	err := validateEventCommand(command)
	if err != nil {
		_, _, err = s.client.api.PostMessage(event.Channel, slack.MsgOptionText(err.Error(), false), slack.MsgOptionTS(ts))
		if err != nil {
			log.Errorf("Failed response on appMentionEventHandler. Error `%s`", err)
		}
		return
	}

	s.actionQueue <- actions.CreateAction(command["Type"], event, s.callBackMessage)
}

func (s *SlackBot) callBackMessage(channel, message, messageTimestamp string) {
	log := logger.FromDefaultContext()
	_, _, err := s.client.api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionTS(messageTimestamp))
	if err != nil {
		log.Errorf("Failed response on appMentionEventHandler. Error `%s`", err)
	}
}

func validateEventCommand(command map[string]string) error {
	if command["Type"] == "" {
		return fmt.Errorf("unknown command")
	}

	if command["Service"] == "" {
		return fmt.Errorf("unknown service")
	}

	if command["Build"] == "" {
		return fmt.Errorf("servie build not specified")
	}

	if command["Environment"] == "" {
		return fmt.Errorf("invalid environment")
	}

	return nil
}

func parseEventText(input string) map[string]string {

	matches := re.FindStringSubmatch(input)
	return map[string]string{
		"Type":        matches[2],
		"Service":     matches[3],
		"Build":       matches[4],
		"Environment": matches[5],
	}
}
