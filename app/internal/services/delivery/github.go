package delivery

import (
	"context"
	"github.com/google/go-github/v58/github"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
}

func NewGithub(token string) *Github {
	log := logger.FromDefaultContext()
	g := &Github{}
	ctx := context.Background()
	githubTokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	githubHttpClient := oauth2.NewClient(ctx, githubTokenSource)
	if githubHttpClient == nil {
		log.Fatal("github oauth2 client init error")
		return nil
	}

	g.client = github.NewClient(githubHttpClient)
	if g.client == nil {
		log.Fatal("github client init error")
	}

	return g
}
