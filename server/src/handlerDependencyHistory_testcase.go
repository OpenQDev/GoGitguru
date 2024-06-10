package server

import (
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
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
			mock.ExpectQuery("-- name: GetRepoDependencies :many").WithArgs("foo", nonExistentRepoUrl, []string{}).WillReturnRows(sqlmock.NewRows([]string{
				"first_use_date",
				"last_use_date",
				"dependency_name",
			}).AddRow(nil, nil, "foo"),
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
			DatesAdded:   []string{"2021-08-25T13:39:56-05:00"},
			DatesRemoved: []string{},
		},

		setupMock: func(mock sqlmock.Sqlmock) {
			mockRows := sqlmock.NewRows([]string{
				"first_use_date",
				"last_use_date",
				"dependency_name",
			})
			row1 := mockRows.AddRow("2021-08-25T13:39:56-05:00", nil, "ethers")

			mock.ExpectQuery("-- name: GetRepoDependencies :many").
				WithArgs("ethers", openqFrontend, []string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"}).
				WillReturnRows(row1)

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
			DatesAdded:   []string{"2024-04-10T18:40:37-05:00"},
			DatesRemoved: []string{},
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("^-- name: GetRepoDependencies :many*").WithArgs("linea", linea, []string{"hardhat.config"}).WillReturnRows(sqlmock.NewRows([]string{
				"first_use_date",
				"last_use_date",
				"dependency_name",
			}).AddRow("2024-04-10T18:40:37-05:00", nil, "linea"),
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
