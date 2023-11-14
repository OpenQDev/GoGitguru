package gitutil

import (
	"fmt"
	"os/exec"
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

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok { // check if err is of type *os.ExitError
			exitCode := exitError.ExitCode()
			fmt.Printf("cmd.Run() failed with exit code: %d\n", exitCode)
		} else {
			fmt.Printf("cmd.Run() failed with %s\n", err)
		}
	}
	fmt.Printf("combined out:\n%s\n", string(output))

	return err
}
