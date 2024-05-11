package data

import (
	"fmt"
	"os"
	gitAuth "promoter/helpers/auth"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

const (
	DefaultPromoterDir = "/promoter-data/repositories/"
)

func CloneRepository() error {

	auth, err := gitAuth.GetSSHAuth()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	manifestRepoUrl := viper.GetString("manifest-repo")
	if manifestRepoUrl == "" {
		fmt.Println("No manifest repo URL given")
		return err
	}

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting home directory:", err)
        return err
    }
    
	fmt.Println("Cloning repository...")

	// Clone the repository
	_, cloneErr := git.PlainClone(homeDir+DefaultPromoterDir+"manifest", false, &git.CloneOptions{
		URL:  manifestRepoUrl,
		Auth: auth,
	})

	return cloneErr
}
