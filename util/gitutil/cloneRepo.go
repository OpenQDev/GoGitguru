package gitutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// CloneRepo clones a git repository from GitHub and places it in a specified directory.
// It takes three parameters:
// - prefixPath: the directory where the cloned repository will be placed. NO TRAILING SLASH
// - organization: the name of the organization on GitHub that owns the repository.
// - repo: the name of the repository to be cloned.
// It returns an error if the cloning process fails.
func CloneRepo(prefixPath string, organization string, repo string) error {
	cloneString := fmt.Sprintf("https://github.com/%s/%s.git", organization, repo)
	cloneDestination := filepath.Join(prefixPath, organization, repo)

	cmd := GitCloneCommand(cloneString, cloneDestination)

	// This allows you to see the stdout and stderr of the command being run on the host machine
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}
