package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/k8s"
)

type ListJob struct {
	botAction  actions.Action
	k8sService *k8s.K8S
}

func (l *ListJob) Launch() {
	// TODO Get from clusters apps states
}

func (l *ListJob) ResponseToBot(message string) {
	l.botAction.ResponseOnAction(message)
}
