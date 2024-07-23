package usersync

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/OpenQDev/GoGitguru/database"
)

func TestPrepareUserDependencies(t *testing.T) {

	usersDependenciesToSync := []database.GetUserDependenciesByUpdatedAtRow{
		{
			UserID:       sql.NullInt32{Int32: 1, Valid: true},
			DependencyID: 1,
			FirstUseDate: sql.NullInt64{Int64: 1231231231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 1231231231203, Valid: true},
		}, {
			UserID:       sql.NullInt32{Int32: 1, Valid: true},
			DependencyID: 2,
			FirstUseDate: sql.NullInt64{Int64: 12331231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 23412312399, Valid: true},
		},
	}

	alreadySyncedUserDependencies := []database.GetUserDependenciesByUserRow{
		{
			UserID:       1,
			DependencyID: 1,
			FirstUseDate: sql.NullInt64{Int64: 34231231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 12312312399, Valid: true},
		}, {
			UserID:       1,
			DependencyID: 2,
			FirstUseDate: sql.NullInt64{Int64: 34231231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 12312312399, Valid: true},
		},
	}
	PrepareUserDependencies(usersDependenciesToSync, alreadySyncedUserDependencies)

	expectedOutput := []database.GetUserDependenciesByUpdatedAtRow{
		{
			UserID:       sql.NullInt32{Int32: 1, Valid: true},
			DependencyID: 1,
			FirstUseDate: sql.NullInt64{Int64: 1231231231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 1231231231203, Valid: true},
		},
		{
			UserID:       sql.NullInt32{Int32: 1, Valid: true},
			DependencyID: 2,
			FirstUseDate: sql.NullInt64{Int64: 12331231203, Valid: true},
			LastUseDate:  sql.NullInt64{Int64: 23412312399, Valid: true},
		},
	}
	if !reflect.DeepEqual(usersDependenciesToSync, expectedOutput) {
		t.Errorf("Expected %v, got %v", expectedOutput, usersDependenciesToSync)
	}
}
