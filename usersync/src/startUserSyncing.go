package usersync

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail string
	RepoUrl     string
}

type Message struct {
	Author_Email string `json:"author_email"`
	Author_Date  string `json:"author_date"`
	Repo_URL     string `json:"repo_url"`
	CommitHash   string `json:"commit_hash"`
}

func setupProducer(environment string, kafkaBrokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	fmt.Printf("Starting producer for environment %s\n", environment)
	if environment == "production" {
		// Set the SASL/OAUTHBEARER configuration
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		config.Net.SASL.TokenProvider = &MSKAccessTokenProvider{}

		// Enable TLS
		tlsConfig := tls.Config{}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tlsConfig
	}

	return sarama.NewSyncProducer(kafkaBrokers, config)
}

func StartUserSyncing(
	db *database.Queries,
	prefixPath string,
	ghAccessToken string,
	batchSize int,
	githubGraphQLUrl string,
	message Message,
) {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	logger.LogBlue("identifying %d new authors", len(message.Author_Email))

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := getRepoToAuthorsMap(message)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := generateBatchAuthors(repoUrlToAuthorsMap, batchSize)

	repoUrlsWithNewUsers := []string{}

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

		upsertRepoToUserByIdParams := database.UpsertRepoToUserByIdParams{
			Url: repoToAuthorBatch.RepoURL,
		}
		for _, commitAuthor := range githubGraphQLCommitAuthors {
			author := commitAuthor.Author

			err := insertIntoRestIdToUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting author RestID %s to Email %s: %s", author.User.GithubRestID, author.Email, err)
			}

			result, err := db.CheckGithubUserIdExists(context.Background(), author.User.GithubRestID)
			if err != nil {
				logger.LogError("error checking if github user exists: %s", err)
			}
			// TODO update their for that specific repo.
			if !result {
				logger.LogBlue("inserting github user %s", author.Name)
				err := insertGithubUser(author, db)
				if err != nil {
					logger.LogError("error occured while inserting github user %s with RestId %s: %s", author.User.Login, author.User.GithubRestID, err)
				} else {
					logger.LogGreen("user %s inserted!", author.Name)
				}

			}

			internal_id, err := db.GetGithubUserByRestId(context.Background(), author.User.GithubRestID)

			if err != nil {
				logger.LogError("error occured while getting GetGithubUserByRestId: %s", err)
			}

			err = GetReposToUsers(db, &upsertRepoToUserByIdParams, internal_id, author)

			if err != nil {
				logger.LogError("error occured while getting repos to users: %s", err)
			}
		}

		if err != nil {
			logger.LogError("error occured while getting repos to users: %s", err)
		}

		if !slices.Contains(repoUrlsWithNewUsers, repoToAuthorBatch.RepoURL) {
			repoUrlsWithNewUsers = append(repoUrlsWithNewUsers, repoToAuthorBatch.RepoURL)
		}

		err = db.UpsertRepoToUserById(context.Background(), upsertRepoToUserByIdParams)
		if err != nil {
			logger.LogError("error occured while upserting repo to user by id: %s", err)
		}
	}

	type UserDependencyKafkaMessage struct {
		RepoUrl string `json:"repo_url"`
	}

	// Create a Kafka producer

	brokers := strings.Split(env.KafkaBrokerUrls, ",")
	producer, err := setupProducer(env.Environment, brokers)
	if err != nil {
		logger.LogError("Failed to create Kafka producer: %s", err)
	} else {
		defer producer.Close()

		// Produce Kafka messages for each email in the list
		for _, repo := range repoUrlsWithNewUsers {

			message := UserDependencyKafkaMessage{
				RepoUrl: repo,
			}

			jsonMessage, err := json.Marshal(message)
			if err != nil {
				logger.LogError("Failed to marshal message to JSON: %s", err)
				continue
			}
			fmt.Println(string(jsonMessage))
			_, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Topic: "user-dependency-sync",
				Value: sarama.StringEncoder(jsonMessage),
			})
			if err != nil {
				logger.LogError("Failed to send message to Kafka: %s", err)
			} else {
				logger.LogGreenDebug("Message sent to Kafka: looking at latest dependency")
			}
		}
	}

}
