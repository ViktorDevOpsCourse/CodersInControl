package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
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
	cluster           map[string]clusters.Cluster
}

func (p *PromoteJob) Launch(ctx context.Context, jobDone chan bool) {
	log := logger.FromContext(ctx)
	// TODO send request to git

	// TODO watch version

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			appEvent, err := p.appsEventsStorage.Get(p.Environment, p.AppName)
			if err != nil {
				if errors.Is(err, storage.NotFoundError) {
					continue
				}
				log.Errorf("Promote Job app `%s`. failed get app event from storage. Err `%s`", p.AppName, err)
			}
			if strings.Contains(appEvent.Image, p.BuildTag) {
				p.ResponseToBot(fmt.Sprintf("App *%s* image: `%s` promoted on %s :tada:", p.AppName, p.BuildTag, p.Environment))
				jobDone <- true
				return
			}

		case <-ctx.Done():
			return
		}
	}

}

func (p *PromoteJob) GetId() string {
	return p.botAction.GetRawCommand()
}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}

func (p *PromoteJob) WaiteDeploymentStatus(status string) {

}
