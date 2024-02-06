package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/delivery"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
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
	ResponseToBot(string)
}

func NewJob(botAction *bot.BotAction,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository,
	clusters clusters.ClustersCopy,
	appUpdater delivery.Updater) (Job, error) {

	switch botAction.GetCommand() {
	case jobList:
		return NewListJob(botAction, clusters)
	case JobDiff:
		return NewDiffJob(botAction, clusters)
	case JobPromote:
		return NewPromoteJob(
			botAction,
			appsStatesStorage,
			appsEventsStorage,
			clusters,
			appUpdater)
	case JobRollBack:
		return NewRollBackJob(
			botAction,
			appsStatesStorage,
			appsEventsStorage,
			clusters,
			appUpdater)
	}
	return nil, fmt.Errorf("unknown job command `%s`", botAction.GetCommand())
}
