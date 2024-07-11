package usersync

import (
	"context"
	"database/sql"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func GetReposToUsers(db *database.Queries, UpsertRepoToUserByIdParams *database.UpsertRepoToUserByIdParams, internal_id int32, author GithubGraphQLAuthor) error {

	userCommits, err := db.GetFirstAndLastCommit(context.Background(), sql.NullString{String: author.Email, Valid: true})

	if err != nil {
		logger.LogError("error occured while getting first and last commit for user %s: %s", author.Email, err)
	}
	alreadySet := slices.Contains(UpsertRepoToUserByIdParams.InternalIds, internal_id)
	if alreadySet {

		//change the first and last commit dates if the current commit is earlier or later than the current first and last commit dates
		for index, id := range UpsertRepoToUserByIdParams.InternalIds {
			if id == internal_id {
				if userCommits.FirstCommitDate.(int64) < UpsertRepoToUserByIdParams.FirstCommitDates[index] {
					UpsertRepoToUserByIdParams.FirstCommitDates[index] = userCommits.FirstCommitDate.(int64)
				}
				if userCommits.LastCommitDate.(int64) > UpsertRepoToUserByIdParams.LastCommitDates[index] {
					UpsertRepoToUserByIdParams.LastCommitDates[index] = userCommits.LastCommitDate.(int64)
				}
			}
		}
	}

	if !alreadySet {
		UpsertRepoToUserByIdParams.InternalIds = append(UpsertRepoToUserByIdParams.InternalIds, internal_id)

		newLastCommitDate, ok := userCommits.LastCommitDate.(int64)
		if !ok {
			newLastCommitDate = 0
		}

		newFirstCommitDate, ok := userCommits.FirstCommitDate.(int64)
		if !ok {
			newFirstCommitDate = 0
		}
		UpsertRepoToUserByIdParams.LastCommitDates = append(UpsertRepoToUserByIdParams.LastCommitDates, newLastCommitDate)
		UpsertRepoToUserByIdParams.FirstCommitDates = append(UpsertRepoToUserByIdParams.FirstCommitDates, newFirstCommitDate)
	}
	return err
}
