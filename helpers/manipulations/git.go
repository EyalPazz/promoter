package manipulations

import (
	"fmt"
	gitAuth "promoter/helpers/auth"
	"promoter/helpers/data"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CommitRepoChange(fileName string, service string, env string, tag string) error {
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

	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	_, commitErr := worktree.Commit(fmt.Sprint("promoting %s/%s in %s to %s", fileName, service, env, tag), &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
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
