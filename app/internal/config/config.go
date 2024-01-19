package config

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
)

var AppVersion = "unknown version"

type Config struct {
	Bot BotConfig
}

type BotConfig struct {
	SlackBotToken string `conf:"env:SLACK_BOT_TOKEN"`
	SlackAppToken string `conf:"env:SLACK_APP_TOKEN"`
}

func ParseAppConfig(version string) (Config, error) {
	var cfg Config

	err := conf.Parse(os.Args[1:], "", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(conf.Usage("", &cfg))
			os.Exit(0)
		}
		if errors.Is(err, conf.ErrVersionWanted) {
			fmt.Printf(`version %s\n`, version)
			os.Exit(0)
		}
		return cfg, err
	}
	return cfg, nil
}
