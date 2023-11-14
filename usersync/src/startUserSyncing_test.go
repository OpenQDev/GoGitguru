package usersync

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestStartUserSync(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../.env")
	debugMode := env.Debug
	targetLiveGithub := env.TargetLiveGithub

	logger.SetDebugMode(debugMode)

	// Open the JSON file
	jsonFile, err := os.Open("../mocks/mockGithubCommitAuthorsResponse_oneAuthor.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// ARRANGE - GLOBAL
	mock, queries := setup.GetMockDatabase()

	mockGithubMux := http.NewServeMux()
	mockGithubMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, jsonFile)
	})
	mockGithubServer := httptest.NewServer(mockGithubMux)
	defer mockGithubServer.Close()

	var serverUrl string
	if targetLiveGithub {
		serverUrl = "https://api.github.com/graphql"
	} else {
		serverUrl = mockGithubServer.URL
	}

	testcases := StartUserSyncingTestCases()

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			tt.setupMock(mock, tt.author)

			// ACT
			StartSyncingUser(
				queries,
				"mock",
				"",
				10,
				serverUrl)

			// ASSERT
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
