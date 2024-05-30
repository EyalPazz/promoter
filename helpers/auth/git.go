package auth

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/term"
)

func GetSSHAuth(hasPassphrase bool) (transport.AuthMethod, error) {
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
	if hasPassphrase {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Print("Enter passphrase (leave blank if none): ")
			passphraseBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, fmt.Errorf("failed to read passphrase: %v", err)
			}
			passphrase = string(passphraseBytes)
			fmt.Println()
		}

	}

	// Create the SSH authentication
	auth, err := ssh.NewPublicKeys("git", key, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH authentication: %v", err)
	}

	return auth, nil
}
