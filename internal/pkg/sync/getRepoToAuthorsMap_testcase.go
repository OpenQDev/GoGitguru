package sync

type GetRepoToAuthorsMapTestCase struct {
	title          string
	input          []UserSync
	expectedOutput map[string][]AuthorCommitTuple
}

func reposToAuthorMap() GetRepoToAuthorsMapTestCase {
	const REPO_TO_AUTHOR_MAP = "REPO_TO_AUTHOR_MAP"
	return GetRepoToAuthorsMapTestCase{
		title: REPO_TO_AUTHOR_MAP,
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
					URL:     "https://github.com/example/repo2",
					NotNull: true,
				},
			},
		},
		expectedOutput: map[string][]AuthorCommitTuple{
			"https://github.com/example/repo":  {{Author: "test@example.com", CommitHash: "abc123"}},
			"https://github.com/example/repo2": {{Author: "otherperson@example.com", CommitHash: "otherCommitHash"}},
		},
	}
}

func GetRepoToAuthorsMapTestCases() []GetRepoToAuthorsMapTestCase {
	return []GetRepoToAuthorsMapTestCase{
		reposToAuthorMap(),
	}
}
