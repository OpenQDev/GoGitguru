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

	commitAuthorsResponse, err := gitutil.GithubGetCommitAuthors(queryString, ghAccessToken)

	if err != nil {
		logger.LogError("error occured while fetching from GraphQL API: %s", err)
	}

	if commitAuthorsResponse.Errors != nil {
		fmt.Printf("GraphQL Error: %v\n", commitAuthorsResponse.Errors)
		fmt.Println("Skipping...")
		return
	}

	if commitAuthorsResponse.Data == nil {
		fmt.Println("Skipping...")
		return
	}

	commits := make(map[string]gitutil.Author, 0)
	for key, value := range commitAuthorsResponse.Data.Repository {
		commits[key] = value
	}
	fmt.Println(commits)

}
