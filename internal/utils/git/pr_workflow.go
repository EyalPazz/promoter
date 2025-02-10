package git

import (
	"fmt"
	"promoter/internal/types"
	"promoter/internal/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type PRGitWorkflow struct {
	BaseGitFlow
	GitProvider  types.IGitProvider
	Title        string
	Body         string
	ChangeBranch string
}

func (gf *PRGitWorkflow) Execute() error {
	worktree, err := gf.CheckoutChange()
	if err != nil {
		return err
	}

	if err := gf.BaseGitFlow.Execute(); err != nil {
		return err
	}

	if err := gf.CreatePR(); err != nil {
		return err
	}

	if err := gf.CheckoutBase(worktree); err != nil {
		return err
	}

	return nil
}

func (gf *PRGitWorkflow) CheckoutBase(worktree *git.Worktree) error {
	config, err := utils.GetConfig()

	if err != nil {
		return err
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + config.PullRequests.BaseBranch),
	}); err != nil {
		return err
	}

	return nil
}

func (gf *PRGitWorkflow) CheckoutChange() (*git.Worktree, error) {
	worktree, err := gf.BaseGitFlow.Checkout()

	if err != nil {
		return nil, err
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.ReferenceName("refs/heads/" + gf.ChangeBranch),
		Keep:   true,
	}); err != nil {
		fmt.Println(plumbing.ReferenceName(gf.ChangeBranch))
		return nil, err
	}

	return worktree, nil
}

func (gf *PRGitWorkflow) CreatePR() error {
	fmt.Print("Creating a PR... \n\n")
	return gf.GitProvider.CreatePR(gf.ChangeBranch, gf.Body, gf.Title)
}
