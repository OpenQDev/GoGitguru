package gitutil

import (
	"fmt"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestDependencyFileExists(t *testing.T) {
	// ARRANGE - GLOBAL
	repoDir := "./mock/openqdev/openq-coinapi"
	dependencyFileExists := "package.json"
	dependencyFileNotExists := "go.mod"

	// ARRANGE - TESTS
	tests := []struct {
		name             string
		dependencyFile   string
		depFilesReturned []string
		wantErr          bool
	}{
		{"Valid", dependencyFileExists, []string{"package.json", "utils/package.json"}, false},
		{"Invalid", dependencyFileNotExists, []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			depFilesReturned, err := GitDependencyFiles(repoDir, tt.dependencyFile)

			if (err != nil) != tt.wantErr {
				fmt.Println("in here")
				t.Errorf("DependencyFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.ElementsMatch(t, depFilesReturned, tt.depFilesReturned) {
				t.Errorf("DependencyFileExists() = %v, want %v", depFilesReturned, tt.depFilesReturned)
			}
		})
	}
}
