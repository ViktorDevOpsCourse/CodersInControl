package jobs

import (
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
	Launch()
	ResponseToBot(message string)
}

func NewJob(botAction actions.Action, cluster *clusters.Cluster) (Job, error) {

	switch botAction.GetCommand() {
	case jobList:
		return &ListJob{
			botAction: botAction,
			cluster:   cluster,
		}, nil
	case JobDiff:
	case JobPromote:
		return &PromoteJob{
			botAction: botAction,
			cluster:   cluster,
		}, nil
	case JobRollBack:
	}
	return nil, fmt.Errorf("unknown job command `%s`", botAction.GetCommand())
}
