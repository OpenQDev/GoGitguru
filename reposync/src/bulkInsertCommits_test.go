package reposync

import (
	"fmt"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/setup"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestBulkInsertCommits(t *testing.T) {
	// Create new Queries instance with mock database
	mock, q := setup.GetMockDatabase()

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
	// go-sqlmock CANNOT accept slices as arguments.
	// Must convert to pq.Array first as is done in databse.BulkInsertCommits
	mock.ExpectExec("^-- name: BulkInsertCommits :exec.*").WithArgs(
		pq.Array(commitHash),
		pq.Array(author),
		pq.Array(authorEmail),
		pq.Array(authorDate),
		pq.Array(committerDate),
		pq.Array(message),
		pq.Array(insertions),
		pq.Array(deletions),
		pq.Array(filesChanged),
		pq.Array(repoUrls),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call function
	err := BulkInsertCommits(q, commitHash, author, authorEmail, authorDate, committerDate, message, insertions, deletions, filesChanged, repoUrls)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
