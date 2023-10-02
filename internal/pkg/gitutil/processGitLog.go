package gitutil

import (
	"fmt"
	"main/internal/pkg/logger"
	"strconv"
	"strings"
)

func ProcessGitLog(log string) (*GitLog, error) {
	lines := strings.Split(log, "\n")
	firstLine := strings.Split(lines[0], "-;-")

	authorDate, err := strconv.ParseInt(firstLine[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing author date: %s", err)
	}

	commitDate, err := strconv.ParseInt(firstLine[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing commit date: %s", err)
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

		logger.LogGreenDebug("file: %v", fileData)

		insertion := int64(0)
		deletion := int64(0)

		if fileData[0] != "-" {
			var err error
			insertion, err = strconv.ParseInt(fileData[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing insertions: %s", err)
			}
		}

		if fileData[1] != "-" {
			var err error
			deletion, err = strconv.ParseInt(fileData[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing deletions: %s", err)
			}
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

	return &output, nil
}
