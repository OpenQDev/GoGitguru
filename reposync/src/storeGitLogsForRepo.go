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

func StoreGitLogsAndDepsHistoryForRepo(params GitLogParams) (int, database.BulkInsertCommitsParams, []string, error) {

	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)

	commitList, err := CreateCommitList(repoDir)

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error getting commit list %s: %s", params.repoUrl, err)
	}

	numberOfCommitsToSync, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error getting number of commits for %s: %s", params.repoUrl, err)
	}

	currentDependencies, err := params.db.GetRepoDependenciesByURL(context.Background(), params.repoUrl)

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error getting current dependencies for %s: %s", params.repoUrl, err)
	}

	dependencyFileRecords, err := params.db.GetAllFilePatterns(context.Background())

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error getting all file patterns for %s: %s", params.repoUrl, err)
	}

	dependencyFiles := []string{}
	for _, dep := range dependencyFileRecords {
		dependencyFiles = append(dependencyFiles, dep.Pattern)
	}

	dependencyHistoryObjects, commitObject, usersToReposObject, numberOfCommitsToSync, err := GetObjectsFromCommitList(params, commitList, numberOfCommitsToSync, currentDependencies, dependencyFiles)

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error getting structs from commit list %s: %s", params.repoUrl, err)
	}

	insertByIdParams := GetUpsertRepoByIdsParams(params, usersToReposObject)

	err = params.db.UpsertRepoToUserById(context.Background(), insertByIdParams)

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error storing users to repo for %s: %s", params.repoUrl, err)
	}

	newRepoDependencies, err := params.db.BatchInsertRepoDependencies(context.Background(), dependencyHistoryObjects)

	if err != nil {
		fmt.Printf("error storing dependency history for %s: %s", params.repoUrl, err)
		fmt.Println("dependencyHistoryObjects", newRepoDependencies)
	}

	err = params.db.BulkInsertCommits(context.Background(), commitObject)

	if err != nil {
		return 0, database.BulkInsertCommitsParams{}, []string{}, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}
	return numberOfCommitsToSync, commitObject, newRepoDependencies, nil
}
