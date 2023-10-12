package sync

type GenerateBatchAuthorsTestCase struct {
	name           string
	input          map[string][]AuthorCommitTuple
	batchSize      int
	expectedOutput [][]interface{}
}

func singleAuthorSingleRepo() GenerateBatchAuthorsTestCase {
	const SINGLE_AUTHOR_ONE_REPO = "SINGLE_AUTHOR_ONE_REPO"
	return GenerateBatchAuthorsTestCase{
		name: SINGLE_AUTHOR_ONE_REPO,
		input: map[string][]AuthorCommitTuple{
			"https://github.com/test/repo": {
				{"test@example.com", "commit1"},
				{"test2@example.com", "commit2"},
				{"test3@example.com", "commit3"},
			},
		},
		batchSize: 2,
		expectedOutput: [][]interface{}{
			{"https://github.com/test/repo", []AuthorCommitTuple{
				{"test@example.com", "commit1"},
				{"test2@example.com", "commit2"},
			}},
			{"https://github.com/test/repo", []AuthorCommitTuple{
				{"test3@example.com", "commit3"},
			}},
		},
	}
}

func singleAuthorTwoRepos() GenerateBatchAuthorsTestCase {
	const SINGLE_AUTHOR_TWO_REPOS = "SINGLE_AUTHOR_TWO_REPOS"
	return GenerateBatchAuthorsTestCase{
		name: SINGLE_AUTHOR_TWO_REPOS,
		input: map[string][]AuthorCommitTuple{
			"https://github.com/test/repo": {
				{"test@example.com", "commit1"},
				{"test2@example.com", "commit2"},
				{"test3@example.com", "commit3"},
			},
			"https://github.com/test/repo2": {
				{"author123@example.com", "commit4"},
				{"author12sdfdsf@example.com", "commit5"},
				{"authosdfsdf@example.com", "commit6"},
			},
		},
		batchSize: 2,
		expectedOutput: [][]interface{}{
			{"https://github.com/test/repo", []AuthorCommitTuple{
				{"test@example.com", "commit1"},
				{"test2@example.com", "commit2"},
			}},
			{"https://github.com/test/repo", []AuthorCommitTuple{
				{"test3@example.com", "commit3"},
			}},
			{"https://github.com/test/repo2", []AuthorCommitTuple{
				{"author123@example.com", "commit4"},
				{"author12sdfdsf@example.com", "commit5"},
			}},
			{"https://github.com/test/repo2", []AuthorCommitTuple{
				{"authosdfsdf@example.com", "commit6"},
			}},
		},
	}
}

func GenerateBatchAuthorsTestCases() []GenerateBatchAuthorsTestCase {
	return []GenerateBatchAuthorsTestCase{
		singleAuthorSingleRepo(),
		singleAuthorTwoRepos(),
	}
}
