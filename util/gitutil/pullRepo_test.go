package gitutil

import (
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestPullRepo(t *testing.T) {
	// ARRANGE - GLOBAL
	prefixPath, _ := os.MkdirTemp("", "repos")

	// ARRANGE - TESTS
	tests := PullRepoTestCases() // Reuse the same test cases

	for _, tt := range tests {
		// Ensure the repository is cloned before trying to pull
		_ = CloneRepo(prefixPath, tt.organization, tt.repo)

		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := PullRepo(prefixPath, tt.organization, tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("PullRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Defer deletion of the repo after each test
			defer DeleteLocalRepo(prefixPath, tt.organization, tt.repo)
		})
	}
}
