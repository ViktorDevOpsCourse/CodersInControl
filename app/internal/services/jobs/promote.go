package jobs

import (
	"context"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type PromoteJob struct {
	botAction   *actions.BotAction
	AppName     string
	BuildTag    string
	Environment string
	cluster     map[string]clusters.Cluster
}

func (p *PromoteJob) Launch(ctx context.Context, jobDone chan bool) {

	// TODO send request to git

	// TODO watch version

	// TODO update state/storage
	jobDone <- true

}

func (p *PromoteJob) GetId() string {
	return p.botAction.GetRawCommand()
}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}

func (p *PromoteJob) WaiteDeploymentStatus(status string) {

}
