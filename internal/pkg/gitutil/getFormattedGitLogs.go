package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"fmt"
	"path/filepath"
)

func GetFormattedGitLogs(prefixPath string, repo string, fromCommitDate string) ([]GitLog, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)

	if !isAGitRepo(fullRepoPath, prefixPath, repo) {
		return nil, fmt.Errorf("%s/%s is not a git repo", prefixPath, repo)
	}

	defaultCommitStartDate := "2020-01-01"
	if fromCommitDate == "" {
		fromCommitDate = defaultCommitStartDate
	}

	cmd := GitLogCommand(fullRepoPath, fromCommitDate)

	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("error running git log in %s: %s", fullRepoPath, err)
	}

	outStr := string(out)
	if outStr != "" && outStr[len(outStr)-1] == '\n' {
		outStr = outStr[:len(outStr)-1]
	}

	gitLogs, err := ProcessGitLogs(outStr)
	if err != nil {
		return nil, fmt.Errorf("error processing git logs for repo %s: %s", fullRepoPath, err)
	}

	return gitLogs, nil
}

func isAGitRepo(fullRepoPath string, prefixPath string, repo string) bool {
	cmdCheck := CheckIsAGitRepo(fullRepoPath)
	err := cmdCheck.Run()
	if err != nil {
		return false
	} else {
		return true
	}
}
