package sync

import (
	"fmt"
	"main/internal/pkg/gitutil"
	"testing"
)

func TestGenerateAuthorBatchGqlQuery(t *testing.T) {
	tests := []struct {
		name           string
		organization   string
		repo           string
		authorList     []string
		expectedOutput string
	}{
		{
			name:         "Single author",
			organization: "testOrg",
			repo:         "testRepo",
			authorList:   []string{"author1"},
			expectedOutput: `{
		rateLimit {
			limit
			used
			resetAt
		}
		repository(owner: "testOrg", name: "testRepo") {
			commit_0: object(oid: "author1") {
				...commitDetails
			}
		}
	}
	` + gitutil.AUTHOR_GRAPHQL_FRAGMENT,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateAuthorBatchGqlQuery(tt.organization, tt.repo, tt.authorList)
			fmt.Println(result)
			if result != tt.expectedOutput {
				t.Errorf("generateAuthorBatchGqlQuery() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
