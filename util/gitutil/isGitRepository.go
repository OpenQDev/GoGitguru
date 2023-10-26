package gitutil

import (
	"os"
	"path/filepath"
)

func IsGitRepository(prefixPath string, organization string, repo string) bool {
	fullRepoPath := filepath.Join(prefixPath, organization, repo)
	_, err := os.Stat(fullRepoPath + "/.git")
	return !os.IsNotExist(err)
}
