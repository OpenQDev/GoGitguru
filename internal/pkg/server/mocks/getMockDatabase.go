package mocks

import (
	"main/internal/database"
	"main/internal/pkg/logger"

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
