package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"regexp"
)

const (
	jobList     = "list"
	JobDiff     = "diff"
	JobPromote  = "promote"
	JobRollBack = "rollback"
)

var rePromote = regexp.MustCompile(`(\S+)@(\S+)\s+to\s+(\w+)`)
var reRollBack = regexp.MustCompile(`(\S+)\s+on\s+(\w+)`)

type Job interface {
	GetId() string
	Launch(context.Context, chan bool)
	ResponseToBot(string)
}

func NewJob(botAction *bot.BotAction,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository,
	clusters map[string]clusters.Cluster) (Job, error) {

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
		matches := rePromote.FindStringSubmatch(botAction.GetCommandArgs())
		// TODO validate matches
		return &PromoteJob{
			AppName:           matches[1],
			BuildTag:          matches[2],
			Environment:       matches[3],
			botAction:         botAction,
			clusters:          clusters,
			appsEventsStorage: appsEventsStorage,
		}, nil
	case JobRollBack:
		matches := reRollBack.FindStringSubmatch(botAction.GetCommandArgs())
		// TODO validate matches
		// TODO overview store state or flow because we can only 1 success state
		prevAppState, err := appsStatesStorage.GetLastSuccessState(matches[2], matches[1])
		if err != nil {
			if errors.Is(err, storage.NotFoundError) {
				return nil, fmt.Errorf("do not find previous service version in database :grimacing: ")
			}
			return nil, fmt.Errorf("error occured while processing rollback: `%s`", err)
		}

		return &PromoteJob{
			AppName:           matches[1],
			BuildTag:          prevAppState.Image,
			Environment:       matches[2],
			botAction:         botAction,
			clusters:          clusters,
			appsEventsStorage: appsEventsStorage,
		}, nil
	}
	return nil, fmt.Errorf("unknown job command `%s`", botAction.GetCommand())
}
