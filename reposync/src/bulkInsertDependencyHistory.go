package reposync

import (
	"context"

	"github.com/OpenQDev/GoGitguru/database"
)

func BulkInsertDependencyHistory(db *database.Queries, url string, dependencyId []int32, firstPresent []int64, lastRemoved []int64) error {

	params := database.BatchInsertRepoDependenciesParams{
		Url:     url,
		Column2: dependencyId,
		Column3: firstPresent,
		Column4: lastRemoved,
	}
	return db.BatchInsertRepoDependencies(context.Background(), params)
}
