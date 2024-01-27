package services

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/jobs"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

type JobDispatcher struct {
	k8sService  *clusters.K8S
	jobs        map[string]jobs.Job
	actionQueue chan actions.Action
}

func NewJobDispatcher(k8sService *clusters.K8S) JobDispatcher {
	return JobDispatcher{
		k8sService:  k8sService,
		jobs:        make(map[string]jobs.Job),
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
			botAction.ResponseOnAction("action already processing, wait please")
			continue
		}

		go func(botAction actions.Action) {
			cluster, err := d.k8sService.GetCluster(botAction.GetEnvironment())
			if err != nil {
				botAction.ResponseOnAction(fmt.Sprintf("Cluster for env `%s` not found", botAction.GetEnvironment()))
				return
			}

			j, err := jobs.NewJob(botAction, cluster)
			if err != nil {
				botAction.ResponseOnAction(fmt.Sprintf("Undefined command `%s`", botAction.GetCommand()))
				return
			}
			j.Launch()
		}(botAction)
	}
}

func (d *JobDispatcher) GetActionProcessorQueue() chan actions.Action {
	return d.actionQueue
}
