package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
)

func checkIndexedReposWithDependency(db *database.Queries, dep string, repos []database.GithubRepo) ([]database.DependenciesToRepo, error) {
	repoIds := make([]int32, 0)
	for _, repo := range repos {
		repoIds = append(repoIds, repo.InternalID)
	}

	result, err := db.QueryBulkRepoDependencyInfo(context.Background(), database.QueryBulkRepoDependencyInfoParams{
		Column1:        repoIds,
		DependencyName: sql.NullString{String: dep, Valid: true},
	})
	return result, err
}

func checkNonIndexedReposWithDependency(db *database.Queries, dep string, repos []database.GithubRepo, datesAddedCommits []int64, datesRemovedCommits []int64) (bool, int64, int64) {
	var lastDateAdded, lastDateRemoved int64 = 0, 0
	var lastDateAddedIso, lastDateRemovedIso string
	datesAddedISO := make([]string, len(datesAddedCommits))
	for i, v := range datesAddedCommits {
		if v >= lastDateAdded {
			lastDateAdded = v
			lastDateAddedIso = time.Unix(v, 0).Format(time.RFC3339)
		}
		datesAddedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
	}

	datesRemovedISO := make([]string, len(datesRemovedCommits))
	for i, v := range datesRemovedCommits {
		if v >= lastDateRemoved {
			lastDateRemoved = v
			lastDateRemovedIso = time.Unix(v, 0).Format(time.RFC3339)
		}

		datesRemovedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
	}
	fmt.Println(lastDateAddedIso, lastDateRemovedIso)
	return lastDateAdded > lastDateRemoved || (lastDateAdded != 0 && lastDateRemoved == 0), lastDateAdded, lastDateRemoved
}

func getHasRepoToDependency(indexedReposWithDependency []database.DependenciesToRepo, repoUrlObject database.GithubRepo) bool {
	for _, indexedRepoWithDependency := range indexedReposWithDependency {
		if indexedRepoWithDependency.InternalID == repoUrlObject.InternalID {
			return true

		}
	}
	return false
}
