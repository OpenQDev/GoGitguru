package server

import "net/http"

type HandlerGithubUserCommitsTestCase struct {
	name           string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
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
	}
}

func HandlerGithubUserCommitsTestCases() []HandlerGithubUserCommitsTestCase {
	return []HandlerGithubUserCommitsTestCase{
		notAuthorized(),
	}
}
