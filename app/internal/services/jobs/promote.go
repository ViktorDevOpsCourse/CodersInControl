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
	"strings"
	"time"
)

type PromoteJob struct {
	botAction         *bot.BotAction
	AppName           string
	BuildTag          string
	Environment       string
	appsEventsStorage storage.EventsRepository
	clusters          map[string]clusters.Cluster
	repo              *delivery.OpsRepo
}

func (p *PromoteJob) Launch(ctx context.Context, jobDone chan bool) {
	log := logger.FromContext(ctx)
	p.ResponseToBot(fmt.Sprintf("image: `%s` promoting :runner:", p.BuildTag))

	if strings.Contains(p.clusters[p.Environment].Applications[p.AppName].Image, p.BuildTag) {
		p.ResponseToBot(fmt.Sprintf("image: `%s` already promoted on %s", p.BuildTag, p.Environment))
		jobDone <- true
		return
	}

	err := p.repo.UpdateImage(fmt.Sprintf("apps/%s/%s-values.yaml", p.Environment, p.AppName), p.BuildTag)
	if err != nil {
		log.Errorf("update image version `%s` failed. Err %s ", p.BuildTag, err)
		p.ResponseToBot(fmt.Sprintf("image: `%s` promoted on %s. github update image version `%[1]s` failed. Err %[3]s ", p.BuildTag, p.Environment, err))
		jobDone <- true
		return
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			log.Infof("check promote finished events for %s in %s", p.AppName, p.Environment)
			appEvent, err := p.appsEventsStorage.GetAndRemove(p.Environment, p.AppName)
			if err != nil {
				if errors.Is(err, storage.NotFoundError) {
					continue
				}
				log.Errorf("Promote Job app `%s`. failed get app event from storage. Err `%s`", p.AppName, err)
			}

			if strings.Contains(appEvent.Image, p.BuildTag) {
				log.Infof("find promote finished event for %s in %s", p.AppName, p.Environment)
				p.ResponseToBot(fmt.Sprintf("image: `%s` promoted on %s :tada:", p.BuildTag, p.Environment))
				jobDone <- true
				return
			}

		case <-ctx.Done():
			return
		}
	}

}

func (p *PromoteJob) GetId() string {
	return fmt.Sprintf("%s%s%s", p.botAction.GetCommand(), p.AppName, p.Environment)
}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}

func (p *PromoteJob) WaiteDeploymentStatus(status string) {

}
