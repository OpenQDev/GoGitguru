package reposync

import (
	"fmt"
	"main/internal/database"
	"main/internal/pkg/gitutil"
)

type GitLogParams struct {
	prefixPath     string
	repo           string
	repoUrl        string
	fromCommitDate string
	db             *database.Queries
}

func StoreGitLogs(params GitLogParams) (int, error) {
	startDate, err := ParseDate(params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	r, err := gitutil.OpenGitRepo(params.prefixPath, params.repo)
	if err != nil {
		return 0, err
	}

	log, err := gitutil.GetCommitHistory(r, startDate)
	if err != nil {
		return 0, err
	}

	numberOfCommits, err := gitutil.GetNumberOfCommits(params.prefixPath, params.repo)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s has %d commits\n", params.repoUrl, numberOfCommits)

	commitCount, err := PrepareCommitHistoryForBulkInsertion(numberOfCommits, log, params)
	if err != nil {
		return 0, err
	}

	return commitCount, nil
}
