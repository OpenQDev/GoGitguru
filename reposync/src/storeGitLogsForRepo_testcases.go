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
				CommitHash:    "09442fceb096a56226fb528368ddf971e776057f",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699383601,
				CommitDate:    1699383601,
				CommitMessage: "Initial commit",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			}, {
				CommitHash:    "a7ce99317e5347735ec5349f303c7036cd007d94",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699383684,
				CommitDate:    1699383684,
				CommitMessage: "Create package.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			}, {
				CommitHash:    "9141d952c3b15d1ad8121527f1f4bfb65f9000c0",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699384512,
				CommitDate:    1699384512,
				CommitMessage: "Create BigFile.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			}, {
				CommitHash:    "a8b0336d4e05acfa79d46beb2442c56c0fb23017",
				AuthorName:    "DRM-Test-User",
				AuthorEmail:   "150183211+DRM-Test-User@users.noreply.github.com",
				AuthorDate:    1699384731,
				CommitDate:    1699384731,
				CommitMessage: "Create BigFile2.json",
				FilesChanged:  0,
				Insertions:    0,
				Deletions:     0,
			}, {
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
			}, {
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
		shouldError: false,
		setupMock: func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string) {
			mock.ExpectQuery("^-- name: GetRepoDependenciesByURL :many.*").WithArgs(repoUrl).WillReturnRows(sqlmock.NewRows([]string{"url"}))

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
			files := []string{"package.json",
				"go.mod"}
			dependencies := []string{"web3", "("}

			firstUseDates := []int64{1699383684, 1699385002}
			lastUseDates := []int64{0, 0}
			rows := sqlmock.NewRows([]string{"id", "pattern", "updated_at", "creator"})
			dependencyFiles := []string{
				"package.json",
				"requirements.txt",
				"pom.xml",
				"Pipfile",
				"go.mod",
				"build.gradle",
				"Gemfile",
				"Cargo.toml",
				".cabal",
				"composer.json",

				"hardhat.config",
				"truffle",
				`\/network\/`,
				`\/deployments\/`,
				"foundry.toml",
			}
			for index, file := range dependencyFiles {
				rows.AddRow(index, file, 1609459200, "GoGitguru")
			}
			mock.ExpectQuery("^-- name: GetAllFilePatterns :many.*").WillReturnRows(rows)
			mock.ExpectQuery("^-- name: GetGithubUserByCommitEmail :many.*").WithArgs(pq.Array([]string{"150183211+DRM-Test-User@users.noreply.github.com", "info@openq.dev"})).WillReturnRows(sqlmock.NewRows([]string{"internal_id", "emails"}))
			mock.ExpectExec("^-- name: UpsertRepoToUserById :exec.*").WithArgs("https://github.com/OpenQDev/OpenQ-DRM-TestRepo", pq.Array([]string{}), pq.Array([]string{}), "{}").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("^-- name: BatchInsertRepoDependencies :exec.*").WithArgs(1609459200, pq.Array(files), pq.Array(dependencies), repoUrl, pq.Array(firstUseDates), pq.Array(lastUseDates)).WillReturnResult(sqlmock.NewResult(1, 1))

			// BULK INSERT COMMITS
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
		},
	}

	return sucessfulGitLogTestCase
}

func StoreGitLogsAndDepsHistoryForRepoTestCases() []StoreGitLogsAndDepsHistoryForRepoTestCase {
	return []StoreGitLogsAndDepsHistoryForRepoTestCase{
		sucessfulGitLog(),
	}
}
