package server

import (
	"errors"
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddHandler(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, _, _ := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	mock, queries := mocks.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	/* ARRANGE - TEST DATA **/

	// "Valid repo URLs"
	targetRepos := []string{"https://github.com/org/repo1", "https://github.com/org/repo2"}

	successReturnBody := HandlerAddResponse{
		Accepted:       targetRepos,
		AlreadyInQueue: []string{},
	}

	twoReposRequest := HandlerAddRequest{
		RepoUrls: targetRepos,
	}

	tests := []struct {
		name             string
		expectedStatus   int
		requestBody      HandlerAddRequest
		expectedResponse HandlerAddResponse
	}{
		{
			name:             "Valid repo URLs",
			expectedStatus:   http.StatusAccepted,
			requestBody:      twoReposRequest,
			expectedResponse: successReturnBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE - LOCAL
			requestBody, err := util.TypeToReader(tt.requestBody)
			if err != nil {
				logger.LogFatalRedAndExit("failed to marshal response to %T: %s", tt.requestBody, err)
			}

			req, _ := http.NewRequest("POST", "", requestBody)
			rr := httptest.NewRecorder()

			// ARRANGE - EXPECT
			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo1").WillReturnError(errors.New("sql: no rows in result set"))
			mock.ExpectExec("^-- name: InsertRepoURL :exec.*").WithArgs("https://github.com/org/repo1").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo2").WillReturnError(errors.New("sql: no rows in result set"))
			mock.ExpectExec("^-- name: InsertRepoURL :exec.*").WithArgs("https://github.com/org/repo2").WillReturnResult(sqlmock.NewResult(1, 1))

			// ACT
			apiCfg.HandlerAdd(rr, req)

			// ARRANGE - EXPECT
			var actualResponse HandlerAddResponse
			err = util.ReaderToType(rr.Result().Body, &actualResponse)
			if err != nil {
				logger.LogFatalRedAndExit("failed to marshal response to %T: %s", actualResponse, err)
			}
			defer rr.Result().Body.Close()

			// ASSERT
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedResponse, actualResponse)

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
