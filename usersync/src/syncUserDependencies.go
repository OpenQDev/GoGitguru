package usersync

import (
	"context"
	"fmt"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

type DepsMessage struct {
	RepoUrl string `json:"repo_url"`
}

func SyncUserDependencies(db *database.Queries, repoUrl string) error {
	fmt.Println("Syncing user dependencies for repo: ", repoUrl)
	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background(), repoUrl)
	if err != nil {
		logger.LogError("error getting user dependencies to sync since: %s", err)
		return err
	}
	fmt.Println("Got user dependencies to sync", len(usersDependenciesToSync))

	// get user dependencies that go wtih the values from above

	getPreviousUserDeps := database.GetUserDependenciesByUserParams{
		UserIds:       []int32{},
		DependencyIds: []int32{},
	}

	for _, userDependency := range usersDependenciesToSync {
		getPreviousUserDeps.UserIds = append(getPreviousUserDeps.UserIds, userDependency.UserID.Int32)
		getPreviousUserDeps.DependencyIds = append(getPreviousUserDeps.DependencyIds, userDependency.DependencyID)
	}
	// get associated user dependencies based that have already been synced
	alreadySyncedUserDependencies, err := db.GetUserDependenciesByUser(context.Background(), getPreviousUserDeps)
	if err != nil {
		logger.LogError("error getting already synced user dependencies: %s", err)
		return err
	}

	bulkInsertUserDependenciesParams := PrepareUserDependencies(usersDependenciesToSync, alreadySyncedUserDependencies)

	err = db.BulkInsertUserDependencies(context.Background(), bulkInsertUserDependenciesParams)
	if err != nil {
		logger.LogError("error inserting user dependencies: %s", err)
		return err
	}
	return nil

}
