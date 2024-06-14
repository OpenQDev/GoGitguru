package usersync

import (
	"context"
	"database/sql"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
)

func SyncUserDependencies(db *database.Queries) error {

	// @flacojones, Any suggestions on how to better handle this?
	tenMinutesAgo := time.Now().Add(-10 * time.Minute).Unix()
	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background(), sql.NullInt64{Int64: tenMinutesAgo, Valid: true})
	if err != nil {
		return err
	}

	//WARNING SETTING UPDATEDAT TO TIME OF REPO SYNC NOT TIME OF USER SYNC
	bulkInsertUserDependenciesParams := database.BulkInsertUserDependenciesParams{
		UpdatedAt: tenMinutesAgo,
	}

	for _, userDependency := range usersDependenciesToSync {
		bulkInsertUserDependenciesParams.UserID = append(bulkInsertUserDependenciesParams.UserID, userDependency.UserID.Int32)
		bulkInsertUserDependenciesParams.DependencyID = append(bulkInsertUserDependenciesParams.DependencyID, userDependency.DependencyID)
		firstUseDate, _ok := userDependency.FirstUseDate.(int64)
		if !_ok {
			firstUseDate = 0
		}
		lastUseDate, _ok := userDependency.LastUseDate.(int64)
		if !_ok {
			lastUseDate = 0
		}
		bulkInsertUserDependenciesParams.FirstUseDate = append(bulkInsertUserDependenciesParams.FirstUseDate, firstUseDate)
		bulkInsertUserDependenciesParams.LastUseDate = append(bulkInsertUserDependenciesParams.LastUseDate, lastUseDate)

		// find repo_deps that exist where user_id -> author -> commit  is shared with a repo
	}
	err = db.BulkInsertUserDependencies(context.Background(), bulkInsertUserDependenciesParams)
	if err != nil {
		return err
	}
	return nil

}
