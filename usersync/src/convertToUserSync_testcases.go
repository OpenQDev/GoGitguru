package usersync

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/database"
)

type ConvertToUserSyncTestCase struct {
	name           string
	input          []database.GetLatestUncheckedCommitPerAuthorRow
	expectedOutput []UserSync
}

func valid() ConvertToUserSyncTestCase {
	const SINGLE_AUTHOR_SINGLE_REPO = "SINGLE_AUTHOR_SINGLE_REPO"
	return ConvertToUserSyncTestCase{
		name: SINGLE_AUTHOR_SINGLE_REPO,
		input: []database.GetLatestUncheckedCommitPerAuthorRow{
			{
				CommitHash: "abc123",
				AuthorEmail: sql.NullString{
					String: "test@example.com",
					Valid:  true,
				},
				RepoUrl: sql.NullString{
					String: "https://github.com/test/repo",
					Valid:  true,
				},
			},
			{
				CommitHash: "abc123",
				AuthorEmail: sql.NullString{
					String: "",
					Valid:  false,
				},
				RepoUrl: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
		},
		expectedOutput: []UserSync{
			{
				CommitHash:  "abc123",
				AuthorEmail: "test@example.com",
				RepoUrl:     "https://github.com/test/repo",
			},
			{
				CommitHash:  "abc123",
				AuthorEmail: "",
				RepoUrl:     "",
			},
		},
	}
}

func ConvertToUserSyncTestCases() []ConvertToUserSyncTestCase {
	return []ConvertToUserSyncTestCase{
		valid(),
	}
}
