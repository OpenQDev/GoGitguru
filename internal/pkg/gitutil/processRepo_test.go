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
	// Initialize a new instance of mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	// Initialize queries with the mocked DB collection.
	queries := database.New(db)

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

			// SET REPO STATUS TO storing_commits
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("storing_commits", tt.repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			commitHash := make([]string, 0)
			author := make([]string, 0)
			authorEmail := make([]string, 0)
			authorDate := make([]int64, 0)
			committerDate := make([]int64, 0)
			message := make([]string, 0)
			insertions := make([]int32, 0)
			deletions := make([]int32, 0)
			filesChanged := make([]int32, 0)
			repoUrls := make([]string, 0)

			for _, commit := range tt.gitLogs {
				commitHash = append(commitHash, commit.CommitHash)
				author = append(author, commit.AuthorName)
				authorEmail = append(authorEmail, commit.AuthorEmail)
				authorDate = append(authorDate, commit.AuthorDate)
				committerDate = append(committerDate, commit.CommitDate)
				message = append(message, commit.CommitMessage)
				insertions = append(insertions, int32(commit.Insertions))
				deletions = append(deletions, int32(commit.Deletions))
				filesChanged = append(filesChanged, int32(commit.FilesChanged))
				repoUrls = append(repoUrls, tt.repoUrl)
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

			// SET REPO STATUS TO synced
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("synced", tt.repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			// Call the ProcessRepo function
			ProcessRepo(prefixPath, tt.repo, tt.repoUrl, queries)
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
