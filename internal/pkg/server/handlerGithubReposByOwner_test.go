package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/database"
	"main/internal/pkg/logger"
	"net/http"
	"net/http/httptest"
	"os"
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

	// Read the JSON file
	jsonFile, err := os.Open("mockReposReturn.json")
	if err != nil {
		t.Errorf("Failed to open mockReposReturn.json: %s", err)
		return
	}
	defer jsonFile.Close()

	// Parse the JSON file into a slice of RestRepo
	byteValue, _ := io.ReadAll(jsonFile)
	var repos []RestRepo
	json.Unmarshal(byteValue, &repos)

	fmt.Println(repos)

	// Create a mock of Github REST API
	mux := http.NewServeMux()
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, jsonFile)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	apiCfg := ApiConfig{
		DB:                   queries,
		GithubRestAPIBaseUrl: server.URL,
	}

	// Define test cases
	tests := []struct {
		name           string
		owner          string
		expectedStatus int
		authorized     bool
		shouldError    bool
	}{
		{
			name:           "should return 401 if no access token",
			owner:          "VipassanaInsight",
			expectedStatus: 401,
			authorized:     false,
			shouldError:    true,
		},
		{
			name:           "should return repos",
			owner:          "TestOwner",
			expectedStatus: 200,
			authorized:     true,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			req, _ := http.NewRequest("GET", "/repos/github/"+tt.owner, nil)

			if tt.authorized {
				req.Header.Add("GH-Authorization", "foo")
			}

			rr := httptest.NewRecorder()

			mock.ExpectExec("^-- name: InsertGithubRepo :one.*").WithArgs(
				repos[0].ID,                              // GithubRestID
				repos[0].NodeID,                          // GithubGraphqlID
				"https://github.com/octocat/Hello-World", // Url
				"Hello-World",                            // Name
				"octocat/Hello-World",                    // FullName
				false,                                    // Private
				"octocat",                                // OwnerLogin
				"https://github.com/images/error/octocat_happy.gif", // OwnerAvatarUrl
				"This your first repo!",                             // Description
				"",                                                  // Homepage
				false,                                               // Fork
				9,                                                   // ForksCount
				false,                                               // Archived
				false,                                               // Disabled
				"mit",                                               // License
				"Ruby",                                              // Language
				80,                                                  // StargazersCount
				80,                                                  // WatchersCount
				9,                                                   // OpenIssuesCount
				true,                                                // HasIssues
				false,                                               // HasDiscussions
				true,                                                // HasProjects
				sql.NullTime{},                                      // CreatedAt
				sql.NullTime{},                                      // UpdatedAt
				sql.NullTime{},                                      // PushedAt
				"public",                                            // Visibility
				108,                                                 // Size
				"master").WillReturnResult(sqlmock.NewResult(1, 1))

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
