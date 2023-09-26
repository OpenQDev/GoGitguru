package s3util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDecompressDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "decompress_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	// Create a directory at tmpDir path
	dirPath := filepath.Join(tmpDir, "dir")
	err = os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %s", err)
	}

	// Create a dummy tarball file in the temp directory
	tarballPath := filepath.Join(tmpDir, "test.tar.gz")
	err = createDummyTarball(tarballPath, dirPath)
	if err != nil {
		t.Fatalf("Failed to create dummy tarball: %s", err)
	}

	// Create a directory to extract to
	extractDir := filepath.Join(tmpDir, "extract")
	err = os.Mkdir(extractDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create extract directory: %s", err)
	}

	// Test the DecompressDirectory function
	err = DecompressDirectory(tarballPath, extractDir)
	if err != nil {
		t.Errorf("DecompressDirectory failed: %s", err)
	}

	// Add assertions to check the contents of the extract directory
	// ...
}

// createDummyTarball creates a dummy tarball file for testing
func createDummyTarball(path string, dirPath string) error {
	fmt.Println(path)
	// Create a new file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write some dummy data to the file
	_, err = file.WriteString("dummy data")
	if err != nil {
		return err
	}

	// Use gzip and tar to create a tarball
	cmd := exec.Command("tar", "-czf", path, dirPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
