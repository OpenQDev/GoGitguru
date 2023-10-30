package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestGitDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/openq-coinapi"
	dependencySearched := "redis"
	depFilePaths := []string{"package.json", "utils/package.json"}

	// ARRANGE - TESTS
	tests := []struct {
		name               string
		dependencySearched string
		depFilePaths       []string
		wantErr            bool
	}{
		{"Valid", dependencySearched, depFilePaths, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			_, err := GitDependencyHistory(repoDir, tt.dependencySearched, tt.depFilePaths)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitDependencyHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}