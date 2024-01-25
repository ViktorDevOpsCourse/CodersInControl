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

	k8sConfig := k8s.Config{
		Clusters: make(map[string]string),
	}
	for environment, clusterOption := range cfg.K8S.Clusters {
		k8sConfig.Clusters[environment] = clusterOption.File
	}

	dispatcher := services.NewJobDispatcher(k8s.NewK8SService(k8sConfig))
	slackBot := bot.NewSlackBot(ctx, bot.SlackOptions{
		ClientOptions: bot.SlackClientOptions{
			SlackBotToken: cfg.Bot.SlackBotToken,
			SlackAppToken: cfg.Bot.SlackAppToken,
			IsDebug:       true,
		},
		BotOptions: bot.SlackBotOptions{
			ActionProcessorQueue: dispatcher.GetActionProcessorQueue(),
		},
		AuthOptions: bot.AuthConfig{
			AllowedUsers: cfg.Bot.Admins,
		},
	})

	slackBot.Run()
	dispatcher.Run()

	return nil
}
