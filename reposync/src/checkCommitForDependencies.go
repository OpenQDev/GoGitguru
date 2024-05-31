package reposync

import (
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommitForDependencies(c *object.Commit, dependencyHistoryObject database.BatchInsertRepoDependenciesParams, dependencies []DependencyWithUpdatedTime, numberOfCommitsToSyncByDependency map[int32]int, commitCount int) (database.BatchInsertRepoDependenciesParams, error) {

	fileContentsLowerMap := make(map[string]string)

	for dependencyIndex, dependencyRecord := range dependencies {
		numberOfCommitsToSync := numberOfCommitsToSyncByDependency[dependencyRecord.InternalID]

		// only execute if we have commits to check
		if commitCount < numberOfCommitsToSync {

			fileContentsLower, ok := fileContentsLowerMap[dependencyRecord.DependencyFile]
			if !ok {
				file, err := c.File(dependencyRecord.DependencyFile)
				if err == nil {
					fileContents, err := file.Contents()
					if err != nil {
						fileContentsLower = ""
					}
					fileContentsLower = strings.ToLower(fileContents)
				}
				if err != nil {
					fileContentsLower = ""
				}
				fileContentsLowerMap[dependencyRecord.DependencyFile] = fileContentsLower
			}
			currentFileContentsLower := fileContentsLowerMap[dependencyRecord.DependencyFile]
			if currentFileContentsLower != "" {
				dependencySearchedLower := strings.ToLower(dependencyRecord.DependencyName)
				currentDateFirstPresent := dependencyHistoryObject.Firstusedates[dependencyIndex]
				currentDateLastRemoved := dependencyHistoryObject.Lastusedates[dependencyIndex]

				if strings.Contains(currentFileContentsLower, dependencySearchedLower) {
					if currentDateFirstPresent > c.Committer.When.Unix() || currentDateFirstPresent == 0 {
						currentDateFirstPresent = c.Committer.When.Unix()
					}
				} else {
					if currentDateFirstPresent != 0 {
						currentDateLastRemoved = c.Committer.When.Unix()
					}
				}
				dependencyHistoryObject.Dependencyids[dependencyIndex] = dependencyRecord.InternalID
				dependencyHistoryObject.Firstusedates[dependencyIndex] = currentDateFirstPresent
				dependencyHistoryObject.Lastusedates[dependencyIndex] = currentDateLastRemoved

			}
		}
	}
	return dependencyHistoryObject, nil
}
