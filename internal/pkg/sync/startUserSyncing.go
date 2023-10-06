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
	repoToAuthorBatches := GenerateBatchAuthors(repoUrlToAuthorsMap, 2)
	logger.LogGreenDebug("repoToAuthorBatches", repoToAuthorBatches)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		commits, err := IdentifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.Tuples, ghAccessToken)
		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		logger.LogGreenDebug("got the following info: %v", commits)
	}
}
