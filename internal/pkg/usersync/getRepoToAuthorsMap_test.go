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
				"REPO_TO_AUTHOR_MAP",
			), tt.title)

			result := GetRepoToAuthorsMap(tt.input)

			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("getRepoToAuthorsMap() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
