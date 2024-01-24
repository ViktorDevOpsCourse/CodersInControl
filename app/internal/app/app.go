package app

import (
	"context"
	"github.com/viktordevopscourse/codersincontrol/app/internal/config"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/bot"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/k8s"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
)

func Run(ctx context.Context) error {
	cfg, err := config.ParseAppConfig(config.AppVersion)
	if err != nil {
		logger.GetDefaultLogger().Fatalln(err)
	}

	// TODO read it from config file or from db. Different envs
	botAuthConfig := bot.AuthConfig{
		AllowedUsers: map[string]struct{}{
			"viktorzhabskiy": struct{}{},
		},
	}

	dispatcher := services.NewJobDispatcher(k8s.NewK8SService())
	slackBot := bot.NewSlackBot(ctx, bot.SlackOptions{
		ClientOptions: bot.SlackClientOptions{
			SlackBotToken: cfg.Bot.SlackBotToken,
			SlackAppToken: cfg.Bot.SlackAppToken,
			IsDebug:       true,
		},
		BotOptions: bot.SlackBotOptions{
			ActionProcessorQueue: dispatcher.GetActionProcessorQueue(),
		},
		AuthOptions: botAuthConfig,
	})

	slackBot.Run()
	dispatcher.Run()

	return nil
}
