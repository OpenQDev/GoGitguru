package gitutil

import "strings"

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
