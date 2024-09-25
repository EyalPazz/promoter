package manipulations

import (
	"fmt"
	"promoter/internal/types"
	gitAuth "promoter/internal/auth"
	"promoter/internal/data"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
)

func HandleDiscard() {
	repoPath, err := data.GetRepoPath()
	if err != nil {
		fmt.Println(err)
	}

	err = discardChanges(repoPath)
	if err != nil {
		fmt.Println(err)
	}
}

func discardChanges(repoPath string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
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

func composeCommitMsg(changes []types.ServiceChanges, env string, project string) string {
	msg := fmt.Sprintf("promotion(%s): %s \n", env, project)
	for _, change := range changes {
		msg += fmt.Sprintf("changed %s to %s \n", change.Name, change.NewTag)
	}
	return msg
}

func CommitRepoChange(project string, changeLog *[]types.ServiceChanges, env string) error {
	repoPath, err := data.GetRepoPath()
	if err != nil {
		return err
	}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(".")
	if err != nil {
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

	_, commitErr := worktree.Commit(composeCommitMsg(*changeLog, env, project), &git.CommitOptions{
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
	repoPath, err := data.GetRepoPath()
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	auth, err := gitAuth.GetSSHAuth(hasPassphrase)
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
