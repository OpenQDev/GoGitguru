package gitutil

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitDependencyHistory(repoDir string, dependenciesSearched []string, depFilePaths []string) (map[string][]int64, map[string][]int64, error) {
	fmt.Println(time.Now().Format(time.RFC3339), "running for", repoDir)
	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, nil, err
	}
	ref, _ := r.Head()
	commits, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	fmt.Println("creating commit list for", repoDir)
	commitList := make([]*object.Commit, 0)
	err = commits.ForEach(func(c *object.Commit) error {
		commitList = append(commitList, c)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("reversing commit list for", repoDir)
	slices.Reverse(commitList)

	datesPresentCommits := make(map[string][]int64)
	datesRemovedCommits := make(map[string][]int64)

	fmt.Println("range over ", len(commitList), "commits", repoDir)
	commitWindow := getCommitWindow(len(commitList))
	for i := 0; i < len(commitList); i += commitWindow {
		c := commitList[i]
		fmt.Printf("Commit number %d: %s\n", i, c.Hash)
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				contents, err := file.Contents()
				if err != nil {
					return nil, nil, err
				}
for _, dependencySearched := range dependenciesSearched {
	dependencyDatePresentCommits := datesPresentCommits[dependencySearched]
	dependencyDateRemovedCommits := datesRemovedCommits[dependencySearched]
	datesPresentCommits[dependencySearched], datesRemovedCommits[dependencySearched]= checkForDependencyInFile(contents, dependencySearched, dependencyDatePresentCommits, c, dependencyDateRemovedCommits)
}
				// Convert both contents and dependencySearched to lowercase for case-insensitive comparison
				
			}
		}
	}

	if len(commitList)%commitWindow != 0 {
		c := commitList[len(commitList)-1]
		fmt.Printf("Commit number %d: %s\n", len(commitList)-1, c.Hash)
		for _, depFilePath := range depFilePaths {
			if file, err := c.File(depFilePath); err == nil {
				contents, err := file.Contents()
				if err != nil {
					return nil, nil, err
				}

				// Convert both contents and dependencySearched to lowercase for case-insensitive comparison
				for _, dependencySearched := range dependenciesSearched {
					dependencyDatePresentCommits := datesPresentCommits[dependencySearched]
					dependencyDateRemovedCommits := datesRemovedCommits[dependencySearched]
					dependencyDatePresentCommits, dependencyDateRemovedCommits= checkForDependencyInFile(contents, dependencySearched, dependencyDatePresentCommits, c, dependencyDateRemovedCommits)
				}
			}
		}
	}

	fmt.Println("assemble arrays for", repoDir)

	dependencyPresentArray := make(map[string][]int64)
	dependencyRemovedArray := make(map[string][]int64)
	for _, dependencySearched := range dependenciesSearched {

		 presentArray :=dependencyPresentArray[dependencySearched]
		if len(datesPresentCommits) > 0 {
			// we only want to know when the dependency was first added
			var earliestPresent = slices.Min(datesPresentCommits[dependencySearched])
			presentArray = []int64{earliestPresent}
		} else {
			presentArray = []int64{}
		}
		dependencyPresentArray[dependencySearched] = presentArray
	
		 removedArray := dependencyRemovedArray[dependencySearched]
		if len(datesRemovedCommits) > 0 {
			// we only consider a dependency "removed" if it is still removed as of "today"
			var latestAbsent = slices.Max(datesRemovedCommits[dependencySearched])
			var latestPresent = slices.Max(datesPresentCommits[dependencySearched])
			if latestAbsent > latestPresent {
				removedArray = []int64{slices.Max(datesRemovedCommits[dependencySearched])}
			}
		} else {
			removedArray = []int64{}
		}
		dependencyRemovedArray[dependencySearched] = removedArray
	
	}

	fmt.Println(time.Now().Format(time.RFC3339), "RAN FOR", repoDir)
	return dependencyPresentArray,dependencyRemovedArray, nil
}

func checkForDependencyInFile(contents string, dependencySearched string, datesPresentCommits []int64, c *object.Commit, datesRemovedCommits []int64) ([]int64, []int64) {
	contentsLower := strings.ToLower(contents)
	dependencySearchedLower := strings.ToLower(dependencySearched)

	if strings.Contains(contentsLower, dependencySearchedLower) {
		datesPresentCommits = append(datesPresentCommits, c.Committer.When.Unix())
	} else {
		if len(datesPresentCommits) != 0 {
			datesRemovedCommits = append(datesRemovedCommits, c.Committer.When.Unix())
		}
	}
	return datesPresentCommits, datesRemovedCommits
}

func getCommitWindow(lenCommitList int) int {
	return int(math.Max(1, math.Floor(float64(lenCommitList)*0.05)))
}
