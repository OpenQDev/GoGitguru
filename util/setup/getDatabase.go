package setup

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/database"

	_ "github.com/lib/pq"
)

func GetDatbase(dbUrl string) (*database.Queries, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	queries := database.New(conn)

	return queries, nil
}

func GetSQLDatbase(dbUrl string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
