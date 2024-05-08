package helpers

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/term"
)

const (
	DefaultPromoterDir = "/promoter-data/repositories/"
)

func getSSHAuth() (transport.AuthMethod, error) {
	sshKeyPath := viper.GetString("ssh-key")
	if sshKeyPath == "" {
		return nil, fmt.Errorf("no SSH key path given")
	}

	// Read the private key
	sshKeyPath = os.ExpandEnv(sshKeyPath)

	key, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %v", err)
	}

	// Ask for passphrase if the key is encrypted
	passphrase := ""
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Enter passphrase (leave blank if none): ")
		passphraseBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, fmt.Errorf("failed to read passphrase: %v", err)
		}
		passphrase = string(passphraseBytes)
		fmt.Println()
	}

	// Create the SSH authentication
	auth, err := ssh.NewPublicKeys("git", key, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH authentication: %v", err)
	}

	return auth, nil
}

func cloneRepository(auth transport.AuthMethod, url, dest string) error {
	fmt.Println("Cloning repository...")

	// Clone the repository
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})

	return err
}

var ManifestRepoClone = &cobra.Command{
	Use:   "clone-manifest",
	Short: "Authenticate with Git provider using SSH key",
	Long:  "Authenticate with Git provider using SSH key to perform Git operations.",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := getSSHAuth()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		manifestRepoUrl := viper.GetString("manifest-repo")
		if manifestRepoUrl == "" {
			fmt.Println("No manifest repo URL given")
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}

		cloneErr := cloneRepository(auth, manifestRepoUrl, homeDir+DefaultPromoterDir+"manifest")
		if cloneErr != nil {
			fmt.Println("Error cloning manifest repo:", cloneErr)
			return
		}

		fmt.Println("Successfully cloned manifest repo")
	},
}
