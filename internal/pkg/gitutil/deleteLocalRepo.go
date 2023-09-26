package gitutil

import (
	"os"
	"path/filepath"
)

func DeleteLocalRepo(prefixPath string, repo string) error {
	err := os.RemoveAll(filepath.Join(prefixPath, repo))

	if err != nil {
		return err
	}

	return nil
}
