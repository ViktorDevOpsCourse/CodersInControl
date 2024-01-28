package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type ListJob struct {
	botAction *actions.BotAction
	clusters  map[string]clusters.Cluster
}

func (l *ListJob) Launch(ctx context.Context, jobDone chan bool) {
	message := ""
	for env, cluster := range l.clusters {
		message = fmt.Sprintf("%s\n\nenvironment: `%s`\n\n", message, env)
		for name, app := range cluster.Applications {
			select {
			case <-ctx.Done():
				return
			default:
				message = fmt.Sprintf("%s *%s* ```version image - %s \nstatus - %s```",
					message, name, app.Image, app.Status.ServiceStatus)
			}
		}
	}

	l.ResponseToBot(message)
	jobDone <- true
}

func (l *ListJob) GetId() string {
	return l.botAction.GetRawCommand()
}

func (l *ListJob) ResponseToBot(message string) {
	l.botAction.ResponseOnAction(message)
}
