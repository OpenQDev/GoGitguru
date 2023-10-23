package usersync

import (
	"time"

	"github.com/OpenQDev/GoGitguru/util/githubGraphQL"
	"github.com/OpenQDev/GoGitguru/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

type StartUserSyncingTestCase struct {
	name        string
	author      githubGraphQL.GithubGraphQLAuthor
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock, repo githubGraphQL.GithubGraphQLAuthor)
}

func startUserSyncingTest1() StartUserSyncingTestCase {
	const SHOULD_STORE_USER = "SHOULD_STORE_USER"
	const email = "andrew@openq.dev"
	const restId = 93455288

	user := githubGraphQL.GithubGraphQLUser{
		GithubRestID:    93455288,
		GithubGraphqlID: "U_kgDOBZIDuA",
		Login:           "FlacoJones",
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
		UpdatedAt: "2023-10-10T15:52:33Z",
	}

	author := githubGraphQL.GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User:  user,
	}

	return StartUserSyncingTestCase{
		name:        SHOULD_STORE_USER,
		author:      author,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, author githubGraphQL.GithubGraphQLAuthor) {
			// EXPECT - GetLatestUncheckedCommitPerAuthor
			rows := sqlmock.NewRows([]string{"commit_hash", "author_email", "repo_url"}).
				AddRow("abc123", "andrew@openq.dev", "https://github.com/OpenQDev/OpenQ-Workflows")
			mock.ExpectQuery("^-- name: GetLatestUncheckedCommitPerAuthor :many.*").WillReturnRows(rows)

			// EXPECT - InsertRestIdToEmail
			rows = sqlmock.NewRows([]string{"rest_id", "email"}).AddRow(restId, email)
			mock.ExpectQuery("^-- name: InsertRestIdToEmail :one.*").WithArgs(restId, email).WillReturnRows(rows)

			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			updatedAt, err := time.Parse(time.RFC3339, author.User.UpdatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			// NOTE - this INTERNALID is generated upon insertion - so it will only appear in the return row
			// it will NOT appear in the call to InsertUser
			rows = sqlmock.NewRows([]string{
				"internal_id",
				"github_rest_id",
				"github_graphql_id",
				"login",
				"name",
				"email",
				"avatar_url",
				"company",
				"location",
				"bio",
				"blog",
				"hireable",
				"twitter_username",
				"followers",
				"following",
				"type",
				"created_at",
				"updated_at",
			}).AddRow(
				0,
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
			)
			mock.ExpectQuery("^-- name: InsertUser :one.*").WithArgs(
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
			).WillReturnRows(rows)
		},
	}
}

func StartUserSyncingTestCases() []StartUserSyncingTestCase {
	return []StartUserSyncingTestCase{
		startUserSyncingTest1(),
	}
}
