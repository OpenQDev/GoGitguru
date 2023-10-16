package usersync

import (
	"context"
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/githubGraphQL"
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
	batchSize int,
) {
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	if err != nil {
		logger.LogError("error in GetLatestUncheckedCommitPerAuthor: %s", err)
	}

	logger.LogGreenDebug("new commit authors to check: %s", newCommitAuthorsRaw)

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("no new authors to process.")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := ConvertToUserSync(newCommitAuthorsRaw)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := GetRepoToAuthorsMap(newCommitAuthors)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := GenerateBatchAuthors(repoUrlToAuthorsMap, 2)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		commits, err := IdentifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.Tuples, ghAccessToken)

		if commits == nil {
			logger.LogError("commits is nil")
			continue
		}

		commitValues := make([]githubGraphQL.Commit, 0, len(*commits))

		for _, value := range *commits {
			commitValues = append(commitValues, value)
		}

		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		logger.LogGreenDebug("got the following info: %v", commits)

		for _, commit := range commitValues {
			restId := commit.Author.User.GithubRestID
			author := commit.Author

			var params database.InsertRestIdToEmailParams
			if restId == 0 {
				params = database.InsertRestIdToEmailParams{
					Email: author.Email,
				}
			} else {
				params = database.InsertRestIdToEmailParams{
					RestID: int32(restId),
					Email:  author.Email,
				}
			}

			_, err = db.InsertRestIdToEmail(context.Background(), params)
			if err != nil {
				logger.LogError("error occured while inserting RestID to Email: %s", err)
			}

			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			updatedAt, err := time.Parse(time.RFC3339, author.User.UpdatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			authorParams := ConvertAuthorToInsertUserParams(author, createdAt, updatedAt)

			_, err = db.InsertUser(context.Background(), authorParams)

			if err != nil {
				logger.LogError("error occured while inserting author: %s", err)
			}

		}
	}
}

func ConvertAuthorToInsertUserParams(author githubGraphQL.Author, createdAt time.Time, updatedAt time.Time) database.InsertUserParams {
	authorParams := database.InsertUserParams{
		GithubRestID:    int32(author.User.GithubRestID),
		GithubGraphqlID: author.User.GithubGraphqlID,
		Login:           author.User.Login,
		Name:            sql.NullString{String: author.User.Name, Valid: true},
		Email:           sql.NullString{String: author.User.Email, Valid: true},
		AvatarUrl:       sql.NullString{String: author.User.AvatarURL, Valid: true},
		Company:         sql.NullString{String: author.User.Company, Valid: true},
		Location:        sql.NullString{String: author.User.Location, Valid: true},
		Bio:             sql.NullString{String: author.User.Bio, Valid: true},
		Blog:            sql.NullString{String: author.User.Blog, Valid: true},
		Hireable:        sql.NullBool{Bool: author.User.Hireable, Valid: true},
		TwitterUsername: sql.NullString{String: author.User.TwitterUsername, Valid: true},
		Followers:       sql.NullInt32{Int32: int32(author.User.Followers.TotalCount), Valid: true},
		Following:       sql.NullInt32{Int32: int32(author.User.Following.TotalCount), Valid: true},
		Type:            "User",
		CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
	}
	return authorParams
}
