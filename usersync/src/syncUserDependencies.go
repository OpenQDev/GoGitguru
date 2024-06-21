package usersync

import (
	"context"
	"database/sql"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
)

func SyncUserDependencies(db *database.Queries) error {

	// get all recent repoDependencies  @flacojones, do you want a smaller time window?
	tenMinutesAgo := time.Now().AddDate(0, 0, -1000).Unix()
	// find user deps based off recent repos
	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background(), sql.NullInt64{Int64: tenMinutesAgo, Valid: true})
	if err != nil {
		return err
	}

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
		return err
	}

	bulkInsertUserDependenciesParams := database.BulkInsertUserDependenciesParams{
		UpdatedAt: tenMinutesAgo,
	}

	for _, userDependency := range usersDependenciesToSync {
		firstUseDate, _ok := userDependency.FirstUseDate.(int64)
		if !_ok {
			firstUseDate = 0
		}
		lastUseDate, _ok := userDependency.LastUseDate.(int64)
		if !_ok {
			lastUseDate = 0
		}
		// reset to already synced vals if necessary
		for _, alreadySynced := range alreadySyncedUserDependencies {
			if userDependency.UserID.Int32 == alreadySynced.UserID && userDependency.DependencyID == alreadySynced.DependencyID {
				if alreadySynced.FirstUseDate.Int64 < firstUseDate && alreadySynced.FirstUseDate.Int64 != 0 {
					firstUseDate = alreadySynced.FirstUseDate.Int64
				}
				if alreadySynced.LastUseDate.Int64 > lastUseDate && alreadySynced.LastUseDate.Int64 != 0 {
					lastUseDate = alreadySynced.LastUseDate.Int64
				}
			}
		}
		bulkInsertUserDependenciesParams.UserID = append(bulkInsertUserDependenciesParams.UserID, userDependency.UserID.Int32)
		bulkInsertUserDependenciesParams.DependencyID = append(bulkInsertUserDependenciesParams.DependencyID, userDependency.DependencyID)

		bulkInsertUserDependenciesParams.FirstUseDate = append(bulkInsertUserDependenciesParams.FirstUseDate, firstUseDate)
		bulkInsertUserDependenciesParams.LastUseDate = append(bulkInsertUserDependenciesParams.LastUseDate, lastUseDate)

	}

	err = db.BulkInsertUserDependencies(context.Background(), bulkInsertUserDependenciesParams)
	if err != nil {
		return err
	}
	return nil

}
