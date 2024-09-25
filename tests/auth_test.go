package tests

import (
	"promoter/internal/auth"

	"testing"

	"github.com/spf13/viper"
)

// Mocking term.ReadPassword
func TestGetSSHAuth(t *testing.T) {

	// Test cases
	tests := []struct {
		name        string
		sshKeyPath  string
		passphrase  string
		expectError bool
	}{
		{"Valid SSH key", "/home/runner/good.ssh", "", false},
		{"Invalid SSH key", "/home/runner/bad.ssh", "", true},
		{"Invalid SSH keypath", "/invalid/path", "", true},
		{"Error reading SSH key file", "/tmp/nonexistent_key", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("ssh-key", tt.sshKeyPath)

			auth, err := auth.GetSSHAuth(false)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if auth == nil {
					t.Errorf("expected valid auth but got nil")
				}
			}
		})
	}
}
