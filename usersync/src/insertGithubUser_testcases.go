package usersync

import (
	"time"

	"github.com/OpenQDev/GoGitguru/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

type InsertGithubUserTestCase struct {
	name        string
	author      GithubGraphQLAuthor
	shouldError bool
	setupMock   func(mock sqlmock.Sqlmock, repo GithubGraphQLAuthor)
}

func insertGithubUserTestCase1() InsertGithubUserTestCase {
	const SHOULD_STORE_GH_USER = "SHOULD_STORE_GH_USER"

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
		UpdatedAt: "2023-10-10T15:52:33Z",
	}

	author := GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "abc123@gmail.com",
		User:  user,
	}

	return InsertGithubUserTestCase{
		name:        SHOULD_STORE_GH_USER,
		author:      author,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, author GithubGraphQLAuthor) {
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

			mock.ExpectExec("^-- name: InsertUser :exec.*").WithArgs(
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
			).WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}
}

func InsertGithubUserTestCases() []InsertGithubUserTestCase {
	return []InsertGithubUserTestCase{
		insertGithubUserTestCase1(),
	}
}
