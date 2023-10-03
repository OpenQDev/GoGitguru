package gitutil

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func ProcessRepo(prefixPath string, repo string, repoUrl string, db *database.Queries) error {
	logger.LogGreenDebug("beginning to process %s", repoUrl)

	// Set repo status to storing_commits
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusStoringCommits,
		Url:    repoUrl,
	})

	commitCount, err := StoreGitLogs(prefixPath, repo, repoUrl, "", db)
	if err != nil {
		// Set repo status to failed
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return err
	}

	// Set repo status to synced
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSynced,
		Url:    repoUrl,
	})

	logger.LogBlue("Successfully stored %d commits for %s in the database.", commitCount, repoUrl)

	return nil
}
