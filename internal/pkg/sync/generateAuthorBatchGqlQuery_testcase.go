package sync

import "main/internal/pkg/gitutil"

type GenerateAuthorBatchGqlQueryTestCase struct {
	name           string
	organization   string
	repo           string
	authorList     []AuthorCommitTuple
	expectedOutput string
}

func singleAuthor() GenerateAuthorBatchGqlQueryTestCase {
	return GenerateAuthorBatchGqlQueryTestCase{
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
	}
}

func GenerateAuthorBatchGqlQueryTestCases() []GenerateAuthorBatchGqlQueryTestCase {
	return []GenerateAuthorBatchGqlQueryTestCase{
		singleAuthor(),
	}
}
