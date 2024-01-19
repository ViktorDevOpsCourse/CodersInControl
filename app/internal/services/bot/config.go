package bot

import "strings"

type SlackOptions struct {
	ClientOptions SlackClientOptions
}

type SlackClientOptions struct {
	SlackBotToken string
	SlackAppToken string
	IsDebug       bool
}

func (c *SlackClientOptions) IsValid() bool {
	if c.SlackAppToken == "" && !strings.HasPrefix(c.SlackAppToken, "xapp-") {
		return false
	}
	if c.SlackBotToken == "" && !strings.HasPrefix(c.SlackBotToken, "xoxb-") {
		return false
	}
	return true
}
