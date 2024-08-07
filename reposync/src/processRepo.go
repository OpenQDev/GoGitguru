package reposync

import (
	"context"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

func ProcessRepo(prefixPath string, organization string, repo string, repoUrl string, startDate time.Time, db *database.Queries, resyncAll bool) error {
	logger.LogGreenDebug("beginning to process %s", repoUrl)

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSyncingRepo,
		Url:    repoUrl,
	})

	commitCount, err := StoreGitLogsAndDepsHistoryForRepo(GitLogParams{prefixPath, organization, repo, repoUrl, startDate, db}, resyncAll)
	if err != nil {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return err
	}

	if commitCount == 0 {
		logger.LogBlue("no new commits in repo %s", repoUrl)
	}
	if resyncAll {
		db.UpdateStatusAndUpdatedAtepoUrlV2(context.Background(), database.UpdateStatusAndUpdatedAtepoUrlV2Params{
			Status: database.RepoStatusSynced,
			Url:    repoUrl,
		})
	} else {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusSynced,
			Url:    repoUrl,
		})
	}
	logger.LogBlue("Successfully stored %d commits for %s in the database.", commitCount, repoUrl)

	return nil
}
