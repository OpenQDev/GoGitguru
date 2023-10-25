package usersync

import (
	"context"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/githubGraphQL"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func insertGithubUser(author githubGraphQL.GithubGraphQLAuthor, db *database.Queries) error {
	createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
	if err != nil && !createdAt.IsZero() {
		logger.LogError("error parsing time: %s", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, author.User.UpdatedAt)
	if err != nil && !createdAt.IsZero() {
		logger.LogError("error parsing time: %s", err)
	}

	authorParams := convertAuthorToInsertUserParams(author, createdAt, updatedAt)

	_, err = db.InsertUser(context.Background(), authorParams)
	return err
}
