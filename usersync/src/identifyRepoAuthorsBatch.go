package usersync

import (
	"fmt"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func identifyRepoAuthorsBatch(repoUrl string, authorCommitList []AuthorCommitTuple, ghAccessToken string, githubGraphQLBaseUrl string) (map[string]GithubGraphQLCommit, error) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorCommitList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	queryString := generateAuthorBatchGqlQuery(organization, repo, authorCommitList)

	result, err := GithubGetCommitAuthors(queryString, ghAccessToken, githubGraphQLBaseUrl)

	if err != nil {
		logger.LogError("error occured while fetching from GraphQL API: %s", err)
	}

	if result.Errors != nil {
		fmt.Printf("skipping due to graphQL error: %v\n", result.Errors)
		return nil, err
	}

	if result.Data == nil {
		logger.LogError("github graphQL api return no data for %s and %s", repoUrl, authorCommitList)
		return nil, nil
	}

	commits := make(map[string]GithubGraphQLCommit, 0)
	for key, value := range result.Data.Repository {
		commits[key] = GithubGraphQLCommit{
			Author: value.Author,
		}
	}

	return commits, nil
}
