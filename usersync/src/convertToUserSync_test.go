package usersync

import (
	"reflect"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestConvertToUserSync(t *testing.T) {
	tests := ConvertToUserSyncTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			result := convertDatabaseObjectToUserSync(tt.input)

			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("convertToUserSync() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
