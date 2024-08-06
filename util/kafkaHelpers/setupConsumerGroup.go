package kafkahelpers

import (
	"crypto/tls"
	"fmt"

	"github.com/IBM/sarama"
)

func SetUpConsumerGroup(environment string, kafkaBrokers []string, group string) (sarama.ConsumerGroup, error) {
	// Set the SASL/OAUTHBEARER configuration
	// Set up the Sarama configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false

	fmt.Printf("Starting consumer for environment %s\n", environment)
	if environment == "production" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		config.Net.SASL.TokenProvider = &MSKAccessTokenProvider{}

		tlsConfig := tls.Config{}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tlsConfig
	}

	return sarama.NewConsumerGroup(kafkaBrokers, group, config)
}
