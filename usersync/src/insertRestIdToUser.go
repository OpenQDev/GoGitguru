package usersync

import (
	"context"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
)

func insertIntoRestIdToUser(author GithubGraphQLAuthor, db *database.Queries) error {
	restId := author.User.GithubRestID

	author.Email = strings.Trim(author.Email, "\"")

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
