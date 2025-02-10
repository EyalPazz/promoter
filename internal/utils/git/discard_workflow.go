package git

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type DiscardGitFlow struct {
	BaseGitFlow
}

func (gf *DiscardGitFlow) Execute() error {
	err := gf.Checkout()
	if err != nil {
		return err
	}

	return nil
}

func (gf *DiscardGitFlow) Checkout() error {
	worktree, err := gf.BaseGitFlow.Checkout()

	if err != nil {
		return err
	}

	repo, err := utils.GetRepo()
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get head reference: %v", err)
	}

	var branch plumbing.ReferenceName
	if headRef.Name().IsBranch() {
		branch = headRef.Name()
	} else {
		return fmt.Errorf("not currently on a branch")
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branch,
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout: %v", err)
	}

	fmt.Println("Reverted All Changes")
	return nil
}
