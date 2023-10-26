package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGitLogsForRepo(t *testing.T) {
	// ARRANGE - GLOBAL
	prefixPath := "mock"
	repo := "OpenQ-DRM-TestRepo"

	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	queries := database.New(db)

	// ARRANGE - TESTS
	tests := StoreGitLogsForRepoTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock, tt.gitLogs, tt.repoUrl)

			commitCount, err := StoreGitLogsForRepo(GitLogParams{prefixPath, organization, repo, tt.repoUrl, tt.fromCommitDate, queries})
			if err != nil && tt.shouldError == false {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", commitCount, err)
			}

			require.Equal(t, 2, commitCount)

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
