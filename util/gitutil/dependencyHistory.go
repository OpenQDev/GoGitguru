package gitutil

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitDependencyHistory(repoDir string, dependencySearched string, depFilePaths []string) ([]object.Commit, error) {
	r, _ := git.PlainOpen(repoDir)
	ref, _ := r.Head()
	commits, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	relevantCommits := []object.Commit{}

	commits.ForEach(func(c *object.Commit) error {
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				fmt.Printf("found file %s\n", depFilePath)
				contents, err := file.Contents()
				if err != nil {
					return err
				}
				if strings.Contains(contents, dependencySearched) {
					fmt.Println("true")
					fmt.Println(c.Author.When)
				} else {
					fmt.Println("false")
					fmt.Println(c.Author.When)
				}
			}
		}
		return nil
	})

	return relevantCommits, nil
}
