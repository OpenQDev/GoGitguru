package usersync

import (
	"fmt"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/lib"
)

func PrepareUserDependencies(usersDependenciesToSync []database.GetUserDependenciesByUpdatedAtRow, alreadySyncedUserDependencies []database.GetUserDependenciesByUserRow) database.BulkInsertUserDependenciesParams {
	fmt.Println("Preparing user dependencies")

	tenMinutesAgo := lib.Now().Add(-10 * time.Minute).Unix()
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
				newDepIsActive := lastUseDate == 0 && firstUseDate != 0
				if alreadySynced.LastUseDate.Int64 > lastUseDate && alreadySynced.LastUseDate.Int64 != 0 && !newDepIsActive {
					lastUseDate = alreadySynced.LastUseDate.Int64
				}
			}
		}
		bulkInsertUserDependenciesParams.UserID = append(bulkInsertUserDependenciesParams.UserID, userDependency.UserID.Int32)
		bulkInsertUserDependenciesParams.DependencyID = append(bulkInsertUserDependenciesParams.DependencyID, userDependency.DependencyID)

		bulkInsertUserDependenciesParams.FirstUseDate = append(bulkInsertUserDependenciesParams.FirstUseDate, firstUseDate)
		bulkInsertUserDependenciesParams.LastUseDate = append(bulkInsertUserDependenciesParams.LastUseDate, lastUseDate)

	}
	return bulkInsertUserDependenciesParams

}
