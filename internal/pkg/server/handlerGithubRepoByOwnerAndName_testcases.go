package server

type HandlerGithubRepoByOwnerAndNameTest struct {
	title          string
	owner          string
	name           string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

func shouldReturn401() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN = "SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN,
		owner:          "DRM-Test-Organization",
		name:           "DRM-Test-Repo",
		expectedStatus: 401,
		authorized:     false,
		shouldError:    true,
	}
}

func shouldReturnRepoForOwnerAndName() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_GET_REPO_FOR_ORG_AND_NAME = "SHOULD_GET_REPO_FOR_ORG_AND_NAME"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_GET_REPO_FOR_ORG_AND_NAME,
		owner:          "DRM-Test-Organization",
		name:           "DRM-Test-Repo",
		expectedStatus: 200,
		authorized:     true,
		shouldError:    false,
	}
}

func HandlerGithubRepoByOwnerAndNameTestCases() []HandlerGithubRepoByOwnerAndNameTest {
	return []HandlerGithubRepoByOwnerAndNameTest{
		shouldReturn401(),
		shouldReturnRepoForOwnerAndName(),
	}
}
