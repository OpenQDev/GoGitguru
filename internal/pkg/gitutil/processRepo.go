package gitutil

import (
	"context"
	"fmt"
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

	gitLogs, err := GetFormattedGitLogs(prefixPath, repo, "")
	if err != nil {
		// Set repo status to failed
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return err
	}

	commit, err := StoreCommits(gitLogs, repoUrl, db)
	if err != nil {
		// Set repo status to failed
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})

		return fmt.Errorf("failed to insert this commit: %v with the following error: %s", commit, err)
	}

	// Set repo status to synced
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSynced,
		Url:    repoUrl,
	})

	return nil
}
