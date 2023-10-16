package server

import (
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserCommitsTestCase struct {
	name           string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock)
}

const login = "DRM-Test-Organization"

func notAuthorized() HandlerGithubUserCommitsTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerGithubUserCommitsTestCase{
		name:           UNAUTHORIZED,
		login:          login,
		expectedStatus: http.StatusBadRequest,
		authorized:     false,
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock) {},
	}
}

func HandlerGithubUserCommitsTestCases() []HandlerGithubUserCommitsTestCase {
	return []HandlerGithubUserCommitsTestCase{
		notAuthorized(),
	}
}
