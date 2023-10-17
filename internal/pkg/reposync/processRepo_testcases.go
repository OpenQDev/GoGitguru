package reposync

import "github.com/DATA-DOG/go-sqlmock"

type ProcessRepoTestCase struct {
	name         string
	organization string
	repo         string
	repoUrl      string
	gitLogs      []GitLog
	setupMock    func(mock sqlmock.Sqlmock, gitLogs []GitLog, repoUrl string)
}

const organization = "OpenQDev"
const repo = "OpenQ-DRM-TestRepo"

func validProcessRepoTest() ProcessRepoTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodProcessRepoTestCase := ProcessRepoTestCase{
		name:         VALID_GIT_LOGS,
		organization: organization,
		repo:         repo,
		repoUrl:      "https://github.com/OpenQ-Dev/OpenQ-DRM-TestRepo",
		gitLogs: []GitLog{
			{
				CommitHash:    "06a12f9c203112a149707ff73e4298749744c358",
				AuthorName:    "FlacoJones",
				AuthorEmail:   "andrew@openq.dev",
				AuthorDate:    1696277247,
				CommitDate:    1696277247,
				CommitMessage: "updates README",
				FilesChanged:  1,
				Insertions:    1,
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
			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("storing_commits", repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			commitHash := make([]string, 0)
			author := make([]string, 0)
			authorEmail := make([]string, 0)
			authorDate := make([]int64, 0)
			committerDate := make([]int64, 0)
			message := make([]string, 0)
			insertions := make([]int32, 0)
			deletions := make([]int32, 0)
			filesChanged := make([]int32, 0)
			repoUrls := make([]string, 0)

			for _, commit := range gitLogs {
				commitHash = append(commitHash, commit.CommitHash)
				author = append(author, commit.AuthorName)
				authorEmail = append(authorEmail, commit.AuthorEmail)
				authorDate = append(authorDate, commit.AuthorDate)
				committerDate = append(committerDate, commit.CommitDate)
				message = append(message, commit.CommitMessage)
				insertions = append(insertions, int32(commit.Insertions))
				deletions = append(deletions, int32(commit.Deletions))
				filesChanged = append(filesChanged, int32(commit.FilesChanged))
				repoUrls = append(repoUrls, repoUrl)
			}

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
			)

			mock.ExpectExec("^-- name: UpdateStatusAndUpdatedAt :exec.*").WithArgs("synced", repoUrl).WillReturnResult(sqlmock.NewResult(1, 1))
		},
	}

	return goodProcessRepoTestCase
}

func ProcessRepoTestCases() []ProcessRepoTestCase {
	return []ProcessRepoTestCase{
		validProcessRepoTest(),
	}

}
