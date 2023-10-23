package usersync

import (
	"context"
	"database/database"
	"util/githubGraphQL"
)

func insertIntoRestIdToUser(author githubGraphQL.GithubGraphQLAuthor, db *database.Queries) error {
	restId := author.User.GithubRestID

	var params database.InsertRestIdToEmailParams
	if restId == 0 {
		params = database.InsertRestIdToEmailParams{
			Email: author.Email,
		}
	} else {
		params = database.InsertRestIdToEmailParams{
			RestID: int32(restId),
			Email:  author.Email,
		}
	}

	_, err := db.InsertRestIdToEmail(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}
