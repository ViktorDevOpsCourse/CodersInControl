package services

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/k8s"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

type JobDispatcher struct {
	k8sService  *k8s.K8S
	jobs        map[string]Job
	actionQueue chan actions.Action
}

func NewJobDispatcher(k8sService *k8s.K8S) JobDispatcher {
	return JobDispatcher{
		k8sService:  k8sService,
		jobs:        make(map[string]Job),
		actionQueue: make(chan actions.Action),
	}
}

func (d *JobDispatcher) Run() {
	log := logger.FromDefaultContext()
	for {
		botAction, ok := <-d.actionQueue
		if !ok {
			log.Fatal("Job dispatcher someone closed active receiver chanel")
			return
		}

		if _, ok := d.jobs[botAction.GetActionID()]; ok {
			botAction.SendResponse("action already processing, wait please")
			continue
		}

		go func(botAction actions.Action) {
			job := NewJob(botAction)
			job.Launch()
		}(botAction)
	}
}

func (d *JobDispatcher) GetActionProcessorQueue() chan actions.Action {
	return d.actionQueue
}
