package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/k8s"
)

type PromoteJob struct {
	botAction  actions.Action
	k8sService *k8s.K8S
}

func (p *PromoteJob) Launch() {
	// TODO send request to run pipeline on update service version
	// TODO receive answer
	// TODO run watcher on update
	// TODO update state/storage

}

func (p *PromoteJob) ResponseToBot(message string) {
	p.botAction.ResponseOnAction(message)
}
