package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
)

type PromoteAction struct {
	event            *slackevents.AppMentionEvent
	callbackResponse func(channel, message, messageTimestamp string)
}

func CreatePromoteAction(event *slackevents.AppMentionEvent, callback func(channel, message, messageTimestamp string)) Action {
	return &PromoteAction{
		event:            event,
		callbackResponse: callback,
	}
}

func (p *PromoteAction) Run() {

}

func (p *PromoteAction) GetActionID() string {
	return fmt.Sprintf("")
}

func (p *PromoteAction) Done() {

}

func (p *PromoteAction) SendResponse(message string) {
	p.callbackResponse(p.event.Channel, message, p.event.ThreadTimeStamp)
}

func (p *PromoteAction) GetName() string {
	// TODO send command name
	return p.event.Channel
}
