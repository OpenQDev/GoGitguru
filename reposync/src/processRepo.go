package reposync

import (
	"context"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

func ProcessRepo(prefixPath string, organization string, repo string, repoUrl string, startDate time.Time, db *database.Queries) error {
	logger.LogGreenDebug("beginning to process %s", repoUrl)

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusStoringCommits,
		Url:    repoUrl,
	})

	commitCount, err := StoreGitLogsForRepo(GitLogParams{prefixPath, organization, repo, repoUrl, startDate, db})
	if err != nil {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return err
	}

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSynced,
		Url:    repoUrl,
	})

	logger.LogBlue("Successfully stored %d commits for %s in the database.", commitCount, repoUrl)

	return nil
}
