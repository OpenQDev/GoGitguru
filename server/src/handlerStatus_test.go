package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug
	logger.SetDebugMode(debugMode)

	// ARRANGE - TESTS
	tests := HandlerStatusTestCases()

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

			tt.setupMock(mock, tt.requestBody.RepoUrls)

			// ACT
			apiCfg.HandlerStatus(rr, req)

			// ASSERT
			if tt.shouldError {
				require.Equal(t, tt.expectedStatus, rr.Code)
				expectedError := `{"error":"repo_urls cannot be empty"}`
				require.JSONEq(t, expectedError, rr.Body.String())
				return
			}

			// ARRANGE - EXPECT
			var actualHandlerStatusResponse []HandlerStatusResponse
			err = marshaller.ReaderToType(rr.Body, &actualHandlerStatusResponse)
			if err != nil {
				fmt.Println(rr.Body.String(), "that other thing", err)
				t.Errorf("Failed to decode rr.Body into HandlerStatusResponse: %s", err)
				return
			}

			require.Equal(t, tt.expectedStatus, rr.Code, rr.Body)
			require.Equal(t, tt.expectedReturnBody, actualHandlerStatusResponse)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
