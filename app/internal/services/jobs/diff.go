package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type DiffJob struct {
	botAction  *actions.BotAction
	clusters   map[string]clusters.Cluster
	currentEnv string
}

func (d *DiffJob) Launch(ctx context.Context, jobDone chan bool) {
	message := ""
	currentApps := d.clusters[d.currentEnv].Applications
	for env, cluster := range d.clusters {
		if env == d.currentEnv {
			continue
		}

		message = fmt.Sprintf("%s\n\ndiff from: `%s`\n\n", message, env)
		for name, app := range cluster.Applications {

			isFoundDisagreements := false
			select {
			case <-ctx.Done():
				return
			default:
				// check if found same app on other env
				if currentApps[name].Name == "" {
					continue
				}

				if *currentApps[name].Replicas != *app.Replicas {
					isFoundDisagreements = true
					message = fmt.Sprintf("%s *%s* ```current replicas - %d \n%s replicas - %d```",
						message, name, currentApps[name].Replicas, env, app.Replicas)
				}

				if currentApps[name].Image != app.Image {
					isFoundDisagreements = true
					message = fmt.Sprintf("%s *%s* ```current image - %d \n%s image - %d```",
						message, name, currentApps[name].Replicas, env, app.Replicas)
				}

				if !isFoundDisagreements {
					message = fmt.Sprintf("%s *%s* ```current and %s same```",
						message, name, env)
				}
			}
		}
	}

	d.ResponseToBot(message)
	jobDone <- true
}

func (d *DiffJob) GetId() string {
	return d.botAction.GetRawCommand()
}

func (d *DiffJob) ResponseToBot(message string) {
	d.botAction.ResponseOnAction(message)
}
