package dbscripts

import (
	"context"
	"fmt"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func ConsolidateFileName() {

	env := setup.ExtractAndVerifyEnvironment("../.env")
	db, _, _ := setup.GetDatbase(env.DbUrl)
	dependencyFiles := []string{
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
		`network`,
		`deployments`,
		"foundry.toml",
	}
	for _, dep := range dependencyFiles {
		params := database.GetDependenciesByFilesParams{
			DependencyFileLike: "%" + dep + "%",
			DependencyFile:     dep,
		}
		fmt.Println("Getting dependencies by files")
		dependencyNames, err := db.GetDependenciesByFiles(context.Background(), params)
		if err != nil {
			logger.LogError("error getting dependencies by files", err)
		}

		fmt.Println("Upserting missing dependencies")
		err = db.UpsertMissingDependencies(context.Background(), database.UpsertMissingDependenciesParams{
			DependencyFile:  dep,
			DependencyNames: dependencyNames,
		})

		if err != nil {
			logger.LogError("error upserting missing dependencies", err)
		}
		fmt.Println("Switching repos relation to simple")
		err = db.SwitchReposRelationToSimple(context.Background(), database.SwitchReposRelationToSimpleParams{
			DependencyFile:     dep,
			DependencyFileLike: "%" + dep + "%",
		})
		if err != nil {
			logger.LogError("error switching repos relation to simple", err)
		}
		err = db.SwitchUsersRelationToSimple(context.Background(), database.SwitchUsersRelationToSimpleParams{
			DependencyFile: dep,
		})

		err = db.DeleteUnusedDependencies(context.Background(), database.DeleteUnusedDependenciesParams{
			DependencyFile:     dep,
			DependencyFileLike: "%" + dep + "%",
		})
		if err != nil {

			logger.LogError("error deleting unused dependencies", err)
		}

	}

}
