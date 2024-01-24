package actions

import (
	"fmt"
)

type PromoteAction struct {
	event            *ActionEvent
	callbackResponse func(channel, message, messageTimestamp string)
}

func CreatePromoteAction(event *ActionEvent, callback func(channel, message, messageTimestamp string)) Action {
	return &PromoteAction{
		event:            event,
		callbackResponse: callback,
	}
}

func (p *PromoteAction) Run() {

}

func (p *PromoteAction) GetActionID() string {
	return fmt.Sprintf("%s:%s:%s:%s", p.event.Name, p.event.Service, p.event.BuildTag, p.event.Environment)
}

func (p *PromoteAction) Done() {

}

func (p *PromoteAction) SendResponse(message string) {
	p.callbackResponse(p.event.botEvent.Channel, message, p.event.botEvent.ThreadTimeStamp)
}

func (p *PromoteAction) GetName() string {
	return p.event.Name
}
