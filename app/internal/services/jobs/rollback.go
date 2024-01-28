package jobs

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type RollbackJob struct {
	botAction  *bot.BotAction
	k8sService *clusters.K8S
}

func (r *RollbackJob) Launch() {
	// TODO get old state from storage
	// TODO send request to run pipeline on update service version
	// TODO receive answer
	// TODO run watcher on update
	// TODO update state/storage

}

func (r *RollbackJob) ResponseToBot(message string) {
	r.botAction.ResponseOnAction(message)
}
