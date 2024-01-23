package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OpenQDev/GoGitguru/util/marshaller"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerGithubUserCommitsTestCase struct {
	name               string
	login              string
	expectedStatus     int
	authorized         bool
	requestBody        HandlerGithubUserCommitsRequest
	expectedReturnBody []CommitWithAuthorInfo
	shouldError        bool
	setupMock          func(mock sqlmock.Sqlmock)
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
		name:               UNAUTHORIZED,
		login:              login,
		expectedStatus:     http.StatusUnauthorized,
		authorized:         false,
		requestBody:        requestBody,
		expectedReturnBody: []CommitWithAuthorInfo{},
		shouldError:        true,
		setupMock:          func(mock sqlmock.Sqlmock) {},
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

	var twoCommitsResponse []CommitWithAuthorInfo
	jsonFile, err := os.Open("./mocks/mockRepoCommitsResponse.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	marshaller.JsonFileToArrayOfType(jsonFile, &twoCommitsResponse)

	expectedReturnBody := twoCommitsResponse

	return HandlerGithubUserCommitsTestCase{
		name:               GET_ALL_USER_COMMITS,
		login:              login,
		expectedStatus:     http.StatusOK,
		authorized:         true,
		requestBody:        requestBody,
		expectedReturnBody: expectedReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock) {
			sinceTime, _ := time.Parse(time.RFC3339, since)
			untilTime, _ := time.Parse(time.RFC3339, until)
			sinceUnix := sinceTime.Unix()
			untilUnix := untilTime.Unix()

			// Define the mock rows
			mockRows := sqlmock.NewRows([]string{
				"commit_hash", "author", "author_email", "author_date", "committer_date", "message", "insertions", "deletions", "lines_changed", "files_changed", "repo_url",
				"rest_id", "email", "internal_id", "github_rest_id", "github_graphql_id", "login", "name", "email_2", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type", "created_at", "updated_at",
			})

			// Add rows to the mock rows
			c1 := twoCommitsResponse[0]
			row1 := mockRows.AddRow(
				c1.CommitHash, c1.Author, c1.AuthorEmail, c1.AuthorDate, c1.CommitterDate, c1.Message, c1.Insertions, c1.Deletions, c1.LinesChanged, c1.FilesChanged, c1.RepoUrl,
				c1.RestID, c1.Email, c1.InternalID, c1.GithubRestID, c1.GithubGraphqlID, c1.Login, c1.Name, c1.Email_2, c1.AvatarUrl, c1.Company, c1.Location, c1.Bio, c1.Blog, c1.Hireable, c1.TwitterUsername, c1.Followers, c1.Following, c1.Type, c1.CreatedAt, c1.UpdatedAt,
			)

			c2 := twoCommitsResponse[1]
			row2 := mockRows.AddRow(
				c2.CommitHash, c2.Author, c2.AuthorEmail, c2.AuthorDate, c2.CommitterDate, c2.Message, c2.Insertions, c2.Deletions, c2.LinesChanged, c2.FilesChanged, c2.RepoUrl,
				c2.RestID, c2.Email, c2.InternalID, c2.GithubRestID, c2.GithubGraphqlID, c2.Login, c2.Name, c2.Email_2, c2.AvatarUrl, c2.Company, c2.Location, c2.Bio, c2.Blog, c2.Hireable, c2.TwitterUsername, c2.Followers, c2.Following, c2.Type, c2.CreatedAt, c2.UpdatedAt,
			)

			// Expect the query with the mock rows
			mock.ExpectQuery("^-- name: GetUserCommitsForRepos :many.*").
				WithArgs(sinceUnix, untilUnix, login).
				WillReturnRows(row1, row2)
		},
	}
}

func HandlerGithubUserCommitsTestCases() []HandlerGithubUserCommitsTestCase {
	return []HandlerGithubUserCommitsTestCase{
		notAuthorized(),
		getAllUserCommits(),
	}
}
