package reposync

import (
	"fmt"
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestAddCommitToCommitObject(t *testing.T) {
	commitObject := database.BulkInsertCommitsParams{
		Commithashes:   []string{},
		Authors:        []string{},
		Authoremails:   []string{},
		Authordates:    []int64{},
		Committerdates: []int64{},
		Messages:       []string{},
		Fileschanged:   []int32{},
	}
	hashString := "d3b07384d113edec49eaa6238ad5ff00"

	// Convert the string to a plumbing.Hash
	hash := plumbing.NewHash(hashString)

	commit := &object.Commit{
		Hash: hash,
		Author: object.Signature{
			Name:  "author",
			Email: "author@email.com",
			When:  time.Now(),
		},
	}
	_ = AddCommitToCommitObject(commit, &commitObject, 0)
	fmt.Println(commitObject.Authors[0], "authors")

	if commitObject.Commithashes[0] != hash.String() {
		t.Errorf("Expected %s, got %s", hashString, commitObject.Commithashes[0])
	}
	if commitObject.Authors[0] != "author" {
		t.Errorf("Expected author, got %s", commitObject.Authors[0])
	}
	if commitObject.Authoremails[0] != "author@email.com" {
		t.Errorf("Expected empty string, got %s", commitObject.Authoremails[0])
	}
}
