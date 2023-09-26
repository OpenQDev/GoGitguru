package gitutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteLocalRepo(prefixPath string, repo string) error {
	path := filepath.Join(prefixPath, repo)
	fmt.Println(path)
	err := os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}
