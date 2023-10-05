package sync

import (
	"context"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
	"time"
)

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {

	newCommitAuthors, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())

	if err != nil {
		logger.LogError("errerrerr", err)
	}

	fmt.Println(newCommitAuthors)
}
