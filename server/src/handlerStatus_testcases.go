package server

import "net/http"

type HandlerStatusTest struct {
	name                    string
	expectedStatus          int
	requestBody             HandlerAddRequest
	expectedSuccessResponse HandlerStatusResponse
	shouldError             bool
}

func statusValidRepoUrls() HandlerStatusTest {
	const VALID_REPO_URLS = "VALID_REPO_URLS"
	targetRepos := []string{"https://github.com/org/repo1", "https://github.com/org/repo2"}

	twoReposRequest := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	successReturnBody := HandlerStatusResponse{}

	validRepoUrls := HandlerStatusTest{
		name:                    VALID_REPO_URLS,
		expectedStatus:          http.StatusAccepted,
		requestBody:             twoReposRequest,
		expectedSuccessResponse: successReturnBody,
		shouldError:             false,
	}

	return validRepoUrls
}

func HandlerStatusTestCases() []HandlerStatusTest {
	return []HandlerStatusTest{
		statusValidRepoUrls(),
	}
}
