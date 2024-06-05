package reposync

import (
	"database/sql"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GetObjectsFromCommitListTestCase struct {
	name                       string
	params                     GitLogParams
	commitList                 []*object.Commit
	numberOfCommits            int
	currentDependencies        []database.GetRepoDependenciesByURLRow
	bulkInsertCommitsParams    database.BulkInsertCommitsParams
	bulkInsertDependencyParams database.BatchInsertRepoDependenciesParams
}

func makeCommitByReference(hashString string, author string, authorEmail string) *object.Commit {
	hash := plumbing.NewHash(hashString)
	return &object.Commit{
		Hash: hash,
		Author: object.Signature{
			Name:  author,
			Email: authorEmail,
			When:  time.Unix(12312381240, 0),
		},
		Committer: object.Signature{
			When: time.Unix(12318351240, 0),
		},
	}
}

func validGetObjectsFromCommitListTest() GetObjectsFromCommitListTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"
	currentDependency := database.GetRepoDependenciesByURLRow{
		DependencyName: "eslint",
		DependencyFile: "package.json",
		FirstUseDate:   sql.NullInt64{Int64: 1620000000, Valid: true},
		LastUseDate:    sql.NullInt64{Int64: 0, Valid: true},
	}
	goodGetObjectsFromCommitListTestCase := GetObjectsFromCommitListTestCase{
		name: VALID_GIT_LOGS,
		commitList: []*object.Commit{
			makeCommitByReference("09442fceb096a56226fb528368ddf971e776057f", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
			makeCommitByReference("a7ce99317e5347735ec5349f303c7036cd007d94", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
			makeCommitByReference("9141d952c3b15d1ad8121527f1f4bfb65f9000c0", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
		},
		numberOfCommits: 2,
		currentDependencies: []database.GetRepoDependenciesByURLRow{
			currentDependency,
		},
		bulkInsertDependencyParams: database.BatchInsertRepoDependenciesParams{
			Url:             "",
			Filenames:       []string{"package.json"},
			Dependencynames: []string{"eslint"},
			Firstusedates:   []int64{1620000000},
			Lastusedates:    []int64{0},
		},
		bulkInsertCommitsParams: database.BulkInsertCommitsParams{
			Commithashes: []string{
				"09442fceb096a56226fb528368ddf971e776057f",
				"a7ce99317e5347735ec5349f303c7036cd007d94",
			},
			Authors: []string{
				"DRM-Test-User",
				"DRM-Test-User",
			},
			Authoremails: []string{
				"150183211+DRM-Test-User@users.noreply.github.com",
				"150183211+DRM-Test-User@users.noreply.github.com",
			},
			Authordates:    []int64{12312381240, 12312381240},
			Committerdates: []int64{12318351240, 12318351240},
			Messages:       []string{"", ""},
			Fileschanged:   []int32{0, 0},
			Repourl: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
	}

	return goodGetObjectsFromCommitListTestCase
}

func GetObjectsFromCommitListTestCases() []GetObjectsFromCommitListTestCase {
	return []GetObjectsFromCommitListTestCase{
		validGetObjectsFromCommitListTest(),
	}

}
