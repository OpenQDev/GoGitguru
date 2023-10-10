package server

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerReadiness(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, _, _ := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	apiCfg := ApiConfig{
		DB: queries,
	}

	successReturnBody := HandlerReadinessResponse{}

	tests := []struct {
		name               string
		expectedStatus     int
		expectedReturnBody HandlerReadinessResponse
	}{
		{
			name:               "should return 200 and empty struct",
			expectedStatus:     200,
			expectedReturnBody: successReturnBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)
			rr := httptest.NewRecorder()

			// ACT
			apiCfg.HandlerReadiness(rr, req)

			// ARRANGE - EXPECT
			var actualReturnBody HandlerReadinessResponse
			err := UnmarshalReader(rr.Result().Body, &actualReturnBody)
			if err != nil {
				t.Errorf("error unmarshalling: %s", err)
			}
			defer rr.Result().Body.Close()

			// ASSERT
			assert.Equal(t, tt.expectedStatus, rr.Result().StatusCode)
			assert.Equal(t, tt.expectedReturnBody, actualReturnBody)
		})
	}
}
