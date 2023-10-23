package usersync

import (
	"fmt"

	"github.com/OpenQDev/GoGitguru/util/githubGraphQL"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func identifyRepoAuthorsBatch(repoUrl string, authorCommitList []AuthorCommitTuple, ghAccessToken string, githubGraphQLBaseUrl string) (map[string]githubGraphQL.GithubGraphQLCommit, error) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorCommitList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	queryString := generateAuthorBatchGqlQuery(organization, repo, authorCommitList)

	result, err := githubGraphQL.GithubGetCommitAuthors(queryString, ghAccessToken, githubGraphQLBaseUrl)

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

	commits := make(map[string]githubGraphQL.GithubGraphQLCommit, 0)
	for key, value := range result.Data.Repository {
		commits[key] = value
	}

	return commits, nil
}
