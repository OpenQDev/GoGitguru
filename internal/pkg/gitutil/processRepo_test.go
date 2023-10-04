package gitutil

import (
	"main/internal/database"
	"main/internal/pkg/logger"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProcessRepo(t *testing.T) {
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

	tmpDir, err := os.MkdirTemp("", "prefixPath")
	if err != nil {
		logger.LogFatalRedAndExit("can't create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	prefixPath := tmpDir

	// Define test cases
	tests := []struct {
		name         string
		organization string
		repo         string
		repoUrl      string
		gitLogs      []GitLog
	}{
		{
			name:         "Valid git logs",
			organization: "OpenQDev",
			repo:         "OpenQ-DRM-TestRepo",
			repoUrl:      "https://github.com/OpenQ-Dev/OpenQ-DRM-TestRepo",
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
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone repo to tmp dir. Will be deleted at end of test
			CloneRepo(prefixPath, tt.organization, tt.repo)

			// Order expected database calls here

			// SET REPO STATUS TO storing_commits
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("storing_commits", tt.repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			// INSERT FIRST COMMIT
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

			// SET REPO STATUS TO synced
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("synced", tt.repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			// Call the ProcessRepo function
			ProcessRepo(prefixPath, tt.repo, tt.repoUrl, apiCfg.DB)
			if err != nil {
				t.Errorf("there was an error processing repo: %s", err)
			}

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// Check if ProcessRepo returned an error
			assert.Nil(t, err)
		})
	}
}
