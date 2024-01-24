package actions

import (
	"fmt"
	"github.com/slack-go/slack/slackevents"
)

type PromoteAction struct {
	event            *slackevents.AppMentionEvent
	callbackResponse chan<- string
}

func CreatePromoteAction(event *slackevents.AppMentionEvent, callbackResponse chan<- string) Action {
	return &PromoteAction{
		event:            event,
		callbackResponse: callbackResponse,
	}
}

func (p *PromoteAction) Run() {

}

func (p *PromoteAction) GetActionID() string {
	return fmt.Sprintf("")
}

func (p *PromoteAction) Done() {
	defer close(p.callbackResponse)
}

func (p *PromoteAction) SendResponse(message string) {
	go func() {
		p.callbackResponse <- message
	}()
}

func (p *PromoteAction) GetName() string {
	// TODO send command name
	return p.event.Channel
}
