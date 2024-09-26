package utils

import (
	"errors"
	"fmt"
	"os"
	"promoter/internal/auth"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

func GetImageRepository(project string, service string, env string, projectFilePath string) (string, error) {

	image, err := GetServiceImage(service, project, env, projectFilePath)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve the input service's image : %s", err)
	}

	imageParts := strings.Split(image, "/")
	if len(imageParts) < 2 {
		return "", errors.New("Invalid Image")
	}

	return imageParts[1], nil
}

func RefreshRepo(hasPassphrase bool) error {

	if val := ManifestRepoExists(); !val {
		if err := cloneRepository(hasPassphrase); err != nil {
			return fmt.Errorf("Error Cloning Git Repo: %s", err)
		}
	}

	auth, err := auth.GetSSHAuth(hasPassphrase)
	if err != nil {
		return fmt.Errorf("Error Authenticating With Git Remote: %s", err)
	}

	repo, err := GetRepo()
	if err != nil {
		return fmt.Errorf("Error Getting manifest repo: %s", err)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("Error fetching from remote: %s\n", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Error getting worktree: %s\n", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("Error pulling from remote: %s\n", err)
	}

	fmt.Println("Successfully Fetched Recent Updates From Manifest")
	return nil
}

func ManifestRepoExists() bool {
	manifestPath, err := GetRepoPath()
	if err != nil {
		fmt.Printf("Error getting repository path: %s\n", err)
		return false
	}

	info, err := os.Stat(manifestPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

func cloneRepository(hasPassphrase bool) error {
	auth, err := auth.GetSSHAuth(hasPassphrase)
	if err != nil {
		return err
	}

	manifestRepoUrl := viper.GetString("manifest-repo")
	if manifestRepoUrl == "" {
		return errors.New("No manifest repo URL given")
	}

	repoPath, err := GetRepoPath()
	if err != nil {
		return err
	}

	fmt.Println("Cloning repository...")

	_, cloneErr := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:  manifestRepoUrl,
		Auth: auth,
	})

	return cloneErr
}
