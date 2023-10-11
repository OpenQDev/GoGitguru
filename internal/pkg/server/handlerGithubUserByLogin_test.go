package server

import (
	"encoding/json"
	"io"
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubUserByLogin(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	jsonFile, err := os.Open("./mocks/mockRepoReturn.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	var repo RestRepo
	err = util.JsonFileToType(jsonFile, &repo)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}
	defer jsonFile.Close()

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
			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)
			// Add {owner} and {name} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = mocks.AppendPathParamToChiContext(req, "owner", tt.owner)
			req = mocks.AppendPathParamToChiContext(req, "name", tt.name)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			// ACT
			apiCfg.HandlerGithubReposByOwner(rr, req)

			// ARRANGE - EXPECT
			var actualRepoReturn RestRepo
			err := json.NewDecoder(rr.Body).Decode(&actualRepoReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into []RestRepo: %s", err)
				return
			}

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			assert.Equal(t, repo, actualRepoReturn)
		})
	}
}
