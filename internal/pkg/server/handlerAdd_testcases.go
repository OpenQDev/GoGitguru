package server

import "net/http"

type HandlerAddTest struct {
	name                          string
	expectedStatus                int
	requestBody                   HandlerAddRequest
	expectedSuccessResponse       HandlerAddResponse
	secondExpectedSuccessResponse HandlerAddResponse
	expectedErrorResponse         ErrorResponse
	shouldError                   bool
}

func validRepoUrls() HandlerAddTest {
	const VALID_REPO_URLS = "VALID_REPO_URLS"
	targetRepos := []string{"https://github.com/org/repo1", "https://github.com/org/repo2"}

	twoReposRequest := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	successReturnBody := HandlerAddResponse{
		Accepted:       targetRepos,
		AlreadyInQueue: []string{},
	}

	secondReturnBody := HandlerAddResponse{
		Accepted:       []string{},
		AlreadyInQueue: targetRepos,
	}

	validRepoUrls := HandlerAddTest{
		name:                          VALID_REPO_URLS,
		expectedStatus:                http.StatusAccepted,
		requestBody:                   twoReposRequest,
		expectedSuccessResponse:       successReturnBody,
		secondExpectedSuccessResponse: secondReturnBody,
		shouldError:                   false,
	}

	return validRepoUrls
}

func emptyRepoUrls() HandlerAddTest {
	const EMPTY_REPO_URLS = "EMPTY_REPO_URLS"

	return HandlerAddTest{
		name:                    EMPTY_REPO_URLS,
		expectedStatus:          http.StatusBadRequest,
		requestBody:             HandlerAddRequest{RepoUrls: []string{}},
		expectedSuccessResponse: HandlerAddResponse{},
		expectedErrorResponse:   ErrorResponse{Error: `error parsing JSON for: {"repo_urls":[]}`},
		shouldError:             true,
	}
}

func HandlerAddTestCases() []HandlerAddTest {
	return []HandlerAddTest{
		validRepoUrls(),
		emptyRepoUrls(),
	}
}
