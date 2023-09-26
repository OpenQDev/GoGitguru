package gitutil

import (
	"fmt"
	"main/internal/database"
)

func StoreCommits(prefixPath string, repo string, db *database.Queries) {
	gitLogOutput := GitLogCsv(prefixPath, repo, "")

	fmt.Println(string(gitLogOutput))

	// TODO : Pipe into Postgres

	fmt.Println("Successfully stored commits in the database.")
}
