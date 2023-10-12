package gitutil

type ProcessRepoTestCase struct {
	name         string
	organization string
	repo         string
	repoUrl      string
	gitLogs      []GitLog
}

const organization = "OpenQDev"
const repo = "OpenQ-DRM-TestRepo"

func test1() ProcessRepoTestCase {
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
	}

	return goodProcessRepoTestCase
}

func ProcessRepoTestCases() []ProcessRepoTestCase {
	return []ProcessRepoTestCase{
		test1(),
	}

}
