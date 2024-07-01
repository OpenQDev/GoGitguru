package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestAddRowToDependencyHistoryObject(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := AddRowToDependencyHistoryObjectTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ARRANGE - LOCAL
			dependencyHistoryObjects := database.BatchInsertRepoDependenciesParams{}
			// ACT
			addRowToDependencyHistoryObject(&dependencyHistoryObjects, tt.currentDependency, tt.currentDependencyFile, tt.firstUseDate)

			assert.Equal(t, tt.expectedResult, dependencyHistoryObjects)

		})
	}
}

func TestSetDateFirstUsed(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := SetDateFirstUsedTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ACT
			setDateFirstUsed(&tt.dependencyHistoryObject, tt.dependencySavedIndex, tt.commit)

			assert.Equal(t, tt.expectedResult, tt.dependencyHistoryObject)

		})
	}
}
func TestSetDateRemoved(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := SetDateRemovedTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ACT
			setDateRemoved(&tt.dependencyHistoryObject, tt.dependenciesThatDoExistCurrentlyIndexes, tt.commit.Committer.When.Unix())

			assert.Equal(t, tt.expectedResult, tt.dependencyHistoryObject)

		})
	}
}

func TestGetPreviousDependenciesInfo(t *testing.T) {
	// BEFORE ALL

	// ARRANGE - TESTS
	tests := GetPreviousDependenciesInfoTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// ACT
			dependencySavedIndex, dependenciesThatDoExistCurrentlyIndexes := getPreviousDependenciesInfo(&tt.dependencyHistoryObject, tt.dependencyName, tt.dependencyFileName, tt.commit)

			assert.Equal(t, tt.dependencySavedIndex, dependencySavedIndex)
			assert.Equal(t, tt.dependenciesThatDoExistCurrentlyIndexes, dependenciesThatDoExistCurrentlyIndexes)

		})
	}
}
