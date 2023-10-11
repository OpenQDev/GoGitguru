package server

import (
	"io"
	"main/internal/database"
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

func TestHandlerRepoCommits(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	mock, queries := mocks.GetMockDatabase()

	jsonFile, err := os.Open("./mocks/mockUserCommitsReturn.json")
	if err != nil {
		t.Errorf("Failed to open ./mocks/mockUserCommitsReturn.json: %s", err)
		return
	}

	var repo []database.Commit
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
			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)

			// Add {login} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = mocks.AppendPathParamToChiContext(req, "login", tt.login)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			// ACT
			apiCfg.HandlerGithubUserCommits(rr, req)

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, repo, "actual")

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
