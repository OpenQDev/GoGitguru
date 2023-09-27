package gitutil

import (
	"context"
	"database/sql"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func StoreCommits(gitLogs []GitLog, repoUrl string, db *database.Queries) {
	for _, gitLog := range gitLogs {
		fmt.Println(gitLog)
		err := db.InsertCommit(context.Background(), database.InsertCommitParams{
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
		})

		if err != nil {
			logger.LogFatalRedAndExit("Failed to insert commit: %s", err)
		}
	}
	logger.LogBlue("Successfully stored commits in the database.")
}
