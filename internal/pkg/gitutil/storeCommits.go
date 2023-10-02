package gitutil

import (
	"context"
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func StoreCommits(gitLogs []GitLog, repoUrl string, db *database.Queries) (*database.InsertCommitParams, error) {
	for _, gitLog := range gitLogs {
		params := database.InsertCommitParams{
			CommitHash:    gitLog.CommitHash,
			Author:        sql.NullString{String: gitLog.AuthorName, Valid: gitLog.AuthorName != ""},
			AuthorEmail:   sql.NullString{String: gitLog.AuthorEmail, Valid: gitLog.AuthorEmail != ""},
			AuthorDate:    sql.NullInt64{Int64: gitLog.AuthorDate, Valid: gitLog.AuthorDate != 0},
			CommitterDate: sql.NullInt64{Int64: gitLog.CommitDate, Valid: gitLog.CommitDate != 0},
			Message:       sql.NullString{String: gitLog.CommitMessage, Valid: gitLog.CommitMessage != ""},
			Insertions:    sql.NullInt32{Int32: 0, Valid: true},
			Deletions:     sql.NullInt32{Int32: 0, Valid: true},
			FilesChanged:  sql.NullInt32{Int32: 0, Valid: true},
			RepoUrl:       sql.NullString{String: repoUrl, Valid: repoUrl != ""},
		}

		err := db.InsertCommit(context.Background(), params)

		if err != nil {
			return &params, err
		}
	}
	logger.LogBlue("Successfully stored commits in the database.")

	return nil, nil
}
