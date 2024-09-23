package data

import (
	"errors"
	"fmt"
	"os"
	gitAuth "promoter/utils/auth"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

const (
	DefaultPromoterDir = "/promoter-data/repositories/"
)

// NOTICE: this determines the convention for image repository names
func GetImageRepository(project string, service string, env string, projectFilePath string) (string, error) {

	image, err := GetServiceImage(service, project, env, projectFilePath, viper.GetString("manifestRepoRoot"))
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve the input service's image : %s", err)
	}

	imageParts := strings.Split(image, "/")
	if len(imageParts) < 2 {
		return "", errors.New("Invalid Image")
	}

	return imageParts[1], nil
}

func GetRepoPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + DefaultPromoterDir + "manifest", nil
}

func RefreshRepo(hasPassphrase bool) {

	if val  := ManifestRepoExists(); !val {
        if err := cloneRepository(hasPassphrase); err != nil {
			fmt.Println("Error Cloning Git Repo", err)
			os.Exit(1)
		}
	}

	auth, err := gitAuth.GetSSHAuth(hasPassphrase)
	if err != nil {
		fmt.Println("Error Authenticating With Git Remote:", err)
		os.Exit(1)
	}

	repoPath, err := GetRepoPath()
	if err != nil {
		fmt.Printf("Error getting repository path: %s\n", err)
		os.Exit(1)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Printf("Error opening repository: %s\n", err)
		os.Exit(1)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("Error fetching from remote: %s\n", err)
		os.Exit(1)
	}

	w, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Error getting worktree: %s\n", err)
		os.Exit(1)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("Error pulling from remote: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully Fetched Recent Updates From Manifest")
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
	auth, err := gitAuth.GetSSHAuth(hasPassphrase)
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
