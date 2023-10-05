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

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("No new authors to process.")
		return
	}

	newCommitAuthors := convertToUserSync(newCommitAuthorsRaw)

	if err != nil {
		logger.LogError("errerrerr", err)
	}

	fmt.Println(newCommitAuthors)

	logger.LogBlue("identifying %d new authors", len(newCommitAuthors))

	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

	fmt.Println(repoUrlToAuthorsMap)
}
