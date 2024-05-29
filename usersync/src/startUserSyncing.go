package usersync

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail string
	RepoUrl     string
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	ghAccessToken string,
	batchSize int,
	githubGraphQLUrl string,
) {

	weekAgo := time.Now()
	usersDependenciesToSync, err := db.GetUserDependenciesByUpdatedAt(context.Background())
	if err != nil {
		logger.LogError("error getting user dependencies by updated at: %s", err)
	}
	// one week in go timevar
	bulkInsertUserDependenciesParams := database.BulkInsertUserDependenciesParams{

		UpdatedAt: sql.NullTime{Time: weekAgo, Valid: true},
	}

	for _, userDependency := range usersDependenciesToSync {
		println(bulkInsertUserDependenciesParams.UpdatedAt.Time.String())
		bulkInsertUserDependenciesParams.Column1 = append(bulkInsertUserDependenciesParams.Column1, userDependency.UserID)
		bulkInsertUserDependenciesParams.Column2 = append(bulkInsertUserDependenciesParams.Column2, userDependency.DependencyID)
		firstUseDate, _ok := userDependency.EarliestFirstUseDate.(int64)
		if !_ok {
			firstUseDate = 0
		}
		lastUseDate, _ok := userDependency.LatestLastUseDate.(int64)
		if !_ok {
			lastUseDate = 0
		}
		bulkInsertUserDependenciesParams.Column3 = append(bulkInsertUserDependenciesParams.Column3, firstUseDate)
		bulkInsertUserDependenciesParams.Column4 = append(bulkInsertUserDependenciesParams.Column4, lastUseDate)

		// find repo_deps that exist where user_id -> author -> commit  is shared with a repo
	}
	_, err = db.BulkInsertUserDependencies(context.Background(), bulkInsertUserDependenciesParams)
	if err != nil {
		logger.LogError("error inserting user dependencies: %s", err)
	}
	newCommitAuthorsRaw, err := getNewCommitAuthors(db)

	if err != nil {
		logger.LogFatalRedAndExit("error getting new commit authors to process: %s", err)
		return
	}
	if newCommitAuthorsRaw == nil {
		logger.LogBlue("no new authors to sync")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := convertDatabaseObjectToUserSync(newCommitAuthorsRaw)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := generateBatchAuthors(repoUrlToAuthorsMap, batchSize)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {

		githubGraphQLCommitAuthorsMap, err := identifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.AuthorCommitTuples, ghAccessToken, githubGraphQLUrl)
		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		if githubGraphQLCommitAuthorsMap == nil {
			logger.LogError("commits is nil")
			continue
		}

		githubGraphQLCommitAuthors := make([]GithubGraphQLCommit, 0, len(githubGraphQLCommitAuthorsMap))

		for _, commitAuthor := range githubGraphQLCommitAuthorsMap {
			githubGraphQLCommitAuthors = append(githubGraphQLCommitAuthors, commitAuthor)
		}

		for _, commitAuthor := range githubGraphQLCommitAuthors {
			author := commitAuthor.Author

			err := insertIntoRestIdToUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting author RestID %s to Email %s: %s", author.User.GithubRestID, author.Email, err)
			}

			exists, err := db.CheckGithubUserExists(context.Background(), strings.ToLower(author.User.Login))
			if err != nil {
				logger.LogError("error checking if github user exists: %s", err)
			}
			// TODO update their for that specific repo.
			if !exists {
				logger.LogBlue("inserting github user %s", author.Name)
				err = insertGithubUser(author, db)
				if err != nil {
					logger.LogError("error occured while inserting github user %s with RestId %s: %s", author.User.Login, author.User.GithubRestID, err)
				} else {
					logger.LogGreen("user %s inserted!", author.Name)
				}
			}

			// attached github user does exist

			if exists {
				continue
			}

		}
	}
}
