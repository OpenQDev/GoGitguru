package server

import (
	"log"
	"net/http"
	"os"

	"github.com/OpenQDev/GoGitguru/util/marshaller"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerFirstCommitTestCase struct {
	name               string
	login              string
	expectedStatus     int
	authorized         bool
	requestBody        HandlerFirstCommitRequest
	expectedReturnBody CommitWithAuthorInfo
	shouldError        bool
	setupMock          func(mock sqlmock.Sqlmock)
}

const firstCommitLogin = "FlacoJones"
const firstCommitRepoUrl = "https://github.com/OpenQDev/OpenQ-Workflows"

func getFirstCommit() HandlerFirstCommitTestCase {
	const GET_FIRST_USER_COMMIT = "GET_FIRST_USER_COMMIT"

	requestBody := HandlerFirstCommitRequest{
		RepoUrl: firstCommitRepoUrl,
		Login:   firstCommitLogin,
	}

	var commitResponse CommitWithAuthorInfo
	jsonFile, err := os.Open("./mocks/mockFirstCommit.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	marshaller.JsonFileToType(jsonFile, &commitResponse)

	expectedReturnBody := commitResponse

	return HandlerFirstCommitTestCase{
		name:               GET_FIRST_USER_COMMIT,
		login:              login,
		expectedStatus:     http.StatusOK,
		authorized:         true,
		requestBody:        requestBody,
		expectedReturnBody: expectedReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock) {
			// Define the mock rows
			mockRows := sqlmock.NewRows([]string{
				"commit_hash", "author", "author_email", "author_date", "committer_date", "message", "insertions", "deletions", "lines_changed", "files_changed", "repo_url",
				"rest_id", "email", "internal_id", "github_rest_id", "github_graphql_id", "login", "name", "email_2", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type", "created_at", "updated_at",
			})

			// Add rows to the mock rows
			row1 := mockRows.AddRow(
				commitResponse.CommitHash, commitResponse.Author, commitResponse.AuthorEmail, commitResponse.AuthorDate, commitResponse.CommitterDate, commitResponse.Message, commitResponse.Insertions, commitResponse.Deletions, commitResponse.LinesChanged, commitResponse.FilesChanged, commitResponse.RepoUrl,
				commitResponse.RestID, commitResponse.Email, commitResponse.InternalID, commitResponse.GithubRestID, commitResponse.GithubGraphqlID, commitResponse.Login, commitResponse.Name, commitResponse.Email_2, commitResponse.AvatarUrl, commitResponse.Company, commitResponse.Location, commitResponse.Bio, commitResponse.Blog, commitResponse.Hireable, commitResponse.TwitterUsername, commitResponse.Followers, commitResponse.Following, commitResponse.Type, commitResponse.CreatedAt, commitResponse.UpdatedAt,
			)

			// Expect the query with the mock rows
			mock.ExpectQuery("^-- name: GetFirstCommit :onoe.*").
				WithArgs(firstCommitRepoUrl, login).
				WillReturnRows(row1)
		},
	}
}

func HandlerFirstCommitTestCases() []HandlerFirstCommitTestCase {
	return []HandlerFirstCommitTestCase{
		getFirstCommit(),
	}
}
