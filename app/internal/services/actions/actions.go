package actions

import "github.com/slack-go/slack/slackevents"

const (
	ListActionName     = "list"
	DiffActionName     = "diff"
	PromoteActionName  = "promote"
	RollbackActionName = "rollback"
)

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

func CreateAction(name string, botEvent *slackevents.AppMentionEvent, responseChan chan<- string) Action {
	var action Action
	switch name {
	case PromoteActionName:
		action = CreatePromoteAction(botEvent, responseChan)
	}

	return action
}
