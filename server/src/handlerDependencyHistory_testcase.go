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
	expectedDependencyHistroyResponse DependencyHistoryResponseMember
	setupMock                         func(mock sqlmock.Sqlmock)
}

func isNotAGitRepository() HandlerDependencyHistoryTestCase {
	nonExistentRepoUrl := []string{"https://github.com/IDont/Exist"}

	NOT_A_GIT_REPOSITORY := "NOT_A_GIT_REPOSITORY"

	return HandlerDependencyHistoryTestCase{

		name:           NOT_A_GIT_REPOSITORY,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrls:           nonExistentRepoUrl,
			FilePaths:          []string{},
			DependencySearched: "foo",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponseMember{
			DatesAdded:   0,
			DatesRemoved: 0,
			RepoUrl:      nonExistentRepoUrl[0],
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("-- name: GetRepoDependencies :many").WithArgs("foo", pq.Array(nonExistentRepoUrl), pq.Array([]string{})).WillReturnRows(sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
				"url",
			}).AddRow("foo", nil, nil, nonExistentRepoUrl[0]),
			)
		},
	}
}

func largeFrontend() HandlerDependencyHistoryTestCase {
	openqFrontend := []string{"https://github.com/OpenQDev/OpenQ-Frontend"}

	LARGE_FRONTEND := "LARGE_FRONTEND"

	return HandlerDependencyHistoryTestCase{
		name:           LARGE_FRONTEND,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrls:           openqFrontend,
			FilePaths:          []string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"},
			DependencySearched: "ethers",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponseMember{
			DatesAdded:   1629919196,
			DatesRemoved: 0,
			RepoUrl:      openqFrontend[0],
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			// TODO: Changed the date format to BIGINT since that's what it is in the schema first_use_date BIGINT DEFAULT NULL,
			mockRows := sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
				"url",
			}).AddRow("ethers", 1629919196, nil, openqFrontend[0])

			mock.ExpectQuery("^-- name: GetRepoDependencies :many*").
				WithArgs("ethers", pq.Array(openqFrontend), pq.Array([]string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"})).
				WillReturnRows(mockRows)
		},
	}
}

func linea() HandlerDependencyHistoryTestCase {
	linea := []string{"https://github.com/compound-finance/comet"}

	LINEA := "LINEA"

	return HandlerDependencyHistoryTestCase{
		name:           LINEA,
		shouldError:    false,
		expectedStatus: http.StatusOK,
		requestBody: DependencyHistoryRequest{
			RepoUrls:           linea,
			FilePaths:          []string{"hardhat.config"},
			DependencySearched: "linea",
		},
		expectedDependencyHistroyResponse: DependencyHistoryResponseMember{
			DatesAdded:   1715595637,
			DatesRemoved: 0,
			RepoUrl:      linea[0],
		},
		setupMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("^-- name: GetRepoDependencies :many*").WithArgs("linea", pq.Array(linea), pq.Array([]string{"hardhat.config"})).WillReturnRows(sqlmock.NewRows([]string{
				"dependency_name",
				"first_use_date",
				"last_use_date",
				"url",
			}).AddRow("linea", 1715595637, nil, linea[0]),
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
