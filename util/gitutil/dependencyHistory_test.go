package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGitDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/dephistory-test-repo"
	dependencySearched := "chai"
	depFilePaths := []string{"package.json"}

	expectedDatesAddedReturn := []int64{1698773760}
	expectedDatesRemovedReturn := []int64{}

	// ARRANGE - TESTS
	tests := []struct {
		name                       string
		dependencySearched         string
		depFilePaths               []string
		expectedDatesAddedReturn   []int64
		expectedDatesRemovedReturn []int64
		wantErr                    bool
	}{
		{"Valid", dependencySearched, depFilePaths, expectedDatesAddedReturn, expectedDatesRemovedReturn, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			datesAdded, datesRemoved, err := GitDependencyHistory(repoDir, tt.dependencySearched, tt.depFilePaths)

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
