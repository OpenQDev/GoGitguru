package reposync

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CreateCommitListTestCase struct {
	name           string
	organization   string
	repo           string
	expectedResult []object.Commit
}

func makeCommit(hashString string, author string, authorEmail string) object.Commit {
	hash := plumbing.NewHash(hashString)
	return object.Commit{
		Hash: hash,
		Author: object.Signature{
			Name:  author,
			Email: authorEmail,
		},
	}
}

func validCreateCommitListTest() CreateCommitListTestCase {
	const VALID_GIT_LOGS = "VALID_GIT_LOGS"

	goodCreateCommitListTestCase := CreateCommitListTestCase{
		name:         VALID_GIT_LOGS,
		organization: organization,
		repo:         repo,
		expectedResult: []object.Commit{
			makeCommit("09442fceb096a56226fb528368ddf971e776057f", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
			makeCommit("a7ce99317e5347735ec5349f303c7036cd007d94", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
			makeCommit("9141d952c3b15d1ad8121527f1f4bfb65f9000c0", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
		}}

	return goodCreateCommitListTestCase
}

func CreateCommitListTestCases() []CreateCommitListTestCase {
	return []CreateCommitListTestCase{
		validCreateCommitListTest(),
	}

}
