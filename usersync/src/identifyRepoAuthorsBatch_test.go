package usersync

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/githubGraphQL"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestIdentifyRepoAuthorsBatch(t *testing.T) {
	// ARRANGE - GLOBAL
	env := setup.ExtractAndVerifyEnvironment("../.env")
	debugMode := env.Debug
	targetLiveGithub := env.TargetLiveGithub

	logger.SetDebugMode(debugMode)

	// Open the JSON file
	jsonFile, err := os.Open("../mocks/mockGithubCommitAuthorsResponse.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type RestRepo
	var commitAuthorsResponse githubGraphQL.GithubGraphQLCommitAuthorsResponse
	err = marshaller.JsonFileToType(jsonFile, &commitAuthorsResponse)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}
	defer jsonFile.Close()

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
		fmt.Println("serverUrl", serverUrl)
	}

	tests := IdentifyRepoAuthorsBatchTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.title)

			// ACT
			resp, err := identifyRepoAuthorsBatch(tt.repoUrl, tt.authorCommitList, "", serverUrl)
			if err != nil {
				t.Fatalf("error in identifyRepoAuthorsBatch test: %s", err)
			}

			if !reflect.DeepEqual(resp, tt.expectedOutput) {
				t.Errorf("Expected output does not match the response. Expected: %v, Got: %v", tt.expectedOutput, resp)
			}

		})
	}
}
