package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddHandler(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug
	logger.SetDebugMode(debugMode)

	// ARRANGE - TESTS
	tests := HandlerAddTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()

			apiCfg := ApiConfig{
				DB: queries,
			}

			// ARRANGE - LOCAL
			requestBody, err := marshaller.TypeToReader(tt.requestBody)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", tt.requestBody, err)
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

			// EXPECT - ERRORS
			if tt.shouldError {
				var actualErrorResponse ErrorResponse
				err = marshaller.ReaderToType(rr.Result().Body, &actualErrorResponse)
				if err != nil {
					t.Errorf("failed to marshal response to %T: %s", actualErrorResponse, err)
				}

				assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
				assert.Equal(t, tt.expectedErrorResponse, actualErrorResponse)
				return
			}

			// EXPECT - SUCCESS
			var actualSuccessResponse HandlerAddResponse
			err = marshaller.ReaderToType(rr.Result().Body, &actualSuccessResponse)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", actualSuccessResponse, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
			assert.Equal(t, tt.expectedSuccessResponse, actualSuccessResponse)

			// --- SECOND CALL --- //

			// ARRANGE - LOCAL
			requestBody, err = marshaller.TypeToReader(tt.requestBody)
			if err != nil {
				logger.LogFatalRedAndExit("failed to marshal response to %T: %s", tt.requestBody, err)
			}

			req, _ = http.NewRequest("POST", "", requestBody)
			rr = httptest.NewRecorder()

			// ARRANGE - EXPECT
			currentTime := time.Now()
			repoURLMockRow1 := sqlmock.NewRows([]string{"url", "status", "created_at", "updated_at"}).AddRow("https://github.com/org/repo1", database.RepoStatusPending, currentTime, currentTime)
			repoURLMockRow2 := sqlmock.NewRows([]string{"url", "status", "created_at", "updated_at"}).AddRow("https://github.com/org/repo2", database.RepoStatusPending, currentTime, currentTime)

			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo1").WillReturnRows(repoURLMockRow1)
			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo2").WillReturnRows(repoURLMockRow2)

			// ACT
			apiCfg.HandlerAdd(rr, req)

			// EXPECT - SUCCESS
			err = marshaller.ReaderToType(rr.Result().Body, &actualSuccessResponse)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", actualSuccessResponse, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
			assert.Equal(t, tt.secondExpectedSuccessResponse, actualSuccessResponse)
		})
	}
}
