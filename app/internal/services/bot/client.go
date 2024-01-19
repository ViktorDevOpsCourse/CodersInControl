package bot

import (
	"context"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

type Client struct {
	ctx   context.Context
	api   *slack.Client
	event *socketmode.Client
}

func NewClient(ctx context.Context, options SlackClientOptions) *Client {
	if !options.IsValid() {
		log := logger.FromContext(ctx)
		log.Fatalf("Bot client options not valid")
		return nil
	}

	apiClient := slack.New(
		options.SlackBotToken,
		slack.OptionDebug(options.IsDebug),
		slack.OptionAppLevelToken(options.SlackAppToken),
	)

	return &Client{
		ctx:   ctx,
		api:   apiClient,
		event: socketmode.New(apiClient),
	}
}
