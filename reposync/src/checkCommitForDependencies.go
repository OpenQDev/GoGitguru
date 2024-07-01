package reposync

import (
	"slices"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommitForDependencies(c *object.Commit, repoDir string, dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, rawDependencyFiles []string) error {

	dependencyFiles, err := gitutil.GitDependencyFiles(repoDir, rawDependencyFiles)

	if err != nil {
		return err
	}

	for _, dependencyFileName := range dependencyFiles {

		// set rawDependencyFileName to the name corresponding to the dependencyFileName
		rawDependencyFileName := ""
		for _, rawDependencyFileName = range rawDependencyFiles {
			if strings.Contains(dependencyFileName, rawDependencyFileName) {
				break
			}
		}

		currentCommitDate := c.Committer.When.Unix()
		file, err := c.File(dependencyFileName)
		if err != nil {
			continue
		}
		if file == nil {
			continue
		}

		dependencies := ParseFile(file)

		// only handle matching file name
		if slices.Contains(dependencyHistoryObject.Filenames, rawDependencyFileName) {
			// should only be package.json
			dependenciesThatDoExistCurrentlyIndexes := []int{}
			// iterate over dependencies that exist in this file and commit within this loop we are looking at actual individual dependencies
			for _, dependency := range dependencies {

				dependencySavedIndex, dependenciesThatDoExistCurrentlyIndexesResult := getPreviousDependenciesInfo(dependencyHistoryObject, dependency, rawDependencyFileName, *c)
				dependenciesThatDoExistCurrentlyIndexes = append(dependenciesThatDoExistCurrentlyIndexes, dependenciesThatDoExistCurrentlyIndexesResult...)

				// handle dependency that doesn't currently exit
				if dependencySavedIndex == -1 {
					addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, rawDependencyFileName, c.Committer.When.Unix())
				} else {
					setDateFirstUsed(dependencyHistoryObject, dependencySavedIndex, *c)
				}

			}

			setDateRemoved(dependencyHistoryObject, dependenciesThatDoExistCurrentlyIndexes, currentCommitDate)

		} else {
			for _, dependency := range dependencies {
				addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, rawDependencyFileName, c.Committer.When.Unix())

			}
		}
	}

	return nil
}
