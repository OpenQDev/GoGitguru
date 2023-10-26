package reposync

import (
	"fmt"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
)

type GitLogParams struct {
	prefixPath     string
	organization   string
	repo           string
	repoUrl        string
	fromCommitDate time.Time
	db             *database.Queries
}

func StoreGitLogsForRepo(params GitLogParams) (int, error) {
	r, err := gitutil.OpenGitRepo(params.prefixPath, params.organization, params.repo)
	if err != nil {
		return 0, err
	}

	log, err := gitutil.GetCommitHistory(r, params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	numberOfCommits, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s has %d commits to sync\n", params.repoUrl, numberOfCommits)

	commitObjects, err := PrepareCommitHistoryForBulkInsertion(numberOfCommits, log, params)

	if err != nil {
		return 0, err
	}

	err = BulkInsertCommits(
		params.db,
		commitObjects.CommitHash,
		commitObjects.Author,
		commitObjects.AuthorEmail,
		commitObjects.AuthorDate,
		commitObjects.CommitterDate,
		commitObjects.Message,
		commitObjects.Insertions,
		commitObjects.Deletions,
		commitObjects.FilesChanged,
		commitObjects.RepoUrls,
	)
	if err != nil {
		return 0, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}

	return numberOfCommits, nil
}
