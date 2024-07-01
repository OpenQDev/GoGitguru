package server

import "net/http"

type HandlerAddDependencyPatternTest struct {
	name                          string
	expectedStatus                int
	firstRequestBody              HandlerAddDependencyPatternRequest
	firstExpectedSuccessResponse  HandlerAddDependencyPatternResponse
	secondRequestBody             HandlerAddDependencyPatternRequest
	secondExpectedSuccessResponse HandlerAddDependencyPatternResponse
	expectedErrorResponse         ErrorResponse
	shouldError                   bool
}

func avalidRepoUrls() HandlerAddDependencyPatternTest {
	const VALID_REPO_URLS = "VALID_REPO_URLS"

	firstFilePatternsAdded := HandlerAddDependencyPatternRequest{
		DependencyPatterns: []string{
			"package.json",
			`\/deployments\/`,
		},
		Creator: "MyGithubUser",
	}
	secondFilePatternsAdded := HandlerAddDependencyPatternRequest{
		DependencyPatterns: []string{
			"package.json",
			`\/deployments\/`,
			"foundry.toml",
			"thomas.tank",
			"thomas.tank2",
		},
		Creator: "MyGithubUser",
	}

	successReturnBody := HandlerAddDependencyPatternResponse{
		Accepted: "true",
	}

	secondReturnBody := HandlerAddDependencyPatternResponse{
		Accepted: "true",
	}

	validRepoUrls := HandlerAddDependencyPatternTest{
		name:                          VALID_REPO_URLS,
		expectedStatus:                http.StatusAccepted,
		firstRequestBody:              firstFilePatternsAdded,
		secondRequestBody:             secondFilePatternsAdded,
		firstExpectedSuccessResponse:  successReturnBody,
		secondExpectedSuccessResponse: secondReturnBody,
		shouldError:                   false,
	}

	return validRepoUrls
}

func aemptyRepoUrls() HandlerAddDependencyPatternTest {
	const EMPTY_REPO_URLS = "EMPTY_REPO_URLS"

	return HandlerAddDependencyPatternTest{
		name:                         EMPTY_REPO_URLS,
		expectedStatus:               http.StatusBadRequest,
		firstRequestBody:             HandlerAddDependencyPatternRequest{DependencyPatterns: []string{}},
		secondRequestBody:            HandlerAddDependencyPatternRequest{DependencyPatterns: []string{}},
		firstExpectedSuccessResponse: HandlerAddDependencyPatternResponse{},
		expectedErrorResponse:        ErrorResponse{Error: `error parsing JSON for: {"repo_urls":[]}`},
		shouldError:                  true,
	}
}

func HandlerAddDependencyPatternTestCases() []HandlerAddDependencyPatternTest {
	return []HandlerAddDependencyPatternTest{
		avalidRepoUrls(),
		aemptyRepoUrls(),
	}
}
