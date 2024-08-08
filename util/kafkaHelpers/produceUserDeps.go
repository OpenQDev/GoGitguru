package kafkahelpers

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type UserDependencyKafkaMessage struct {
	RepoUrl string `json:"repo_url"`
}

func ProduceUserDependencyMessage(producer sarama.SyncProducer, env setup.EnvConfig, repoUrlsWithNewUsers []string) {

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
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: env.UserDepsSyncTopic,
			Value: sarama.StringEncoder(jsonMessage),
		})
		if err != nil {
			logger.LogError("Failed to send message to Kafka: %s", err)
		} else {
			logger.LogGreenDebug("Message sent to Kafka: looking at latest dependency")
		}
	}

}
