package services

import (
	"github.com/slack-go/slack/slackevents"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot/actions"
)

type JobScheduler struct {
}

type Job struct {
	event  slackevents.EventsAPIEvent
	action actions.Action
}
