package gitprovider

import (
	"context"
	"promoter/internal/utils"

	"github.com/google/go-github/v69/github"
)

type GitHub struct {
	client *github.Client
}

func NewGitHubClient() *GitHub {
	return &GitHub{
		client: github.NewClient(nil).WithAuthToken(utils.GetGitProviderToken()),
	}
}

func (gh *GitHub) CreatePR(head, body, title string) error {

	ctx := context.Background()
	config, _ := utils.GetConfig()

	if _, _, err := gh.client.PullRequests.Create(ctx, config.PullRequests.Org, config.PullRequests.RepoName, &github.NewPullRequest{
		Head:  utils.PtrTo(head),
		Base:  &config.PullRequests.BaseBranch,
		Body:  utils.PtrTo(body),
		Title: utils.PtrTo(title),
	}); err != nil {
		return err
	}

	return nil
}
