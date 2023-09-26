package gitutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteLocalRepo(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}

	// Create a test repo directory inside the temp directory
	testRepo := filepath.Join(tempDir, "testRepo")
	err = os.Mkdir(testRepo, 0755)
	if err != nil {
		t.Fatalf("Failed to create test repo directory: %s", err)
	}

	// Call the function to test
	err = DeleteLocalRepo(tempDir, "testRepo")
	if err != nil {
		t.Fatalf("DeleteLocalRepoAndTarball failed: %s", err)
	}

	// Check if the test repo directory has been deleted
	if _, err := os.Stat(testRepo); !os.IsNotExist(err) {
		t.Fatalf("Test repo directory was not deleted")
	}
}
