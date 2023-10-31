package gitutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitDependencyHistory(repoDir string, dependencySearched string, depFilePaths []string) ([]time.Time, []time.Time, error) {
	r, _ := git.PlainOpen(repoDir)
	ref, _ := r.Head()
	commits, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	datesAddedCommits := []time.Time{}
	datesRemovedCommits := []time.Time{}

	commits.ForEach(func(c *object.Commit) error {
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				contents, err := file.Contents()
				if err != nil {
					return err
				}
				if strings.Contains(contents, dependencySearched) {
					fmt.Println("true")
					fmt.Println(c.Committer.When)
					datesAddedCommits = append(datesAddedCommits, c.Committer.When)
				} else {
					fmt.Println("false")
					fmt.Println(c.Committer.When)
					datesRemovedCommits = append(datesRemovedCommits, c.Committer.When)
				}
			}
		}
		return nil
	})

	return datesAddedCommits, datesRemovedCommits, nil
}
