package reposync

import (
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/testhelpers"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestStoreGitLogs(t *testing.T) {
	// ARRANGE - GLOBAL
	tempDir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	CloneRepo(tempDir, "OpenQDev", "OpenQ-DRM-TestRepo")

	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	queries := database.New(db)

	// ARRANGE - TESTS
	tests := GitLogTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock, tt.gitLogs, tt.repoUrl)

			commit, err := StoreGitLogs(tempDir, "OpenQ-DRM-TestRepo", tt.repoUrl, "", queries)
			if err != nil && tt.shouldError == false {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", commit, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if tt.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
