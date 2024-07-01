package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/util/lib"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddDependencyPatternHandler(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug
	logger.SetDebugMode(debugMode)

	// ARRANGE - TESTS
	tests := HandlerAddDependencyPatternTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"VALID_REPO_URLS",
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()
			apiCfg := ApiConfig{
				DB: queries,
			}

			// ARRANGE - LOCAL
			requestBody, err := marshaller.TypeToReader(tt.firstRequestBody)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", tt.firstRequestBody, err)
			}

			req, _ := http.NewRequest("POST", "", requestBody)
			rr := httptest.NewRecorder()

			mockTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			lib.Now = func() time.Time { return mockTime }

			// ARRANGE - EXPECT
			mock.ExpectExec("^-- name: BulkUpsertFilePatterns :exec.*").WithArgs(pq.Array(tt.firstRequestBody.DependencyPatterns), mockTime.Unix(), tt.firstRequestBody.Creator).WillReturnResult(sqlmock.NewResult(1, 1))

			// ACT
			apiCfg.HandlerAddDependencyPattern(rr, req)

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
			var actualSuccessResponse HandlerAddDependencyPatternResponse
			err = marshaller.ReaderToType(rr.Result().Body, &actualSuccessResponse)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", actualSuccessResponse, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
			assert.Equal(t, tt.firstExpectedSuccessResponse, actualSuccessResponse)

			// --- SECOND CALL --- //

			// ARRANGE - LOCAL
			requestBody, err = marshaller.TypeToReader(tt.secondRequestBody)
			if err != nil {
				logger.LogFatalRedAndExit("failed to marshal response to %T: %s", tt.secondRequestBody, err)
			}

			req, _ = http.NewRequest("POST", "", requestBody)
			rr = httptest.NewRecorder()

			// ARRANGE - EXPECT
			mock.ExpectExec("^-- name: BulkUpsertFilePatterns :exec.*").WithArgs(pq.Array(tt.secondRequestBody.DependencyPatterns), mockTime.Unix(), tt.secondRequestBody.Creator).WillReturnResult(sqlmock.NewResult(1, 1))

			// ACT
			apiCfg.HandlerAddDependencyPattern(rr, req)

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
