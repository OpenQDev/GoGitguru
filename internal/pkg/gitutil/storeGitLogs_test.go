package gitutil

import (
	"errors"
	"main/internal/database"
	"main/internal/pkg/logger"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
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

	// Initialize a new instance of mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	// Initialize queries with the mocked DB collection.
	queries := database.New(db)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError {
				mock.ExpectExec("^-- name: BulkInsertCommits :exec.*").WillReturnError(errors.New("mock error"))
			} else {
				numberOfCommits := len(tt.gitLogs)
				var (
					commitHash    = make([]string, numberOfCommits)
					author        = make([]string, numberOfCommits)
					authorEmail   = make([]string, numberOfCommits)
					authorDate    = make([]int64, numberOfCommits)
					committerDate = make([]int64, numberOfCommits)
					message       = make([]string, numberOfCommits)
					insertions    = make([]int32, numberOfCommits)
					deletions     = make([]int32, numberOfCommits)
					filesChanged  = make([]int32, numberOfCommits)
					repoUrls      = make([]string, numberOfCommits)
				)

				for i, commit := range tt.gitLogs {
					commitHash[i] = commit.CommitHash
					author[i] = commit.AuthorName
					authorEmail[i] = commit.AuthorEmail
					authorDate[i] = commit.AuthorDate
					committerDate[i] = commit.CommitDate
					message[i] = commit.CommitMessage
					insertions[i] = int32(commit.Insertions)
					deletions[i] = int32(commit.Deletions)
					filesChanged[i] = int32(commit.FilesChanged)
					repoUrls[i] = tt.repoUrl
				}

				// BULK INSERT COMMITS
				mock.ExpectExec("^-- name: BulkInsertCommits :exec.*").WithArgs(
					commitHash,
					author,
					authorEmail,
					authorDate,
					committerDate,
					message,
					insertions,
					deletions,
					filesChanged,
					repoUrls,
				)
			}

			commit, err := StoreGitLogs(tempDir, "OpenQ-DRM-TestRepo", tt.repoUrl, "", queries)
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
