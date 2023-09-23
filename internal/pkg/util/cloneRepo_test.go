package util

import (
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	tests := []struct {
		name         string
		repo         string
		organization string
		wantErr      bool
	}{
		{
			name:         "Valid repo and organization",
			repo:         "OpenQ-Workflows",
			organization: "OpenQDev",
			wantErr:      false,
		},
		{
			name:         "Invalid repo",
			repo:         "invalid-repo",
			organization: "valid-org",
			wantErr:      true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloneRepo(tt.organization, tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Defer deletion of the repo after each test
			defer os.RemoveAll(tt.repo)
		})
	}
}
