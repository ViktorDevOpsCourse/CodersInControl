package jobs

import (
	"context"
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

var re = regexp.MustCompile(`(\S+)@(\S+)\s+to\s+(\w+)`)

type Job interface {
	GetId() string
	Launch(context.Context, chan bool)
	ResponseToBot(string)
}

func NewJob(botAction *bot.BotAction,
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
		matches := re.FindStringSubmatch(botAction.GetCommandArgs())
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
	}
	return nil, fmt.Errorf("unknown job command `%s`", botAction.GetCommand())
}
