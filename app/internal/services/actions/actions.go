package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^\<\@(?P<botUserId>\w*)\>\s+(?P<Command>\w+)(?P<RawArgs>\s+.*)?$`)

type ReceivedEvent struct {
	Command   string
	RawArgs   string
	ChannelId string
	rawEvent  *slackevents.AppMentionEvent
}

type BotAction struct {
	event            *ReceivedEvent
	callbackResponse func(channel, message, messageTimestamp string)
}

func NewAction(event *slackevents.AppMentionEvent, callback func(channel, message, messageTimestamp string)) (*BotAction, error) {
	actionEvent, err := parseArgs(event)
	if err != nil {
		callback(event.Channel, err.Error(), event.ThreadTimeStamp)
		return nil, err
	}

	actionEvent.rawEvent = event

	return &BotAction{
		event:            actionEvent,
		callbackResponse: callback,
	}, nil
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

func parseArgs(event *slackevents.AppMentionEvent) (*ReceivedEvent, error) {
	match := re.FindStringSubmatch(event.Text)
	result := map[string]string{}
	for keyIndex, value := range match {
		if keyIndex > 0 {
			result[re.SubexpNames()[keyIndex]] = strings.TrimSpace(value)
		}
	}
	a := &ReceivedEvent{
		Command:   result["Command"],
		RawArgs:   result["RawArgs"],
		ChannelId: event.Channel,
	}
	return a, nil
}
