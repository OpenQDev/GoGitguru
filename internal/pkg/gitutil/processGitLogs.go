package gitutil

import (
	"strings"
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

func ProcessGitLogs(logs string) ([]GitLog, error) {
	gitLogs := make([]GitLog, 0)

	emptyNewline := "\n\n"
	logEntries := strings.Split(logs, emptyNewline)

	for _, logEntry := range logEntries {
		gitLog, err := ProcessGitLog(logEntry)
		if err != nil {
			return nil, err
		}
		gitLogs = append(gitLogs, *gitLog)
	}

	return gitLogs, nil
}
