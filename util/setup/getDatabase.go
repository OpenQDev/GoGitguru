package setup

import (
	"database/database"
	"database/sql"

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
