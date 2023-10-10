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

	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubUserCommits(t *testing.T) {
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	// Initialize a new instance of ApiConfig with mocked DB
	ghAccessToken, targetLiveGithub, mock, queries := GetMockDatabase(ghAccessToken, targetLiveGithub)

	// Read the JSON file
	jsonFile, err := os.Open("mockUserCommitsReturn.json")
	if err != nil {
		t.Errorf("Failed to open mockReposReturn.json: %s", err)
		return
	}
	defer jsonFile.Close()

	// Parse the JSON file into a slice of RestRepo
	repos := ParseJsonFileToCommits(jsonFile)

	// Create a mock of Github REST API
	mux := http.NewServeMux()

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
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
		login          string
		expectedStatus int
		authorized     bool
		shouldError    bool
	}{
		// {
		// 	name:           "should return 401 if no access token",
		// 	login:          "DRM-Test-Organization",
		// 	expectedStatus: 401,
		// 	authorized:     false,
		// 	shouldError:    true,
		// },
		{
			name:           "should store repos for organization",
			login:          "IAmATestUserForDRM",
			expectedStatus: 200,
			authorized:     true,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/users/github/%s/commits", tt.login), nil)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			// Add {owner} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "login", tt.login)

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
			// impl here
			)

			// Call the handler function
			apiCfg.HandlerGithubUserCommits(rr, req)

			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.Equal(t, repos, "actual")

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func ParseJsonFileToCommits(jsonFile *os.File) []database.Commit {
	byteValue, _ := io.ReadAll(jsonFile)
	var commits []database.Commit
	json.Unmarshal(byteValue, &commits)
	jsonFile.Seek(0, 0)
	return commits
}
