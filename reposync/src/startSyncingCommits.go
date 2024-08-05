package reposync

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func StartSyncingCommits(
	db *database.Queries,
	conn *sql.DB,
	prefixPath string,
	gitguruUrl string,
	repoUrl string,
) {
	if repoUrl == "" {
		logger.LogBlue("no new repo urls to sync. exiting...")
		return
	}

	logger.LogGreenDebug("beginning sync for the following repo:\n%s", repoUrl)

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)
	logger.LogGreenDebug("processing %s/%s...", organization, repo)
	JAN_1_2020 := time.Unix(1577858400, 0)

	startDate := JAN_1_2020 // Unix time for Jan 1, 2020

	// Check if the repo is present in the repos directory
	if gitutil.IsGitRepository(prefixPath, organization, repo) {
		// If it is, pull the latest changes
		logger.LogBlue("repository %s exists. pulling...", repoUrl)
		err := gitutil.PullRepo(prefixPath, organization, repo)
		if err != nil {
			logger.LogError("error pulling repo %s/%s: %s", organization, repo, err)

			err = db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
				Status: database.RepoStatusFailed,
				Url:    repoUrl,
			})

			if err != nil {
				logger.LogError("error setting to synced for existing repository %s: %s", repoUrl, err)
			}
		}
		logger.LogBlue("repository %s pulled!", repoUrl)
		//	 there are cases where the repository may exist in local, but hasn't been synced
		//	 no rows in result set just means it didn't have any commit entries for that repo
		latestCommitterDate, err := db.GetLatestCommitterDate(context.Background(), repoUrl)
		if err != nil {
			if !strings.Contains(err.Error(), "sql: no rows in result set") {
				logger.LogFatalRedAndExit("error getting latest committer date: %s ", err)
			}
		}

		latestCommitterDateTime := time.Unix(int64(latestCommitterDate), 0)
		/// Unsure why but sometimes commits before JAN_1_2020 were being stored after initia clone-sync, causing issues
		if latestCommitterDateTime.After(JAN_1_2020) {
			startDate = latestCommitterDateTime
		}
	} else {
		// If not, clone it
		logger.LogBlue("repository %s does not exist. cloning...", repoUrl)
		err := gitutil.CloneRepo(prefixPath, organization, repo)
		if err != nil {
			logger.LogError("error cloning repo %s/%s: %s", organization, repo, err)

			err = db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
				Status: database.RepoStatusFailed,
				Url:    repoUrl,
			})
			if err != nil {
				logger.LogError("error setting repo status to failed", err)
			}
		}
		logger.LogBlue("repository %s cloned!", repoUrl)
	}

	emailList, repoWithUpdatedDeps, err := ProcessRepo(prefixPath, organization, repo, repoUrl, startDate, db)
	if err != nil {
		logger.LogError("error processing repo %s: %s", repoUrl, err)
		return
	}

	type GithubUserKafkaMessage struct {
		AuthorEmail string    `json:"author_email"`
		AuthorDate  time.Time `json:"author_date"`
		RepoUrl     string    `json:"repo_url"`
		CommitHash  string    `json:"commit_hash"`
	}

	type UserDependencyKafkaMessage struct {
		RepoUrl      string `json:"repo_url"`
		DependencyId int32  `json:"dependency_id"`
	}

	// Create a Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		logger.LogError("Failed to create Kafka producer: %s", err)
	} else {
		defer producer.Close()

		// Produce Kafka messages for each email in the list
		for _, user := range emailList {
			// Check if the email already exists in the rest_id_to_email table
			exists, err := db.CheckGithubUserRestIdAuthorEmailExists(context.Background(), user.AuthorEmail)
			if err != nil {
				logger.LogError("Failed to check if email exists: %s", err)
				continue
			}

			// If the email doesn't exist, produce a Kafka message
			if !exists {
				message := GithubUserKafkaMessage{
					AuthorEmail: user.AuthorEmail,
					AuthorDate:  user.AuthorDate,
					RepoUrl:     user.RepoUrl,
					CommitHash:  user.CommitHash,
				}

				jsonMessage, err := json.Marshal(message)
				if err != nil {
					logger.LogError("Failed to marshal message to JSON: %s", err)
					continue
				}
				fmt.Println(string(jsonMessage))
				_, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: "user-sync",
					Value: sarama.StringEncoder(jsonMessage),
				})
				if err != nil {
					logger.LogError("Failed to send message to Kafka: %s", err)
				} else {
					logger.LogGreenDebug("Message sent to Kafka: %s", user.AuthorEmail)
				}
			} else {
				logger.LogGreenDebug("Email already exists, skipping Kafka message: %s", user.AuthorEmail)
			}
		}

		for _, repo := range repoWithUpdatedDeps {

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

	if err != nil {
		logger.LogFatalRedAndExit("error while processing repository %s: %s", repoUrl, err)
	}
}
