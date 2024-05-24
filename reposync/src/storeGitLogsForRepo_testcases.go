package reposync

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

type StoreGitLogsAndDepsHistoryForRepoTestCase struct {
	name           string
	repoUrl        string
	repo           string
	fromCommitDate time.Time
	gitLogs        []GitLog
	shouldError    bool
	setupMock      func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string)
}

func sucessfulGitLog() StoreGitLogsAndDepsHistoryForRepoTestCase {
	sucessfulGitLogTestCase := StoreGitLogsAndDepsHistoryForRepoTestCase{
		name:           "Valid git logs",
		repoUrl:        "https://github.com/OpenQDev/OpenQ-DRM-TestRepo",
		repo:           "OpenQ-DRM-TestRepo",
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
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string) {
			numberOfCommits := len(gitLogs)
			var (
				commitHash    = make([]string, numberOfCommits)
				author        = make([]string, numberOfCommits)
				authorEmail   = make([]string, numberOfCommits)
				authorDate    = make([]int64, numberOfCommits)
				committerDate = make([]int64, numberOfCommits)
				message       = make([]string, numberOfCommits)
				insertions    = make([]int32, numberOfCommits)
				deletions     = make([]int32, numberOfCommits)
				filesChanged  = make([]int32, numberOfCommits)
				repoUrls      = make([]string, numberOfCommits)
			)

			for i, commit := range gitLogs {
				commitHash[i] = commit.CommitHash
				author[i] = commit.AuthorName
				authorEmail[i] = commit.AuthorEmail
				authorDate[i] = commit.AuthorDate
				committerDate[i] = commit.CommitDate
				message[i] = commit.CommitMessage
				insertions[i] = int32(commit.Insertions)
				deletions[i] = int32(commit.Deletions)
				filesChanged[i] = int32(commit.FilesChanged)
				repoUrls[i] = repoUrl
			}

			// BULK INSERT COMMITS
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
		},
	}

	return sucessfulGitLogTestCase
}

func StoreGitLogsAndDepsHistoryForRepoTestCases() []StoreGitLogsAndDepsHistoryForRepoTestCase {
	return []StoreGitLogsAndDepsHistoryForRepoTestCase{
		sucessfulGitLog(),
	}
}
