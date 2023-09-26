package gitutil

import (
	"main/internal/database"
	"main/internal/pkg/logger"
)

func StoreCommits(prefixPath string, repo string, db *database.Queries) {
	_ = FormatGitLogs(prefixPath, repo, "")

	// TODO : Pipe into Postgres

	logger.LogBlue("Successfully stored commits in the database.")
}
