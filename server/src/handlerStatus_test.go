package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment(".env")
	debugMode := env.Debug
	logger.SetDebugMode(debugMode)

	_, queries := setup.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	// ARRANGE - TESTS
	tests := HandlerStatusTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ARRANGE - LOCAL
			requestBody, err := marshaller.TypeToReader(tt.requestBody)
			if err != nil {
				t.Errorf("failed to marshal response to %T: %s", tt.requestBody, err)
			}

			req, _ := http.NewRequest("POST", "", requestBody)
			rr := httptest.NewRecorder()

			// ACT
			apiCfg.HandlerStatus(rr, req)

			// EXPECT
			assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
		})
	}
}
