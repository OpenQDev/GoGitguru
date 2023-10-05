package sync

import (
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
)

func identifyRepoAuthorsBatch(repoUrl string, authorList []string) {
	logger.LogBlue("Identifying %d authors for repo %s", len(authorList), repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

	generateAuthorBatchGqlQuery(organization, repo, authorList)
}
