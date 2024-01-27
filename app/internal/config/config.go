package config

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"gopkg.in/yaml.v2"
	"os"
)

const AppConfigPath = "./config.example.yaml"

var AppVersion = "unknown version"

type Config struct {
	Bot                   BotConfig
	K8S                   K8SConfig
	ServiceConfigFilePath string `conf:"env:SERVICE_CONFIG_FILE_PATH"`
}

type FileConfig struct {
	Clusters map[string]Cluster `yaml:"clusters"`
	Bot      Bot                `yaml:"bot"`
}
type Cluster struct {
	File string `yaml:"file"`
}

type K8SConfig struct {
	Clusters map[string]Cluster
}

type Bot struct {
	Admins []string `yaml:"admins"`
}

type BotConfig struct {
	Admins        map[string]struct{}
	SlackBotToken string `conf:"env:SLACK_BOT_TOKEN"`
	SlackAppToken string `conf:"env:SLACK_APP_TOKEN"`
}

func ParseAppConfig(version string) (Config, error) {
	cfg := Config{
		Bot: BotConfig{
			Admins: make(map[string]struct{}),
		},
		K8S: K8SConfig{
			Clusters: make(map[string]Cluster),
		},
	}

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

	if cfg.ServiceConfigFilePath == "" {
		cfg.ServiceConfigFilePath = AppConfigPath
	}

	yamlFile, err := os.ReadFile(cfg.ServiceConfigFilePath)
	if err != nil {
		fmt.Printf("Failed read config file: %v", err)
		os.Exit(0)
	}

	var fc FileConfig

	err = yaml.Unmarshal(yamlFile, &fc)
	if err != nil {
		fmt.Printf("Failed unpacke config yaml file: %v", err)
		os.Exit(0)
	}

	cfg.K8S = K8SConfig{
		Clusters: fc.Clusters,
	}

	for _, admin := range fc.Bot.Admins {
		cfg.Bot.Admins[admin] = struct{}{}
	}

	return cfg, nil
}
