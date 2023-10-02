package gitutil

import (
	"errors"
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStoreCommits(t *testing.T) {
	// Initialize a new instance of ApiConfig with mocked DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	// Initialize queries with the mocked DB collection.
	queries := database.New(db)

	apiCfg := handlers.ApiConfig{
		DB: queries,
	}

	// Define test cases
	tests := []struct {
		name        string
		repoUrl     string
		gitLogs     []GitLog
		shouldError bool
	}{
		{
			name:    "Valid git logs",
			repoUrl: "https://github.com/OpenQDev/OpenQ-Workflows",
			gitLogs: []GitLog{
				{
					CommitHash:    "abc123",
					AuthorName:    "John Doe",
					AuthorEmail:   "john.doe@example.com",
					AuthorDate:    1633027200,
					CommitDate:    1633027200,
					CommitMessage: "Initial commit",
					FilesChanged:  3,
					Insertions:    42,
					Deletions:     0,
				},
				{
					CommitHash:    "def456",
					AuthorName:    "Jane Doe",
					AuthorEmail:   "jane.doe@example.com",
					AuthorDate:    1633113600,
					CommitDate:    1633113600,
					CommitMessage: "Add feature",
					FilesChanged:  1,
					Insertions:    10,
					Deletions:     2,
				},
			},
			shouldError: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expectations and actions for the mock DB can be defined here
			if tt.shouldError {
				mock.ExpectExec("^-- name: InsertCommit :one.*").WillReturnError(errors.New("mock error"))
			} else {
				for _, gitLog := range tt.gitLogs {
					linesChanged := gitLog.Insertions + gitLog.Deletions
					mock.ExpectQuery("^-- name: InsertCommit :one.*").WithArgs(
						gitLog.CommitHash,
						gitLog.AuthorName,
						gitLog.AuthorEmail,
						gitLog.AuthorDate,
						gitLog.CommitDate,
						gitLog.CommitMessage,
						gitLog.Insertions,
						gitLog.Deletions,
						gitLog.FilesChanged,
						tt.repoUrl,
					).WillReturnRows(sqlmock.NewRows([]string{
						"commit_hash",
						"author",
						"author_email",
						"author_date",
						"committer_date",
						"message",
						"insertions",
						"deletions",
						"lines_changed",
						"files_changed",
						"repo_url",
					}).AddRow(
						gitLog.CommitHash,
						gitLog.AuthorName,
						gitLog.AuthorEmail,
						gitLog.AuthorDate,
						gitLog.CommitDate,
						gitLog.CommitMessage,
						gitLog.Insertions,
						gitLog.Deletions,
						linesChanged,
						gitLog.FilesChanged,
						tt.repoUrl,
					))
				}
			}

			commit, err := StoreCommits(tt.gitLogs, tt.repoUrl, apiCfg.DB)
			if err != nil && tt.shouldError == false {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", commit, err)
			}

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// Check if StoreCommits returned an error
			if tt.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
