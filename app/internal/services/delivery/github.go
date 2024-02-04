package delivery

import (
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type FluxRepo struct {
	client        *github.Client
	repoOwner     string
	repoName      string
	workingBranch string
}

type Repo struct {
	Owner         string
	Name          string
	WorkingBranch string
}

type Application struct {
	FilePath string
	Version  string
}

type Config struct {
	Repo Repo
}

func NewFluxRepo(token string, config Config) Updater {

	g := &FluxRepo{
		repoOwner:     config.Repo.Owner,
		repoName:      config.Repo.Name,
		workingBranch: config.Repo.WorkingBranch,
	}

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

func (g *FluxRepo) Update(args interface{}) error {
	appConfig := Application{}
	var ok bool
	if appConfig, ok = args.(Application); !ok {
		return fmt.Errorf("invalid arguments to update flux, got %#v, expected %#v", args, Application{})
	}

	fileContent, _, _, err := g.client.Repositories.GetContents(
		context.Background(),
		g.repoOwner,
		g.repoName,
		appConfig.FilePath,
		&github.RepositoryContentGetOptions{
			Ref: g.workingBranch,
		})
	if err != nil {
		return fmt.Errorf("github could not get file content: %v", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return fmt.Errorf("github not decode file content: %v", err)
	}

	newContent := strings.ReplaceAll(content, `version: "*"`, fmt.Sprintf(`version: "%s"`, appConfig.Version))

	sha := fileContent.GetSHA()

	commitMessage := fmt.Sprintf("Update version to %s", appConfig.Version)
	commit := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(newContent),
		SHA:     &sha,
	}

	_, _, err = g.client.Repositories.UpdateFile(
		context.Background(),
		g.repoOwner,
		g.repoName,
		appConfig.Version,
		commit)
	if err != nil {
		return fmt.Errorf("github could not update file: %v", err)
	}

	return nil
}
