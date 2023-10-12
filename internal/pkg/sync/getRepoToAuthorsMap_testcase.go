package sync

type GetRepoToAuthorsMapTestCase struct {
	name           string
	input          []UserSync
	expectedOutput map[string][]AuthorCommitTuple
}

func fooo() GetRepoToAuthorsMapTestCase {
	const FOO = "FOO"
	return GetRepoToAuthorsMapTestCase{
		name: FOO,
		input: []UserSync{
			{
				CommitHash: "abc123",
				Author: struct {
					Email   string
					NotNull bool
				}{
					Email:   "test@example.com",
					NotNull: true,
				},
				Repo: struct {
					URL     string
					NotNull bool
				}{
					URL:     "https://github.com/example/repo",
					NotNull: true,
				},
			},
			{
				CommitHash: "otherCommitHash",
				Author: struct {
					Email   string
					NotNull bool
				}{
					Email:   "otherperson@example.com",
					NotNull: true,
				},
				Repo: struct {
					URL     string
					NotNull bool
				}{
					URL:     "https://github.com/example/other-repo",
					NotNull: true,
				},
			},
		},
		expectedOutput: map[string][]AuthorCommitTuple{
			"https://github.com/example/repo":  []AuthorCommitTuple{AuthorCommitTuple{Author: "test@example.com", CommitHash: "abc123"}},
			"https://github.com/example/repo2": []AuthorCommitTuple{AuthorCommitTuple{Author: "otherperson@example.com", CommitHash: "otherCommitHash"}},
		},
	}
}

func GetRepoToAuthorsMapTestCases() []GetRepoToAuthorsMapTestCase {
	return []GetRepoToAuthorsMapTestCase{
		fooo(),
	}
}
