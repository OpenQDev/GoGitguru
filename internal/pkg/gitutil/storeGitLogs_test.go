package gitutil

import (
	"errors"
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/server"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStoreGitLogs(t *testing.T) {
	// Setup a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	CloneRepo(tempDir, "OpenQDev", "OpenQ-DRM-TestRepo")

	// Initialize a new instance of ApiConfig with mocked DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	// Initialize queries with the mocked DB collection.
	queries := database.New(db)

	apiCfg := server.ApiConfig{
		DB: queries,
	}

	// Define test cases
	tests := []struct {
		name        string
		repoUrl     string
		repo        string
		gitLogs     []GitLog
		shouldError bool
	}{
		{
			name:    "Valid git logs",
			repoUrl: "https://github.com/OpenQDev/OpenQ-DRM-TestRepo",
			repo:    "OpenQ-DRM-TestRepo",
			gitLogs: []GitLog{
				{
					CommitHash:    "06a12f9c203112a149707ff73e4298749744c358",
					AuthorName:    "FlacoJones",
					AuthorEmail:   "andrew@openq.dev",
					AuthorDate:    1696277247,
					CommitDate:    1696277247,
					CommitMessage: "updates README",
					FilesChanged:  1,
					Insertions:    1,
					Deletions:     0,
				},
				{
					CommitHash:    "9fae86bc8e89895b961d81bd7e9e4e897501c8bb",
					AuthorName:    "FlacoJones",
					AuthorEmail:   "andrew@openq.dev",
					AuthorDate:    1696277205,
					CommitDate:    1696277205,
					CommitMessage: "initial commit",
					FilesChanged:  0,
					Insertions:    0,
					Deletions:     0,
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

			commit, err := StoreGitLogs(tempDir, "OpenQ-DRM-TestRepo", tt.repoUrl, "", apiCfg.DB)
			if err != nil && tt.shouldError == false {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", commit, err)
			}

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// Check if GetGitLogs returned an error
			if tt.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
