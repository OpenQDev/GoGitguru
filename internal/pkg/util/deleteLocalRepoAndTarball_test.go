package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteLocalRepoAndTarball(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a temporary repo and tarball
	repo := "testRepo"
	repoPath := filepath.Join(tempDir, "repos", repo)
	tarballPath := filepath.Join(tempDir, "repos", repo+".tar.gz")

	// Create the repo and tarball
	os.MkdirAll(repoPath, 0755)
	os.WriteFile(tarballPath, []byte("test"), 0644)

	// Run the function
	err = DeleteLocalRepoAndTarball(tempDir+"repos", repo)
	if err != nil {
		t.Fatalf("Failed to delete repo and tarball: %v", err)
	}

	// Check if the repo and tarball have been deleted
	if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
		t.Errorf("Repo was not deleted")
	}
	if _, err := os.Stat(tarballPath); !os.IsNotExist(err) {
		t.Errorf("Tarball was not deleted")
	}
}
