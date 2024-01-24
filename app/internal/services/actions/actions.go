package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
	"regexp"
)

const (
	ListActionName     = "list"
	DiffActionName     = "diff"
	PromoteActionName  = "promote"
	RollbackActionName = "rollback"
)

// TODO read from config command and envs
var re = regexp.MustCompile(`<@(\S+)> (promote|list|diff|rollback) (\S+)@(\S+) to (stage|qa|prod)`)

type ActionMetadata interface {
	GetName() string
}

type Action interface {
	ActionMetadata

	Run()
	GetActionID() string
	SendResponse(message string)
	Done()
}

type ActionEvent struct {
	Name        string
	Service     string
	BuildTag    string
	Environment string
	botEvent    *slackevents.AppMentionEvent
}

func NewAction(event *slackevents.AppMentionEvent, callback func(channel, message, messageTimestamp string)) (Action, error) {
	actionEvent, err := createActionEvent(event)
	if err != nil {
		callback(event.Channel, err.Error(), event.ThreadTimeStamp)
		return nil, err
	}

	actionEvent.botEvent = event

	var action Action

	switch actionEvent.Name {
	case PromoteActionName:
		action = CreatePromoteAction(actionEvent, callback)
	}

	return action, nil
}

func createActionEvent(event *slackevents.AppMentionEvent) (*ActionEvent, error) {
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
