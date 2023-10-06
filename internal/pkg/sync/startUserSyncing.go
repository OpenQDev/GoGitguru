package sync

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/logger"
	"time"
)

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration,
	ghAccessToken string,
	batchSize int,
) {
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	if err != nil {
		logger.LogError("errerrerr", err)
	}

	logger.LogGreenDebug("new commit authors to check: %s", newCommitAuthorsRaw)

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("No new authors to process.")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := ConvertToUserSync(newCommitAuthorsRaw)
	logger.LogGreenDebug("newCommitAuthors", newCommitAuthors)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := GetRepoToAuthorsMap(newCommitAuthors)
	logger.LogGreenDebug("repoUrlToAuthorsMap", repoUrlToAuthorsMap)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := BatchAuthors(repoUrlToAuthorsMap, 2)
	logger.LogGreenDebug("repoToAuthorBatches", repoToAuthorBatches)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		repoUrl, ok := repoToAuthorBatch[0].(string)
		if !ok {
			logger.LogError("Unable to cast repoToAuthorBatch[0] to string")
			continue
		}
		authors, ok := repoToAuthorBatch[1].([]string)
		if !ok {
			logger.LogError("Unable to cast repoToAuthorBatch[1] to []string")
			continue
		}

		logger.LogGreenDebug("%s: %v", repoUrl, authors)

		commits, err := IdentifyRepoAuthorsBatch(repoUrl, authors, ghAccessToken)
		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s -> %v", repoUrl, authors)

		logger.LogGreenDebug("got the following info: %v", commits)
	}
}
