package s3util

import (
	"fmt"
	"os"
	"os/exec"
)

// Create a tarball
func CompressDirectory(archiveFilename string, directoryToArchive string) (string, error) {
	cmd := exec.Command("tar", "-czf", archiveFilename, directoryToArchive)

	// Print any errors from running tar
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("error tarring %s to %s. failed with error: %s", directoryToArchive, archiveFilename, err))
	}

	return fmt.Sprintf("%s.tar.gz", directoryToArchive), nil
}
