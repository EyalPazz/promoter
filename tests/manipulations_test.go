package tests

import (
	"promoter/helpers/manipulations"
    "os"
    "path"
	"testing"
)

func TestFileExists(t *testing.T) {

	// Test cases
	tests := []struct {
		name          string
		filePath      string
		expectedValue bool
	}{
		{"Valid Path", "good.ssh", true},
		{"Invalid Path", "I_DONT_EXIST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

            home, _ := os.UserHomeDir()
			exists := manipulations.FileExists(path.Join(home, tt.filePath))

			if tt.expectedValue != exists {
				t.Errorf("expected %v but got %v", tt.expectedValue, exists)
			}
		})
	}
}
