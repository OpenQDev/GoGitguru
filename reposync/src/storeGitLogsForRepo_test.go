package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGitLogsAndDepsHistoryForRepo(t *testing.T) {
	// ARRANGE - GLOBAL
	prefixPath := "../mock"
	repo := "OpenQ-DRM-TestRepo"

	// ARRANGE - TESTS
	tests := StoreGitLogsAndDepsHistoryForRepoTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()

			tt.setupMock(mock, tt.gitLogs, tt.repoUrl)

			commitCount, err := StoreGitLogsAndDepsHistoryForRepo(GitLogParams{prefixPath, organization, repo, tt.repoUrl, tt.fromCommitDate, queries})
			if err != nil && tt.shouldError == false {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", commitCount, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			require.Equal(t, 8, commitCount)

			if tt.shouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
