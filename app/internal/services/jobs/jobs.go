package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

const (
	jobList     = "list"
	JobDiff     = "diff"
	JobPromote  = "promote"
	JobRollBack = "rollback"
)

type Job interface {
	GetId() string
	Launch(context.Context, chan bool)
	ResponseToBot(message string)
}

func NewJob(botAction actions.Action, clusters map[string]clusters.Cluster) (Job, error) {

	switch botAction.GetCommand() {
	case jobList:
		return &ListJob{
			botAction: botAction,
			clusters:  clusters,
		}, nil
	case JobDiff:
		currentEnv := botAction.GetCommandArgs()
		if _, ok := clusters[currentEnv]; !ok {
			return nil, fmt.Errorf("invalid command or unknow environment. Accept `@bot diff environment`")
		}
		return &DiffJob{
			botAction:  botAction,
			clusters:   clusters,
			currentEnv: botAction.GetCommandArgs(),
		}, nil
	case JobPromote:
		return &PromoteJob{
			botAction: botAction,
			//cluster:   clusters,
		}, nil
	case JobRollBack:
	}
	return nil, fmt.Errorf("unknown job command `%s`", botAction.GetCommand())
}
