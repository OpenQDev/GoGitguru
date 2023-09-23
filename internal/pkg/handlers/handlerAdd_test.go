package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"main/internal/database"
	"main/internal/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddHandler(t *testing.T) {
	// Initialize a new instance of ApiConfig with mocked DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB:", err)
	}

	// Expectations and actions for the mock DB can be defined here
	// For example:
	// mock.ExpectQuery("^SELECT (.+) FROM users$").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John"))

	// Initialize queries with the mocked DB collection. All other queries and models are identical
	// Only now we can use db as a spy to ensure methods were called
	queries := database.New(db)

	apiCfg := ApiConfig{
		DB: queries,
	}

	// Define test cases
	tests := []struct {
		name           string
		repoUrls       []string
		expectedStatus int
	}{
		{
			name:           "Valid repo URLs",
			repoUrls:       []string{"https://github.com/org/repo1", "https://github.com/org/repo2"},
			expectedStatus: http.StatusAccepted,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			body, _ := json.Marshal(map[string][]string{
				"repo_urls": tt.repoUrls,
			})
			req, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo1").WillReturnError(errors.New("sql: no rows in result set"))
			mock.ExpectExec("^-- name: InsertRepoURL :exec.*").WithArgs("https://github.com/org/repo1").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("^-- name: GetRepoURL :one.*").WithArgs("https://github.com/org/repo2").WillReturnError(errors.New("sql: no rows in result set"))
			mock.ExpectExec("^-- name: InsertRepoURL :exec.*").WithArgs("https://github.com/org/repo2").WillReturnResult(sqlmock.NewResult(1, 1))

			// Call the handler function
			apiCfg.HandlerAdd(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
