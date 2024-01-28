package services

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/jobs"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"sync"
	"time"
)

type JobDispatcher struct {
	k8sService     *clusters.K8S
	jobs           sync.Map
	actionReceiver chan actions.Action
}

func NewJobDispatcher(k8sService *clusters.K8S) JobDispatcher {
	return JobDispatcher{
		k8sService:     k8sService,
		actionReceiver: make(chan actions.Action),
	}
}

func (d *JobDispatcher) Run() {
	log := logger.FromDefaultContext()
	for {
		botAction, ok := <-d.actionReceiver
		if !ok {
			log.Fatal("Job dispatcher someone closed active receiver chanel")
			return
		}

		j, err := jobs.NewJob(botAction, d.k8sService.GetClustersCopy())
		if err != nil {
			botAction.ResponseOnAction(fmt.Sprintf("Undefined command `%s`. Error %s",
				botAction.GetRawCommand(), err))
			return
		}
		if d.isJobExist(j.GetId()) {
			botAction.ResponseOnAction("action already processing, wait please")
			return
		}

		d.proceedJob(j)

	}
}

func (d *JobDispatcher) isJobExist(jobId string) bool {
	_, ok := d.jobs.Load(jobId)
	if !ok {
		return false
	}

	return true
}

func (d *JobDispatcher) proceedJob(job jobs.Job) {
	d.jobs.Store(job.GetId(), job)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	jobDone := make(chan bool)
	go job.Launch(ctx, jobDone)

	select {
	case <-ctx.Done():
		job.ResponseToBot("timeout exceeded")
	case <-jobDone:
		d.jobs.Delete(job.GetId())
	}
}

func (d *JobDispatcher) GetActionQueueReceiver() chan actions.Action {
	return d.actionReceiver
}
