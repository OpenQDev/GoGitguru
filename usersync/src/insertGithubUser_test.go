package usersync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/setup"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestInsertGithubUser(t *testing.T) {
	// ARRANGE - GLOBAL
	tests := InsertGithubUserTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH
			mock, queries := setup.GetMockDatabase()

			// ARRANGE - LOCAL
			tt.setupMock(mock, tt.author)

			// ACT
			err := insertGithubUser(tt.author, queries)

			// ASSERT
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
