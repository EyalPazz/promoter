package manipulations

import (
	"fmt"
	gitAuth "promoter/helpers/auth"
	"promoter/helpers/data"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
)

func CommitRepoChange(project string, service string, env string, tag string) error {
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
		fmt.Println("No git name was given in config")
		return err
	}

	gitEmail := viper.GetString("git-email")
	if gitEmail == "" {
		fmt.Println("No git email was given in config")
		return err
	}

	commitMessage := fmt.Sprintf("promoting %s/%s in %s to %s", project, service, env, tag)
	_, commitErr := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})

	if commitErr != nil {
		return err
	}

	auth, err := gitAuth.GetSSHAuth()
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
