package s3util

import (
	"fmt"
	"os/exec"
)

// Create a tarball
func TarAndGzip(directoryToArchive string, tarballOutputPath string) (string, error) {
	err := exec.Command("tar", "-czf", directoryToArchive, tarballOutputPath).Run()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.tar.gz", tarballOutputPath), nil
}
