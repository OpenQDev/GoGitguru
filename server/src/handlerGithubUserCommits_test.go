package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerGithubUserCommits(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug
	ghAccessToken := env.GhAccessToken

	logger.SetDebugMode(debugMode)

	mock, queries := setup.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	tests := HandlerGithubUserCommitsTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)

			// Add {login} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "login", tt.login)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			tt.setupMock(mock)

			// ACT
			apiCfg.HandlerGithubUserCommits(rr, req)

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// ARRANGE - EXPECT
			var actualRepoCommitsReturn []CommitWithAuthorInfo
			err := marshaller.ReaderToType(rr.Body, &actualRepoCommitsReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into []RestRepo: %s", err)
				return
			}

			require.Equal(t, tt.expectedStatus, rr.Code, rr.Body)
			require.Equal(t, tt.expectedReturnBody, actualRepoCommitsReturn)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
