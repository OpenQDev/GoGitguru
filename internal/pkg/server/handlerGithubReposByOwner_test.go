package server

import (
	"main/internal/database"
	"main/internal/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubReposByOwner(t *testing.T) {
	// Initialize a new instance of ApiConfig with mocked DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	// Initialize queries with the mocked DB collection.
	queries := database.New(db)

	apiCfg := ApiConfig{
		DB: queries,
	}

	// Define test cases
	tests := []struct {
		name           string
		owner          string
		expectedStatus int
		shouldError    bool
		// Add more fields as needed
	}{
		{
			name:           "should return 401 if no access token",
			owner:          "VipassanaInsight",
			expectedStatus: 401,
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			req, _ := http.NewRequest("GET", "/repos/github/"+tt.owner, nil)
			rr := httptest.NewRecorder()

			// Define your expectations here. For example, if you expect the InsertGithubRepo function to be called,
			// you can set it up like this:
			// mock.ExpectExec("^-- name: InsertGithubRepo :one.*").WithArgs(...).WillReturnResult(sqlmock.NewResult(1, 1))

			// Call the handler function
			apiCfg.HandlerGithubReposByOwner(rr, req)

			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			expectedResponse := struct{}{}
			assert.Equal(t, expectedResponse, rr.Body.String())

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
