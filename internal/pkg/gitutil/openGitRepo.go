package gitutil

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func OpenGitRepo(prefixPath string, organization string, repo string) (*git.Repository, error) {
	fullRepoPath := filepath.Join(prefixPath, organization, repo)
	return git.PlainOpen(fullRepoPath)
}
