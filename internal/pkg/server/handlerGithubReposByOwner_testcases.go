package server

type HandlerGithubReposByOwnerTestCase struct {
	name           string
	owner          string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

func unauthorized() HandlerGithubReposByOwnerTestCase {
	const SHOULD_401 = "SHOULD_401"
	return HandlerGithubReposByOwnerTestCase{
		name:           SHOULD_401,
		owner:          "DRM-Test-Organization",
		expectedStatus: 401,
		authorized:     false,
		shouldError:    true,
	}
}

func sucess() HandlerGithubReposByOwnerTestCase {
	const SHOULD_STORE_ALL_REPOS_FOR_ORG = "SHOULD_STORE_ALL_REPOS_FOR_ORG"
	return HandlerGithubReposByOwnerTestCase{
		name:           SHOULD_STORE_ALL_REPOS_FOR_ORG,
		owner:          "DRM-Test-Organization",
		expectedStatus: 200,
		authorized:     true,
		shouldError:    false,
	}
}

func HandlerGithubReposByOwnerTestCases() []HandlerGithubReposByOwnerTestCase {
	return []HandlerGithubReposByOwnerTestCase{
		unauthorized(),
		sucess(),
	}
}
