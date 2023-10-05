package gitutil

import (
	"os/exec"
)

/*
git clone https://github.com/OpenQDev/OpenQ-DRM repos/OpenQ-DRM
*/
func GitCloneCommand(cloneString string, cloneDestination string) *exec.Cmd {
	return exec.Command("git", "clone", cloneString, cloneDestination)
}

/*
git -C repos/OpenQ-Workflows rev-list --all --count
*/
func GitCommitCount(repoDir string) *exec.Cmd {
	return exec.Command("git", "-C", repoDir, "rev-list", "--all", "--count")
}
