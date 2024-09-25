package utils

import (
	"fmt"
	"os"

	"time"

	"promoter/internal/types"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func GetLatestRevisions(project string, env string, dayInterval int) ([]string, error) {
	repo, err := GetRepo()
	if err != nil {
		return nil, err
	}

    projectFile, err := GetProjectFile(project, env, true)
	if err != nil {
		return nil, err
	}

    sinceTime := time.Now().AddDate(0,0, dayInterval)
    logs, err := repo.Log(&git.LogOptions{
            PathFilter: func(s string) bool {
            return s == projectFile
        },
            Since: &sinceTime ,
    })

	if err != nil {
		return nil, err
	}

    var res []string;

    err = logs.ForEach(func(c *object.Commit) error {
        res = append(res, c.Hash.String())
        return nil
    })

	return res, nil
}
