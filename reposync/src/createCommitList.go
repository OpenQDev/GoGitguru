package reposync

import (
	"slices"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateCommitList(repoDir string) ([]*object.Commit, error) {

	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, err
	}
	ref, _ := r.Head()
	commits, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	commitList := make([]*object.Commit, 0)
	err = commits.ForEach(func(c *object.Commit) error {
		commitList = append(commitList, c)
		return nil
	})

	// we want to look at commits from oldest  to newest
	slices.Reverse(commitList)
	return commitList, err
}
