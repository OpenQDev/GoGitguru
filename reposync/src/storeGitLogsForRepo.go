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
	rawDependencies, err := params.db.GetDependencies(context.Background(), params.repoUrl)
	if err != nil {
		println("error getting rawDependencies")
	}
	dependenciesFiles := make([]string, len(rawDependencies))
	dependencyNames := make([]string, len(rawDependencies))
	dependencies := []database.Dependency{}
	for i, rawDependency := range rawDependencies {
		dependency := database.Dependency{
			DependencyName: rawDependency.DependencyName.String,
			DependencyFile: rawDependency.DependencyFile.String,
			InternalID:     rawDependency.InternalID.Int32,
		}
		dependencies = append(dependencies, dependency)

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
	dependencyHistory, commitList, err := gitutil.GitDependencyHistory(repoDir, dependencies)
	if err != nil {
		return 0, err
	}

	dependencyHistoryObjects, err := PrepareDependencyHistoryForBulkInsertion(dependencyHistory, dependencies, params.repoUrl)

	if err != nil {
		return 0, err
	}

	numberOfCommits, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	err = BulkInsertDependencyHistory(params.db, params.repoUrl, dependencyHistoryObjects.DependencyId, dependencyHistoryObjects.DateFirstPresent, dependencyHistoryObjects.DateLastRemoved)
	if err != nil {
		return 0, fmt.Errorf("error storing dependency history for %s: %s", params.repoUrl, err)
	}

	fmt.Printf("%s has %d commits to sync\n", params.repoUrl, numberOfCommits)

	if numberOfCommits == 0 {
		return 0, nil
	}

	commitObjects, err := PrepareCommitHistoryForBulkInsertion(numberOfCommits, commitList, params)

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
