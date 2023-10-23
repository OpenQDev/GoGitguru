package usersync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGetRepoToAuthorsMap(t *testing.T) {
	tests := GetRepoToAuthorsMapTestCases()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.title)

			result := getRepoToAuthorsMap(tt.input)

			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("getRepoToAuthorsMap() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
