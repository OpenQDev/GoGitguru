package server

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, _, _ := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	successResponse := HandlerVersionResponse{Version: "1.0.0"}

	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody HandlerVersionResponse
	}{
		{
			name:                 "should return 200 and version 1.0.0",
			expectedStatusCode:   200,
			expectedResponseBody: successResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
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
