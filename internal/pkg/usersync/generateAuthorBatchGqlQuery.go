package usersync

import (
	"fmt"
	"main/internal/pkg/githubGraphQL"
)

func generateAuthorBatchGqlQuery(organization string, repo string, authorList []AuthorCommitTuple) string {
	gqlQuery := fmt.Sprintf(`{
		rateLimit {
			limit
			used
			resetAt
		}
		repository(owner: "%s", name: "%s") {`, organization, repo)

	// gql_query operates on the repository level, ordered by repositoryUrl
	// prepares 100 of these objects
	for i, commit := range authorList {
		gqlQuery += fmt.Sprintf(`
			commit_%d: object(oid: "%s") {
				...commitDetails
			}`, i, commit.CommitHash)
	}

	gqlQuery += `
		}
	}
	`
	// author_graphql_fragment is not defined in the original code, assuming it's a string
	gqlQuery += githubGraphQL.AUTHOR_GRAPHQL_FRAGMENT

	return gqlQuery
}
