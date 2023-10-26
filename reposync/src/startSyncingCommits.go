package reposync

import (
	"context"
	"strconv"
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

		// Check if the repo is present in the repos directory
		if gitutil.IsGitRepository(prefixPath, organization, repo) {
			// If it is, pull the latest changes
			logger.LogBlue("repository %s exists. pulling...", repoUrl)
			gitutil.PullRepo(prefixPath, organization, repo)
			logger.LogBlue("repository %s pulled!", repoUrl)
		} else {
			// If not, clone it
			logger.LogBlue("repository %s does not exist. cloning...", repoUrl)
			gitutil.CloneRepo(prefixPath, organization, repo)
			logger.LogBlue("repository %s cloned!", repoUrl)
		}

		latestCommitterDate, err := db.GetLatestCommitterDate(context.Background(), repoUrl)
		if err != nil {
			logger.LogFatalRedAndExit("error getting latest committer date: %s ", err)
		}

		latestCommitterDateString := strconv.FormatInt(int64(latestCommitterDate), 10)
		startDate, err := ParseDate(latestCommitterDateString)
		if err != nil {
			logger.LogFatalRedAndExit("failed to parse latest commiter date %s to string", latestCommitterDate)
		}

		err = ProcessRepo(prefixPath, organization, repo, repoUrl, startDate, db)
		if err != nil {
			logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
		}
	}
}
