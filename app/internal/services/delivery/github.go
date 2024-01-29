package delivery

import (
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type OpsRepo struct {
	client     *github.Client
	repoOwner  string
	repoName   string
	branchName string
}

type RepoConfig struct {
	RepoOwner  string
	RepoName   string
	BranchName string
}

func NewOpsRepo(token string, config RepoConfig) *OpsRepo {

	g := &OpsRepo{
		repoOwner:  config.RepoOwner,
		repoName:   config.RepoName,
		branchName: config.BranchName,
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

func (g *OpsRepo) UpdateImage(filePath string, version string) error {

	fileContent, _, _, err := g.client.Repositories.GetContents(context.Background(), g.repoOwner, g.repoName, filePath, &github.RepositoryContentGetOptions{
		Ref: g.branchName,
	})
	if err != nil {
		return fmt.Errorf("github could not get file content: %v", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return fmt.Errorf("github not decode file content: %v", err)
	}

	newContent := strings.ReplaceAll(content, `version: "*"`, fmt.Sprintf(`version: "%s"`, version))

	sha := fileContent.GetSHA()

	commitMessage := fmt.Sprintf("Update version to %s", version)
	commit := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(newContent),
		SHA:     &sha,
	}

	_, _, err = g.client.Repositories.UpdateFile(context.Background(), g.repoOwner, g.repoName, filePath, commit)
	if err != nil {
		return fmt.Errorf("github could not update file: %v", err)
	}

	return nil
}
