package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type ListJob struct {
	botAction actions.Action
	cluster   *clusters.Cluster
}

func (l *ListJob) Launch() {
	// TODO Get from clusters apps states
}

func (l *ListJob) ResponseToBot(message string) {
	l.botAction.ResponseOnAction(message)
}
