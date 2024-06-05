package reposync

import (
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type AddRowToDependencyHistoryObjectTestCase struct {
	name                  string
	currentDependency     string
	currentDependencyFile string
	firstUseDate          int64
	expectedResult        database.BatchInsertRepoDependenciesParams
}

func validAddRowToDependencyHistoryObjectTest() AddRowToDependencyHistoryObjectTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodAddRowToDependencyHistoryObjectTestCase := AddRowToDependencyHistoryObjectTestCase{
		name:                  VALID_GIT_LOGS,
		currentDependency:     "eslint",
		currentDependencyFile: "package.json",
		firstUseDate:          1620000000,
		expectedResult: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{0},
		}}

	return goodAddRowToDependencyHistoryObjectTestCase
}

func AddRowToDependencyHistoryObjectTestCases() []AddRowToDependencyHistoryObjectTestCase {
	return []AddRowToDependencyHistoryObjectTestCase{
		validAddRowToDependencyHistoryObjectTest(),
	}

}

type SetDateFirstUsedTestCase struct {
	name                    string
	dependencyHistoryObject database.BatchInsertRepoDependenciesParams
	dependencySavedIndex    int
	commit                  object.Commit
	expectedResult          database.BatchInsertRepoDependenciesParams
}

func validSetDateFirstUsedTest() SetDateFirstUsedTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodSetDateFirstUsedTestCase := SetDateFirstUsedTestCase{
		name: VALID_GIT_LOGS,
		dependencyHistoryObject: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{0},
		},
		dependencySavedIndex: 0,
		commit: object.Commit{
			Committer: object.Signature{
				When: time.Unix(1420000000, 0),
			},
		},
		expectedResult: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1420000000},
			Lastusedates:    []int64{0},
		}}

	return goodSetDateFirstUsedTestCase
}

func SetDateFirstUsedTestCases() []SetDateFirstUsedTestCase {
	return []SetDateFirstUsedTestCase{
		validSetDateFirstUsedTest(),
	}

}

type SetDateRemovedTestCase struct {
	name                                    string
	dependencyHistoryObject                 database.BatchInsertRepoDependenciesParams
	dependenciesThatDoExistCurrentlyIndexes []int
	commit                                  object.Commit
	expectedResult                          database.BatchInsertRepoDependenciesParams
}

func validSetDateRemovedTest() SetDateRemovedTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodSetDateRemovedTestCase := SetDateRemovedTestCase{
		name: VALID_GIT_LOGS,
		dependencyHistoryObject: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{4},
		},
		dependenciesThatDoExistCurrentlyIndexes: []int{},
		commit: object.Commit{
			Committer: object.Signature{
				When: time.Unix(1420000000, 0),
			},
		},
		expectedResult: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{1420000000},
		}}

	return goodSetDateRemovedTestCase
}

func SetDateRemovedTestCases() []SetDateRemovedTestCase {
	return []SetDateRemovedTestCase{
		validSetDateRemovedTest(),
	}

}

type GetPreviousDependenciesInfoTestCase struct {
	name                                    string
	dependencyHistoryObject                 database.BatchInsertRepoDependenciesParams
	dependencyName                          string
	dependencyFileName                      string
	commit                                  object.Commit
	dependencySavedIndex                    int
	dependenciesThatDoExistCurrentlyIndexes []int
	expectedResult                          database.BatchInsertRepoDependenciesParams
}

func validGetPreviousDependenciesInfoTest() GetPreviousDependenciesInfoTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodGetPreviousDependenciesInfoTestCase := GetPreviousDependenciesInfoTestCase{
		name: VALID_GIT_LOGS,
		dependencyHistoryObject: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{4},
		},
		commit: object.Commit{
			Committer: object.Signature{
				When: time.Unix(1420000000, 0),
			},
		},
		dependencyName:                          "bslint",
		dependencyFileName:                      "package.json",
		dependencySavedIndex:                    -1,
		dependenciesThatDoExistCurrentlyIndexes: []int{},
		expectedResult: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{1420000000},
		}}

	return goodGetPreviousDependenciesInfoTestCase
}

func otherGetPreviousDependenciesInfoTest() GetPreviousDependenciesInfoTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodGetPreviousDependenciesInfoTestCase := GetPreviousDependenciesInfoTestCase{
		name: VALID_GIT_LOGS,
		dependencyHistoryObject: database.BatchInsertRepoDependenciesParams{
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{4},
		},
		commit: object.Commit{
			Committer: object.Signature{
				When: time.Unix(1420000000, 0),
			},
		},
		dependencyName:                          "eslint",
		dependencyFileName:                      "package.json",
		dependencySavedIndex:                    0,
		dependenciesThatDoExistCurrentlyIndexes: []int{0}}

	return goodGetPreviousDependenciesInfoTestCase
}

func GetPreviousDependenciesInfoTestCases() []GetPreviousDependenciesInfoTestCase {
	return []GetPreviousDependenciesInfoTestCase{
		validGetPreviousDependenciesInfoTest(),
		otherGetPreviousDependenciesInfoTest(),
	}

}
