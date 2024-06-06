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
				CommitHash:    "09442fceb096a56226fb528368ddf971e776057f",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699383601,
				CommitDate:    1699383601,
				CommitMessage: "Initial commit",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "a7ce99317e5347735ec5349f303c7036cd007d94",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699383684,
				CommitDate:    1699383684,
				CommitMessage: "Create package.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "9141d952c3b15d1ad8121527f1f4bfb65f9000c0",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699384512,
				CommitDate:    1699384512,
				CommitMessage: "Create BigFile.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "a8b0336d4e05acfa79d46beb2442c56c0fb23017",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699384731,
				CommitDate:    1699384731,
				CommitMessage: "Create BigFile2.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "0cc787dfb7f6a5808c54b5654e7bf871f004b890",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699385002,
				CommitDate:    1699385002,
				CommitMessage: "Add files via upload",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "70488e2cc8ef84edaab39aafda542b0ac2cee092",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699385069,
				CommitDate:    1699385069,
				CommitMessage: "Add files via upload",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "a4b132ba0fac0380bc7479730a4216218c39b716",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699385139,
				CommitDate:    1699385139,
				CommitMessage: "Add files via upload",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
			{
				CommitHash:    "32f8b288406652840a600e18d562a51661d64d99",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "info@openq.dev",
				AuthorDate:    1699386034,
				CommitDate:    1699386034,
				CommitMessage: "updates",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			},
		},
		setupMock: func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string) {

			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs(database.RepoStatusSyncingRepo, repoUrl).WillReturnResult(sqlmock.NewResult(1, 2))
			mock.ExpectQuery("^-- name: GetRepoDependenciesByURL :many.*").WithArgs(repoUrl).WillReturnRows(sqlmock.NewRows([]string{"url"}))
			// Define test data
			commitCount := 8
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

			files := []string{"package.json",
				"reposync/go.mod",
				"reposync/go.mod",
				"reposync/go.mod",
				"server/go.mod",
				"server/go.mod",
				"server/go.mod",
				"usersync/go.mod",
				"usersync/go.mod",
				"usersync/go.mod"}
			dependencies := []string{"web3", "base", "ton", "ergo", "base", "ton", "ergo", "base", "ton", "ergo"}
			firstUseDates := []int64{1699383684, 1699385002, 1699385002, 1699385002, 1699385069, 1699385069, 1699385069, 1699385139, 1699385139, 1699385139}
			lastUseDates := []int64{1699386034, 1699386034, 1699386034, 1699386034, 1699386034, 1699386034, 1699386034, 1699386034, 1699386034, 1699386034}

			mock.ExpectExec("^-- name: BatchInsertRepoDependencies :exec.*").WithArgs(pq.Array(files), pq.Array(dependencies), repoUrl, pq.Array(firstUseDates), pq.Array(lastUseDates)).WillReturnResult(sqlmock.NewResult(1, 1))
			// Define expected SQL statement
			// go-sqlmock CANNOT accept slices as arguments. Must convert to pq.Array first as is done in databse.BulkInsertCommits
			mock.ExpectExec("^-- name: BulkInsertCommits :exec.*").WithArgs(
				pq.Array(commitHash),
				pq.Array(author),
				pq.Array(authorEmail),
				pq.Array(authorDate),
				pq.Array(committerDate),
				pq.Array(message),
				pq.Array(filesChanged),
				repoUrl,
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
