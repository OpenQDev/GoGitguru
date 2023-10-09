package server

import (
	"encoding/json"
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
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

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
		owner          string
		expectedStatus int
		authorized     bool
		shouldError    bool
	}{
		// {
		// 	name:           "should return 401 if no access token",
		// 	owner:          "DRM-Test-Organization",
		// 	expectedStatus: 401,
		// 	authorized:     false,
		// 	shouldError:    true,
		// },
		{
			name:           "should store repos for organization",
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

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			// Add {owner} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "owner", tt.owner)

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
				repos[0].ID,              // 0 - GithubRestID
				repos[0].NodeID,          // 1 - GithubGraphqlID
				repos[0].URL,             // 2 - Url
				repos[0].Name,            // 3 - Name
				repos[0].FullName,        // 4 - FullName
				repos[0].Private,         // 5 - Private
				repos[0].Owner.Login,     // 6 - OwnerLogin
				repos[0].Owner.AvatarURL, // 7 - OwnerAvatarUrl
				repos[0].Description,     // 8 - Description
				repos[0].Homepage,        // 9 - Homepage
				repos[0].Fork,            // 10 - Fork
				repos[0].ForksCount,      // 11 - ForksCount
				repos[0].Archived,        // 12 - Archived
				repos[0].Disabled,        // 13 - Disabled
				repos[0].License,         // 14 - License
				repos[0].Language,        // 15 - Language
				repos[0].StargazersCount, // 16 - StargazersCount
				repos[0].WatchersCount,   // 17 - WatchersCount
				repos[0].OpenIssuesCount, // 18 - OpenIssuesCount
				repos[0].HasIssues,       // 19 - HasIssues
				repos[0].HasDiscussions,  // 20 - HasDiscussions
				repos[0].HasProjects,     // 21 - HasProjects
				repos[0].CreatedAt,       // 22 - CreatedAt
				repos[0].UpdatedAt,       // 23 - UpdatedAt
				repos[0].PushedAt,        // 24 - PushedAt
				repos[0].Visibility,      // 25 - Visibility
				repos[0].Size,            // 26 - Size
				repos[0].DefaultBranch,   // 27 - DefaultBranch
			)

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
