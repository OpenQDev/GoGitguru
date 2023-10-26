package gitutil

import (
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestIsGitRepository(t *testing.T) {
	// ARRANGE - GLOBAL
	prefixPath, _ := os.MkdirTemp("", "repos")

	// ARRANGE - TESTS
	tests := IsGitRepositoryTestCases() // You will need to define this function

	for _, tt := range tests {
		// Ensure the repository is cloned before trying to check
		_ = CloneRepo(prefixPath, tt.organization, tt.repo)

		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			got := IsGitRepository(prefixPath, tt.organization, tt.repo)
			if got != tt.want {
				t.Errorf("IsGitRepository() = %v, want %v", got, tt.want)
			}
			// Defer deletion of the repo after each test
			defer DeleteLocalRepo(prefixPath, tt.organization, tt.repo)
		})
	}
}
