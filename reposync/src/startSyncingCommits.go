package reposync

import (
	"context"
	"database/database"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func StartSyncingCommits(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {

	repoUrlObjects, err := db.GetRepoURLs(context.Background())

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
		logger.LogGreenDebug("processing %s/%s...", organization, repo)

		defer gitutil.DeleteLocalRepo(prefixPath, organization, repo)

		// just returns an error and continues if already there. otherwise clones
		// no need to even check for "isGitRepo"
		gitutil.CloneRepo(prefixPath, organization, repo)

		err := ProcessRepo(prefixPath, organization, repo, repoUrl, db)
		if err != nil {
			logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
		}
	}
}
