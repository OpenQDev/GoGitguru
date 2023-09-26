package s3util

import (
	"fmt"
	"os"
	"os/exec"
)

// Decompress a tarball
func DecompressDirectory(archiveFilename string, directoryToExtractTo string) error {
	cmd := exec.Command("tar", "-xzf", archiveFilename, "-C", directoryToExtractTo)

	// Print any errors from running tar
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("error untarring %s to %s. failed with error: %s", archiveFilename, directoryToExtractTo, err)
	}

	return nil
}
