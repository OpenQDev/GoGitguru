package gitutil

import (
	"os"
	"testing"
)

func TestFormatGitLogs(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	CloneRepo(tempDir, "OpenQDev", "OpenQ-Workflows")

	// Call the function under test
	_ = GetFormattedGitLogs(tempDir, "OpenQ-Workflows", "")
}
