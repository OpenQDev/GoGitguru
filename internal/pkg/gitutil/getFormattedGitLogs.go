package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"fmt"
	"path/filepath"
)

type GitLog struct {
	CommitHash    string
	AuthorName    string
	AuthorEmail   string
	AuthorDate    int64
	CommitDate    int64
	CommitMessage string
	FilesChanged  int64
	Insertions    int64
	Deletions     int64
}

func GetFormattedGitLogs(prefixPath string, repo string, fromCommitDate string) ([]GitLog, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)

	if !isAGitRepo(fullRepoPath, prefixPath, repo) {
		return nil, fmt.Errorf("%s/%s is not a git repo", prefixPath, repo)
	}

	defaultCommitStartDate := "2020-01-01"
	if fromCommitDate == "" {
		fromCommitDate = defaultCommitStartDate
	}

	return nil, nil
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
