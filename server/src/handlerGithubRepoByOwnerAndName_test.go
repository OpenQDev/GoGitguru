package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/githubRest"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubRepoByOwnerAndName(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment(".env")
	debugMode := env.Debug
	ghAccessToken := env.GhAccessToken
	targetLiveGithub := env.TargetLiveGithub

	logger.SetDebugMode(debugMode)

	mock, queries := setup.GetMockDatabase()

	// Open the JSON file
	jsonFile, err := os.Open("./mocks/mockGithubRepoReturn.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type RestRepo
	var repo githubRest.GithubRestRepo
	err = marshaller.JsonFileToType(jsonFile, &repo)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}
	defer jsonFile.Close()
	expectedGitguruRepo := RestRepoToGitguruRepo(repo)

	mockGithubMux := http.NewServeMux()
	mockGithubMux.HandleFunc("/repos/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, jsonFile)
	})
	mockGithubServer := httptest.NewServer(mockGithubMux)
	defer mockGithubServer.Close()

	var serverUrl string
	if targetLiveGithub {
		serverUrl = "https://api.github.com"
	} else {
		serverUrl = mockGithubServer.URL
	}

	apiCfg := ApiConfig{
		DB:                   queries,
		GithubRestAPIBaseUrl: serverUrl,
	}

	tests := HandlerGithubRepoByOwnerAndNameTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.title)

			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)

			// Add {owner} and {name} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "name", tt.name)
			req = AppendPathParamToChiContext(req, "owner", tt.owner)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			tt.setupMock(mock, repo)

			// ACT
			apiCfg.HandlerGithubRepoByOwnerAndName(rr, req)

			// ASSERT - ERROR
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			} else if rr.Code < 200 || rr.Code >= 300 {
				t.Errorf("Unexpected HTTP status code: %d. Response: %s", rr.Code, rr.Body.String())
				return
			}

			// ARRANGE - EXPECT
			var actualRepoReturn GitguruRepo
			marshaller.ReaderToType(rr.Result().Body, &actualRepoReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into GithubRestRepo: %s", err)
				return
			}

			assert.Equal(t, expectedGitguruRepo, actualRepoReturn)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
