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
	gitLogs := make([]GitLog, 100)

	emptyNewline := "\n\n"
	logEntries := strings.Split(logs, emptyNewline)

	for _, logEntry := range logEntries {
		gitLog := ProcessGitLog(logEntry)
		gitLogs = append(gitLogs, gitLog)
	}

	return gitLogs
}

func ProcessGitLog(log string) GitLog {
	lines := strings.Split(log, "\n")
	firstLine := strings.Split(lines[0], "-;-")

	output := GitLog{
		CommitHash:    firstLine[0],
		AuthorName:    firstLine[1],
		AuthorEmail:   firstLine[2],
		AuthorData:    firstLine[3],
		CommitDate:    firstLine[4],
		CommitMessage: lines[1],
	}

	return output
}
