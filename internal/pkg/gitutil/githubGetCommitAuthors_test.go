package gitutil

import (
	"fmt"
	"main/internal/pkg/setup"
	"testing"
)

func TestGithubGetCommitAuthors(t *testing.T) {
	_, _, _, _, _, _, _, _, ghAccessToken := setup.ExtractAndVerifyEnvironment("../../../.env")

	tests := []struct {
		name    string
		owner   string
		repo    string
		commits []string
		wantErr bool
	}{
		{
			name:    "TestGithubGetCommitAuthors with valid query",
			owner:   "OpenQDev",
			repo:    "OpenQ-Workflows",
			commits: []string{"8799411585c826b577f632f1ef5c0415914267ed", "657bd8b7f7d83e8b842411cbf65666901d65431c"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			resp, err := GithubGetCommitAuthors(query, ghAccessToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("GithubGetCommitAuthors() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualRestId := resp.Data.Repository["commit_1"].Author.User.GithubRestID
			expectedRestId := 93455288
			if actualRestId != expectedRestId {
				t.Errorf("GithubGetCommitAuthors() unexpected return. expect rest ID of %d but got %d", expectedRestId, actualRestId)
			}
		})
	}
}
