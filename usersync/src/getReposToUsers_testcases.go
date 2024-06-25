package usersync

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/OpenQDev/GoGitguru/database"
)

type GetReposToUsersTestCase struct {
	name           string
	initialParams  database.UpsertRepoToUserByIdParams
	expectedParams database.UpsertRepoToUserByIdParams
	internal_id    int32
	author         GithubGraphQLAuthor
	setupMock      func(mock sqlmock.Sqlmock, author GithubGraphQLAuthor)
}

func getReposToUsersTest1() GetReposToUsersTestCase {
	const SHOULD_STORE_USER = "SHOULD_STORE_USER"
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

	return GetReposToUsersTestCase{
		name:   SHOULD_STORE_USER,
		author: author,
		initialParams: database.UpsertRepoToUserByIdParams{
			Url:              "https//:github.com/AndrewOBrien/GoGitguru",
			InternalIds:      []int32{},
			FirstCommitDates: []int64{},
			LastCommitDates:  []int64{},
		},
		expectedParams: database.UpsertRepoToUserByIdParams{
			Url:              "https//:github.com/AndrewOBrien/GoGitguru",
			InternalIds:      []int32{1},
			FirstCommitDates: []int64{2},
			LastCommitDates:  []int64{2},
		},
		internal_id: 1,
		setupMock: func(mock sqlmock.Sqlmock, author GithubGraphQLAuthor) {

			mockTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			now = func() time.Time { return mockTime }

			rows := sqlmock.NewRows([]string{
				"first_commit_date", "last_commit_date"})
			rows.AddRow(1, 2)
			mock.ExpectQuery(("^-- name: GetFirstAndLastCommit :one.*")).WithArgs(author.Email).WillReturnRows(rows)

		},
	}
}

func GetReposToUsersTestCases() []GetReposToUsersTestCase {
	return []GetReposToUsersTestCase{
		getReposToUsersTest1(),
	}
}
