package server

import (
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserByLoginTestCase struct {
	title          string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock)
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
		setupMock:      func(mock sqlmock.Sqlmock) {},
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
		setupMock:      func(mock sqlmock.Sqlmock) {},
	}
}

func HandlerGithubUserByLoginTestCases() []HandlerGithubUserByLoginTestCase {
	return []HandlerGithubUserByLoginTestCase{
		should401(),
		valid(),
	}
}
