package services

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/actions"
)

type Job struct {
	botAction actions.Action
}

func NewJob(botAction actions.Action) Job {
	return Job{
		botAction: botAction,
	}
}

func (j *Job) Launch() {
	// TODO
	// TODO call to github
	// TODO watch on k8s services
	// TODO send response in bot chanel
	j.botAction.SendResponse(fmt.Sprintf("Job `%s` in chanel `%s`", j.botAction.GetName(), j.botAction.GetName()))
	j.botAction.Done()
}
