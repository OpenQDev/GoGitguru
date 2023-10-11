package server

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"main/internal/pkg/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerVersion(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, _, _ := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	// ARRANGE - TESTS
	type HandlerVersionTest struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody HandlerVersionResponse
	}

	// TEST 1
	const SHOULD_RETURN_200_AND_CORRECT_VERSION = "should return 200 and version 1.0.0"
	successResponse := HandlerVersionResponse{Version: "1.0.0"}

	shouldReturn200 := HandlerVersionTest{
		name:                 SHOULD_RETURN_200_AND_CORRECT_VERSION,
		expectedStatusCode:   200,
		expectedResponseBody: successResponse,
	}

	tests := []HandlerVersionTest{
		shouldReturn200,
	}

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
			err := util.ReaderToType(rr.Result().Body, &actualResponse)
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
