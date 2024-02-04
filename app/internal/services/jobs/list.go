package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type ListJob struct {
	botAction *bot.BotAction
	clusters  clusters.ClustersCopy
}

func NewListJob(botAction *bot.BotAction,
	clusters clusters.ClustersCopy) (*ListJob, error) {
	return &ListJob{
		botAction: botAction,
		clusters:  clusters,
	}, nil
}

func (l *ListJob) Launch(ctx context.Context, jobDone chan bool) {
	message := ""
	for clusterName, cluster := range l.clusters {
		message = fmt.Sprintf("%s\n\nenvironment: `%s`\n\n", message, clusterName)
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
	return fmt.Sprintf("%s %s", l.botAction.Event.ChannelId, l.botAction.GetRawCommand())
}

func (l *ListJob) ResponseToBot(message string) {
	l.botAction.ResponseOnAction(message)
}
