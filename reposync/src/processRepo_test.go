package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestProcessRepo(t *testing.T) {
	// BEFORE ALL
	prefixPath := "mock"

	// ARRANGE - TESTS
	tests := ProcessRepoTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()

			// ARRANGE - LOCAL
			tt.setupMock(mock, tt.gitLogs, tt.repoUrl)

			// ACT
			err := ProcessRepo(prefixPath, tt.organization, tt.repo, tt.repoUrl, tt.fromCommitDate, queries)
			if err != nil {
				t.Errorf("there was an error processing repo: %s", err)
			}

			// ASSERT

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			assert.Nil(t, err)
		})
	}
}
