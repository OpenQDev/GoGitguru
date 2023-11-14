package gitutil

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "clone", cloneString, cloneDestination, "--single-branch")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out")
	}

	if err != nil {
		return err
	}

	if strings.Contains(out.String(), "empty repository") {
		return fmt.Errorf("%s/%s is an empty repository", organization, repo)
	}

	return nil
}
