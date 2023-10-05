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
	newCommitAuthors := convertToUserSync(newCommitAuthorsRaw)
	fmt.Println("newCommitAuthors", newCommitAuthors)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)
	fmt.Println("repoUrlToAuthorsMap", repoUrlToAuthorsMap)

	// Create batches of repos for GraphQL query
	authorBatches := batchAuthors(repoUrlToAuthorsMap, 2)
	fmt.Println("authorBatches", authorBatches)

	// Get info for each batch
	// for _, authorBatch := range authorBatches {
	// 	identifyRepoAuthorsBatch(authorBatch.repoUrl, authorBatch.authorList, ghAccessToken)
	// }
}
