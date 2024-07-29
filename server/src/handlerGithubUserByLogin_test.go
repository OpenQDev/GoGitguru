package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerGithubUserByLogin(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	debugMode := env.Debug
	ghAccessToken := env.GhAccessTokens
	targetLiveGithub := env.TargetLiveGithub

	logger.SetDebugMode(debugMode)

	jsonFile, err := os.Open("./mocks/mockGithubUserResponse.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	var user User
	err = marshaller.JsonFileToType(jsonFile, &user)
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

	tests := HandlerGithubUserByLoginTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.title)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()
			apiCfg := ApiConfig{
				DB:                   queries,
				GithubRestAPIBaseUrl: serverUrl,
			}

			// ARRANGE - LOCAL
			req, _ := http.NewRequest("GET", "", nil)
			// Add {owner} and {name} to the httptest.ResponseRecorder context since we are NOT calling this via Chi router
			req = AppendPathParamToChiContext(req, "login", tt.login)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			tt.setupMock(mock, user)

			// ACT
			apiCfg.HandlerGithubUserByLogin(rr, req)

			if rr.Code < 200 || rr.Code >= 300 {
				fmt.Println(rr.Body.String())
			}

			// ASSERT
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			}

			// ARRANGE - EXPECT
			var actualUserResponse ReturnUser
			err := json.NewDecoder(rr.Body).Decode(&actualUserResponse)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into User: %s", err)
				return
			}

			require.Equal(t, ConvertToReturnUser(user), actualUserResponse)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			require.Equal(t, tt.expectedStatus, rr.Code)

		})
	}
}
