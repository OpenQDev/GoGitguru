package gitutil

import "time"

type GetCommitHistoryTestCase struct {
	name      string
	startDate time.Time
	wantErr   bool
}

func validGitCommitHistory() GetCommitHistoryTestCase {
	const VALID = "VALID"
	return GetCommitHistoryTestCase{
		name:      VALID,
		startDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		wantErr:   false,
	}
}

func invalidGitCommitHistory() GetCommitHistoryTestCase {
	const INVALID = "INVALID"
	return GetCommitHistoryTestCase{
		name:      INVALID,
		startDate: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC), // Future date
		wantErr:   true,
	}
}

func GetCommitHistoryTestCases() []GetCommitHistoryTestCase {
	return []GetCommitHistoryTestCase{
		validGitCommitHistory(),
		invalidGitCommitHistory(),
	}
}
