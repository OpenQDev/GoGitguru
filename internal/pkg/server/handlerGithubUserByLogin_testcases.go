package server

import "net/http"

type HandlerGithubUserByLoginTestCase struct {
	title          string
	owner          string
	name           string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

const owner = "DRM-Test-Organization"
const repo = "DRM-Test-Repo"

func should401() HandlerGithubUserByLoginTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerGithubUserByLoginTestCase{
		title:          UNAUTHORIZED,
		owner:          owner,
		name:           repo,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
	}
}

func valid() HandlerGithubUserByLoginTestCase {
	const VALID = "VALID"
	return HandlerGithubUserByLoginTestCase{
		title:          VALID,
		owner:          owner,
		name:           repo,
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
