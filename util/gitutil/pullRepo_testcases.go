package gitutil

type PullRepoTestCase struct {
	name         string
	repo         string
	organization string
	wantErr      bool
}

func validPull() PullRepoTestCase {
	const VALID = "VALID"
	return PullRepoTestCase{
		name:         VALID,
		repo:         "OpenQ-Workflows",
		organization: "OpenQDev",
		wantErr:      false,
	}
}

func invalidPull() PullRepoTestCase {
	const INVALID = "INVALID"
	return PullRepoTestCase{
		name:         INVALID,
		repo:         "invalid-repo",
		organization: "valid-org",
		wantErr:      true,
	}
}

func PullRepoTestCases() []PullRepoTestCase {
	return []PullRepoTestCase{
		validPull(),
		invalidPull(),
	}
}
