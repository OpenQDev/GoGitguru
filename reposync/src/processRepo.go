package reposync

import (
	"context"
	"database/database"
	"util/logger"
)

func ProcessRepo(prefixPath string, organization string, repo string, repoUrl string, db *database.Queries) error {
	logger.LogGreenDebug("beginning to process %s", repoUrl)

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusStoringCommits,
		Url:    repoUrl,
	})

	commitCount, err := StoreGitLogsForRepo(GitLogParams{prefixPath, organization, repo, repoUrl, "", db})
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