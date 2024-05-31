package reposync

import (
	"context"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
)

// from commitDate should be the date of the last commit that was synced for the repository or any of the dependencies.

func GetRepoDependencies(db *database.Queries, repoUrl string) ([]DependencyWithUpdatedTime, error) {
	rawDependencies, err := db.GetDependencies(context.Background(), repoUrl)
	if err != nil {
		println("error getting rawDependencies")
	}
	dependenciesFiles := make([]string, len(rawDependencies))
	dependencyNames := make([]string, len(rawDependencies))
	dependencies := []DependencyWithUpdatedTime{}
	for i, rawDependency := range rawDependencies {
		dependency := DependencyWithUpdatedTime{
			DependencyName: rawDependency.DependencyName.String,
			DependencyFile: rawDependency.DependencyFile.String,
			UpdatedAt:      rawDependency.UpdatedAt.Int64,
			InternalID:     rawDependency.InternalID.Int32,
		}
		dependencies = append(dependencies, dependency)

		if slices.Contains(dependenciesFiles, dependency.DependencyFile) {
			dependenciesFiles[i] = dependency.DependencyFile
		}
		if slices.Contains(dependencyNames, dependency.DependencyName) {
			dependencyNames[i] = dependency.DependencyName
		}
	}
	return dependencies, err
}
