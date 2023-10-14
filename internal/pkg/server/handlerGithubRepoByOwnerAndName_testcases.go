package server

import "net/http"

type HandlerGithubRepoByOwnerAndNameTest struct {
	title          string
	owner          string
	name           string
	expectedStatus int
	authorized     bool
	shouldError    bool
}

const drmTestOrg = "DRM-Test-Organization"
const drmTestRepo = "DRM-Test-Repo"

func shouldReturn401() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN = "SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN,
		owner:          drmTestOrg,
		name:           drmTestRepo,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
	}
}

func shouldReturnRepoIfExistsInDb() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_RETURN_REPO_IF_EXISTS_IN_DB = "SHOULD_RETURN_REPO_IF_EXISTS_IN_DB"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_RETURN_REPO_IF_EXISTS_IN_DB,
		owner:          drmTestOrg,
		name:           drmTestRepo,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
	}
}

func shouldReturnRepoForOwnerAndName() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_GET_REPO_FOR_ORG_AND_NAME = "SHOULD_GET_REPO_FOR_ORG_AND_NAME"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_GET_REPO_FOR_ORG_AND_NAME,
		owner:          drmTestOrg,
		name:           drmTestRepo,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
	}
}

func HandlerGithubRepoByOwnerAndNameTestCases() []HandlerGithubRepoByOwnerAndNameTest {
	return []HandlerGithubRepoByOwnerAndNameTest{
		shouldReturn401(),
		shouldReturnRepoForOwnerAndName(),
		shouldReturnRepoIfExistsInDb(),
	}
}
