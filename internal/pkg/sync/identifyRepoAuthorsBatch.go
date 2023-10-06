package sync

import (
	"fmt"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
)

func IdentifyRepoAuthorsBatch(repoUrl string, authorCommitList []AuthorCommitTuple, ghAccessToken string) (*map[string]gitutil.Commit, error) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorCommitList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	queryString := GenerateAuthorBatchGqlQuery(organization, repo, authorCommitList)

	commitAuthorsResponse, err := gitutil.GithubGetCommitAuthors(queryString, ghAccessToken)

	logger.LogGreenDebug("commit authors response: %v", commitAuthorsResponse)

	if err != nil {
		logger.LogError("error occured while fetching from GraphQL API: %s", err)
	}

	if commitAuthorsResponse.Errors != nil {
		fmt.Printf("skipping due to graphQL error: %v\n", commitAuthorsResponse.Errors)
		fmt.Println()
		return nil, err
	}

	if commitAuthorsResponse.Data == nil {
		logger.LogError("github graphQL api return no data for %s and %s", repoUrl, authorCommitList)
		return nil, nil
	}

	commits := make(map[string]gitutil.Commit, 0)
	for key, value := range commitAuthorsResponse.Data.Repository {
		commits[key] = value
	}
	fmt.Println(commits)

	return &commits, nil
}
