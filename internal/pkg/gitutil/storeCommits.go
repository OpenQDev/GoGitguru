package gitutil

import (
	"context"
	"database/sql"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
)

func StoreCommits(gitLogs []GitLog, repoUrl string, db *database.Queries) (*database.Commit, error) {
	fmt.Println(gitLogs)
	for _, gitLog := range gitLogs {
		params := database.InsertCommitParams{
			CommitHash:    gitLog.CommitHash,
			Author:        sql.NullString{String: gitLog.AuthorName, Valid: gitLog.AuthorName != ""},
			AuthorEmail:   sql.NullString{String: gitLog.AuthorEmail, Valid: gitLog.AuthorEmail != ""},
			AuthorDate:    sql.NullInt64{Int64: gitLog.AuthorDate, Valid: gitLog.AuthorDate != 0},
			CommitterDate: sql.NullInt64{Int64: gitLog.CommitDate, Valid: gitLog.CommitDate != 0},
			Message:       sql.NullString{String: gitLog.CommitMessage, Valid: gitLog.CommitMessage != ""},
			Insertions:    sql.NullInt32{Int32: int32(gitLog.Insertions), Valid: true},
			Deletions:     sql.NullInt32{Int32: int32(gitLog.Deletions), Valid: true},
			FilesChanged:  sql.NullInt32{Int32: int32(gitLog.FilesChanged), Valid: true},
			RepoUrl:       sql.NullString{String: repoUrl, Valid: repoUrl != ""},
		}

		commit, err := db.InsertCommit(context.Background(), params)

		if err != nil {
			return &commit, err
		}
	}

	logger.LogBlue("Successfully stored %d commits in the database.", len(gitLogs))

	return nil, nil
}
