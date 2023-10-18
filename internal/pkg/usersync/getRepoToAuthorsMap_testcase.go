package usersync

type GetRepoToAuthorsMapTestCase struct {
	title          string
	input          []UserSync
	expectedOutput RepoToAuthorCommitTuples
}

func reposToAuthorMap() GetRepoToAuthorsMapTestCase {
	const REPO_TO_AUTHOR_MAP = "REPO_TO_AUTHOR_MAP"
	return GetRepoToAuthorsMapTestCase{
		title: REPO_TO_AUTHOR_MAP,
		input: []UserSync{
			{
				CommitHash:  "abc123",
				AuthorEmail: "test@example.com",
				RepoUrl:     "https://github.com/example/repo",
			},
			{
				CommitHash:  "otherCommitHash",
				AuthorEmail: "otherperson@example.com",
				RepoUrl:     "https://github.com/example/repo2",
			},
		},
		expectedOutput: RepoToAuthorCommitTuples{
			Repos: map[string][]AuthorCommitTuple{
				"https://github.com/example/repo":  {{Author: "test@example.com", CommitHash: "abc123"}},
				"https://github.com/example/repo2": {{Author: "otherperson@example.com", CommitHash: "otherCommitHash"}},
			},
		},
	}
}

func GetRepoToAuthorsMapTestCases() []GetRepoToAuthorsMapTestCase {
	return []GetRepoToAuthorsMapTestCase{
		reposToAuthorMap(),
	}
}
