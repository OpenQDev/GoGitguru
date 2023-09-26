package s3util

import (
	"os"
	"testing"
)

func TestCompressDirectory(t *testing.T) {
	// Setup: Create a temporary directory for testing
	testDir, err := os.MkdirTemp("", "testDir")

	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	tarballPath := "testTar.tar.gz"

	defer os.RemoveAll(testDir)
	defer os.RemoveAll(tarballPath)

	// Test: Call the function with the test directory
	_, err = CompressDirectory(tarballPath, testDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check: Verify that the tarball was created
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		t.Errorf("Tarball was not created")
	}
}
