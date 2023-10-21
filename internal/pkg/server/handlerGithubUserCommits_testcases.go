package server

import (
	"net/http"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserCommitsTestCase struct {
	name           string
	login          string
	expectedStatus int
	authorized     bool
	requestBody    HandlerGithubUserCommitsRequest
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock)
}

const login = "DRM-Test-Organization"

func notAuthorized() HandlerGithubUserCommitsTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"

	requestBody := HandlerGithubUserCommitsRequest{
		RepoUrls: []string{"https://github.com/openqdev/openq-workflows"},
		Since:    time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
		Until:    time.Now().AddDate(0, 0, 0).Format(time.RFC3339),
	}

	return HandlerGithubUserCommitsTestCase{
		name:           UNAUTHORIZED,
		login:          login,
		expectedStatus: http.StatusUnauthorized,
		authorized:     false,
		requestBody:    requestBody,
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock) {},
	}
}

func getAllUserCommits() HandlerGithubUserCommitsTestCase {
	const GET_ALL_USER_COMMITS = "GET_ALL_USER_COMMITS"
	since := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	until := time.Now().AddDate(0, 0, 0).Format(time.RFC3339)

	requestBody := HandlerGithubUserCommitsRequest{
		RepoUrls: []string{},
		Since:    since,
		Until:    until,
	}

	return HandlerGithubUserCommitsTestCase{
		name:           GET_ALL_USER_COMMITS,
		login:          login,
		expectedStatus: http.StatusOK,
		authorized:     true,
		requestBody:    requestBody,
		shouldError:    false,
		setupMock: func(mock sqlmock.Sqlmock) {

			sinceTime, _ := time.Parse(time.RFC3339, since)
			untilTime, _ := time.Parse(time.RFC3339, until)
			sinceUnix := sinceTime.Unix()
			untilUnix := untilTime.Unix()

			mock.ExpectQuery("^-- name: GetAllUserCommits :many.*").
				WithArgs(sinceUnix, untilUnix, login).
				WillReturnRows(sqlmock.NewRows([]string{
					"commit_hash", "author", "author_email", "author_date", "committer_date", "message", "insertions", "deletions", "lines_changed", "files_changed", "repo_url",
					"internal_id", "github_rest_id", "github_graphql_id", "login", "name", "email", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type", "created_at", "updated_at",
				}))

		},
	}
}

func HandlerGithubUserCommitsTestCases() []HandlerGithubUserCommitsTestCase {
	return []HandlerGithubUserCommitsTestCase{
		notAuthorized(),
		getAllUserCommits(),
	}
}
