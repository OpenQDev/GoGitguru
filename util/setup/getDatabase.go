package setup

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/database"

	_ "github.com/lib/pq"
)

func GetDatbase(dbUrl string) (*database.Queries, *sql.DB, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, nil, err
	}

	queries := database.New(conn)

	return queries, conn, nil
}
