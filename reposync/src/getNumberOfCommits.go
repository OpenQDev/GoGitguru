package reposync

import (
	"time"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
)

type NumberOfCommits struct {
	ToCheckByDependency          map[int32]int
	GreatestToCheckForDependency int
	ToSync                       int
}

func GetNumberOfCommitsPerDependency(dependencies []DependencyWithUpdatedTime, params GitLogParams) (NumberOfCommits, error) {
	numberOfCommitsToSyncByDependency := make(map[int32]int)
	var err error
	greatestNumberOfCommitsToCheckForDependency := 0
	for _, dependency := range dependencies {
		lastUpdatedTime := time.Unix(dependency.UpdatedAt, 0)
		numberCommits, localError := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, dependency.DependencyName, lastUpdatedTime)

		if localError != nil {
			err = localError
			break
		}
		numberOfCommitsToSyncByDependency[dependency.InternalID] = numberCommits
		if numberCommits > greatestNumberOfCommitsToCheckForDependency {
			greatestNumberOfCommitsToCheckForDependency = numberCommits
		}
	}

	numberOfCommitsToSync, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)

	result := NumberOfCommits{
		ToCheckByDependency:          numberOfCommitsToSyncByDependency,
		GreatestToCheckForDependency: greatestNumberOfCommitsToCheckForDependency,
		ToSync:                       numberOfCommitsToSync,
	}
	return result, err

}
