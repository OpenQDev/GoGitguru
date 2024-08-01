package usersync

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/lib"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func SyncUserDependencies(db *database.Queries) error {

	// get all recent repoDependencies  @flacojones, do you want a smaller time window?
	fiveMinutesAgo := lib.Now().Add(-5 * time.Minute).Unix()
	fifteenMinutesAgo := lib.Now().Add(-15 * time.Minute).Unix()
	// find user deps based off recent repos

	usersToDependenciesParams := database.GetUserDependenciesByUpdatedAtParams{
		Since: sql.NullInt64{Int64: fifteenMinutesAgo, Valid: true},
		Until: sql.NullInt64{Int64: fiveMinutesAgo, Valid: true},
	}

	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background(), usersToDependenciesParams)
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
