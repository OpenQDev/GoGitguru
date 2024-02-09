package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OpenQDev/GoGitguru/util/githubRest"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubRepoByOwnerAndNameTest struct {
	title          string
	owner          string
	name           string
	expectedStatus int
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo)
}

const drmTestOrg = "drm-test-organization"
const drmTestRepo = "drm-test-repo"

func shouldReturn401() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN = "SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN"
	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_RETURN_401_IF_NO_ACCESS_TOKEN,
		owner:          drmTestOrg,
		name:           drmTestRepo,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		shouldError:    true,
		setupMock: func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo) {
			fullName := fmt.Sprintf("%s/%s", drmTestOrg, drmTestRepo)
			mock.ExpectQuery("-- name: CheckGithubRepoExists :one").WithArgs(fullName).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			createdAtUnix := createdAt.Unix()
			updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
			updatedAtUnix := updatedAt.Unix()
			pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)
			pushedAtUnix := pushedAt.Unix()

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repo.GithubRestID, repo.GithubGraphqlID, repo.URL, repo.Name, fullName, repo.Private, repo.Owner.Login, repo.Owner.AvatarURL, repo.Description, repo.Homepage, repo.Fork, repo.ForksCount, repo.Archived, repo.Disabled, repo.License.Name, repo.Language, repo.StargazersCount, repo.WatchersCount, repo.OpenIssuesCount, repo.HasIssues, repo.HasDiscussions, repo.HasProjects, createdAtUnix, updatedAtUnix, pushedAtUnix, repo.Visibility, repo.Size, repo.DefaultBranch)

			mock.ExpectQuery("-- name: GetGithubRepo :one").WithArgs(fullName).WillReturnRows(rows)
		},
	}
}

func shouldStoreRepoIfNotInDb() HandlerGithubRepoByOwnerAndNameTest {
	const SHOULD_STORE_REPO_IF_NOT_IN_DB = "SHOULD_STORE_REPO_IF_NOT_IN_DB"

	return HandlerGithubRepoByOwnerAndNameTest{
		title:          SHOULD_STORE_REPO_IF_NOT_IN_DB,
		owner:          drmTestOrg,
		name:           drmTestRepo,
		expectedStatus: http.StatusOK,
		authorized:     true,
		shouldError:    false,
		setupMock: func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo) {
			fullName := fmt.Sprintf("%s/%s", drmTestOrg, drmTestRepo)

			createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			createdAtUnix := createdAt.Unix()
			updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
			updatedAtUnix := updatedAt.Unix()
			pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)
			pushedAtUnix := pushedAt.Unix()

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repo.GithubRestID, repo.GithubGraphqlID, repo.URL, repo.Name, fullName, repo.Private, repo.Owner.Login, repo.Owner.AvatarURL, repo.Description, repo.Homepage, repo.Fork, repo.ForksCount, repo.Archived, repo.Disabled, repo.License.Name, repo.Language, repo.StargazersCount, repo.WatchersCount, repo.OpenIssuesCount, repo.HasIssues, repo.HasDiscussions, repo.HasProjects, createdAtUnix, updatedAtUnix, pushedAtUnix, repo.Visibility, repo.Size, repo.DefaultBranch)

			mock.ExpectQuery("-- name: CheckGithubRepoExists :one").WithArgs(fullName).WillReturnRows(sqlmock.NewRows([]string{"not_exists"}).AddRow(false))

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
				repo.GithubRestID,    // 0 - GithubRestID
				repo.GithubGraphqlID, // 1 - GithubGraphqlID
				repo.URL,             // 2 - Url
				repo.Name,            // 3 - Name
				fullName,             // 4 - FullName
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
				createdAtUnix,        // 22 - CreatedAt
				updatedAtUnix,        // 23 - UpdatedAt
				pushedAtUnix,         // 24 - PushedAt
				repo.Visibility,      // 25 - Visibility
				repo.Size,            // 26 - Size
				repo.DefaultBranch,   // 27 - DefaultBranch
			).WillReturnRows(rows)
		},
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
		setupMock: func(mock sqlmock.Sqlmock, repo githubRest.GithubRestRepo) {
			fullName := fmt.Sprintf("%s/%s", drmTestOrg, drmTestRepo)

			createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			createdAtUnix := createdAt.Unix()
			updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
			updatedAtUnix := updatedAt.Unix()
			pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)
			pushedAtUnix := pushedAt.Unix()

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repo.GithubRestID, repo.GithubGraphqlID, repo.URL, repo.Name, fullName, repo.Private, repo.Owner.Login, repo.Owner.AvatarURL, repo.Description, repo.Homepage, repo.Fork, repo.ForksCount, repo.Archived, repo.Disabled, repo.License.Name, repo.Language, repo.StargazersCount, repo.WatchersCount, repo.OpenIssuesCount, repo.HasIssues, repo.HasDiscussions, repo.HasProjects, createdAtUnix, updatedAtUnix, pushedAtUnix, repo.Visibility, repo.Size, repo.DefaultBranch)

			mock.ExpectQuery("-- name: CheckGithubRepoExists :one").WithArgs(fullName).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			mock.ExpectQuery("-- name: GetGithubRepo :one").WithArgs(fullName).WillReturnRows(rows)
		},
	}
}

func HandlerGithubRepoByOwnerAndNameTestCases() []HandlerGithubRepoByOwnerAndNameTest {
	return []HandlerGithubRepoByOwnerAndNameTest{
		shouldReturn401(),
		shouldStoreRepoIfNotInDb(),
		shouldReturnRepoIfExistsInDb(),
	}
}
