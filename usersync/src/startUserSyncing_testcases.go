package usersync

import (
	"time"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/lib/pq"

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
			now = func() time.Time { return mockTime }
			// EXPECT - GetLatestUncheckedCommitPerAuthor
			rows := sqlmock.NewRows([]string{"commit_hash", "author_email", "author_date", "repo_url"}).
				AddRow("65062be663cc004b77ca8a3b13255bc5efa42f25", "andrew@openq.dev", 0, "https://github.com/OpenQDev/OpenQ-Workflows")
			mock.ExpectQuery("^-- name: GetLatestUncheckedCommitPerAuthor :many.*").WillReturnRows(rows)

			// EXPECT - InsertRestIdToEmail
			rows = sqlmock.NewRows([]string{"rest_id", "email"}).AddRow(restId, email)
			mock.ExpectQuery("^-- name: InsertRestIdToEmail :one.*").WithArgs(restId, email).WillReturnRows(rows)

			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}
			mock.ExpectQuery("^-- name: CheckGithubUserExists :one.*").WithArgs(author.User.Login).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
			// NOTE - this INTERNALID is generated upon insertion - so it will only appear in the return row
			// it will NOT appear in the call to InsertUser
			rows = sqlmock.NewRows([]string{
				"user_id", "first_use_date", "last_use_date", "dependency_id"})
			rows.AddRow(1, 2, 2, 1)
			mock.ExpectQuery(("^-- name: GetFirstAndLastCommit :one.*")).WithArgs().WillReturnRows(rows)
			rows = sqlmock.NewRows([]string{"first_commit_date", "last_commit_date", "user_id", "dependency_id"}).AddRow(1, 2, 1, 2)
			mock.ExpectQuery("^-- name: GetUserDependenciesByUpdatedAt :many.*").WithArgs(1609458600).WillReturnRows(rows)

			rows = sqlmock.NewRows([]string{
				"first_use_date", "last_use_date", "dependency_id", "user_id"}).AddRow(0, 1, 2, 3)
			mock.ExpectQuery("^-- name: GetUserDependenciesByUser :many.*").WithArgs(pq.Array([]int{1}), pq.Array([]int{2})).WillReturnRows(rows)

			column1 := []int{1}
			column2 := []int{2}
			column3 := []uint64{2}
			column4 := []uint64{1}

			mock.ExpectExec("^-- name: BulkInsertUserDependencies.").WithArgs(
				pq.Array(column1),
				pq.Array(column2),
				pq.Array(column3),
				pq.Array(column4),
				1609458600,
			).WillReturnResult(sqlmock.NewResult(1, 1))
			rows = sqlmock.NewRows([]string{
				"internal_id"}).AddRow(0)
			mock.ExpectExec("^-- name: UpsertRepoToUserById :exec.*").WithArgs("https://github.com/OpenQDev/OpenQ-Workflows", "{1}", "{1}", "{2}").WillReturnResult(sqlmock.NewResult(1, 1))
			/*	mock.ExpectQuery("^-- name: InsertUser :one.*").WithArgs(
				author.User.GithubRestID,
				author.User.GithubGraphqlID,
				author.User.Login,
				author.User.Name,
				author.User.Email,
				author.User.AvatarURL,
				author.User.Company,
				author.User.Location,
				author.User.Bio,
				author.User.Blog,
				author.User.Hireable,
				author.User.TwitterUsername,
				author.User.Followers.TotalCount,
				author.User.Following.TotalCount,
				"User",
				createdAt,
				updatedAt,
			).WillReturnRows(rows)*/
		},
	}
}

func StartUserSyncingTestCases() []StartUserSyncingTestCase {
	return []StartUserSyncingTestCase{
		startUserSyncingTest1(),
	}
}
