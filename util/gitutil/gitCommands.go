package gitutil

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

/*
git clone https://github.com/OpenQDev/OpenQ-DRM repos/OpenQ-DRM
*/
func GitCloneCommand(cloneString string, cloneDestination string) *exec.Cmd {
	return exec.Command("git", "clone", cloneString, cloneDestination, "--single-branch")
}

/*
git clone with a 1 minute timeout
*/
func GitCloneCommandWithTimeout(cloneString string, cloneDestination string) *exec.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "git", "clone", cloneString, cloneDestination, "--single-branch")
}

/*
git pull
*/
func GitPullCommand(pullDestination string) *exec.Cmd {
	return exec.Command("git", "-C", pullDestination, "pull")
}

/*
git -C . ls-files | grep -E 'go.mod'
git -C <repoDir> ls-files | grep -E '<gitGrepExists>'
*/
func LogDependencyFiles(repoDir string, gitGrepExists string) *exec.Cmd {
	cmd := fmt.Sprintf("git -C %s ls-files | grep -E '%s'", repoDir, gitGrepExists)
	return exec.Command("/bin/sh", "-c", cmd)
}

/*
git -C ./mock/openqdev/openq-coinapi log -p --raw --unified=0 -i -S'github.com/lib/pq v1.10.9' go.mod
git log -p --raw --unified=0 -i -S'<dependency-name>' <path-to-dependency-file>
git -C ./mock/openqdev/openq-coinapi log -p --raw --unified=0 -i -S'axios' '**package.json**' '**utils/package.json**'
*/
func GitDepFileHistory(repoDir string, dependencyName string, dependencyFilePaths string) *exec.Cmd {
	cmd := fmt.Sprintf("git -C %s log -p --raw --unified=0 -i -S'%s' %s", repoDir, dependencyName, dependencyFilePaths)
	return exec.Command("/bin/sh", "-c", cmd)
}

/*
git -C {repo_dir} grep -rli '{dependency_searched}' -- {files_paths_formatted}
git -C ./mock/openqdev/openq-coinapi grep -rli 'axios' -- '**package.json**' '**utils/package.json**'
*/
func GitDepFileToday(repoDir string, dependencyName string, dependencyFilePaths string) *exec.Cmd {
	cmd := fmt.Sprintf("git -C %s grep -rli '%s' -- %s", repoDir, dependencyName, dependencyFilePaths)
	return exec.Command("/bin/sh", "-c", cmd)
}

/*
if filepaths don't exists -> errors
if dependency doesn't exist -> empty return
*
*/
func GitDependencySearch(repoDir string, dependencySearched string, filesPathsFormatted string) *exec.Cmd {
	cmd := fmt.Sprintf("git -C %s grep -rli '%s' -- %s", repoDir, dependencySearched, filesPathsFormatted)
	return exec.Command("/bin/sh", "-c", cmd)
}
