package sync

import (
	"main/internal/pkg/testhelpers"
	"testing"
)

func TestGenerateAuthorBatchGqlQuery(t *testing.T) {
	// ARRANGE - TESTS
	tests := GenerateAuthorBatchGqlQueryTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GenerateAuthorBatchGqlQuery(tt.organization, tt.repo, tt.authorList)
			if result != tt.expectedOutput {
				t.Errorf("generateAuthorBatchGqlQuery() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
