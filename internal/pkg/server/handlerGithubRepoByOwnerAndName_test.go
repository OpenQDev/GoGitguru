package server

import (
	"encoding/json"
	"io"
	"main/internal/pkg/logger"
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"main/internal/pkg/testhelpers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubRepoByOwnerAndName(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	_, queries := mocks.GetMockDatabase()

	// Open the JSON file
	jsonFile, err := os.Open("./mocks/mockRepoReturn.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type RestRepo
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

	tests := HandlerGithubRepoByOwnerAndNameTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"SHOULD_GET_REPO_FOR_ORG_AND_NAME",
			), tt.title)

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

			// ASSERT - ERROR
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// ARRANGE - EXPECT
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
