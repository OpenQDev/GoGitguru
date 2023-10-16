package sync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGenerateBatchAuthors(t *testing.T) {
	tests := GenerateBatchAuthorsTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"SINGLE_AUTHOR_ONE_REPO",
			), tt.name)

			result := GenerateBatchAuthors(tt.input, tt.batchSize)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("batchAuthors() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
