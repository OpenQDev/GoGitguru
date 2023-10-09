package server

import (
	"encoding/json"
	"fmt"
	"io"
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubReposByOwner(t *testing.T) {
	_, _, _, _, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")

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

	// Create a mock of Github REST API
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("jsonFile", jsonFile)
		io.Copy(w, jsonFile)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	var serverUrl string
	if targetLiveGithub {
		serverUrl = "https://api.github.com"
	} else {
		serverUrl = server.URL
	}

	apiCfg := ApiConfig{
		DB:                   queries,
		GithubRestAPIBaseUrl: serverUrl,
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
			owner:          "DRM-Test-Organization",
			expectedStatus: 401,
			authorized:     false,
			shouldError:    true,
		},
		{
			name:           "should return repos",
			owner:          "DRM-Test-Organization",
			expectedStatus: 200,
			authorized:     true,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			req, _ := http.NewRequest("GET", "/repos/github/"+tt.owner, nil)
			fmt.Println("req.URL", req.URL)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			mock.ExpectExec("^-- name: InsertGithubRepo :one.*").WithArgs(
				repos[0].ID,              // GithubRestID
				repos[0].NodeID,          // GithubGraphqlID
				repos[0].URL,             // Url
				repos[0].Name,            // Name
				repos[0].FullName,        // FullName
				repos[0].Private,         // Private
				repos[0].Owner.Login,     // OwnerLogin
				repos[0].Owner.AvatarURL, // OwnerAvatarUrl
				repos[0].Description,     // Description
				repos[0].Homepage,        // Homepage
				repos[0].Fork,            // Fork
				repos[0].ForksCount,      // ForksCount
				repos[0].Archived,        // Archived
				repos[0].Disabled,        // Disabled
				repos[0].License,         // License
				repos[0].Language,        // Language
				repos[0].StargazersCount, // StargazersCount
				repos[0].WatchersCount,   // WatchersCount
				repos[0].OpenIssuesCount, // OpenIssuesCount
				repos[0].HasIssues,       // HasIssues
				repos[0].HasDiscussions,  // HasDiscussions
				repos[0].HasProjects,     // HasProjects
				repos[0].CreatedAt,       // CreatedAt
				repos[0].UpdatedAt,       // UpdatedAt
				repos[0].PushedAt,        // PushedAt
				repos[0].Visibility,      // Visibility
				repos[0].Size,            // Size
				repos[0].DefaultBranch,   // DefaultBranch
			).WillReturnResult(sqlmock.NewResult(1, 1))

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
