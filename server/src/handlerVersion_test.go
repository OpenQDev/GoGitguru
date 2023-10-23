package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"util/logger"
	"util/marshaller"
	"util/setup"
	"util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestHandlerVersion(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment(".env")
	debugMode := env.Debug

	logger.SetDebugMode(debugMode)

	_, queries := setup.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	// ARRANGE - TESTS
	tests := HandlerVersionTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)
			rr := httptest.NewRecorder()

			// ACT
			apiCfg.HandlerVersion(rr, req)

			// ARRANGE - EXPECT
			var actualResponse HandlerVersionResponse
			err := marshaller.ReaderToType(rr.Result().Body, &actualResponse)
			if err != nil {
				logger.LogFatalRedAndExit("failed to marshal response to %T: %s", actualResponse, err)
			}
			defer rr.Result().Body.Close()

			// ASSERT
			assert.Equal(t, tt.expectedStatusCode, rr.Result().StatusCode)
			assert.Equal(t, tt.expectedResponseBody, actualResponse)
		})
	}
}
