package gitutil

import (
	"os"
	"testing"
	"util/testhelpers"
)

func TestCloneRepo(t *testing.T) {
	// ARRANGE - GLOBAL
	prefixPath, _ := os.MkdirTemp("", "repos")

	// ARRANGE - TESTS
	tests := CloneRepoTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := CloneRepo(prefixPath, tt.organization, tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Defer deletion of the repo after each test
			defer DeleteLocalRepo(prefixPath, tt.organization, tt.repo)
		})
	}
}
