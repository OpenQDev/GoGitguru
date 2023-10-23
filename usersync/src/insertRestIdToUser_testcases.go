package usersync

import (
	"util/githubGraphQL"

	"github.com/DATA-DOG/go-sqlmock"
)

type InsertIntoRestIdToUserTestCase struct {
	name        string
	author      githubGraphQL.GithubGraphQLAuthor
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock, repo githubGraphQL.GithubGraphQLAuthor)
}

func fooo() InsertIntoRestIdToUserTestCase {
	const SHOULD_STORE_USER_TO_REPO_ID = "SHOULD_STORE_USER_TO_REPO_ID"
	const email = "abc123@gmail.com"
	const restId = 123

	user := githubGraphQL.GithubGraphQLUser{GithubRestID: restId}

	author := githubGraphQL.GithubGraphQLAuthor{
		Email: email,
		User:  user,
	}

	return InsertIntoRestIdToUserTestCase{
		name:        SHOULD_STORE_USER_TO_REPO_ID,
		author:      author,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, author githubGraphQL.GithubGraphQLAuthor) {
			rows := sqlmock.NewRows([]string{"rest_id", "email"}).AddRow(restId, email)
			mock.ExpectQuery("^-- name: InsertRestIdToEmail :one.*").WithArgs(restId, email).WillReturnRows(rows)
		},
	}
}

func InsertIntoRestIdToUserTestCases() []InsertIntoRestIdToUserTestCase {
	return []InsertIntoRestIdToUserTestCase{
		fooo(),
	}
}
