package usersync

import (
	"fmt"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestGithubGetCommitAuthors(t *testing.T) {
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	ghAccessToken := env.GhAccessTokens
	targetLiveGithub := env.TargetLiveGithub

	var serverUrl string
	if targetLiveGithub {
		serverUrl = "https://api.github.com/graphql"
	} else {
		serverUrl = "https://api.github.com/graphql"
	}

	tests := GithubGetCommitAuthorsTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			var commitDetails string

			for i, commit := range tt.commits {
				commitDetails += fmt.Sprintf(`commit_%d: object(oid: "%s") {
					...commitDetails
				}
				`, i, commit)
			}

			query := fmt.Sprintf(`{
				rateLimit {
					limit
					used
					resetAt
				}
				repository(owner: "%s", name: "%s") {
					%s
				}
			}
			`, tt.owner, tt.repo, commitDetails) + AUTHOR_GRAPHQL_FRAGMENT

			resp, err := GithubGetCommitAuthors(query, ghAccessToken, serverUrl)

			if (err != nil) != tt.wantErr {
				t.Errorf("GithubGetCommitAuthors() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualRestId := resp.Data.Repository["commit_1"].Author.User.GithubRestID
			expectedRestId := 93455288
			if actualRestId != int32(expectedRestId) {
				t.Errorf("GithubGetCommitAuthors() unexpected return. expect rest ID of %d but got %d", expectedRestId, actualRestId)
			}
		})
	}
}
