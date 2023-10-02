package gitutil

import (
	"main/internal/database"
	"main/internal/pkg/handlers"
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

	apiCfg := handlers.ApiConfig{
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
	}{
		{
			name:         "Valid git logs",
			organization: "OpenQDev",
			repo:         "OpenQ-Workflows",
			repoUrl:      "https://github.com/OpenQ-Dev/OpenQ-Workflows",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone repo to tmp dir. Will be deleted at end of test
			CloneRepo(prefixPath, "OpenQDev", "OpenQ-Workflows")

			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs(
				"storing_commits",
				tt.repoUrl,
			).WillReturnResult(sqlmock.NewResult(1, 1))

			// Call the ProcessRepo function
			ProcessRepo(prefixPath, tt.repo, tt.repoUrl, apiCfg.DB)
			if err != nil {
				t.Errorf("there was an error processing repo: %s", err)
			}

			// Check if there were any unexpected calls to the mock DB
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// Check if StoreCommits returned an error
			assert.Nil(t, err)
		})
	}
}
