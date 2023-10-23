package server

import (
	"net/http"
	"time"
	"util/githubRest"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubReposByOwnerTestCase struct {
	name           string
	owner          string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo)
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
		setupMock:      func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo) {},
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
		setupMock: func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo) {
			createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
			pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repo.GithubRestID, repo.GithubGraphqlID, repo.URL, repo.Name, repo.FullName, repo.Private, repo.Owner.Login, repo.Owner.AvatarURL, repo.Description, repo.Homepage, repo.Fork, repo.ForksCount, repo.Archived, repo.Disabled, repo.License.Name, repo.Language, repo.StargazersCount, repo.WatchersCount, repo.OpenIssuesCount, repo.HasIssues, repo.HasDiscussions, repo.HasProjects, createdAt, updatedAt, pushedAt, repo.Visibility, repo.Size, repo.DefaultBranch)

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
				repo.GithubRestID,    // 0 - GithubRestID
				repo.GithubGraphqlID, // 1 - GithubGraphqlID
				repo.URL,             // 2 - Url
				repo.Name,            // 3 - Name
				repo.FullName,        // 4 - FullName
				repo.Private,         // 5 - Private
				repo.Owner.Login,     // 6 - OwnerLogin
				repo.Owner.AvatarURL, // 7 - OwnerAvatarUrl
				repo.Description,     // 8 - Description
				repo.Homepage,        // 9 - Homepage
				repo.Fork,            // 10 - Fork
				repo.ForksCount,      // 11 - ForksCount
				repo.Archived,        // 12 - Archived
				repo.Disabled,        // 13 - Disabled
				repo.License.Name,    // 14 - License
				repo.Language,        // 15 - Language
				repo.StargazersCount, // 16 - StargazersCount
				repo.WatchersCount,   // 17 - WatchersCount
				repo.OpenIssuesCount, // 18 - OpenIssuesCount
				repo.HasIssues,       // 19 - HasIssues
				repo.HasDiscussions,  // 20 - HasDiscussions
				repo.HasProjects,     // 21 - HasProjects
				createdAt,            // 22 - CreatedAt
				updatedAt,            // 23 - UpdatedAt
				pushedAt,             // 24 - PushedAt
				repo.Visibility,      // 25 - Visibility
				repo.Size,            // 26 - Size
				repo.DefaultBranch,   // 27 - DefaultBranch
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
