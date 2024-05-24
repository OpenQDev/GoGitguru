package reposync

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
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

func StoreGitLogsAndDepsHistoryForRepo(params GitLogParams) (int, error) {
	

	dependencies, err := params.db.GetDependencies(context.Background())
	if err != nil {
		println("error getting dependencies")
	}
	dependenciesFiles := make([]string, len(dependencies))
	dependencyNames := make([]string, len(dependencies))
	for i, dependency := range dependencies {

		if slices.Contains(dependenciesFiles, dependency.DependencyFile) {
			dependenciesFiles[i] = dependency.DependencyFile
		}
		if slices.Contains(dependencyNames, dependency.DependencyName) {
			dependencyNames[i] = dependency.DependencyName
		}
	}

	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)
	if err != nil {
		return 0, err
	}

	dependencyHistory, log, err := gitutil.GitDependencyHistory(repoDir,  dependencies)
	if(err != nil) {
		return 0, err
	}
	println("dependency history", len(dependencyHistory))

	dependencyHistoryObjects, err := PrepareDependencyHistoryForBulkInsertion(dependencyHistory, dependencies, params.repoUrl)

	if err != nil {
		return 0, err
	}

	numberOfCommits, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s has %d commits to sync\n", params.repoUrl, numberOfCommits)

	if numberOfCommits == 0 {
		return 0, nil
	}

	commitObjects, err := PrepareCommitHistoryForBulkInsertion(numberOfCommits, log, params)

	if err != nil {
		return 0, err
	}
	err = BulkInsertDependencyHistory(params.db, dependencyHistoryObjects.RepoUrls, dependencyHistoryObjects.DependencyId, dependencyHistoryObjects.DateFirstPresent, dependencyHistoryObjects.DateLastRemoved)
	if err != nil {
		return 0, fmt.Errorf("error storing dependency history for %s: %s", params.repoUrl, err)
	}
	for commitIndex := 0; commitIndex < numberOfCommits; commitIndex++ {
		println(commitObjects.RepoUrls[commitIndex])
	}



	return 8, nil
}
