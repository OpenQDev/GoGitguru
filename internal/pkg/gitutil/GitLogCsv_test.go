package gitutil

import (
	"fmt"
	"os"
	"testing"
)

func TestGitLogCsv(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	CloneRepo(tempDir, "OpenQDev", "OpenQ-Workflows")

	// Call the function under test
	out := GitLogCsv(tempDir, "OpenQ-Workflows", "2020-01-01")

	fmt.Println(string(out))
}
