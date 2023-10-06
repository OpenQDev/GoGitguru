package usersync

import (
	"fmt"
	"main/internal/pkg/gitutil"
)

func GenerateAuthorBatchGqlQuery(organization string, repo string, authorList []string) string {
	gqlQuery := fmt.Sprintf(`{
		rateLimit {
			limit
			used
			resetAt
		}
		repository(owner: "%s", name: "%s") {`, organization, repo)

	// gql_query operates on the repository level, ordered by repositoryUrl
	// prepares 100 of these objects
	for i, author := range authorList {
		gqlQuery += fmt.Sprintf(`
			commit_%d: object(oid: "%s") {
				...commitDetails
			}`, i, author)
	}

	gqlQuery += `
		}
	}
	`
	// author_graphql_fragment is not defined in the original code, assuming it's a string
	gqlQuery += gitutil.AUTHOR_GRAPHQL_FRAGMENT

	return gqlQuery
}
