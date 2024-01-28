package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
	"regexp"
	"strings"
)

// TODO read from config command and envs
var re = regexp.MustCompile(`^\<\@(?P<botUserId>\w*)\>\s+(?P<Command>\w+)(?P<RawArgs>\s+.*)?$`)

type Action interface {
	GetCommand() string
	GetCommandArgs() string
	GetRawCommand() string
	ResponseOnAction(message string)
}

type ActionEvent struct {
	Command  string
	RawArgs  string
	rawEvent *slackevents.AppMentionEvent
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

func (p *BotAction) ResponseOnAction(message string) {
	p.callbackResponse(p.event.rawEvent.Channel, message, p.event.rawEvent.ThreadTimeStamp)
}

func (p *BotAction) GetCommand() string {
	return p.event.Command
}

func (p *BotAction) GetCommandArgs() string {
	return p.event.RawArgs
}

func (p *BotAction) GetRawCommand() string {
	return fmt.Sprintf("%s %s", p.event.Command, p.event.RawArgs)
}

func parseArgs(event *slackevents.AppMentionEvent) (*ActionEvent, error) {
	match := re.FindStringSubmatch(event.Text)
	result := map[string]string{}
	for keyIndex, value := range match {
		if keyIndex > 0 {
			result[re.SubexpNames()[keyIndex]] = strings.TrimSpace(value)
		}
	}
	a := &ActionEvent{
		Command: result["Command"],
		RawArgs: result["RawArgs"],
	}
	return a, nil
}
