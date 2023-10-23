package usersync

import (
	"context"
	"database/database"
)

func getNewCommitAuthors(db *database.Queries) ([]database.GetLatestUncheckedCommitPerAuthorRow, error) {
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	if err != nil {
		return nil, err
	}

	noNewCommitAuthors := len(newCommitAuthorsRaw) == 0
	if noNewCommitAuthors {
		return nil, nil
	}

	return newCommitAuthorsRaw, nil
}
