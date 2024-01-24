package bot

import "github.com/slack-go/slack"

type Auth struct {
	api  *slack.Client
	conf AuthConfig
}

func NewAuth(api *slack.Client, conf AuthConfig) *Auth {
	return &Auth{api: api, conf: conf}
}

func (a *Auth) hasPermissions(userId string) (bool, error) {
	user, err := a.api.GetUserInfo(userId)
	if err != nil {
		return false, err
	}

	if _, ok := a.conf.AllowedUsers[user.Name]; !ok {
		return false, nil
	}

	return true, nil
}
