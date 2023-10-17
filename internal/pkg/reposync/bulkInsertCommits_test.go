package reposync

package reposync_test

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"main/internal/database"
	"main/internal/pkg/reposync"
)

func TestBulkInsertCommits(t *testing.T) {
	// Create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create new Queries instance with mock database
	q := database.New(db)

// Define test data
commitCount := 2
commitHash := make([]string, commitCount)
author := make([]string, commitCount)
authorEmail := make([]string, commitCount)
authorDate := make([]int64, commitCount)
committerDate := make([]int64, commitCount)
message := make([]string, commitCount)
insertions := make([]int32, commitCount)
deletions := make([]int32, commitCount)
filesChanged := make([]int32, commitCount)
repoUrls := make([]string, commitCount)

// Fill the arrays
for i := 0; i < commitCount; i++ {
	commitHash[i] = fmt.Sprintf("hash%d", i+1)
	author[i] = fmt.Sprintf("author%d", i+1)
	authorEmail[i] = fmt.Sprintf("email%d", i+1)
	authorDate[i] = int64(i + 1)
	committerDate[i] = int64(i + 3)
	message[i] = fmt.Sprintf("message%d", i+1)
	insertions[i] = int32(i + 5)
	deletions[i] = int32(i + 7)
	filesChanged[i] = int32(i + 9)
	repoUrls[i] = fmt.Sprintf("url%d", i+1)
}

	// Define expected SQL statement
	mock.ExpectExec("^-- name: BulkInsertCommits :exec.*").WithArgs(
		commitHash,
		author,
		authorEmail,
		authorDate,
		committerDate,
		message,
		insertions,
		deletions,
		filesChanged,
		repoUrls,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call function
	err = reposync.BulkInsertCommits(q, commitHash, author, authorEmail, authorDate, committerDate, message, insertions, deletions, filesChanged, repoUrls)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}