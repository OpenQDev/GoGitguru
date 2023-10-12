package gitutil

type CloneRepoTestCase struct {
	name         string
	repo         string
	organization string
	wantErr      bool
}

func valid() CloneRepoTestCase {
	const VALID = "VALID"
	return CloneRepoTestCase{
		name:         VALID,
		repo:         "OpenQ-Workflows",
		organization: "OpenQDev",
		wantErr:      false,
	}
}

func invalid() CloneRepoTestCase {
	const INVALID = "INVALID"
	return CloneRepoTestCase{
		name:         INVALID,
		repo:         "invalid-repo",
		organization: "valid-org",
		wantErr:      true,
	}
}

func CloneRepoTestCases() []CloneRepoTestCase {
	return []CloneRepoTestCase{
		valid(),
		invalid(),
	}
}
