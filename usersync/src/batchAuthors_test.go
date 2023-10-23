package usersync

import (
	"reflect"
	"testing"
	"util/testhelpers"
)

func TestGenerateBatchAuthors(t *testing.T) {
	tests := GenerateBatchAuthorsTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			result := generateBatchAuthors(tt.input, tt.batchSize)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("batchAuthors() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
