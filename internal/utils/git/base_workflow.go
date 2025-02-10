package git

import (
	"fmt"
	"promoter/internal/auth"
	"promoter/internal/utils"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
)

type BaseGitFlow struct {
	CommitMsg  string
	Passphrase bool
}

func (gf *BaseGitFlow) Execute() error {
	worktree, err := gf.Checkout()
	if err != nil {
		return err
	}

	if err := gf.Add(worktree); err != nil {
		return err
	}

	if err := gf.Commit(worktree, gf.CommitMsg); err != nil {
		return err
	}

	if err := gf.Push(gf.Passphrase); err != nil {
		return err
	}

	return nil
}

func (gf *BaseGitFlow) Checkout() (*git.Worktree, error) {
	repo, err := utils.GetRepo()
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return worktree, nil
}

func (gf *BaseGitFlow) Add(worktree *git.Worktree) error {
	if _, err := worktree.Add("."); err != nil {
		return err
	}

	return nil
}

func (gf *BaseGitFlow) Commit(worktree *git.Worktree, commitMsg string) error {
	gitName := viper.GetString("git-name")
	if gitName == "" {
		return fmt.Errorf("no git name was given in config")
	}

	gitEmail := viper.GetString("git-email")
	if gitEmail == "" {
		return fmt.Errorf("no git email was given in config")
	}

	_, err := worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})

	return err
}

func (gf *BaseGitFlow) Push(hasPassphrase bool) error {
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

	return err
}
