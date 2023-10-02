package gitutil

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func ProcessRepo(prefixPath string, repo string, repoUrl string, db *database.Queries) {
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusStoringCommits,
		Url:    repoUrl,
	})

	err := StoreCommits(GetFormattedGitLogs(prefixPath, repo, ""), repoUrl, db)
	if err != nil {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		logger.LogFatalRedAndExit("Failed to insert commit: %s", err)
	}
}
