package gitutil

type IsGitRepositoryTestCase struct {
	name         string
	repo         string
	organization string
	want         bool
}

func validRepo() IsGitRepositoryTestCase {
	return IsGitRepositoryTestCase{
		name:         "Valid Repository",
		repo:         "OpenQ-Workflows",
		organization: "OpenQDev",
		want:         true,
	}
}

func invalidRepo() IsGitRepositoryTestCase {
	return IsGitRepositoryTestCase{
		name:         "Invalid Repository",
		repo:         "invalid-repo",
		organization: "valid-org",
		want:         false,
	}
}

func IsGitRepositoryTestCases() []IsGitRepositoryTestCase {
	return []IsGitRepositoryTestCase{
		validRepo(),
		invalidRepo(),
	}
}
