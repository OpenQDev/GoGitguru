package usersync

import (
	"database/sql"
	"time"

	"github.com/OpenQDev/GoGitguru/util/lib"
	"github.com/OpenQDev/GoGitguru/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

type StartUserSyncingTestCase struct {
	name        string
	author      GithubGraphQLAuthor
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock, repo GithubGraphQLAuthor)
}

func startUserSyncingTest1() StartUserSyncingTestCase {
	const SHOULD_STORE_USER = "SHOULD_STORE_USER"
	const email = "andrew@openq.dev"
	const restId = 93455288

	user := GithubGraphQLUser{
		GithubRestID:    93455288,
		GithubGraphqlID: "U_kgDOBZIDuA",
		Login:           "flacojones",
		Name:            "AndrewOBrien",
		Email:           "",
		AvatarURL:       "https://avatars.githubusercontent.com/u/93455288?u=fd1fb04b6ff2bf397f8353eafffc3bfb4bd66e84\u0026v=4",
		Company:         "",
		Location:        "",
		Hireable:        false,
		Bio:             "builder at OpenQ",
		Blog:            "",
		TwitterUsername: "",
		Followers: struct {
			TotalCount int `json:"totalCount"`
		}{
			TotalCount: 12,
		},
		Following: struct {
			TotalCount int `json:"totalCount"`
		}{
			TotalCount: 0,
		},
		CreatedAt: "2021-10-30T23:43:10Z",
		UpdatedAt: "2024-06-19T19:22:09Z",
	}

	author := GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User:  user,
	}

	return StartUserSyncingTestCase{
		name:        SHOULD_STORE_USER,
		author:      author,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, author GithubGraphQLAuthor) {

			mockTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			lib.Now = func() time.Time { return mockTime }
			// EXPECT - GetLatestUncheckedCommitPerAuthor
			rows := sqlmock.NewRows([]string{"commit_hash", "author_email", "author_date", "repo_url", "github_user_email"}).
				AddRow("65062be663cc004b77ca8a3b13255bc5efa42f25", "andrew@openq.dev", 0, "https://github.com/OpenQDev/OpenQ-Workflows", sql.NullString{Valid: false})
			mock.ExpectQuery("^-- name: GetLatestUncheckedCommitPerAuthor :many.*").WillReturnRows(rows)
			// NOTE - this INTERNALID is generated upon insertion - so it will only appear in the return row
			// it will NOT appear in the call to InsertUser
			// EXPECT - InsertRestIdToEmail			//
			mock.ExpectExec("^-- name: InsertRestIdToEmail :exec.*").WithArgs(restId, email).WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectQuery("^-- name: CheckGithubUserIdExists :one.*").WithArgs(93455288).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

			mock.ExpectQuery("^-- name: GetGithubUserByRestId :one.*").WithArgs(93455288).WillReturnRows(sqlmock.NewRows([]string{"internal_id"}).AddRow(1))
			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}
			rows = sqlmock.NewRows([]string{
				"user_id", "first_use_date"})
			rows.AddRow(1, 2)
			mock.ExpectQuery(("^-- name: GetFirstAndLastCommit :one.*")).WithArgs().WillReturnRows(rows)

			mock.ExpectExec("^-- name: UpsertRepoToUserById :exec.*").WithArgs("https://github.com/OpenQDev/OpenQ-Workflows", "{1}", "{1}", "{2}").WillReturnResult(sqlmock.NewResult(1, 1))

			// Expect - SetAllCommitsToChecked
			mock.ExpectExec("^-- name: SetAllCommitsToChecked :exec.*").WillReturnResult(sqlmock.NewResult(1, 1))

		},
	}
}

func StartUserSyncingTestCases() []StartUserSyncingTestCase {
	return []StartUserSyncingTestCase{
		startUserSyncingTest1(),
	}
}
