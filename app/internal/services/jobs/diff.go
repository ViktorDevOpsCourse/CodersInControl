package jobs

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters"
)

type DiffJob struct {
	botAction          *bot.BotAction
	clusters           clusters.ClustersCopy
	currentClusterName string
}

func NewDiffJob(botAction *bot.BotAction,
	clusters clusters.ClustersCopy) (*DiffJob, error) {

	currentClusterName := botAction.GetCommandArgs()
	if _, err := clusters.GetCluster(currentClusterName); err != nil {
		return nil, fmt.Errorf("invalid command or unknow environment. Accept `@bot diff environment`")
	}
	return &DiffJob{
		botAction:          botAction,
		clusters:           clusters,
		currentClusterName: botAction.GetCommandArgs(),
	}, nil
}

func (d *DiffJob) Launch(ctx context.Context, jobDone chan bool) {
	if len(d.clusters) <= 1 {
		d.ResponseToBot("you have 1 or less clusters")
		jobDone <- true
		return
	}

	d.ResponseToBot(d.compareApps(ctx))
	jobDone <- true
}

func (d *DiffJob) GetId() string {
	return fmt.Sprintf("%s %s", d.botAction.Event.ChannelId, d.botAction.GetRawCommand())
}

func (d *DiffJob) ResponseToBot(message string) {
	d.botAction.ResponseOnAction(message)
}

func (d *DiffJob) compareApps(ctx context.Context) string {

	cluster, _ := d.clusters.GetCluster(d.currentClusterName)
	baseApps := cluster.Applications
	message := ""

	for clusterName, cluster := range d.clusters {
		if clusterName == d.currentClusterName {
			continue
		}

		message = fmt.Sprintf("%s\n\ndifference between: `%s` and `%s` \n\n", message, d.currentClusterName, clusterName)
		for name, app := range cluster.Applications {

			select {
			case <-ctx.Done():
				return ""
			default:
				baseApp := baseApps[name]

				// check if found same app on other env
				if baseApp.Name == "" {
					continue
				}

				differenceMessage := d.compareApp(baseApp, app, clusterName)

				message = fmt.Sprintf("%s app: *%s* ```%s```",
					message, name, differenceMessage)
			}
		}
	}

	return message
}

func (d *DiffJob) compareApp(currentApp, otherApp clusters.Application, otherClusterName string) string {
	compareMessage := ""
	isFoundDisagreements := false

	if currentApp.Replicas != nil && otherApp.Replicas != nil {
		if *currentApp.Replicas != *otherApp.Replicas {
			isFoundDisagreements = true
			compareMessage = fmt.Sprintf("%s%s replicas - %d \n%s replicas - %d\n",
				compareMessage, d.currentClusterName, *currentApp.Replicas, otherClusterName, *otherApp.Replicas)
		}
	}

	if currentApp.Image != otherApp.Image {
		isFoundDisagreements = true
		compareMessage = fmt.Sprintf("%s%s image - %s \n%s image - %s\n",
			compareMessage, d.currentClusterName, currentApp.Image, otherClusterName, otherApp.Image)
	}

	if !isFoundDisagreements {
		compareMessage = fmt.Sprintf("%s and %s same", d.currentClusterName, otherClusterName)
	}

	return compareMessage
}
