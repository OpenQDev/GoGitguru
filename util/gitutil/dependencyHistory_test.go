package gitutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestGitDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/openq-coinapi"
	dependencySearched := "redis"
	depFilePaths := []string{"package.json"}
	expectedDatesAddedReturn := []time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}
	expectedDatesRemovedReturn := []time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}

	// ARRANGE - TESTS
	tests := []struct {
		name                       string
		dependencySearched         string
		depFilePaths               []string
		expectedDatesAddedReturn   []time.Time
		expectedDatesRemovedReturn []time.Time
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

			if !reflect.DeepEqual(datesAdded, tt.expectedDatesAddedReturn) {
				t.Errorf("GitDependencyHistory() datesAdded = %v, want %v", datesAdded, tt.expectedDatesAddedReturn)
			}

			if !reflect.DeepEqual(datesRemoved, tt.expectedDatesRemovedReturn) {
				t.Errorf("GitDependencyHistory() datesRemoved = %v, want %v", datesRemoved, tt.expectedDatesRemovedReturn)
			}
		})
	}
}
