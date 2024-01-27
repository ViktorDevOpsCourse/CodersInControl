package jobs

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type PromoteJob struct {
	botAction actions.Action
	cluster   *clusters.Cluster
}

func (p *PromoteJob) Launch() {
	// TODO send request to run pipeline on update service version
	// TODO receive answer
	// TODO run watcher on update
	// TODO update state/storage
	fmt.Println(p.cluster.GetApplications())

}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}
