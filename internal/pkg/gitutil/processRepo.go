package gitutil

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func ProcessRepo(prefixPath string, repo string, repoUrl string, db *database.Queries) {
	// Set repo status to storing_commits
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusStoringCommits,
		Url:    repoUrl,
	})

	commit, err := StoreCommits(GetFormattedGitLogs(prefixPath, repo, ""), repoUrl, db)
	if err != nil {
		// Set repo status to failed
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})

		logger.LogFatalRedAndExit("for repository %s, failed to insert this commit: %s with the following error: ", repoUrl, commit, err)
	}
}
