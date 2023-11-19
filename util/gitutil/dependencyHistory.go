package gitutil

import (
	"slices"
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

	commitList := make([]*object.Commit, 0)
	err = commits.ForEach(func(c *object.Commit) error {
		commitList = append(commitList, c)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	slices.Reverse(commitList)

	datesPresentCommits := []int64{}
	datesRemovedCommits := []int64{}

	for _, c := range commitList {
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				contents, err := file.Contents()
				if err != nil {
					return nil, nil, err
				}
				if strings.Contains(contents, dependencySearched) {
					datesPresentCommits = append(datesPresentCommits, c.Committer.When.Unix())
					break
				} else {
					if len(datesPresentCommits) != 0 {
						datesRemovedCommits = append(datesRemovedCommits, c.Committer.When.Unix())
					}
					break
				}
			}
		}
	}

	var presentArray []int64
	if len(datesPresentCommits) > 0 {
		// we only want to know when the dependency was first added
		var earliestPresent = slices.Min(datesPresentCommits)
		presentArray = []int64{earliestPresent}
	} else {
		presentArray = []int64{}
	}

	var removedArray []int64
	if len(datesRemovedCommits) > 0 {
		// we only consider a dependency "removed" if it is still removed as of "today"
		var latestAbsent = slices.Max(datesRemovedCommits)
		var latestPresent = slices.Max(datesPresentCommits)
		if latestAbsent > latestPresent {
			removedArray = []int64{slices.Max(datesRemovedCommits)}
		}
	} else {
		removedArray = []int64{}
	}

	return presentArray, removedArray, nil
}
