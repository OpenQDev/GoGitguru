package gitutil

import (
	"os"
	"path/filepath"
)

func DeleteLocalRepo(prefixPath string, organization string, repo string) error {
	path := filepath.Join(prefixPath, organization, repo)

	err := os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}
