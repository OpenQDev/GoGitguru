package gitutil

import (
	"os"
	"path/filepath"
)

// CloneRepo clones a git repository from GitHub and places it in a specified directory.
// It takes three parameters:
// - prefixPath: the directory where the cloned repository will be placed. NO TRAILING SLASH
// - organization: the name of the organization on GitHub that owns the repository.
// - repo: the name of the repository to be cloned.
// It returns an error if the cloning process fails.
func PullRepo(prefixPath string, organization string, repo string) error {
	pullDestination := filepath.Join(prefixPath, organization, repo)

	cmd := GitPullCommand(pullDestination)

	// This allows you to see the stdout and stderr of the command being run on the host machine
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
