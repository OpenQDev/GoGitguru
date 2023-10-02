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

	files := lines[2:]
	filesChanged := int64(len(files))

	insertions := int64(0)
	deletions := int64(0)

	for _, file := range files {
		if file == "" {
			continue
		}

		fileData := strings.Fields(file)

		insertion, err := strconv.ParseInt(fileData[0], 10, 64)
		if err != nil {
			logger.LogError("error parsing insertions", err)
		}

		deletion, err := strconv.ParseInt(fileData[1], 10, 64)
		if err != nil {
			logger.LogError("error parsing deletions", err)
		}

		insertions += insertion
		deletions += deletion
	}

	output := GitLog{
		CommitHash:    firstLine[0],
		AuthorName:    firstLine[1],
		AuthorEmail:   firstLine[2],
		AuthorDate:    authorDate,
		CommitDate:    commitDate,
		CommitMessage: lines[1],
		FilesChanged:  filesChanged,
		Insertions:    insertions,
		Deletions:     deletions,
	}

	return output
}
