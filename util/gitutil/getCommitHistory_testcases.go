package gitutil

import "time"

type GetCommitHistoryTestCase struct {
	name                string
	startDate           time.Time
	expectedCommitCount int
	wantErr             bool
}

func validGitCommitHistory() GetCommitHistoryTestCase {
	const VALID = "VALID"
	return GetCommitHistoryTestCase{
		name:                VALID,
		startDate:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		expectedCommitCount: 2,
		wantErr:             false,
	}
}

func validGitCommitHistory_onlyOne() GetCommitHistoryTestCase {
	const VALID_ONE_COMMIT = "VALID_ONE_COMMIT"
	return GetCommitHistoryTestCase{
		name:                VALID_ONE_COMMIT,
		startDate:           time.Date(2023, 10, 2, 15, 6, 48, 0, time.FixedZone("", -5*60*60)),
		expectedCommitCount: 1,
		wantErr:             false,
	}
}

func invalidGitCommitHistory() GetCommitHistoryTestCase {
	const INVALID = "INVALID"
	return GetCommitHistoryTestCase{
		name:                INVALID,
		startDate:           time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC), // Future date
		expectedCommitCount: 0,
		wantErr:             true,
	}
}

func GetCommitHistoryTestCases() []GetCommitHistoryTestCase {
	return []GetCommitHistoryTestCase{
		validGitCommitHistory(),
		invalidGitCommitHistory(),
		validGitCommitHistory_onlyOne(),
	}
}
