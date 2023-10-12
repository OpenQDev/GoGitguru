package sync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGenerateBatchAuthors(t *testing.T) {
	tests := GenerateBatchAuthorsTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBatchAuthors(tt.input, tt.batchSize)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("batchAuthors() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
