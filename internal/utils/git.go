package utils

import (
	"os"
    "fmt"

	"promoter/internal/types"

	"github.com/go-git/go-git/v5"

)

const (
	DEFAULT_PROMOTER_DIR = "/promoter-data/repositories/"
)

func GetRepoPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + DEFAULT_PROMOTER_DIR + "manifest", nil
}

func GetRepo() (*git.Repository, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
    }

    return repo, nil
}

func ComposeCommitMsg(changes []types.ServiceChanges, env string, project string) string {
	msg := fmt.Sprintf("promotion(%s): %s \n", env, project)
	for _, change := range changes {
		msg += fmt.Sprintf("changed %s to %s \n", change.Name, change.NewTag)
	}
	return msg
}
