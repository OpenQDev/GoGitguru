package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestHandlerDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug

	logger.SetDebugMode(debugMode)

	tests := HandlerDependencyHistoryTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"LINEA",
			), tt.name)

			// BEFORE EACH
			_, queries := setup.GetMockDatabase()
			apiCfg := ApiConfig{
				DB: queries,
			}

			// ARRANGE - LOCAL
			req, _ := http.NewRequest("POST", "", nil)

			rr := httptest.NewRecorder()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// ACT
			apiCfg.HandlerDependencyHistory(rr, req)

			// ARRANGE - EXPECT
			var actualDependencyHistroyResponse DependencyHistoryResponse
			err := json.NewDecoder(rr.Body).Decode(&actualDependencyHistroyResponse)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into DependencyHistoryResponse: %s", err)
				return
			}

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			assert.Equal(t, tt.expectedDependencyHistroyResponse, actualDependencyHistroyResponse)
		})
	}
}
