package gitutil

import (
	"fmt"
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

func ProcessGitLogs(logs string) []GitLog {
	gitLogs := make([]GitLog, 0)

	emptyNewline := "\n\n"
	logEntries := strings.Split(logs, emptyNewline)

	for _, logEntry := range logEntries {
		gitLog := ProcessGitLog(logEntry)
		gitLogs = append(gitLogs, gitLog)
	}

	for _, gitLog := range gitLogs {
		fmt.Printf("CommitHash: %s\n", gitLog.CommitHash)
		fmt.Printf("AuthorName: %s\n", gitLog.AuthorName)
		fmt.Printf("AuthorEmail: %s\n", gitLog.AuthorEmail)
		fmt.Printf("AuthorDate: %d\n", gitLog.AuthorDate)
		fmt.Printf("CommitDate: %d\n", gitLog.CommitDate)
		fmt.Printf("CommitMessage: %s\n", gitLog.CommitMessage)
		fmt.Printf("FilesChanged: %d\n", gitLog.FilesChanged)
		fmt.Printf("Insertions: %d\n", gitLog.Insertions)
		fmt.Printf("Deletions: %d\n", gitLog.Deletions)
		fmt.Println("-----------------------------")
	}

	return gitLogs
}
