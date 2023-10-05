package sync

import (
	"context"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
	"time"
)

type UserSync struct {
	CommitHash string
	Author     struct {
		Email   string
		NotNull bool
	}
	Repo struct {
		URL     string
		NotNull bool
	}
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {
	/*
		[
			{
				65062be663cc004b77ca8a3b13255bc5efa42f25
				{andrew@openq.dev true}
				{https://github.com/openqdev/openq-workflows true}
			}
		]
		**/

	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	newCommitAuthors := convertToUserSync(newCommitAuthorsRaw)

	if err != nil {
		logger.LogError("errerrerr", err)
	}

	fmt.Println(newCommitAuthors)
}

func convertToUserSync(newCommitAuthorsRaw []database.GetLatestUncheckedCommitPerAuthorRow) []UserSync {
	var newCommitAuthors []UserSync

	for _, author := range newCommitAuthorsRaw {
		newCommitAuthors = append(newCommitAuthors, UserSync{
			CommitHash: author.CommitHash,
			Author: struct {
				Email   string
				NotNull bool
			}{
				Email:   author.AuthorEmail.String,
				NotNull: author.AuthorEmail.Valid,
			},
			Repo: struct {
				URL     string
				NotNull bool
			}{
				URL:     author.RepoUrl.String,
				NotNull: author.RepoUrl.Valid,
			},
		})
	}

	return newCommitAuthors
}
