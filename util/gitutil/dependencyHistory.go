package gitutil

import (
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitDependencyHistory(repoDir string, dependencySearched string, depFilePaths []string) ([]int64, []int64, error) {
	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, nil, err
	}
	ref, _ := r.Head()
	commits, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	datesPresentCommits := []int64{}
	datesRemovedCommits := []int64{}

	commits.ForEach(func(c *object.Commit) error {
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				contents, err := file.Contents()
				if err != nil {
					return err
				}
				if strings.Contains(contents, dependencySearched) {
					datesPresentCommits = append(datesPresentCommits, c.Committer.When.Unix())
					break
				} else {
					datesRemovedCommits = append(datesRemovedCommits, c.Committer.When.Unix())
					break
				}
			}
		}
		return nil
	})

	sort.Slice(datesPresentCommits, func(i, j int) bool { return datesPresentCommits[i] < datesPresentCommits[j] })
	sort.Slice(datesRemovedCommits, func(i, j int) bool { return datesRemovedCommits[i] < datesRemovedCommits[j] })

	return datesPresentCommits, datesRemovedCommits, nil
}
