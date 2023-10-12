package gitutil

import (
	"fmt"
	"main/internal/pkg/setup"
	"testing"
)

func TestGithubGetCommitAuthors(t *testing.T) {
	_, _, _, _, _, _, _, _, ghAccessToken, _ := setup.ExtractAndVerifyEnvironment("../../../.env")

	tests := GithubGetCommitAuthorsTestCases()

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
