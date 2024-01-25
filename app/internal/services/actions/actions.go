package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
	"regexp"
)

// TODO read from config command and envs
var re = regexp.MustCompile(`<@(\S+)> (promote|list|diff|rollback) (\S+)@(\S+) to (stage|qa|prod)`)

type ActionMetadata interface {
	GetName() string
	GetServiceName() string
	GetBuildTag() string
	GetEnvironment() string
}

type Action interface {
	ActionMetadata

	GetActionID() string
	ResponseOnAction(message string)
}

type ActionEvent struct {
	Name        string
	Service     string
	BuildTag    string
	Environment string
	rawEvent    *slackevents.AppMentionEvent
}

func NewAction(event *slackevents.AppMentionEvent, callback func(channel, message, messageTimestamp string)) (Action, error) {
	actionEvent, err := parseArgs(event)
	if err != nil {
		callback(event.Channel, err.Error(), event.ThreadTimeStamp)
		return nil, err
	}

	actionEvent.rawEvent = event

	return createAction(actionEvent, callback), nil
}

func (e *ActionEvent) validate() error {
	if e.Name == "" {
		return fmt.Errorf("unknown command")
	}

	if e.Service == "" {
		return fmt.Errorf("unknown service")
	}

	if e.BuildTag == "" {
		return fmt.Errorf("servie build not specified")
	}

	if e.Environment == "" {
		return fmt.Errorf("invalid environment")
	}

	return nil
}

type BotAction struct {
	event            *ActionEvent
	callbackResponse func(channel, message, messageTimestamp string)
}

func createAction(event *ActionEvent, callback func(channel, message, messageTimestamp string)) Action {
	return &BotAction{
		event:            event,
		callbackResponse: callback,
	}
}

func (p *BotAction) GetActionID() string {
	return fmt.Sprintf("%s:%s:%s:%s", p.event.Name, p.event.Service, p.event.BuildTag, p.event.Environment)
}

func (p *BotAction) ResponseOnAction(message string) {
	p.callbackResponse(p.event.rawEvent.Channel, message, p.event.rawEvent.ThreadTimeStamp)
}

func (p *BotAction) GetName() string {
	return p.event.Name
}

func (p *BotAction) GetServiceName() string {
	return p.event.Service
}

func (p *BotAction) GetBuildTag() string {
	return p.event.BuildTag
}

func (p *BotAction) GetEnvironment() string {
	return p.event.Environment
}

func parseArgs(event *slackevents.AppMentionEvent) (*ActionEvent, error) {
	matches := re.FindStringSubmatch(event.Text)
	a := &ActionEvent{
		Name:        matches[2],
		Service:     matches[3],
		BuildTag:    matches[4],
		Environment: matches[5],
	}
	err := a.validate()
	if err != nil {
		return nil, err
	}
	return a, nil
}
