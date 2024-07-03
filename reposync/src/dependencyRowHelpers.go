package reposync

import (
	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type DependencyToAdd struct {
	DependencyName string
	DependencyFile string
	DateFirstUsed  int64
	DateLastUsed   int64
}

// a file of helpers that mutate dependencyHistoryObject within checkCommitForDependencies.go

// finds where current dependencies sit in the dependencyHistoryObject
func getPreviousDependenciesInfo(dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, currentDependency string, currentDependencyFile string, commit object.Commit) int {
	dependencySavedIndex := -1

	for savedDependencyNameIndex, savedDependencyName := range dependencyHistoryObject.Dependencynames {
		// only filename at matching index
		for savedDependencyFileNameIndex, savedDependencyFileName := range dependencyHistoryObject.Filenames {
			fileNameIsCorrect := savedDependencyFileName == currentDependencyFile
			dependencyNameIsCorrect := savedDependencyName == currentDependency
			indexIsCorrect := savedDependencyNameIndex == savedDependencyFileNameIndex
			if fileNameIsCorrect && dependencyNameIsCorrect && indexIsCorrect {
				// set for outside of loop

				dependencySavedIndex = savedDependencyNameIndex

			}
		}
	}
	return dependencySavedIndex
}
func addRowToDependencyHistoryObject(dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, currentDependency string, currentDependencyFile string, firstUseDate int64) {
	dependencyHistoryObject.Filenames = append(dependencyHistoryObject.Filenames, currentDependencyFile)
	dependencyHistoryObject.Dependencynames = append(dependencyHistoryObject.Dependencynames, currentDependency)
	dependencyHistoryObject.Firstusedates = append(dependencyHistoryObject.Firstusedates, firstUseDate)
	dependencyHistoryObject.Lastusedates = append(dependencyHistoryObject.Lastusedates, 0)
}

func setDateRemoved(dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, currentCommitDate int64) {

	for savedDependencyIndex := range dependencyHistoryObject.Dependencynames {
		currentRemovedDate := dependencyHistoryObject.Lastusedates[savedDependencyIndex]
		if currentRemovedDate < currentCommitDate || currentRemovedDate == 0 {
			dependencyHistoryObject.Lastusedates[savedDependencyIndex] = currentCommitDate
		}

	}
}

func setDateFirstUsed(dependencyHistoryObject *database.BatchInsertRepoDependenciesParams, dependencySavedIndex int, c object.Commit) {
	currentDateFirstPresent := dependencyHistoryObject.Firstusedates[dependencySavedIndex]
	if currentDateFirstPresent > c.Committer.When.Unix() || currentDateFirstPresent == 0 {
		dependencyHistoryObject.Firstusedates[dependencySavedIndex] = c.Committer.When.Unix()
	}
}
