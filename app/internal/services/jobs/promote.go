package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/delivery"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"regexp"
	"strings"
	"time"
)

var rePromote = regexp.MustCompile(`(\S+)@(\S+)\s+to\s+(\w+)`)

type UpdateAppJob struct {
	botAction          *bot.BotAction
	AppName            string
	BuildTag           string
	ClusterName        string
	appsStatesStorage  storage.StateRepository
	appsEventsStorage  storage.EventsRepository
	clusters           clusters.ClustersCopy
	ApplicationUpdater delivery.Updater
	currentAppState    clusters.Application
}

func NewUpdateAppJob(
	botAction *bot.BotAction,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository,
	clusters clusters.ClustersCopy,
	appUpdater delivery.Updater) (*UpdateAppJob, error) {

	matches := rePromote.FindStringSubmatch(botAction.GetRawCommand())
	err := isValidPromote(matches)
	if err != nil {
		return nil, fmt.Errorf("failed processing command with args `%s`. Reason `%s`", botAction.GetRawCommand(), err)
	}

	return &UpdateAppJob{
		AppName:            matches[1],
		BuildTag:           matches[2],
		ClusterName:        matches[3],
		botAction:          botAction,
		clusters:           clusters,
		appsStatesStorage:  appsStatesStorage,
		appsEventsStorage:  appsEventsStorage,
		ApplicationUpdater: appUpdater,
	}, nil
}

func (p *UpdateAppJob) Launch(ctx context.Context, jobDone chan bool) {
	log := logger.FromContext(ctx)
	defer func() {
		jobDone <- true
	}()

	p.ResponseToBot(fmt.Sprintf("image: `%s` promoting :runner:", p.BuildTag))

	cluster, err := p.clusters.GetCluster(p.ClusterName)
	if err != nil {
		p.ResponseToBot(fmt.Sprintf("error occured while get cluster `%s`", p.ClusterName))
		return
	}

	p.currentAppState = cluster.GetApplicationByName(p.AppName)

	if strings.Contains(p.currentAppState.Image, p.BuildTag) {
		p.ResponseToBot(fmt.Sprintf("image: `%s` already promoted on %s", p.BuildTag, p.ClusterName))
		return
	}

	err = p.ApplicationUpdater.Update(delivery.Application{
		FilePath: fmt.Sprintf("apps/%s/%s-values.yaml", p.ClusterName, p.AppName),
		Version:  p.BuildTag,
	})
	if err != nil {
		log.Errorf("update image version `%s` failed. Err %s ", p.BuildTag, err)
		p.ResponseToBot(fmt.Sprintf("image: `%s` promoted on %s. github update image version `%[1]s` failed. Err %[3]s ", p.BuildTag, p.ClusterName, err))
		return
	}

	p.waitDeploymentEvent(ctx)

}

func (p *UpdateAppJob) GetId() string {
	return fmt.Sprintf("%s%s%s", p.botAction.GetCommand(), p.AppName, p.ClusterName)
}

func (p *UpdateAppJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}

func (p *UpdateAppJob) waitDeploymentEvent(ctx context.Context) {
	log := logger.FromContext(ctx)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			log.Infof("check promote finished events for %s in %s", p.AppName, p.ClusterName)
			appEvent, err := p.appsEventsStorage.GetAndRemove(p.ClusterName, p.AppName)
			if err != nil {
				if errors.Is(err, storage.NotFoundError) {
					continue
				}
				log.Errorf("Promote Job app `%s`. failed get app event from storage. Err `%s`", p.AppName, err)
			}

			if !strings.Contains(appEvent.Image, p.BuildTag) {
				log.Errorf("found some other event %#v in event storage. for job %#v", appEvent, p)
				continue
			}

			log.Infof("find promote finished event for %s in %s. Status `%s`", p.AppName, p.ClusterName, appEvent.Status)
			if clusters.Status(appEvent.Status) != clusters.RunningStatus {
				p.ResponseToBot(fmt.Sprintf("image: `%s` promoted on %s with status %s", p.BuildTag, p.ClusterName, appEvent.Status))
				return
			}

			// TODO can be do it beautiful
			if p.botAction.GetCommand() == "promote" {
				err = p.appsStatesStorage.Save(p.ClusterName, p.currentAppState.GetName(), storage.State{
					Image: p.currentAppState.Image,
				})
				if err != nil {
					log.Error(err)
				}
			}

			p.ResponseToBot(fmt.Sprintf("image: `%s` promoted on %s. Status %s :tada:", p.BuildTag, p.ClusterName, appEvent.Status))
			return

		case <-ctx.Done():
			return
		}
	}
}

func isValidPromote(matches []string) error {
	if len(matches) < 3 {
		return fmt.Errorf("invalid command. Accept `@bot promote service@version to environment`")
	}
	if matches[1] == "" {
		return fmt.Errorf("invalid application name")
	}
	if matches[2] == "" {
		return fmt.Errorf("invalid application build tag")
	}
	if matches[3] == "" {
		return fmt.Errorf("invalid application environment")
	}

	return nil
}
