package usersync

import (
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
)

type StartUserDepsSyncingTestCase struct {
	name        string
	setupMock   func(mock sqlmock.Sqlmock)
	shouldError bool
}

func startUserDepsSyncingTest1() StartUserDepsSyncingTestCase {
	const SHOULD_STORE_USER_DEPS = "SHOULD_STORE_USER_DEPS"

	return StartUserDepsSyncingTestCase{
		name:        SHOULD_STORE_USER_DEPS,
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"first_commit_date", "last_commit_date", "user_id", "dependency_id"}).AddRow(1, 2, 1, 2)
			mock.ExpectQuery("^-- name: GetUserDependenciesByUpdatedAt :many.*").WithArgs(1609458600).WillReturnRows(rows)

			rows = sqlmock.NewRows([]string{
				"first_use_date", "last_use_date", "dependency_id", "user_id"}).AddRow(0, 1, 2, 3)
			mock.ExpectQuery("^-- name: GetUserDependenciesByUser :many.*").WithArgs(pq.Array([]int{1}), pq.Array([]int{2})).WillReturnRows(rows)

			column1 := []int{1}
			column2 := []int{2}
			column3 := []uint64{2}
			column4 := []uint64{1}

			mock.ExpectExec("^-- name: BulkInsertUserDependencies.").WithArgs(
				pq.Array(column1),
				pq.Array(column2),
				pq.Array(column3),
				pq.Array(column4),
				1609458600,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		},
	}

}

func StartUserDepsSyncingTestCases() []StartUserDepsSyncingTestCase {
	return []StartUserDepsSyncingTestCase{
		startUserDepsSyncingTest1(),
	}
}
