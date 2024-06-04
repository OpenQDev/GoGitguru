package usersync

import (
	"context"

	"github.com/OpenQDev/GoGitguru/database"
)

func SyncUserDependencies(db *database.Queries) error {
	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background())
	if err != nil {
		return err
	}
	bulkInsertUserDependenciesParams := database.BulkInsertUserDependenciesParams{}

	for _, userDependency := range usersDependenciesToSync {
		bulkInsertUserDependenciesParams.Column1 = append(bulkInsertUserDependenciesParams.Column1, userDependency.UserID.Int32)
		bulkInsertUserDependenciesParams.Column2 = append(bulkInsertUserDependenciesParams.Column2, userDependency.DependencyID)
		firstUseDate, _ok := userDependency.FirstUseDate.(int64)
		if !_ok {
			firstUseDate = 0
		}
		lastUseDate, _ok := userDependency.LastUseDate.(int64)
		if !_ok {
			lastUseDate = 0
		}
		bulkInsertUserDependenciesParams.Column3 = append(bulkInsertUserDependenciesParams.Column3, firstUseDate)
		bulkInsertUserDependenciesParams.Column4 = append(bulkInsertUserDependenciesParams.Column4, lastUseDate)

		// find repo_deps that exist where user_id -> author -> commit  is shared with a repo
	}
	_, err = db.BulkInsertUserDependencies(context.Background(), bulkInsertUserDependenciesParams)
	if err != nil {
		return err
	}
	return nil

}
