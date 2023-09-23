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
	defer os.RemoveAll(testDir)

	// Test: Call the function with the test directory
	_, err = CompressDirectory(testDir, "/tmp/testTar")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check: Verify that the tarball was created
	if _, err := os.Stat("/tmp/testTar.tar.gz"); os.IsNotExist(err) {
		t.Errorf("Tarball was not created")
	}
}
