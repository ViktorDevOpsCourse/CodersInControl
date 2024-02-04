package jobs

import (
	"errors"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/delivery"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"regexp"
)

var reRollBack = regexp.MustCompile(`(\S+)\s+on\s+(\w+)`)

func NewRollBackJob(
	botAction *bot.BotAction,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository,
	clusters clusters.ClustersCopy,
	appUpdater delivery.Updater) (*UpdateAppJob, error) {

	matches := reRollBack.FindStringSubmatch(botAction.GetRawCommand())
	err := isValidRollback(matches)
	if err != nil {
		return nil, fmt.Errorf("failed processing command with args `%s`. Reason `%s`", botAction.GetRawCommand(), err)
	}

	prevAppState, err := appsStatesStorage.GetLastSuccessState(matches[2], matches[1])
	if err != nil {
		if errors.Is(err, storage.NotFoundError) {
			return nil, fmt.Errorf("do not find previous service version in database :grimacing: ")
		}
		return nil, fmt.Errorf("error occured while processing rollback: `%s`", err)
	}

	return &UpdateAppJob{
		AppName:            matches[1],
		BuildTag:           prevAppState.Image,
		ClusterName:        matches[2],
		botAction:          botAction,
		clusters:           clusters,
		appsEventsStorage:  appsEventsStorage,
		ApplicationUpdater: appUpdater,
	}, nil
}

func isValidRollback(matches []string) error {
	if len(matches) < 2 {
		return fmt.Errorf("invalid command. Accept `@bot promote service@version to environment`")
	}
	if matches[1] == "" {
		return fmt.Errorf("invalid application name")
	}
	if matches[2] == "" {
		return fmt.Errorf("invalid application environment")
	}

	return nil
}
