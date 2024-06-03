package reposync

import (
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type DependencyToAdd struct {
	DependencyName string
	DependencyFile string
	DateFirstUsed  int64
	DateLastUsed   int64
}

func CheckCommitForDependencies(c *object.Commit, dependencyHistoryObject *database.BatchInsertRepoDependenciesParams) error {
	dependencyFiles := []string{
		"package.json",
		"requirements.txt",
		"pom.xml",
		"Pipfile",
		"go.mod",
		"build.gradle",
		"Gemfile",
		"Cargo.toml",
		".cabal", // Adjusted from "cabal" to ".cabal" to match the correct file name
		"composer.json",
	}

	for _, dependencyFileName := range dependencyFiles {
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
		if slices.Contains(dependencyHistoryObject.Filenames, dependencyFileName) {
			// should only be package.json
			dependenciesThatDoExistCurrentlyIndexes := []int{}

			// iterate over dependencies that exist in this file and commit within this loop we are looking at actual individual dependencies
			for _, dependency := range dependencies {

				dependencySavedIndex := -1
				for savedDependencyNameIndex, savedDependencyName := range dependencyHistoryObject.Dependencynames {
					// only filename at matching index
					for savedDependencyFileNameIndex, savedDependencyFileName := range dependencyHistoryObject.Filenames {
						if savedDependencyFileName == dependencyFileName {
							if savedDependencyNameIndex == savedDependencyFileNameIndex {
								if savedDependencyName == dependency {
									// set for outside of loop
									dependenciesThatDoExistCurrentlyIndexes = append(dependenciesThatDoExistCurrentlyIndexes, savedDependencyNameIndex)
									dependencySavedIndex = savedDependencyNameIndex
									// if it does exist also clear lastRemovedDate
									if c.Committer.When.Unix() > dependencyHistoryObject.Lastusedates[savedDependencyNameIndex] {

										dependencyHistoryObject.Lastusedates[savedDependencyNameIndex] = 0
									}

									break
								}
							}
						}
					}
				}

				// dependencySavedIndex is the index of the dependency <> filename in the list of dependencies and -1 if it doesn't exist as that
				//handle dependency never existed before - if never existed before add firstUseDate of right now and lastRemovedDate of null
				if dependencySavedIndex == -1 {

					dependencyHistoryObject.Filenames = append(dependencyHistoryObject.Filenames, dependencyFileName)
					dependencyHistoryObject.Dependencynames = append(dependencyHistoryObject.Dependencynames, dependency)
					dependencyHistoryObject.Firstusedates = append(dependencyHistoryObject.Firstusedates, c.Committer.When.Unix())
					dependencyHistoryObject.Lastusedates = append(dependencyHistoryObject.Lastusedates, 0)

					//  if existed before update the value at that index where applicable
				} else {
					currentDateFirstPresent := dependencyHistoryObject.Firstusedates[dependencySavedIndex]
					if currentDateFirstPresent > c.Committer.When.Unix() || currentDateFirstPresent == 0 {
						dependencyHistoryObject.Firstusedates[dependencySavedIndex] = c.Committer.When.Unix()
					}
				}

			}
			// iterate through saved dependencies and if they haven't been added. set lastRemovedDate to current commit
			for savedDependencyIndex := range dependencyHistoryObject.Dependencynames {
				if !slices.Contains(dependenciesThatDoExistCurrentlyIndexes, savedDependencyIndex) {
					currentRemovedDate := dependencyHistoryObject.Lastusedates[savedDependencyIndex]
					if currentRemovedDate > currentCommitDate || currentRemovedDate == 0 {
						println("setting last used date")
						dependencyHistoryObject.Lastusedates[savedDependencyIndex] = currentCommitDate
					}
				}

			}

			// what if it doesn't exist?

			// if it doesn't exist set set firstUseDate to now

		} else {
			for _, dependency := range dependencies {
				dependencyHistoryObject.Filenames = append(dependencyHistoryObject.Filenames, dependencyFileName)
				dependencyHistoryObject.Dependencynames = append(dependencyHistoryObject.Dependencynames, dependency)
				dependencyHistoryObject.Firstusedates = append(dependencyHistoryObject.Firstusedates, c.Committer.When.Unix())
				dependencyHistoryObject.Lastusedates = append(dependencyHistoryObject.Lastusedates, 0)

			}
		}

	}

	return nil
}
