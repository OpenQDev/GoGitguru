package server

import "net/http"

type HandlerDependencyHistoryTestCase struct {
	name                              string
	shouldError                       bool
	expectedStatus                    int
	requestBody                       DependencyHistoryRequest
	expectedDependencyHistroyResponse DependencyHistoryResponse
}

func isNotAGitRepository() HandlerDependencyHistoryTestCase {
	const nonExistentRepoUrl = "https://github.com/IDont/Exist"

	NOT_A_GIT_REPOSITORY := "NOT_A_GIT_REPOSITORY"

	return HandlerDependencyHistoryTestCase{
		name:           NOT_A_GIT_REPOSITORY,
		shouldError:    true,
		expectedStatus: http.StatusNotFound,
		requestBody: DependencyHistoryRequest{
			RepoUrl:            nonExistentRepoUrl,
			FilePaths:          []string{},
			DependencySearched: "foo",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponse{},
	}
}

func largeFrontend() HandlerDependencyHistoryTestCase {
	const openqFrontend = "https://github.com/OpenQDev/OpenQ-Frontend"

	LARGE_FRONTEND := "LARGE_FRONTEND"

	return HandlerDependencyHistoryTestCase{
		name:           LARGE_FRONTEND,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrl:            openqFrontend,
			FilePaths:          []string{},
			DependencySearched: "web3",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponse{},
	}
}

func HandlerDependencyHistoryTestCases() []HandlerDependencyHistoryTestCase {
	return []HandlerDependencyHistoryTestCase{
		isNotAGitRepository(),
		largeFrontend(),
	}
}
