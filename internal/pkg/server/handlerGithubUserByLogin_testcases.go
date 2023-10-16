package server

import (
	"net/http"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserByLoginTestCase struct {
	title          string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, user User)
}

const userLogin = "FlacoJones"

func should401() HandlerGithubUserByLoginTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerGithubUserByLoginTestCase{
		title:          UNAUTHORIZED,
		login:          userLogin,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock, user User) {},
	}
}

func valid() HandlerGithubUserByLoginTestCase {
	const VALID = "VALID"
	return HandlerGithubUserByLoginTestCase{
		title:          VALID,
		login:          login,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
		setupMock: func(mock sqlmock.Sqlmock, user User) {
			createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, user.UpdatedAt)

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "login", "name", "email", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type", "created_at", "updated_at"}).
				AddRow(1, user.GithubRestID, user.GithubGraphqlID, user.Login, user.Name, user.Email, user.AvatarURL, user.Company, user.Location, user.Bio, user.Blog, user.Hireable, user.TwitterUsername, user.Followers, user.Following, user.Type, createdAt, updatedAt)

			mock.ExpectQuery("-- name: CheckGithubUserExists :one").WithArgs(user.Login).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			mock.ExpectQuery("^-- name: InsertUser :one.*").WithArgs(
				user.GithubRestID,
				user.GithubGraphqlID,
				user.Login,
				user.Name,
				user.Email,
				user.AvatarURL,
				user.Company,
				user.Location,
				user.Bio,
				user.Blog,
				user.Hireable,
				user.TwitterUsername,
				user.Followers,
				user.Following,
				user.Type,
				user.CreatedAt,
				user.UpdatedAt,
			).WillReturnRows(rows)
		},
	}
}

func HandlerGithubUserByLoginTestCases() []HandlerGithubUserByLoginTestCase {
	return []HandlerGithubUserByLoginTestCase{
		should401(),
		valid(),
	}
}
