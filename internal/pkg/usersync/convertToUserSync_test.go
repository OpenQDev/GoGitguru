package usersync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestConvertToUserSync(t *testing.T) {
	tests := ConvertToUserSyncTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToUserSync(tt.input)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("convertToUserSync() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
