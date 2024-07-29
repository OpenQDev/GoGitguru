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
	expectedReturnBody HandlerFirstCommitResponse
	shouldError        bool
	setupMock          func(mock sqlmock.Sqlmock)
}

func getFirstCommit() HandlerFirstCommitTestCase {
	const GET_FIRST_USER_COMMIT = "GET_FIRST_USER_COMMIT"

	const firstCommitRepoUrl = "https://github.com/OpenQDev/OpenQ-Workflows"
	const firstCommitLogin = "mktcode"

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
	err = marshaller.JsonFileToType(jsonFile, &commitResponse)
	if err != nil {
		log.Fatal(err)
	}

	expectedReturnBody := HandlerFirstCommitResponse{AuthorDate: 1697637308}

	return HandlerFirstCommitTestCase{
		name:               GET_FIRST_USER_COMMIT,
		login:              firstCommitLogin,
		expectedStatus:     http.StatusOK,
		authorized:         true,
		requestBody:        requestBody,
		expectedReturnBody: expectedReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock) {
			// Define the mock rows
			mockRows := sqlmock.NewRows([]string{
				"author_date",
			})

			// Add rows to the mock rows
			row1 := mockRows.AddRow(
				commitResponse.AuthorDate,
			)

			// Expect the query with the mock rows
			mock.ExpectQuery("^-- name: GetFirstCommit :one.*").
				WithArgs(firstCommitRepoUrl, firstCommitLogin).
				WillReturnRows(row1)
		},
	}
}

func HandlerFirstCommitTestCases() []HandlerFirstCommitTestCase {
	return []HandlerFirstCommitTestCase{
		getFirstCommit(),
	}
}
