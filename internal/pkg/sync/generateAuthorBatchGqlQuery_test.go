package sync

import (
	"main/internal/pkg/testhelpers"
	"strings"
	"testing"
)

func TestGenerateAuthorBatchGqlQuery(t *testing.T) {
	// ARRANGE - TESTS
	tests := GenerateAuthorBatchGqlQueryTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			"SINGLE_AUTHOR",
		), tt.title)

		t.Run(tt.title, func(t *testing.T) {
			result := GenerateAuthorBatchGqlQuery(tt.organization, tt.repo, tt.authorList)

			sanitizedResult := sanitizeString(result)
			sanitizedExpectedOutput := sanitizeString(tt.expectedOutput)

			if sanitizedResult != sanitizedExpectedOutput {
				t.Errorf("generateAuthorBatchGqlQuery() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

func sanitizeString(str string) string {
	noSpaces := strings.ReplaceAll(str, " ", "")
	noNewLines := strings.ReplaceAll(noSpaces, "\n", "")
	noTabs := strings.ReplaceAll(noNewLines, "\t", "")
	return noTabs
}
