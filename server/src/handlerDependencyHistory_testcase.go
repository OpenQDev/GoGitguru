package server

import (
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

type HandlerDependencyHistoryTestCase struct {
	name                              string
	shouldError                       bool
	expectedStatus                    int
	requestBody                       DependencyHistoryRequest
	expectedDependencyHistroyResponse DependencyHistoryResponse
	setupMock                         func(mock sqlmock.Sqlmock)
}

func isNotAGitRepository() HandlerDependencyHistoryTestCase {
	const nonExistentRepoUrl = "https://github.com/IDont/Exist"

	NOT_A_GIT_REPOSITORY := "NOT_A_GIT_REPOSITORY"

	return HandlerDependencyHistoryTestCase{
		name:           NOT_A_GIT_REPOSITORY,
		shouldError:    true,
		expectedStatus: http.StatusNotFound,
		requestBody: DependencyHistoryRequest{
			RepoUrl:            nonExistentRepoUrl,
			FilePaths:          []string{},
			DependencySearched: "foo",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponse{},
		setupMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("-- name: GetRepoDependencies :many").WithArgs("foo", nonExistentRepoUrl, pq.Array([]string{})).WillReturnRows(sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
			}).AddRow("foo", nil, nil),
			)
		},
	}
}

func largeFrontend() HandlerDependencyHistoryTestCase {
	const openqFrontend = "https://github.com/OpenQDev/OpenQ-Frontend"

	LARGE_FRONTEND := "LARGE_FRONTEND"

	return HandlerDependencyHistoryTestCase{
		name:           LARGE_FRONTEND,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrl:            openqFrontend,
			FilePaths:          []string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"},
			DependencySearched: "ethers",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponse{
			DatesAdded:   []string{"2021-08-25T19:19:56Z"},
			DatesRemoved: []string{},
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			// TODO: Changed the date format to BIGINT since that's what it is in the schema first_use_date BIGINT DEFAULT NULL,
			mockRows := sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
			}).AddRow("ethers", 1629919196, nil)

			mock.ExpectQuery("^-- name: GetRepoDependencies :many*").
				WithArgs("ethers", openqFrontend, pq.Array([]string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"})).
				WillReturnRows(mockRows)
		},
	}
}

func linea() HandlerDependencyHistoryTestCase {
	const linea = "https://github.com/compound-finance/comet"

	LINEA := "LINEA"

	return HandlerDependencyHistoryTestCase{
		name:           LINEA,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrl:            linea,
			FilePaths:          []string{"hardhat.config"},
			DependencySearched: "linea",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponse{
			DatesAdded:   []string{"2024-05-13T10:20:37Z"},
			DatesRemoved: []string{},
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("^-- name: GetRepoDependencies :many*").WithArgs("linea", linea, pq.Array([]string{"hardhat.config"})).WillReturnRows(sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
			}).AddRow("linea", 1715595637, nil),
			)
		},
	}
}

func HandlerDependencyHistoryTestCases() []HandlerDependencyHistoryTestCase {
	return []HandlerDependencyHistoryTestCase{
		isNotAGitRepository(),
		largeFrontend(),
		linea(),
	}
}
