package gitutil

import (
	"fmt"
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"testing"
)

func TestGithubGetCommitAuthors(t *testing.T) {
	_, _, _, _, _, _, _, _, ghAccessToken := setup.ExtractAndVerifyEnvironment("../../../.env")

	query := `{
		rateLimit {
			limit
			used
			resetAt
		}
		repository(owner: "OpenQDev", name: "OpenQ-Workflows") {
			commit_0: object(oid: "8799411585c826b577f632f1ef5c0415914267ed") {
				...commitDetails
			}
		}
	}
	` + AUTHOR_GRAPHQL_FRAGMENT

	// Use server.URL as the graphql API endpoint
	res, err := GithubGetCommitAuthors(query, ghAccessToken)
	if err != nil {
		logger.LogError("error in GithubGetCommitAuthors", err)
	}

	fmt.Println(res)
}
