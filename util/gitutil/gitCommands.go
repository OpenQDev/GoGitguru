package gitutil

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	return exec.Command("git", "-C", pullDestination, "pull", "--strategy-option=theirs")
}

/*
git -C . ls-files | grep -E 'go.mod'
git -C <repoDir> ls-files | grep -E '<gitGrepExists>'
*/
func LogDependencyFiles(repoDir string, gitGrepExists string) ([]string, error) {
	// Open the repository
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	// Get the file tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get file tree: %w", err)
	}

	// Compile the regular expression
	re, err := regexp.Compile(gitGrepExists)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}

	// List the files and filter them based on the regular expression
	var files []string
	err = tree.Files().ForEach(func(f *object.File) error {
		if re.MatchString(f.Name) {
			files = append(files, f.Name)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
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
