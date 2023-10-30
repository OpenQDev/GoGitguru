package gitutil

import (
	"regexp"
	"strings"
	"time"
)

type CommitSummary struct {
	CommitHash   string
	Author       string
	Date         time.Time
	Message      string
	AddedLines   int
	RemovedLines int
	FullContent  string
}

type DiffHistoryResult struct {
	CommitsSummary []CommitSummary
	DatesAdded     []time.Time
	DatesRemoved   []time.Time
}

func DiffHistoryObject(dependencyHistoryLogs string, dependencySearched string, dependencyTodayOutput string) DiffHistoryResult {
	patternCommit := regexp.MustCompile(`commit [a-f0-9]{40}`)
	commits := SplitWithDelimiters(dependencyHistoryLogs, patternCommit)[1:]

	patternAdded := regexp.MustCompile(`(?m)^\+\s.*` + regexp.QuoteMeta(dependencySearched) + `.*$`)
	patternRemoved := regexp.MustCompile(`(?m)^-\s.*` + regexp.QuoteMeta(dependencySearched) + `.*$`)
	results := DiffHistoryResult{}
	commitsSummary := make([]CommitSummary, 0)
	var commitHash string
	datesAdded := make([]time.Time, 0)
	datesRemoved := make([]time.Time, 0)

	for _, commit := range commits {
		if strings.HasPrefix(commit, "commit") {
			commitHash = strings.TrimPrefix(strings.TrimSpace(commit), "commit ")
			continue
		}

		lines := strings.Split(strings.TrimSpace(commit), "\n")
		authorLine := strings.TrimPrefix(lines[0], "Author: ")
		dateLine := strings.TrimPrefix(lines[1], "Date: ")
		messageLine := lines[3]
		addedLines := len(patternAdded.FindAllString(commit, -1))
		removedLines := len(patternRemoved.FindAllString(commit, -1))

		if dateLine == "" {
			continue
		}

		formattedDate, _ := time.Parse(time.RFC3339, dateLine)

		if addedLines > removedLines {
			datesAdded = append(datesAdded, formattedDate)
		} else if removedLines > addedLines {
			datesRemoved = append(datesRemoved, formattedDate)
		}

		commitsSummary = append(commitsSummary, CommitSummary{
			CommitHash:   commitHash,
			Author:       authorLine,
			Date:         formattedDate,
			Message:      messageLine,
			AddedLines:   addedLines,
			RemovedLines: removedLines,
			FullContent:  commit,
		})
	}

	results.CommitsSummary = commitsSummary
	results.DatesAdded = datesAdded
	if dependencyTodayOutput == "" {
		results.DatesRemoved = datesRemoved
	}

	return results
}
