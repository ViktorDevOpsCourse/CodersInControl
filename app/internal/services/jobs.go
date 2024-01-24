package services

import (
	"fmt"
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
			botAction.SendResponse("action already processing waite please")
			botAction.Done()
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

type Job struct {
	botAction actions.Action
}

func NewJob(botAction actions.Action) Job {
	return Job{
		botAction: botAction,
	}
}

func (j *Job) Launch() {
	// TODO
	// TODO call to github
	// TODO watch on k8s services
	// TODO send response in bot chanel
	j.botAction.SendResponse(fmt.Sprintf("Job `%s` in chanel `%s`", j.botAction.GetName(), j.botAction.GetName()))
}
