package util

import (
	"os"
	"path/filepath"
)

func DeleteLocalRepoAndTarball(prefixPath string, repo string) error {
	err := os.RemoveAll(filepath.Join(prefixPath, repo))
	if err != nil {
		return err
	}
	err = os.RemoveAll(filepath.Join(prefixPath, repo+".tar.gz"))
	if err != nil {
		return err
	}
	return nil
}
