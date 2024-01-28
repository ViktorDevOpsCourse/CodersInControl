package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type PromoteJob struct {
	botAction *actions.BotAction
	cluster   *clusters.Cluster
}

func (p *PromoteJob) Launch(ctx context.Context, jobDone chan bool) {
	// TODO send request to run pipeline on update service version
	// TODO receive answer
	// TODO run watcher on update
	// TODO update state/storage
	fmt.Println(p.cluster.GetApplications())
	jobDone <- true

}

func (p *PromoteJob) GetId() string {
	return ""
}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}
