package usersync

import (
	"context"
	"fmt"

	"github.com/OpenQDev/GoGitguru/database"
)

func getNewCommitAuthors(db *database.Queries) ([]database.GetLatestUncheckedCommitPerAuthorRow, error) {
	fmt.Println("in getNewCommitAuthors")
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	fmt.Println("read getNewCommitAuthors from db")

	if err != nil {
		return nil, err
	}

	noNewCommitAuthors := len(newCommitAuthorsRaw) == 0
	if noNewCommitAuthors {
		return nil, nil
	}

	return newCommitAuthorsRaw, nil
}
