package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type DiffJob struct {
	botAction  *bot.BotAction
	clusters   map[string]clusters.Cluster
	currentEnv string
}

func NewDiffJob(botAction *bot.BotAction,
	clusters map[string]clusters.Cluster) (*DiffJob, error) {

	currentEnv := botAction.GetCommandArgs()
	if _, ok := clusters[currentEnv]; !ok {
		return nil, fmt.Errorf("invalid command or unknow environment. Accept `@bot diff environment`")
	}
	return &DiffJob{
		botAction:  botAction,
		clusters:   clusters,
		currentEnv: botAction.GetCommandArgs(),
	}, nil
}

func (d *DiffJob) Launch(ctx context.Context, jobDone chan bool) {
	if len(d.clusters) <= 1 {
		d.ResponseToBot("you have 1 or less clusters")
		jobDone <- true
		return
	}

	message := ""
	currentApps := d.clusters[d.currentEnv].Applications
	for env, cluster := range d.clusters {
		if env == d.currentEnv {
			continue
		}

		message = fmt.Sprintf("%s\n\ndifference between: `%s` and `%s` \n\n", message, d.currentEnv, env)
		for name, app := range cluster.Applications {

			select {
			case <-ctx.Done():
				return
			default:
				currentApp := currentApps[name]

				// check if found same app on other env
				if currentApp.Name == "" {
					continue
				}

				differenceMessage := d.compareApps(currentApp, app, env)

				message = fmt.Sprintf("%s app: *%s* ```%s```",
					message, name, differenceMessage)
			}
		}
	}

	d.ResponseToBot(message)
	jobDone <- true
}

func (d *DiffJob) GetId() string {
	return fmt.Sprintf("%s %s", d.botAction.Event.ChannelId, d.botAction.GetRawCommand())
}

func (d *DiffJob) ResponseToBot(message string) {
	d.botAction.ResponseOnAction(message)
}

func (d *DiffJob) compareApps(currentApp, otherApp clusters.Application, otherEnv string) string {
	compareMessage := ""
	isFoundDisagreements := false

	if currentApp.Replicas != nil && otherApp.Replicas != nil {
		if *currentApp.Replicas != *otherApp.Replicas {
			isFoundDisagreements = true
			compareMessage = fmt.Sprintf("%s%s replicas - %d \n%s replicas - %d\n",
				compareMessage, d.currentEnv, *currentApp.Replicas, otherEnv, *otherApp.Replicas)
		}
	}

	if currentApp.Image != otherApp.Image {
		isFoundDisagreements = true
		compareMessage = fmt.Sprintf("%s%s image - %s \n%s image - %s\n",
			compareMessage, d.currentEnv, currentApp.Image, otherEnv, otherApp.Image)
	}

	if !isFoundDisagreements {
		compareMessage = fmt.Sprintf("%s and %s same", d.currentEnv, otherEnv)
	}

	return compareMessage
}
