package reposync

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
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

func StoreGitLogsAndDepsHistoryForRepo(params GitLogParams, resyncAll bool) (int, error) {

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

	dependencyFileRecords, err := params.db.GetAllFilePatterns(context.Background())

	if err != nil {
		return 0, fmt.Errorf("error getting all file patterns for %s: %s", params.repoUrl, err)
	}

	dependencyFiles := []string{}
	for _, dep := range dependencyFileRecords {
		dependencyFiles = append(dependencyFiles, dep.Pattern)
	}

	dependencyHistoryObjects, commitObject, usersToReposObject, numberOfCommitsToSync, err := GetObjectsFromCommitList(params, commitList, numberOfCommitsToSync, currentDependencies, dependencyFiles, resyncAll)
	if resyncAll {
		logger.LogGreenDebug("resyncing all deps for %s", dependencyHistoryObjects.Url, dependencyHistoryObjects.Dependencynames, dependencyHistoryObjects.UpdatedAt.Int64)
	}

	if err != nil {
		return 0, fmt.Errorf("error getting structs from commit list %s: %s", params.repoUrl, err)
	}

	insertByIdParams := GetUpsertRepoByIdsParams(params, usersToReposObject)

	err = params.db.UpsertRepoToUserById(context.Background(), insertByIdParams)

	if err != nil {
		return 0, fmt.Errorf("error storing users to repo for %s: %s", params.repoUrl, err)
	}

	err = params.db.BatchInsertRepoDependencies(context.Background(), dependencyHistoryObjects)

	if err != nil {
		fmt.Printf("error storing dependency history for %s: %s", params.repoUrl, err)
		fmt.Println("dependencyHistoryObjects", dependencyHistoryObjects)
	}

	err = params.db.BulkInsertCommits(context.Background(), commitObject)

	if err != nil {
		return 0, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}
	return numberOfCommitsToSync, nil
}
