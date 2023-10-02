package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"main/internal/pkg/logger"
	"path/filepath"
)

func GetFormattedGitLogs(prefixPath string, repo string, fromCommitDate string) []GitLog {
	fullRepoPath := filepath.Join(prefixPath, repo)

	checkIsAGitRepo(fullRepoPath, prefixPath, repo)

	defaultCommitStartDate := "2020-01-01"
	if fromCommitDate == "" {
		fromCommitDate = defaultCommitStartDate
	}

	cmd := GitLogCommand(fullRepoPath, fromCommitDate)

	out, err := cmd.Output()

	if err != nil {
		logger.LogFatalRedAndExit("error running git log in %s: %s", fullRepoPath, err)
	}

	outStr := string(out)
	if outStr != "" && outStr[len(outStr)-1] == '\n' {
		outStr = outStr[:len(outStr)-1]
	}

	gitLogs := ProcessGitLogs(outStr)

	return gitLogs
}

func checkIsAGitRepo(fullRepoPath string, prefixPath string, repo string) {
	cmdCheck := CheckIsAGitRepo(fullRepoPath)
	err := cmdCheck.Run()
	if err != nil {
		logger.LogFatalRedAndExit("%s/%s is not a git repository: %s", prefixPath, repo, err)
	}
}
