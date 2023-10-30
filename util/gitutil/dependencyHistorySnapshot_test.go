package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGitDependencyHistorySnapshot(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/openq-coinapi"
	dependencySearched := "axios"
	depFilePaths := "'**package.json**' '**utils/package.json**'"
	expectedReturnListValid := []string{"package.json"}

	// ARRANGE - TESTS
	tests := []struct {
		name               string
		dependencySearched string
		depFilePaths       string
		expectedReturnList []string
		wantErr            bool
	}{
		{"Valid", dependencySearched, depFilePaths, expectedReturnListValid, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			out, err := GitDependencyHistorySnapshot(repoDir, tt.dependencySearched, tt.depFilePaths)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitDependencyHistorySnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.ElementsMatch(t, out, tt.expectedReturnList) {
				t.Errorf("GitDependencyHistorySnapshot() = %v, want %v", out, tt.expectedReturnList)
			}
		})
	}
}
