package gitutil

import (
	"main/internal/pkg/logger"
	"strconv"
	"strings"
)

func ProcessGitLog(log string) GitLog {
	lines := strings.Split(log, "\n")
	firstLine := strings.Split(lines[0], "-;-")

	authorDate, err := strconv.ParseInt(firstLine[3], 10, 64)
	if err != nil {
		logger.LogError("error parsing author date", err)
	}

	commitDate, err := strconv.ParseInt(firstLine[4], 10, 64)
	if err != nil {
		logger.LogError("error parsing commit date", err)
	}

	output := GitLog{
		CommitHash:    firstLine[0],
		AuthorName:    firstLine[1],
		AuthorEmail:   firstLine[2],
		AuthorDate:    authorDate,
		CommitDate:    commitDate,
		CommitMessage: lines[1],
	}

	return output
}
