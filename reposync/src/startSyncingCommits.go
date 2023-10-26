package reposync

import (
	"context"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func StartSyncingCommits(
	db *database.Queries,
	prefixPath string) {

	repoUrlObjects, err := db.GetRepoURLs(context.Background())

	repoUrls := sortRepoUrls(repoUrlObjects)
	logger.LogGreenDebug("beginning sync for the following repos:\n%v", strings.Join(repoUrls, "\n"))

	if err != nil {
		logger.LogFatalRedAndExit("error getting repo urls: %s ", err)
	}

	for _, repoUrl := range repoUrls {
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)
		logger.LogGreenDebug("processing %s/%s...", organization, repo)

		defer gitutil.DeleteLocalRepo(prefixPath, organization, repo)

		// Check if the repo is present in the repos directory
		if !gitutil.IsGitRepo(prefixPath, organization, repo) {
			// If not, clone it
			logger.LogBlue("repository %s does not exist. cloning...", repoUrl)
			gitutil.CloneRepo(prefixPath, organization, repo)
			logger.LogBlue("repository %s cloned!", repoUrl)
		} else {
			// If it is, pull the latest changes
			logger.LogBlue("repository %s exists. pulling...", repoUrl)
			gitutil.PullRepo(prefixPath, organization, repo)
			logger.LogBlue("repository %s pulled!", repoUrl)
		}

		err := ProcessRepo(prefixPath, organization, repo, repoUrl, db)
		if err != nil {
			logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
		}
	}
}
