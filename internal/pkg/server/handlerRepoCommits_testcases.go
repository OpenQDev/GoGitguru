package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerRepoCommitsTestCase struct {
	name           string
	login          string
	expectedStatus int
	requestBody    HandlerRepoCommitsRequest
	authorized     bool
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock)
}

func foo() HandlerRepoCommitsTestCase {
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerRepoCommitsTestCase{
		name:           UNAUTHORIZED,
		login:          login,
		expectedStatus: http.StatusBadRequest,
		authorized:     false,
		requestBody:    HandlerRepoCommitsRequest{},
		shouldError:    true,
		setupMock:      func(mock sqlmock.Sqlmock) {},
	}
}

func getAllRepoCommits() HandlerRepoCommitsTestCase {
	const GET_ALL_REPO_COMMITS = "GET_ALL_REPO_COMMITS"

	since := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	fmt.Println("since", since)
	until := time.Now().AddDate(0, 0, 0).Format(time.RFC3339)

	requestBody := HandlerRepoCommitsRequest{
		RepoURL: "https://github.com/openqdev/openq-workflows",
		Since:   since,
		Until:   until,
	}

	return HandlerRepoCommitsTestCase{
		name:           GET_ALL_REPO_COMMITS,
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

			mock.ExpectQuery("-- name: GetCommitsWithAuthorInfo :many").
				WithArgs(requestBody.RepoURL, sinceUnix, untilUnix).
				WillReturnRows(sqlmock.NewRows([]string{
					"commit_hash", "author", "author_email", "author_date", "committer_date", "message", "insertions", "deletions", "lines_changed", "files_changed", "repo_url",
					"rest_id", "gure.email", "internal_id", "github_rest_id", "github_graphql_id", "login", "name", "gu.email", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type", "created_at", "updated_at",
				}))
		},
	}
}

func HandlerRepoCommitsTestCases() []HandlerRepoCommitsTestCase {
	return []HandlerRepoCommitsTestCase{
		foo(),
		getAllRepoCommits(),
	}
}
