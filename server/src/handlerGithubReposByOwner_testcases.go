package server

import (
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubReposByOwnerTestCase struct {
	name           string
	owner          string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, repo database.InsertGithubRepoParams)
}

const owner = "DRM-Test-Organization"

func unauthorized() HandlerGithubReposByOwnerTestCase {
	const SHOULD_401 = "SHOULD_401"
	return HandlerGithubReposByOwnerTestCase{
		name:           SHOULD_401,
		owner:          owner,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock, repo database.InsertGithubRepoParams) {},
	}
}

func sucess() HandlerGithubReposByOwnerTestCase {
	const SHOULD_STORE_ALL_REPOS_FOR_ORG = "SHOULD_STORE_ALL_REPOS_FOR_ORG"
	return HandlerGithubReposByOwnerTestCase{
		name:           SHOULD_STORE_ALL_REPOS_FOR_ORG,
		owner:          owner,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
		setupMock: func(mock sqlmock.Sqlmock, repo database.InsertGithubRepoParams) {

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
				repo.GithubRestID,          // 0 - GithubRestID
				repo.GithubGraphqlID,       // 1 - GithubGraphqlID
				repo.Url,                   // 2 - Url
				repo.Name,                  // 3 - Name
				repo.FullName,              // 4 - FullName
				repo.Private.Bool,          // 5 - Private
				repo.OwnerLogin,            // 6 - OwnerLogin
				repo.OwnerAvatarUrl,        // 7 - OwnerAvatarUrl
				repo.Description.String,    // 8 - Description
				repo.Homepage.String,       // 9 - Homepage
				repo.Fork.Bool,             // 10 - Fork
				repo.ForksCount.Int32,      // 11 - ForksCount
				repo.Archived.Bool,         // 12 - Archived
				repo.Disabled.Bool,         // 13 - Disabled
				repo.License.String,        // 14 - License
				repo.Language.String,       // 15 - Language
				repo.StargazersCount.Int32, // 16 - StargazersCount
				repo.WatchersCount.Int32,   // 17 - WatchersCount
				repo.OpenIssuesCount.Int32, // 18 - OpenIssuesCount
				repo.HasIssues.Bool,        // 19 - HasIssues
				repo.HasDiscussions.Bool,   // 20 - HasDiscussions
				repo.HasProjects.Bool,      // 21 - HasProjects
				repo.CreatedAt.Int32,       // 22 - CreatedAt
				repo.UpdatedAt.Int32,       // 23 - UpdatedAt
				repo.PushedAt.Int32,        // 24 - PushedAt
				repo.Visibility.String,     // 25 - Visibility
				repo.Size.Int32,            // 26 - Size
				repo.DefaultBranch.String,  // 27 - DefaultBranch
			).WillReturnRows(rows)
		},
	}
}

func HandlerGithubReposByOwnerTestCases() []HandlerGithubReposByOwnerTestCase {
	return []HandlerGithubReposByOwnerTestCase{
		unauthorized(),
		sucess(),
	}
}
