package gitutil

import (
	"encoding/json"
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
			commit_1: object(oid: "657bd8b7f7d83e8b842411cbf65666901d65431c") {
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

	resBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		logger.LogError("error in json.MarshalIndent", err)
	}

	fmt.Println(string(resBytes))
}
