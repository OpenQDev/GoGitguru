package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestParseDependencyFiles(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := ParseFileTestCases()

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.fileName)

			// BEFORE EACH

			// ACT

			result := ParseFile(tt.file, tt.fileName)

			// ASSERT

			assert.ElementsMatch(t, tt.dependencies, result)

		})
	}
}
