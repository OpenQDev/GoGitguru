package usersync

type GithubGetCommitAuthorsTest struct {
	name    string
	owner   string
	repo    string
	commits []string
	wantErr bool
}

func success() GithubGetCommitAuthorsTest {
	const VALID_QUERY = "VALID_QUERY"

	return GithubGetCommitAuthorsTest{
		name:    VALID_QUERY,
		owner:   "OpenQDev",
		repo:    "OpenQ-Workflows",
		commits: []string{"8799411585c826b577f632f1ef5c0415914267ed", "657bd8b7f7d83e8b842411cbf65666901d65431c"},
		wantErr: false,
	}
}

func GithubGetCommitAuthorsTestCases() []GithubGetCommitAuthorsTest {
	return []GithubGetCommitAuthorsTest{
		success(),
	}
}
