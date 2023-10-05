package sync

import (
	"fmt"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
)

func identifyRepoAuthorsBatch(repoUrl string, authorList []string, ghAccessToken string) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	queryString := generateAuthorBatchGqlQuery(organization, repo, authorList)

	data, err := gitutil.GithubGraphQL(queryString, ghAccessToken)

	if err != nil {
		logger.LogError("error occured while fetching from GraphQL API: %s", err)
	}

	if _, ok := data["errors"]; ok {
		fmt.Printf("GraphQL Error: %v\n", data["errors"])
		fmt.Println("Skipping...")
		return
	}

	if data == nil {
		fmt.Println("Skipping...")
		return
	}

}
