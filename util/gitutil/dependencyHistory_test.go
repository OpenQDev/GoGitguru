package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGitDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/dephistory-test-repo"
	repoDirChurned := "./mock/openqdev/dephistory-test-repo-churned"
	repoDirChurnedReturned := "./mock/openqdev/dephistory-test-repo-churned-returned"

	dependencySearched := "chai"
	depFilePaths := []string{"package.json"}

	expectedDatesAddedReturn := []int64{1698773760}
	expectedDatesRemovedReturn := []int64{}

	expectedDatesAddedReturnChurned := []int64{1698773760}
	expectedDatesRemovedReturnChurned := []int64{1699555209}

	expectedDatesAddedReturnChurnedReturned := []int64{1698773760, 1699983952}
	expectedDatesRemovedReturnChurnedReturned := []int64{1699555209}

	// ARRANGE - TESTS
	tests := []struct {
		name                       string
		dependencySearched         string
		depFilePaths               []string
		expectedDatesAddedReturn   []int64
		expectedDatesRemovedReturn []int64
		repoDir                    string
		wantErr                    bool
	}{
		{"Added after init, never removed", dependencySearched, depFilePaths, expectedDatesAddedReturn, expectedDatesRemovedReturn, repoDir, false},
		{"Added after init, then removed", dependencySearched, depFilePaths, expectedDatesAddedReturnChurned, expectedDatesRemovedReturnChurned, repoDirChurned, false},
		{"Added after init, then removed, then re-added", dependencySearched, depFilePaths, expectedDatesAddedReturnChurnedReturned, expectedDatesRemovedReturnChurnedReturned, repoDirChurnedReturned, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"Added after init, then removed",
			), tt.name)

			datesAdded, datesRemoved, err := GitDependencyHistory(tt.repoDir, tt.dependencySearched, tt.depFilePaths)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitDependencyHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.ElementsMatch(t, datesAdded, tt.expectedDatesAddedReturn) {
				t.Errorf("DependencyFileExists() = %v, want %v", datesAdded, tt.expectedDatesAddedReturn)
			}

			if !assert.ElementsMatch(t, datesRemoved, tt.expectedDatesRemovedReturn) {
				t.Errorf("DependencyFileExists() = %v, want %v", datesRemoved, tt.expectedDatesRemovedReturn)
			}
		})
	}
}
