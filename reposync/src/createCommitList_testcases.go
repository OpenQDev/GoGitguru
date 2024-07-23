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
		expectedResult: []object.Commit{makeCommit("32f8b288406652840a600e18d562a51661d64d99", "DRM-Test-User", "info@openq.dev"),

			makeCommit("a4b132ba0fac0380bc7479730a4216218c39b716", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
			makeCommit("70488e2cc8ef84edaab39aafda542b0ac2cee092", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com"),
		}}

	return goodCreateCommitListTestCase
}

func CreateCommitListTestCases() []CreateCommitListTestCase {
	return []CreateCommitListTestCase{
		validCreateCommitListTest(),
	}

}
