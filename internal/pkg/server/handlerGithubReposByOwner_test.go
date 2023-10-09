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
	"time"

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
	repos := ParseJsonFileToRestRepos(jsonFile)

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

			createdAt, _ := time.Parse(time.RFC3339, repos[0].CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, repos[0].UpdatedAt)
			pushedAt, _ := time.Parse(time.RFC3339, repos[0].PushedAt)

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repos[0].ID, repos[0].NodeID, repos[0].URL, repos[0].Name, repos[0].FullName, repos[0].Private, repos[0].Owner.Login, repos[0].Owner.AvatarURL, repos[0].Description, "homepage", repos[0].Fork, repos[0].ForksCount, repos[0].Archived, repos[0].Disabled, "license", "language", repos[0].StargazersCount, repos[0].WatchersCount, repos[0].OpenIssuesCount, repos[0].HasIssues, repos[0].HasDiscussions, repos[0].HasProjects, createdAt, updatedAt, pushedAt, repos[0].Visibility, repos[0].Size, repos[0].DefaultBranch)

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
				"homepage",               // 9 - Homepage
				repos[0].Fork,            // 10 - Fork
				repos[0].ForksCount,      // 11 - ForksCount
				repos[0].Archived,        // 12 - Archived
				repos[0].Disabled,        // 13 - Disabled
				"license",                // 14 - License
				"language",               // 15 - Language
				repos[0].StargazersCount, // 16 - StargazersCount
				repos[0].WatchersCount,   // 17 - WatchersCount
				repos[0].OpenIssuesCount, // 18 - OpenIssuesCount
				repos[0].HasIssues,       // 19 - HasIssues
				repos[0].HasDiscussions,  // 20 - HasDiscussions
				repos[0].HasProjects,     // 21 - HasProjects
				createdAt,                // 22 - CreatedAt
				updatedAt,                // 23 - UpdatedAt
				pushedAt,                 // 24 - PushedAt
				repos[0].Visibility,      // 25 - Visibility
				repos[0].Size,            // 26 - Size
				repos[0].DefaultBranch,   // 27 - DefaultBranch
			).WillReturnRows(rows)

			// Call the handler function
			apiCfg.HandlerGithubReposByOwner(rr, req)

			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			_ = struct{}{}
			fmt.Println(rr.Body.String())
			// assert.Equal(t, expectedResponse, rr.Body.String())

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func ParseJsonFileToRestRepos(jsonFile *os.File) []RestRepo {
	byteValue, _ := io.ReadAll(jsonFile)
	var repos []RestRepo
	json.Unmarshal(byteValue, &repos)
	jsonFile.Seek(0, 0)
	return repos
}
