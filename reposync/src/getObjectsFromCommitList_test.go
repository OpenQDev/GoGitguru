package reposync

import (
	"fmt"
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGetObjectsFromCommitList(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := GetObjectsFromCommitListTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)
			mockTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			now = func() time.Time { return mockTime }
			// ACT
			bulkInsertDependencyParams, bulkInsertCommitParams, usersToRepoUrls, _, err := GetObjectsFromCommitList(tt.params, tt.commitList, tt.numberOfCommits, tt.currentDependencies)
			if err != nil {
				t.Errorf("there was an error storing this commit: %v - the error was: %s", bulkInsertCommitParams, err)
			}
			fmt.Println(usersToRepoUrls)
			assert.Equal(t, tt.bulkInsertCommitsParams, bulkInsertCommitParams)
			assert.Equal(t, tt.bulkInsertDependencyParams, bulkInsertDependencyParams)
		})
	}
}