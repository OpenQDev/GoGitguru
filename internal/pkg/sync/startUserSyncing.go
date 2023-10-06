package sync

import (
	"context"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
	usersync "main/internal/pkg/sync/user"
	"time"
)

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration,
	ghAccessToken string,
) {
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	if err != nil {
		logger.LogError("errerrerr", err)
	}

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("No new authors to process.")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := usersync.ConvertToUserSync(newCommitAuthorsRaw)
	fmt.Println("newCommitAuthors", newCommitAuthors)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := usersync.GetRepoToAuthorsMap(newCommitAuthors)
	fmt.Println("repoUrlToAuthorsMap", repoUrlToAuthorsMap)

	// Create batches of repos for GraphQL query
	authorBatches := usersync.BatchAuthors(repoUrlToAuthorsMap, 2)
	fmt.Println("authorBatches", authorBatches)

	// Get info for each batch
	// for _, authorBatch := range authorBatches {
	// 	identifyRepoAuthorsBatch(authorBatch.repoUrl, authorBatch.authorList, ghAccessToken)
	// }
}
