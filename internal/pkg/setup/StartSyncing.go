package setup

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"sort"
	"strings"
	"time"
)

func StartSyncing(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {
	// Fetch all repository URLs
	repoUrlObjects, err := db.GetRepoURLs(context.Background())

	// Prepare an alphabetical, lowercased list of only repoUrls
	// Log them with new lines
	repoUrls := sortRepoUrls(repoUrlObjects)
	logger.LogGreenDebug("beginning sync for the following repos:\n%v", strings.Join(repoUrls, "\n"))

	if err != nil {
		logger.LogFatalRedAndExit("error getting repo urls: %s ", err)
	}

	if err != nil {
		logger.LogFatalRedAndExit("Failed to insert repo url: %s", err)
	}

	for _, repoUrl := range repoUrls {
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

		defer gitutil.DeleteLocalRepo(prefixPath, repo)

		// just returns an error and continues if already there. otherwise clones
		// no need to even check for "isGitRepo"
		gitutil.CloneRepo(prefixPath, organization, repo)

		err := gitutil.ProcessRepo(prefixPath, repo, repoUrl, db)
		if err != nil {
			logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
		}
	}
}

func sortRepoUrls(repoUrlObjects []database.RepoUrl) []string {
	repoUrls := make([]string, len(repoUrlObjects))

	for i, repo := range repoUrlObjects {
		// since sort.Strings uses case-sensitive lexicographic ordering, we must lowercase
		repoUrls[i] = strings.ToLower(repo.Url)
	}

	sort.Strings(repoUrls)
	return repoUrls
}
