package gitutil

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitDependencyHistory(repoDir string, dependenciesSearched []string, depFilePaths []string, dependencies []database.Dependency) (map[int32]DependencyResult, object.CommitIter, error) {
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

	dependenciesResults := make(map[int32]DependencyResult)

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

				// Convert both contents and dependencySearched to lowercase for case-insensitive comparison
				dependenciesResults = checkForDependenciesInFile(contents, dependenciesSearched, c, dependenciesResults, depFilePath, dependencies)
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
				dependenciesResults = checkForDependenciesInFile(contents, dependenciesSearched, c, dependenciesResults, depFilePath, dependencies)
			}
		}
	}

	fmt.Println("assemble arrays for", repoDir)

	return dependenciesResults, commits, err
}

type DependencyResult struct {
	DateFirstPresent int64
	DateLastRemoved  int64
}

func checkForDependenciesInFile(contents string, dependenciesSearched []string, c *object.Commit, currentDependenciesResult map[int32]DependencyResult, filePath string, dependencies []database.Dependency) map[int32]DependencyResult {
	contentsLower := strings.ToLower(contents)
	results := make(map[int32]DependencyResult)
	for _, dependencySearched := range dependenciesSearched {
		dependency := findDependency(dependencies, filePath, dependencySearched)
		dependencySearchedLower := strings.ToLower(dependencySearched)
		currentDateFirstPresent := currentDependenciesResult[dependency.InternalID].DateFirstPresent
		currentDateLastRemoved := currentDependenciesResult[dependency.InternalID].DateLastRemoved
		if strings.Contains(contentsLower, dependencySearchedLower) {
			currentDateFirstPresent = c.Committer.When.Unix()
		} else {
			if currentDateFirstPresent != 0 {
				currentDateLastRemoved = c.Committer.When.Unix()
			}
		}
		dependencyResult := DependencyResult{
			DateFirstPresent: currentDateFirstPresent,
			DateLastRemoved:  currentDateLastRemoved,
		}
		results[dependency.InternalID] = dependencyResult
	}
	return results
}

func getCommitWindow(lenCommitList int) int {
	return int(math.Max(1, math.Floor(float64(lenCommitList)*0.05)))
}
func findDependency(dependencies []database.Dependency, filename, dependencyName string) *database.Dependency {
	for _, dep := range dependencies {
		if dep.DependencyFile == filename && dep.DependencyName == dependencyName {
			return &dep // Return a pointer to the matching dependency
		}
	}
	return nil // Return nil if no matching dependency is found
}
