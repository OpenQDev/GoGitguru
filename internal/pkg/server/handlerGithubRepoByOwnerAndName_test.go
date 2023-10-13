package server

import (
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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGithubRepoByOwnerAndName(t *testing.T) {
	// ARRANGE - GLOBAL
	_, _, _, debugMode, _, _, _, _, ghAccessToken, targetLiveGithub := setup.ExtractAndVerifyEnvironment("../../../.env")
	logger.SetDebugMode(debugMode)

	mock, queries := mocks.GetMockDatabase()

	// Open the JSON file
	jsonFile, err := os.Open("./mocks/mockGithubRepoReturn.json")
	if err != nil {
		t.Errorf("error opening json file: %s", err)
	}

	// Decode the JSON file to type RestRepo
	var repo GithubRestRepo
	err = util.JsonFileToType(jsonFile, &repo)
	if err != nil {
		t.Errorf("Failed to read JSON file: %s", err)
	}
	defer jsonFile.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/", func(w http.ResponseWriter, r *http.Request) {
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
			req = mocks.AppendPathParamToChiContext(req, "name", tt.name)
			req = mocks.AppendPathParamToChiContext(req, "owner", tt.owner)

			if tt.authorized {
				req.Header.Add("GH-Authorization", ghAccessToken)
			}

			rr := httptest.NewRecorder()

			createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
			updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
			pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)

			rows := sqlmock.NewRows([]string{"internal_id", "github_rest_id", "github_graphql_id", "url", "name", "full_name", "private", "owner_login", "owner_avatar_url", "description", "homepage", "fork", "forks_count", "archived", "disabled", "license", "language", "stargazers_count", "watchers_count", "open_issues_count", "has_issues", "has_discussions", "has_projects", "created_at", "updated_at", "pushed_at", "visibility", "size", "default_branch"}).
				AddRow(1, repo.GithubRestID, repo.GithubGraphqlID, repo.URL, repo.Name, repo.FullName, repo.Private, repo.Owner.Login, repo.Owner.AvatarURL, repo.Description, "homepage", repo.Fork, repo.ForksCount, repo.Archived, repo.Disabled, "license", "language", repo.StargazersCount, repo.WatchersCount, repo.OpenIssuesCount, repo.HasIssues, repo.HasDiscussions, repo.HasProjects, createdAt, updatedAt, pushedAt, repo.Visibility, repo.Size, repo.DefaultBranch)

			mock.ExpectQuery("^-- name: InsertGithubRepo :one.*").WithArgs(
				repo.GithubRestID,    // 0 - GithubRestID
				repo.GithubGraphqlID, // 1 - GithubGraphqlID
				repo.URL,             // 2 - Url
				repo.Name,            // 3 - Name
				repo.FullName,        // 4 - FullName
				repo.Private,         // 5 - Private
				repo.Owner.Login,     // 6 - OwnerLogin
				repo.Owner.AvatarURL, // 7 - OwnerAvatarUrl
				repo.Description,     // 8 - Description
				"",                   // 9 - Homepage
				repo.Fork,            // 10 - Fork
				repo.ForksCount,      // 11 - ForksCount
				repo.Archived,        // 12 - Archived
				repo.Disabled,        // 13 - Disabled
				"",                   // 14 - License
				"",                   // 15 - Language
				repo.StargazersCount, // 16 - StargazersCount
				repo.WatchersCount,   // 17 - WatchersCount
				repo.OpenIssuesCount, // 18 - OpenIssuesCount
				repo.HasIssues,       // 19 - HasIssues
				repo.HasDiscussions,  // 20 - HasDiscussions
				repo.HasProjects,     // 21 - HasProjects
				createdAt,            // 22 - CreatedAt
				updatedAt,            // 23 - UpdatedAt
				pushedAt,             // 24 - PushedAt
				repo.Visibility,      // 25 - Visibility
				repo.Size,            // 26 - Size
				repo.DefaultBranch,   // 27 - DefaultBranch
			).WillReturnRows(rows)

			// ACT
			apiCfg.HandlerGithubRepoByOwnerAndName(rr, req)

			// ASSERT - ERROR
			if tt.shouldError {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				return
			} else if rr.Code < 200 || rr.Code >= 300 {
				t.Errorf("Unexpected HTTP status code: %d. Error: %s", rr.Code, err)
				return
			}

			// ARRANGE - EXPECT
			var actualRepoReturn GithubRestRepo
			util.ReaderToType(rr.Result().Body, &actualRepoReturn)
			if err != nil {
				t.Errorf("Failed to decode rr.Body into []RestRepo: %s", err)
				return
			}

			assert.Equal(t, repo, actualRepoReturn)
		})
	}
}
