package reposync

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/OpenQDev/GoGitguru/database"
	"github.com/lib/pq"
)

type ProcessRepoTestCase struct {
	name           string
	organization   string
	repo           string
	repoUrl        string
	gitLogs        []GitLog
	fromCommitDate time.Time
	setupMock      func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string)
}

const organization = "openqdev"
const repo = "openq-drm-testrepo"

func validProcessRepoTest() ProcessRepoTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodProcessRepoTestCase := ProcessRepoTestCase{
		name:           VALID_GIT_LOGS,
		organization:   organization,
		repo:           repo,
		repoUrl:        "https://github.com/OpenQ-Dev/OpenQ-DRM-TestRepo",
		fromCommitDate: time.Unix(1696277204, 0),
		gitLogs: []GitLog{
			{
				CommitHash:    "06a12f9c203112a149707ff73e4298749744c358",
				AuthorName:    "FlacoJones",
				AuthorEmail:   "andrew@openq.dev",
				AuthorDate:    1696277247,
				CommitDate:    1696277247,
				CommitMessage: "updates README",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "9fae86bc8e89895b961d81bd7e9e4e897501c8bb",
				AuthorName:    "FlacoJones",
				AuthorEmail:   "andrew@openq.dev",
				AuthorDate:    1696277205,
				CommitDate:    1696277205,
				CommitMessage: "initial commit",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
		},
		setupMock: func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string) {
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs(database.RepoStatusSyncingRepo, repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

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
				commitHash[i] = gitLogs[i].CommitHash
				author[i] = gitLogs[i].AuthorName
				authorEmail[i] = gitLogs[i].AuthorEmail
				authorDate[i] = gitLogs[i].AuthorDate
				committerDate[i] = gitLogs[i].CommitDate
				message[i] = gitLogs[i].CommitMessage
				insertions[i] = int32(gitLogs[i].Insertions)
				deletions[i] = int32(gitLogs[i].Deletions)
				filesChanged[i] = int32(gitLogs[i].FilesChanged)
				repoUrls[i] = repoUrl
			}

			// Define expected SQL statement
			// go-sqlmock CANNOT accept slices as arguments. Must convert to pq.Array first as is done in databse.BulkInsertCommits
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

			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs(database.RepoStatusSynced, repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}

	return goodProcessRepoTestCase
}

func ProcessRepoTestCases() []ProcessRepoTestCase {
	return []ProcessRepoTestCase{
		validProcessRepoTest(),
	}

}
