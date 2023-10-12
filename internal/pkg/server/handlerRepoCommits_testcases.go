package server

import "net/http"

type HandlerRepoCommitsTestCase struct {
	name           string
	login          string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

func foo() HandlerRepoCommitsTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerRepoCommitsTestCase{
		name:           UNAUTHORIZED,
		login:          login,
		expectedStatus: http.StatusBadRequest,
		authorized:     false,
		shouldError:    true,
	}
}

func HandlerRepoCommitsTestCases() []HandlerRepoCommitsTestCase {
	return []HandlerRepoCommitsTestCase{
		foo(),
	}
}
