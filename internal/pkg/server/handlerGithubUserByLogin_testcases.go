package server

import "net/http"

type HandlerGithubUserByLoginTestCase struct {
	title          string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

const owner = "FlacoJones"

func should401() HandlerGithubUserByLoginTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerGithubUserByLoginTestCase{
		title:          UNAUTHORIZED,
		login:          login,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
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
	}
}

func HandlerGithubUserByLoginTestCases() []HandlerGithubUserByLoginTestCase {
	return []HandlerGithubUserByLoginTestCase{
		should401(),
		valid(),
	}
}
