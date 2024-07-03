package reposync

import (
	"slices"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommitForDependencies(c *object.Commit, repoDir string, dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, rawDependencyFiles []string) error {

	dependencyFileNamesWithPath, err := gitutil.GitDependencyFiles(repoDir, rawDependencyFiles)
	if err != nil {
		return err
	}
	// set date removed at start and unset later if necessary
	setDateRemoved(dependencyHistoryObject, c.Committer.When.Unix())

	for _, dependencyFileNameWithPath := range dependencyFileNamesWithPath {

		// set indexableDependencyFileName to the name corresponding to the dependencyFileNameWithPath
		indexableDependencyFileName := ""
		for _, indexableDependencyFileName = range rawDependencyFiles {
			if strings.Contains(dependencyFileNameWithPath, indexableDependencyFileName) {
				break
			}
		}

		file, err := c.File(dependencyFileNameWithPath)

		if err != nil {
			continue
		}
		if file == nil {
			continue
		}
		dependencies := ParseFile(file, indexableDependencyFileName)

		// only handle matching file name
		if slices.Contains(dependencyHistoryObject.Filenames, indexableDependencyFileName) {

			// iterate over dependencies that exist in this file and commit within this loop we are looking at actual individual dependencies
			for _, dependency := range dependencies {

				dependencySavedIndex := getPreviousDependenciesInfo(dependencyHistoryObject, dependency, indexableDependencyFileName, *c)

				// handle dependency that doesn't currently exit
				if dependencySavedIndex == -1 {
					addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, indexableDependencyFileName, c.Committer.When.Unix())
				} else {
					setDateFirstUsed(dependencyHistoryObject, dependencySavedIndex, *c)

					dependencyHistoryObject.Lastusedates[dependencySavedIndex] = 0
				}

			}

		} else {
			// handle new file
			for _, dependency := range dependencies {
				addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, indexableDependencyFileName, c.Committer.When.Unix())

			}
		}
	}

	return nil
}
