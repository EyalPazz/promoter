package tests

import (
	"promoter/helpers/manipulations"

	"testing"
)

func TestFileExists(t *testing.T) {

	// Test cases
	tests := []struct {
		name          string
		filePath      string
		expectedValue bool
	}{
		{"Valid Path", "/home/hilma/good.ssh", true},
		{"Invalid Path", "/homa/hilma/I_DONT_EXIST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			exists := manipulations.FileExists(tt.filePath)

			if tt.expectedValue != exists {
				t.Errorf("expected %v but got %v", tt.expectedValue, exists)
			}
		})
	}
}
