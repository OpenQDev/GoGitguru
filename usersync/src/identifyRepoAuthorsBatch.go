package usersync

import (
	"fmt"
	"strconv"
	"strings"

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
		fmt.Println()
		return nil, err
	}

	if result.Data == nil {
		logger.LogError("github graphQL api return no data for %s and %s", repoUrl, authorCommitList)
		return nil, nil
	}

	commits := make(map[string]GithubGraphQLCommit, 0)
	for key, value := range result.Data.Repository {
		if value.Author.Email != "" {
			commits[key] = value
		}
		if value.Author.Email == "" {
			numStr := strings.TrimPrefix(key, "commit_")
			// Convert the remaining string to an integer
			correctIndex, err := strconv.Atoi(numStr)
			for i := range authorCommitList {
				if err != nil {
					return nil, err
				}
				if i == correctIndex {
					commits[key] = GithubGraphQLCommit{
						Author: GithubGraphQLAuthor{
							Name:  authorCommitList[i].Author,
							Email: authorCommitList[i].Author,
						},
					}
				}
			}
		}

	}

	return commits, nil
}
