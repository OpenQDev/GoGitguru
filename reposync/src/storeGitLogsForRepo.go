package reposync

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
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
	dependencies, err := GetRepoDependencies(params.db, params.repoUrl)
	if err != nil {
		println("error getting rawDependencies")
	}
	commitList, err := CreateStartWithLatestCommitList(repoDir)
	if err != nil {
		return 0, fmt.Errorf("error getting commit list %s: %s", params.repoUrl, err)
	}

	numberOfCommits, err := GetNumberOfCommitsPerDependency(dependencies, params)
	if err != nil {
		return 0, fmt.Errorf("error getting number of commits for %s: %s", params.repoUrl, err)
	}

	dependencyHistoryObjects, commitObject, err := GetObjectsFromCommitList(params, dependencies, commitList, numberOfCommits)
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

	return numberOfCommits.ToSync, nil
}
