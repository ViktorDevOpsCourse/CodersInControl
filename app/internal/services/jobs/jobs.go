package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/k8s"
)

const (
	ListJobName     = "list"
	DiffJobName     = "diff"
	PromoteJobName  = "promote"
	RollbackJobName = "rollback"
)

type Job interface {
	Launch()
	ResponseToBot(message string)
}

func NewJob(botAction actions.Action, k8sService *k8s.K8S) Job {
	var job Job

	switch botAction.GetName() {
	case ListJobName:
	case DiffJobName:
	case PromoteJobName:
		job = &PromoteJob{
			botAction:  botAction,
			k8sService: k8sService,
		}
	case RollbackJobName:
	default:

	}
	return job
}
