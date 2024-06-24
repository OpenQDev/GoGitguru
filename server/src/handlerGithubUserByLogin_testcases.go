package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserByLoginTestCase struct {
	title          string
	login          string
	loginLower     string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, user User)
}

func should401() HandlerGithubUserByLoginTestCase {
	const userLogin = "FlacoJones"
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerGithubUserByLoginTestCase{
		title:          UNAUTHORIZED,
		login:          userLogin,
		loginLower:     strings.ToLower(userLogin),
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock, user User) {},
	}
}

func valid() HandlerGithubUserByLoginTestCase {
	const VALID = "VALID"
	const login = "DRM-Test-Organization"
	return HandlerGithubUserByLoginTestCase{
		title:          VALID,
		login:          login,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
		setupMock: func(mock sqlmock.Sqlmock, user User) {
			createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, user.UpdatedAt)

			rows := sqlmock.NewRows([]string{"internal_id"}).
				AddRow(1)

			mock.ExpectQuery("-- name: CheckGithubUserExists :one").WithArgs(user.Login).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))
			mock.ExpectQuery("^-- name: InsertUser :one.*").WithArgs(
				user.GithubRestID,
				user.GithubGraphqlID,
				strings.ToLower(user.Login),
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
				createdAt,
				updatedAt,
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
