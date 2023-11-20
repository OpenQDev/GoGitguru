package reposync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

type GetDueUrlResponse struct {
	RepoUrl string `json:"repo_url"`
}

func StartSyncingCommits(
	db *database.Queries,
	prefixPath string,
	gitguruUrl string,
	getDueRepoUrlExpiration int,
) {

	for {
		url := fmt.Sprintf("%s/get-next-repo-url", gitguruUrl)
		resp, err := http.Get(url)
		if err != nil {
			logger.LogFatalRedAndExit("error getting repo url: %s ", err)
		}
		defer resp.Body.Close()

		var getDueUrlResponse GetDueUrlResponse
		err = json.NewDecoder(resp.Body).Decode(&getDueUrlResponse)
		if err != nil {
			logger.LogFatalRedAndExit("error decoding repo url: %s ", err)
		}

		repoUrl := getDueUrlResponse.RepoUrl

		if repoUrl == "" {
			break
		}

		logger.LogGreenDebug("beginning sync for the following repo:\n%s", repoUrl)

		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)
		logger.LogGreenDebug("processing %s/%s...", organization, repo)
		JAN_1_2020 := time.Unix(1577858400, 0)

		startDate := JAN_1_2020 // Unix time for Jan 1, 2020

		// Check if the repo is present in the repos directory
		if gitutil.IsGitRepository(prefixPath, organization, repo) {
			// If it is, pull the latest changes
			logger.LogBlue("repository %s exists. pulling...", repoUrl)
			err := gitutil.PullRepo(prefixPath, organization, repo)
			if err != nil {
				logger.LogError("error pulling repo %s/%s: %s", organization, repo, err)

				logger.LogError("deleting repo url %s/%s since it does not exist, is private, too large, or is empty", organization, repo)
				err := db.DeleteRepoURL(context.Background(), repoUrl)
				if err != nil {
					logger.LogError("error deleting repo url %s: %s", repoUrl, err)
				}
				logger.LogError("repo url %s/%s deleted!", organization, repo)

				continue
			}
			logger.LogBlue("repository %s pulled!", repoUrl)

			// there are cases where the repository may exist in local, but hasn't been synced
			// no rows in result set just means it didn't have any commit entries for that repo
			latestCommitterDate, err := db.GetLatestCommitterDate(context.Background(), repoUrl)

			if err != nil {
				if !strings.Contains(err.Error(), "sql: no rows in result set") {
					logger.LogFatalRedAndExit("error getting latest committer date: %s ", err)
				}
			}

			latestCommitterDateTime := time.Unix(int64(latestCommitterDate), 0)

			// Unsure why but sometimes commits before JAN_1_2020 were being stored after initia clone-sync, causing issues
			if latestCommitterDateTime.After(JAN_1_2020) {
				startDate = latestCommitterDateTime
			}
		} else {
			// If not, clone it
			logger.LogBlue("repository %s does not exist. cloning...", repoUrl)
			err := gitutil.CloneRepo(prefixPath, organization, repo)
			if err != nil {
				logger.LogError("error cloning repo %s/%s: %s", organization, repo, err)

				logger.LogError("deleting repo url %s/%s since it does not exist, is private, too large, or is empty", organization, repo)
				err := db.DeleteRepoURL(context.Background(), repoUrl)
				if err != nil {
					logger.LogError("error deleting repo url %s: %s", repoUrl, err)
				}
				logger.LogError("repo url %s/%s deleted!", organization, repo)

				continue
			}
			logger.LogBlue("repository %s cloned!", repoUrl)
		}

		err = ProcessRepo(prefixPath, organization, repo, repoUrl, startDate, db)
		if err != nil {
			logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
		}
	}
}
