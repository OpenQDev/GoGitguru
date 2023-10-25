package server

import (
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

type HandlerStatusTest struct {
	name               string
	expectedStatus     int
	requestBody        HandlerAddRequest
	expectedReturnBody []HandlerStatusResponse
	shouldError        bool
	setupMock          func(mock sqlmock.Sqlmock, repos []string)
}

func noRepoUrls() HandlerStatusTest {
	const NO_REPO_URLS = "NO_REPO_URLS"
	targetRepos := []string{}

	noRequestsBody := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	noRepoUrls := HandlerStatusTest{
		name:               NO_REPO_URLS,
		expectedStatus:     http.StatusBadRequest,
		requestBody:        noRequestsBody,
		expectedReturnBody: []HandlerStatusResponse{},
		shouldError:        true,
		setupMock:          func(mock sqlmock.Sqlmock, repos []string) {},
	}

	return noRepoUrls
}

func statusValidRepoUrls() HandlerStatusTest {
	const VALID_REPO_URLS = "VALID_REPO_URLS"
	repo1Url := "https://github.com/org/repo1"
	repo2Url := "https://github.com/org/repo2"
	targetRepos := []string{repo1Url, repo2Url}

	twoReposRequest := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	successReturnBody := []HandlerStatusResponse{
		{
			Url:            "https://github.com/org/repo1",
			Status:         "pending",
			PendingAuthors: 1,
		},
		{
			Url:            "https://github.com/org/repo2",
			Status:         "pending",
			PendingAuthors: 2,
		},
	}

	rows := sqlmock.NewRows([]string{"url", "status", "pending_authors"})
	rows.AddRow(repo1Url, "pending", 1)
	rows.AddRow(repo2Url, "pending", 2)

	validRepoUrls := HandlerStatusTest{
		name:               VALID_REPO_URLS,
		expectedStatus:     http.StatusAccepted,
		requestBody:        twoReposRequest,
		expectedReturnBody: successReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock, repos []string) {
			mock.ExpectQuery("^-- name: GetReposStatus :many.*").
				WithArgs(pq.Array(targetRepos)).
				WillReturnRows(rows)
		},
	}

	return validRepoUrls
}

func missingRepoUrl() HandlerStatusTest {
	const MISSING_REPO_URL = "MISSING_REPO_URL"
	repo1Url := "https://github.com/org/repo1"
	missingRepoUrl := "https://github.com/org/iDontExist"
	targetRepos := []string{repo1Url, missingRepoUrl}

	twoReposRequest := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	successReturnBody := []HandlerStatusResponse{
		{
			Url:            repo1Url,
			Status:         "pending",
			PendingAuthors: 1,
		},
		{
			Url:            missingRepoUrl,
			Status:         "not_listed",
			PendingAuthors: 0,
		},
	}

	rows := sqlmock.NewRows([]string{"url", "status", "pending_authors"})
	rows.AddRow(repo1Url, "pending", 1)

	oneMissingRepoUrl := HandlerStatusTest{
		name:               MISSING_REPO_URL,
		expectedStatus:     http.StatusAccepted,
		requestBody:        twoReposRequest,
		expectedReturnBody: successReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock, repos []string) {
			mock.ExpectQuery("^-- name: GetReposStatus :many.*").
				WithArgs(pq.Array(targetRepos)).
				WillReturnRows(rows)
		},
	}

	return oneMissingRepoUrl
}

func HandlerStatusTestCases() []HandlerStatusTest {
	return []HandlerStatusTest{
		statusValidRepoUrls(),
		missingRepoUrl(),
		noRepoUrls(),
	}
}
