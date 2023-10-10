package server

import (
	"encoding/json"
	"fmt"
	"io"
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubRepoByOwnerAndName(t *testing.T) {
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	// Read the JSON file
	jsonFile, err := os.Open("./mocks/mockRepoReturn.json")
	if err != nil {
		t.Errorf("Failed to open ./mocks/mockReposReturn.json: %s", err)
		return
	}
	defer jsonFile.Close()

	// Parse the JSON file into a slice of RestRepo
	repo := ParseJsonFileToRestRepo(jsonFile)

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
		title          string
		owner          string
		name           string
		expectedStatus int
		authorized     bool
		shouldError    bool
	}{
		// {
		// 	title:          "should return 401 if no access token",
		// 	owner:          "DRM-Test-Organization",
		// 	name:           "DRM-Test-Repo",
		// 	expectedStatus: 401,
		// 	authorized:     false,
		// 	shouldError:    true,
		// },
		{
			title:          "should store repos for organization",
			owner:          "DRM-Test-Organization",
			name:           "DRM-Test-Repo",
			expectedStatus: 200,
			authorized:     true,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the HTTP request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/repos/github/%s%s", tt.owner, tt.name), nil)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			// Add {owner} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = mocks.AppendPathParamToChiContext(req, "owner", tt.owner)
			req = mocks.AppendPathParamToChiContext(req, "name", tt.name)

			// Call the handler function
			apiCfg.HandlerGithubReposByOwner(rr, req)

			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			var actualRepoReturn RestRepo
			err := json.NewDecoder(rr.Body).Decode(&actualRepoReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into []RestRepo: %s", err)
				return
			}

			assert.Equal(t, repo, actualRepoReturn)
		})
	}
}

func ParseJsonFileToRestRepo(jsonFile *os.File) RestRepo {
	byteValue, _ := io.ReadAll(jsonFile)
	var repo RestRepo
	json.Unmarshal(byteValue, &repo)
	jsonFile.Seek(0, 0)
	return repo
}
