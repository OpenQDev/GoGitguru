package sync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGetRepoToAuthorsMap(t *testing.T) {
	tests := GetRepoToAuthorsMapTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoToAuthorsMap(tt.input)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("getRepoToAuthorsMap() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
