package gitutil

type StoreGitLogsTestCase struct {
	name        string
	repoUrl     string
	repo        string
	gitLogs     []GitLog
	shouldError bool
}

func sucessfulGitLog() StoreGitLogsTestCase {
	foo := StoreGitLogsTestCase{
		name:    "Valid git logs",
		repoUrl: "https://github.com/OpenQDev/OpenQ-DRM-TestRepo",
		repo:    "OpenQ-DRM-TestRepo",
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
		shouldError: false,
	}

	return foo
}

func GitLogTestCases() []StoreGitLogsTestCase {
	return []StoreGitLogsTestCase{
		sucessfulGitLog(),
	}
}
