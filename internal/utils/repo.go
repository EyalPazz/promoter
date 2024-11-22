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

func GetImageRepository(project, service, env, projectFilePath string) (string, error) {

	image, err := GetServiceImage(service, project, env, projectFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the input service's image : %s", err)
	}

	imageParts := strings.Split(image, "/")
	if len(imageParts) < 2 {
		return "", errors.New("invalid image")
	}

	return imageParts[1], nil
}

func RefreshRepo(hasPassphrase bool) error {

	if val := ManifestRepoExists(); !val {
		if err := cloneRepository(hasPassphrase); err != nil {
			return fmt.Errorf("error cloning git repo: %s", err)
		}
	}

	auth, err := auth.GetSSHAuth(hasPassphrase)
	if err != nil {
		return fmt.Errorf("error authenticating with git remote: %s", err)
	}

	repo, err := GetRepo()
	if err != nil {
		return fmt.Errorf("error getting manifest repo: %s", err)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching from remote: %s", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %s", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pulling from remote: %s", err)
	}

	fmt.Print("Successfully Fetched Recent Updates From Manifest \n\n")
	return nil
}

func ManifestRepoExists() bool {
	manifestPath, err := GetRepoPath()
	if err != nil {
		fmt.Printf("error getting repository path: %s", err)
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
		return errors.New("no manifest repo URL given")
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
