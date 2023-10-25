package server

import "net/http"

type HandlerStatusTest struct {
	name                          string
	expectedStatus                int
	requestBody                   HandlerAddRequest
	expectedSuccessResponse       HandlerAddResponse
	secondExpectedSuccessResponse HandlerAddResponse
	expectedErrorResponse         ErrorResponse
	shouldError                   bool
}

func statusValidRepoUrls() HandlerStatusTest {
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

	validRepoUrls := HandlerStatusTest{
		name:                          VALID_REPO_URLS,
		expectedStatus:                http.StatusAccepted,
		requestBody:                   twoReposRequest,
		expectedSuccessResponse:       successReturnBody,
		secondExpectedSuccessResponse: secondReturnBody,
		shouldError:                   false,
	}

	return validRepoUrls
}

func HandlerStatusTestCases() []HandlerStatusTest {
	return []HandlerStatusTest{
		statusValidRepoUrls(),
	}
}
