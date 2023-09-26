package gitutil

import (
	"strings"
)

type GitLog struct {
	CommitHash    string
	AuthorName    string
	AuthorEmail   string
	AuthorData    string
	CommitDate    string
	CommitMessage string
}

func ProcessGitLogs(logs string) []GitLog {
	gitLogs := make([]GitLog, 0)

	emptyNewline := "\n\n"
	logEntries := strings.Split(logs, emptyNewline)

	for _, logEntry := range logEntries {
		gitLog := ProcessGitLog(logEntry)
		gitLogs = append(gitLogs, gitLog)
	}

	return gitLogs
}
