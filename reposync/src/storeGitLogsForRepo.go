package reposync

import (
	"context"
	"fmt"
	"path/filepath"
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

// from commitDate should be the date of the last commit that was synced for the repository or any of the dependencies.

func StoreGitLogsAndDepsHistoryForRepo(params GitLogParams) (int, error) {

	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)

	commitList, err := CreateCommitList(repoDir)
	if err != nil {
		return 0, fmt.Errorf("error getting commit list %s: %s", params.repoUrl, err)
	}

	numberOfCommitsToSync, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, fmt.Errorf("error getting number of commits for %s: %s", params.repoUrl, err)
	}

	currentDependencies, err := params.db.GetRepoDependenciesByURL(context.Background(), params.repoUrl)

	if err != nil {
		return 0, fmt.Errorf("error getting current dependencies for %s: %s", params.repoUrl, err)
	}

	dependencyHistoryObjects, commitObject, numberOfCommitsToSync, err := GetObjectsFromCommitList(params, commitList, numberOfCommitsToSync, currentDependencies)

	if err != nil {
		return 0, err
	}
	err = params.db.BatchInsertRepoDependencies(context.Background(), dependencyHistoryObjects)
	if err != nil {
		return 0, fmt.Errorf("error storing dependency history for %s: %s", params.repoUrl, err)
	}
	err = params.db.BulkInsertCommits(context.Background(), commitObject)

	if err != nil {
		return 0, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}
	return numberOfCommitsToSync, nil
}
