package server

import (
	"context"

	"github.com/OpenQDev/GoGitguru/database"
)

func BulkInsertRepoToDependencies(
	db *database.Queries,
	GithubRepoId []int32,
	DependencyName string,
	FirstCommitDate []int64,
	DateAdded []int64,
	DateRemoved []int64,

) error {
	// create an array with same lenght as the other arrays and make each memmber have value of DependencyName

	DependencyNameArray := make([]string, len(GithubRepoId))
	for i := 0; i < len(GithubRepoId); i++ {
		DependencyNameArray[i] = DependencyName
	}
	params := database.BulkInsertRepoDependencyInfoParams{
		Column1: GithubRepoId,
		Column2: DependencyNameArray,
		Column3: FirstCommitDate,
		Column4: DateAdded,
		Column5: DateRemoved,
	}

	err := db.BulkInsertRepoDependencyInfo(context.Background(), params)
	return err
}
