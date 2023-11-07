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

func TestHandlerGithubReposByOwner(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment(".env")
	debugMode := env.Debug
	ghAccessToken := env.GhAccessToken
	targetLiveGithub := env.TargetLiveGithub

	logger.SetDebugMode(debugMode)

	mock, queries := setup.GetMockDatabase()

	// ARRANGE - TEST DATA

	// Open the JSON file
	jsonFile, err := os.Open("./mocks/mockGithubReposReturn.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type []GithubRestRepo
	var repos []githubRest.GithubRestRepo
	err = marshaller.JsonFileToType(jsonFile, &repos)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}

	defer jsonFile.Close()

	mockGithubMux := http.NewServeMux()

	mockGithubMux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
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

	// ARRANGE - TESTS
	tests := HandlerGithubReposByOwnerTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)
			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)
			// Add {owner} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "owner", tt.owner)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			tt.setupMock(mock, repos[0])

			// ACT
			apiCfg.HandlerGithubReposByOwner(rr, req)

			// ARRANGE - EXPECT
			var actualReposReturn []GitguruRepo
			marshaller.ReaderToType(rr.Body, &actualReposReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into []GithubRestRepo: %s", err)
				return
			}

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			assert.Equal(t, tt.expectedStatus, rr.Code)

			convertedRepos := []GitguruRepo{}
			for _, repo := range repos {
				convertedRepos = append(convertedRepos, RestRepoToGitguruRepo(repo))
			}
			actualReposReturn = convertedRepos

			assert.Equal(t, convertedRepos, actualReposReturn)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
