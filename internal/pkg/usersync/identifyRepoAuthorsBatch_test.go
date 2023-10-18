package usersync

import (
	"io"
	"main/internal/pkg/githubGraphQL"
	"main/internal/pkg/logger"
	"main/internal/pkg/server"
	"main/internal/pkg/server/util"
	"main/internal/pkg/setup"
	"main/internal/pkg/testhelpers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIdentifyRepoAuthorsBatch(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, _, targetLiveGithub, _ := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	// Open the JSON file
	jsonFile, err := os.Open("./mocks/mockGithubCommitAuthorsResponse.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type RestRepo
	var commitAuthorsResponse githubGraphQL.CommitAuthorsResponse
	err = util.JsonFileToType(jsonFile, &commitAuthorsResponse)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}
	defer jsonFile.Close()

	mockGithubMux := http.NewServeMux()
	mockGithubMux.HandleFunc("/graphql/", func(w http.ResponseWriter, r *http.Request) {
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

	apiCfg := server.ApiConfig{
		GithubGraphQLBaseUrl: serverUrl,
	}

	tests := IdentifyRepoAuthorsBatchTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.title)

			// ACT
			_, err = identifyRepoAuthorsBatch(tt.repoUrl, tt.authorCommitList, "", apiCfg)
		})
	}
}
