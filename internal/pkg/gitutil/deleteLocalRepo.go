package gitutil

import (
	"os"
	"path/filepath"
)

func DeleteLocalRepo(prefixPath string, repo string) error {
	path := filepath.Join(prefixPath, repo)

	err := os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}
