package usersync

import (
	"github.com/DATA-DOG/go-sqlmock"
)

type InsertIntoRestIdToUserTestCase struct {
	name        string
	author      GithubGraphQLAuthor
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock, repo GithubGraphQLAuthor)
}

func fooo() InsertIntoRestIdToUserTestCase {
	const SHOULD_STORE_USER_TO_REPO_ID = "SHOULD_STORE_USER_TO_REPO_ID"
	const email = "abc123@gmail.com"
	const restId = 123

	user := GithubGraphQLUser{GithubRestID: restId}

	author := GithubGraphQLAuthor{
		Email: email,
		User:  user,
	}

	return InsertIntoRestIdToUserTestCase{
		name:        SHOULD_STORE_USER_TO_REPO_ID,
		author:      author,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, author GithubGraphQLAuthor) {

			mock.ExpectExec("^-- name: InsertRestIdToEmail :exec.*").WithArgs(restId, email).WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}
}

func InsertIntoRestIdToUserTestCases() []InsertIntoRestIdToUserTestCase {
	return []InsertIntoRestIdToUserTestCase{
		fooo(),
	}
}
