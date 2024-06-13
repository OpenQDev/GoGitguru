package usersync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertAuthorToInsertUserParams(t *testing.T) {
	// ARRANGE
	user := GithubGraphQLUser{
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

	author := GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User:  user,
	}

	createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, user.UpdatedAt)

	// ACT
	result := convertAuthorToInsertUserParams(author, createdAt, updatedAt)

	// ASSERT
	assert.Equal(t, int32(93455288), result.GithubRestID)
	assert.Equal(t, "U_kgDOBZIDuA", result.GithubGraphqlID)
	assert.Equal(t, "flacojones", result.Login)
	assert.Equal(t, "AndrewOBrien", result.Name.String)
	assert.Equal(t, "", result.Email.String)
	assert.Equal(t, "https://avatars.githubusercontent.com/u/93455288?u=fd1fb04b6ff2bf397f8353eafffc3bfb4bd66e84\u0026v=4", result.AvatarUrl.String)
	assert.Equal(t, "", result.Company.String)
	assert.Equal(t, "", result.Location.String)
	assert.Equal(t, false, result.Hireable.Bool)
	assert.Equal(t, "builder at OpenQ", result.Bio.String)
	assert.Equal(t, "", result.Blog.String)
	assert.Equal(t, "", result.TwitterUsername.String)
	assert.Equal(t, int32(12), result.Followers.Int32)
	assert.Equal(t, int32(0), result.Following.Int32)
	assert.Equal(t, "User", result.Type)
	assert.Equal(t, createdAt, result.CreatedAt.Time)
	assert.Equal(t, updatedAt, result.UpdatedAt.Time)
}
