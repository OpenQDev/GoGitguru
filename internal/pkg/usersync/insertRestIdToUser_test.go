package usersync

import (
	"main/internal/pkg/server/mocks"
	"main/internal/pkg/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertIntoRestIdToUser(t *testing.T) {
	// ARRANGE - GLOBAL
	mock, queries := mocks.GetMockDatabase()

	tests := InsertIntoRestIdToUserTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			tt.setupMock(mock, tt.author)

			// ACT
			err := insertIntoRestIdToUser(tt.author, queries)

			// ASSERT
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
