package reposync

import (
	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
)

func PrepareDependencyHistoryForBulkInsertion(dependencyHistory map[int32]gitutil.DependencyResult, dependencies []database.Dependency, repoUrl string) (RepoDependencyHistoryObject, error) {
	numberOfDependencies := len(dependencies)
	var (
		dependencyId     = make([]int32, numberOfDependencies)
		dateFirstPresent = make([]int64, numberOfDependencies)
		dateLastRemoved  = make([]int64, numberOfDependencies)
	)
	for i, dependency := range dependencies {
		dependencyId[i] = dependency.InternalID
		dateFirstPresent[i] = dependencyHistory[dependency.InternalID].DateFirstPresent
		dateLastRemoved[i] = dependencyHistory[dependency.InternalID].DateLastRemoved
	}
	depenencyHistoryObject := RepoDependencyHistoryObject{
		DependencyId:     dependencyId,
		DateFirstPresent: dateFirstPresent,
		DateLastRemoved:  dateLastRemoved,
	}
	return depenencyHistoryObject, nil
}
