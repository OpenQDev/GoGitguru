package reposync

import (
	"slices"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckCommitForDependencies(c *object.Commit, repoDir string, dependencyHistoryObject *database.BatchInsertRepoDependenciesParams) error {
	rawDependencyFiles := []string{
		"package.json",
		"requirements.txt",
		"pom.xml",
		"Pipfile",
		"go.mod",
		"build.gradle",
		"Gemfile",
		"Cargo.toml",
		".cabal",
		"composer.json",

		"hardhat.config",
		"truffle",
		"network",
		"deployments",
		"foundry.toml",
	}

	dependencyFiles, err := gitutil.GitDependencyFiles(repoDir, rawDependencyFiles)
	if err != nil {
		return err
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
		time := time.Now().Unix()
		dependencies := ParseFile(file)
		printExecutionTime(time, "ParseFile", dependencyFileName+" "+repoDir)

		// only handle matching file name
		if slices.Contains(dependencyHistoryObject.Filenames, dependencyFileName) {
			// should only be package.json
			dependenciesThatDoExistCurrentlyIndexes := []int{}
			// iterate over dependencies that exist in this file and commit within this loop we are looking at actual individual dependencies
			for _, dependency := range dependencies {

				dependencySavedIndex, dependenciesThatDoExistCurrentlyIndexesResult := getPreviousDependenciesInfo(dependencyHistoryObject, dependency, dependencyFileName, *c)
				dependenciesThatDoExistCurrentlyIndexes = append(dependenciesThatDoExistCurrentlyIndexes, dependenciesThatDoExistCurrentlyIndexesResult...)

				// handle dependency that doesn't currently exit
				if dependencySavedIndex == -1 {
					addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, dependencyFileName, c.Committer.When.Unix())
				} else {
					setDateFirstUsed(dependencyHistoryObject, dependencySavedIndex, *c)
				}

			}

			setDateRemoved(dependencyHistoryObject, dependenciesThatDoExistCurrentlyIndexes, currentCommitDate)

		} else {
			for _, dependency := range dependencies {
				addRowToDependencyHistoryObject(dependencyHistoryObject, dependency, dependencyFileName, c.Committer.When.Unix())

			}
		}
	}

	return nil
}
