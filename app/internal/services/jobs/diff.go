package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type DiffJob struct {
	botAction  actions.Action
	k8sService *clusters.K8S
}

func (d *DiffJob) Launch() {
	// // TODO Get from clusters apps states
	// TODO compare it and show promote

}

func (d *DiffJob) ResponseToBot(message string) {
	d.botAction.ResponseOnAction(message)
}
