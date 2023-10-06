package gitutil

import (
	"fmt"
	"os/exec"
)

/*
git clone https://github.com/OpenQDev/OpenQ-DRM repos/OpenQ-DRM
*/
func GitCloneCommand(cloneString string, cloneDestination string) *exec.Cmd {
	return exec.Command("git", "clone", cloneString, cloneDestination, "--single-branch")
}

/*
git -C repos/OpenQ-Workflows rev-list --all --count
*/
func GitCommitCount(repoDir string) *exec.Cmd {
	return exec.Command("git", "-C", repoDir, "rev-list", "--all", "--count")
}

/*
git log -p --raw --unified=0 -i -S'github.com/lib/pq v1.10.9' go.mod
git log -p --raw --unified=0 -i -S'<dependency-name>' <path-to-dependency-file>
*/
func GitDepFileHistory(repoDir string, dependencyName string, dependencyFilePath string) *exec.Cmd {
	gitLogCommand := fmt.Sprintf("git log -p --raw --unified=0 -i -S'%s' %s", dependencyName, dependencyFilePath)
	return exec.Command("sh", "-c", gitLogCommand)
}
