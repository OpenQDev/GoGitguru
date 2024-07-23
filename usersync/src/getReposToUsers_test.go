package usersync

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestGetReposToUsers(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug

	logger.SetDebugMode(debugMode)

	// Open the JSON file
	jsonFile, err := os.Open("../mocks/mockGithubCommitAuthorsResponse_oneAuthor.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// ARRANGE - GLOBAL

	mockGithubMux := http.NewServeMux()
	mockGithubMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, jsonFile)
	})
	mockGithubServer := httptest.NewServer(mockGithubMux)
	defer mockGithubServer.Close()

	testcases := GetReposToUsersTestCases()

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()

			// ARRANGE - LOCAL
			tt.setupMock(mock, tt.author)

			// ACT
			GetReposToUsers(
				queries,
				&tt.initialParams, tt.internal_id, tt.author)

			// ASSERT

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if reflect.DeepEqual(tt.expectedParams, tt.initialParams) {
				t.Errorf("expected: %v, got: %v", tt.expectedParams, tt.initialParams)
			}

		})
	}
}
