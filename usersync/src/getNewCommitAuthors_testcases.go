package usersync

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type GetNewCommitAuthorsTestCase struct {
	name        string
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock)
}

func getNewCommitAuthorsTestCase1() GetNewCommitAuthorsTestCase {
	return GetNewCommitAuthorsTestCase{
		name:        "Test Case 1",
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"author_email", "commit_hash", "repo_url", "author_date", "github_user_email"}).
				AddRow("andrew@openq.dev", "abc123", 0, "https://github.com/OpenQDev/OpenQ-Workflows", sql.NullString{Valid: false})
			mock.ExpectQuery("^-- name: GetLatestUncheckedCommitPerAuthor :many.*").WillReturnRows(rows)
		},
	}
}

func GetNewCommitAuthorsTestCases() []GetNewCommitAuthorsTestCase {
	return []GetNewCommitAuthorsTestCase{
		getNewCommitAuthorsTestCase1(),
	}
}
