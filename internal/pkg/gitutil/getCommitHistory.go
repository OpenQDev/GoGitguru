package gitutil

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetCommitHistory(r *git.Repository, startDate time.Time) (object.CommitIter, error) {
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}
	return r.Log(&git.LogOptions{From: ref.Hash(), Since: &startDate})
}
