package usersync

import (
	"fmt"
	"main/internal/pkg/githubGraphQL"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
)

func IdentifyRepoAuthorsBatch(repoUrl string, authorCommitList []AuthorCommitTuple, ghAccessToken string) (*map[string]githubGraphQL.Commit, error) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorCommitList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	queryString := GenerateAuthorBatchGqlQuery(organization, repo, authorCommitList)

	result, err := githubGraphQL.GithubGetCommitAuthors(queryString, ghAccessToken)

	logger.LogGreenDebug("GithubGetCommitAuthors response: %v", result)

	if err != nil {
		logger.LogError("error occured while fetching from GraphQL API: %s", err)
	}

	if result.Errors != nil {
		fmt.Printf("skipping due to graphQL error: %v\n", result.Errors)
		fmt.Println()
		return nil, err
	}

	if result.Data == nil {
		logger.LogError("github graphQL api return no data for %s and %s", repoUrl, authorCommitList)
		return nil, nil
	}

	commits := make(map[string]githubGraphQL.Commit, 0)
	for key, value := range result.Data.Repository {
		commits[key] = value
	}

	return &commits, nil
}
