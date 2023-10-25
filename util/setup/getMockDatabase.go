package setup

import (
	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

func GetMockDatabase() (sqlmock.Sqlmock, *database.Queries) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.LogFatalRedAndExit("can't create mock DB: %s", err)
	}

	queries := database.New(db)

	return mock, queries
}
