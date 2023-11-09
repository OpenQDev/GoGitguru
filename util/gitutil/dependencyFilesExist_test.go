package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestDependencyFileExists(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/openq-coinapi"
	dependencyFileExists := []string{"package.json"}
	dependencyFileNotExists := []string{"go.mod"}

	// ARRANGE - TESTS
	tests := []struct {
		name             string
		dependencyFiles  []string
		depFilesReturned []string
		wantErr          bool
	}{
		{"Valid", dependencyFileExists, []string{"package.json", "utils/package.json"}, false},
		{"Invalid", dependencyFileNotExists, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"Invalid",
			), tt.name)

			depFilesReturned, err := GitDependencyFiles(repoDir, tt.dependencyFiles)

			if (err != nil) != tt.wantErr {
				t.Errorf("DependencyFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.ElementsMatch(t, depFilesReturned, tt.depFilesReturned) {
				t.Errorf("DependencyFileExists() = %v, want %v", depFilesReturned, tt.depFilesReturned)
			}
		})
	}
}
