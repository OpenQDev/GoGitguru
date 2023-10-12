package gitutil

import (
	"testing"
)

func TestExtractOrganizationAndRepositoryFromUrl(t *testing.T) {
	// ARRANGE - TESTS

	tests := ExtractOrganizationAndRepositoryFromUrlTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, repo := ExtractOrganizationAndRepositoryFromUrl(tt.url)
			if org != tt.org || repo != tt.repo {
				t.Errorf("got %v %v, want %v %v", org, repo, tt.org, tt.repo)
			}
		})
	}
}
