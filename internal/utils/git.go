package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"promoter/internal/types"

	"gopkg.in/yaml.v3"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	DEFAULT_PROMOTER_DIR = "/promoter/"
)

func GetRepoPath() (string, error) {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user cache directory: %w", err)
	}
	return baseDir + DEFAULT_PROMOTER_DIR, nil
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

func ComposeCommitMsg(changes *[]types.ServiceChanges, env, project string) string {
	msg := fmt.Sprintf("promotion(%s): %s \n", env, project)
	for _, change := range *changes {
		msg += fmt.Sprintf("changed %s to %s \n", change.Name, change.NewTag)
	}
	return msg
}

func GetLatestRevisions(project, env string, dayInterval int) ([]*object.Commit, error) {
	repo, err := GetRepo()
	if err != nil {
		return nil, err
	}

	projectFile, err := GetProjectFile(project, env, true)
	if err != nil {
		return nil, err
	}

	sinceTime := time.Now().AddDate(0, 0, -dayInterval)
	logs, err := repo.Log(&git.LogOptions{
		PathFilter: func(s string) bool {
			return s == projectFile
		},
		Order: git.LogOrderCommitterTime,
		Since: &sinceTime,
	})

	if err != nil {
		return nil, err
	}

	var res []*object.Commit

	err = logs.ForEach(func(c *object.Commit) error {
		res = append(res, c)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetFileFromCommit(commit *object.Commit, filePath string) (*Config, error) {
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	file, err := tree.File(filePath)
	if err != nil {
		return nil, err
	}

	reader, err := file.Blob.Reader()
	if err != nil {
		return nil, err
	}

	fileData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return nil, err
	}

	return &config, nil

}
