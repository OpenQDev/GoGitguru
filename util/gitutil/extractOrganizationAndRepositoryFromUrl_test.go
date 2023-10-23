package gitutil

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestExtractOrganizationAndRepositoryFromUrl(t *testing.T) {
	// ARRANGE - TESTS

	tests := ExtractOrganizationAndRepositoryFromUrlTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			org, repo := ExtractOrganizationAndRepositoryFromUrl(tt.url)
			if org != tt.org || repo != tt.repo {
				t.Errorf("got %v %v, want %v %v", org, repo, tt.org, tt.repo)
			}
		})
	}
}
