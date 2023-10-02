package gitutil

import (
	"os"
	"reflect"
	"testing"
)

func TestFormatGitLogs(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	CloneRepo(tempDir, "OpenQDev", "OpenQ-DRM-TestRepo")

	// Call the function under test
	gitLogs := GetFormattedGitLogs(tempDir, "OpenQ-DRM-TestRepo", "")

	// Create a test data
	expectedGitLogs := []GitLog{
		{
			CommitHash:    "9fae86bc8e89895b961d81bd7e9e4e897501c8bb",
			AuthorName:    "FlacoJones",
			AuthorEmail:   "andrew@openq.dev",
			AuthorDate:    1696277205,
			CommitDate:    1696277205,
			CommitMessage: "initial commit",
			FilesChanged:  1,
			Insertions:    0,
			Deletions:     0,
		},
		{
			CommitHash:    "06a12f9c203112a149707ff73e4298749744c358",
			AuthorName:    "FlacoJones",
			AuthorEmail:   "andrew@openq.dev",
			AuthorDate:    1696277247,
			CommitDate:    1696277247,
			CommitMessage: "updates README",
			FilesChanged:  1,
			Insertions:    1,
			Deletions:     0,
		},
	}

	// Check if the gitLogs matches the expectedGitLogs
	if !reflect.DeepEqual(gitLogs, expectedGitLogs) {
		t.Errorf("Expected %v, but got %v", expectedGitLogs, gitLogs)
	}
}
