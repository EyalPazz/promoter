package manipulations

import (
	"fmt"
	"promoter/internal/auth"
	"promoter/internal/utils"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
)

func DiscardChanges() error {

	repo, err := utils.GetRepo()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
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

func CommitRepoChange(commitMsg string) error {

	repo, err := utils.GetRepo()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	if _, err := worktree.Add("."); err != nil {
		return err
	}

	gitName := viper.GetString("git-name")
	if gitName == "" {
		return fmt.Errorf("No git name was given in config")
	}

	gitEmail := viper.GetString("git-email")
	if gitEmail == "" {
		return fmt.Errorf("No git email was given in config")
	}

	_, commitErr := worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})

	if commitErr != nil {
		return err
	}

	return nil
}

func PushToManifest(hasPassphrase bool) error {

	repo, err := utils.GetRepo()
	if err != nil {
		return err
	}

	auth, err := auth.GetSSHAuth(hasPassphrase)
	if err != nil {
		fmt.Println("Error Authenticating With Git Remote:", err)
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	return nil
}

func HandleRepoActions(commitMsg string, passphrase bool) error {
	if err := CommitRepoChange(commitMsg); err != nil {
		return err
	}

	if err := PushToManifest(passphrase); err != nil {
		return err
	}

	return nil
}
