package sync

import (
	"main/internal/pkg/gitutil"
	"main/internal/pkg/testhelpers"
	"testing"
)

func TestGenerateAuthorBatchGqlQuery(t *testing.T) {
	tests := []struct {
		name           string
		organization   string
		repo           string
		authorList     []AuthorCommitTuple
		expectedOutput string
	}{
		{
			name:         "Single author",
			organization: "testOrg",
			repo:         "testRepo",
			authorList:   []AuthorCommitTuple{{Author: "author1", CommitHash: "commit1"}},
			expectedOutput: `{
		rateLimit {
			limit
			used
			resetAt
		}
		repository(owner: "testOrg", name: "testRepo") {
			commit_0: object(oid: "commit1") {
				...commitDetails
			}
		}
	}
	` + gitutil.AUTHOR_GRAPHQL_FRAGMENT,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GenerateAuthorBatchGqlQuery(tt.organization, tt.repo, tt.authorList)
			if result != tt.expectedOutput {
				t.Errorf("generateAuthorBatchGqlQuery() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
