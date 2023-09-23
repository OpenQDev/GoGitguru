package util

import (
	"fmt"
	"os"
	"os/exec"
)

func CloneRepo(organization string, repo string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", organization, repo))

	// This allows you to see the stdout and stderr of the command being run on the host machine
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
