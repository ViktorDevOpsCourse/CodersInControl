package services

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/delivery"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/jobs"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"sync"
	"time"
)

type JobDispatcher struct {
	k8sService        *clusters.K8S
	jobs              sync.Map
	actionReceiver    chan *bot.BotAction
	appsStatesStorage storage.StateRepository
	appsEventsStorage storage.EventsRepository
	github            *delivery.OpsRepo
}

func NewJobDispatcher(k8sService *clusters.K8S,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository,
	github *delivery.OpsRepo) JobDispatcher {
	return JobDispatcher{
		k8sService:        k8sService,
		actionReceiver:    make(chan *bot.BotAction),
		appsStatesStorage: appsStatesStorage,
		appsEventsStorage: appsEventsStorage,
		github:            github,
	}
}

// Run launch dispatch jobs. Blocked operation
func (d *JobDispatcher) Run() {
	log := logger.FromDefaultContext()
	for {
		botAction, ok := <-d.actionReceiver
		if !ok {
			log.Fatal("Job dispatcher someone closed active receiver chanel")
			return
		}

		j, err := jobs.NewJob(botAction, d.appsStatesStorage, d.appsEventsStorage, d.k8sService.GetClustersCopy(), d.github)
		if err != nil {
			botAction.ResponseOnAction(fmt.Sprintf("Something went wrong with command `%s`. Error %s",
				botAction.GetRawCommand(), err))
			continue
		}

		if d.isJobExist(j.GetId()) {
			botAction.ResponseOnAction("action already processing, wait please")
			continue
		}

		go d.proceedJob(j)

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
	log := logger.FromDefaultContext()

	d.jobs.Store(job.GetId(), true)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	jobDone := make(chan bool)
	go job.Launch(ctx, jobDone)

	defer d.jobs.Delete(job.GetId())

	select {
	case <-ctx.Done():
		job.ResponseToBot("timeout exceeded")
	case <-jobDone:
		log.Infof("job %s done", job.GetId())
	}
}

func (d *JobDispatcher) GetActionQueueReceiver() chan *bot.BotAction {
	return d.actionReceiver
}

func (d *JobDispatcher) Register() chan *bot.BotAction {
	return d.actionReceiver
}
