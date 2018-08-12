/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/Shopify/sarama"
	localconfig "github.com/mcc-github/blockchain/orderer/common/localconfig"
)

func newBrokerConfig(
	tlsConfig localconfig.TLS,
	saslPlain localconfig.SASLPlain,
	retryOptions localconfig.Retry,
	kafkaVersion sarama.KafkaVersion,
	chosenStaticPartition int32) *sarama.Config {

	
	paddingDelta := 1 * 1024 * 1024

	brokerConfig := sarama.NewConfig()

	brokerConfig.Consumer.Retry.Backoff = retryOptions.Consumer.RetryBackoff

	
	brokerConfig.Consumer.Return.Errors = true

	brokerConfig.Metadata.Retry.Backoff = retryOptions.Metadata.RetryBackoff
	brokerConfig.Metadata.Retry.Max = retryOptions.Metadata.RetryMax

	brokerConfig.Net.DialTimeout = retryOptions.NetworkTimeouts.DialTimeout
	brokerConfig.Net.ReadTimeout = retryOptions.NetworkTimeouts.ReadTimeout
	brokerConfig.Net.WriteTimeout = retryOptions.NetworkTimeouts.WriteTimeout

	brokerConfig.Net.TLS.Enable = tlsConfig.Enabled
	if brokerConfig.Net.TLS.Enable {
		
		keyPair, err := tls.X509KeyPair([]byte(tlsConfig.Certificate), []byte(tlsConfig.PrivateKey))
		if err != nil {
			logger.Panic("Unable to decode public/private key pair:", err)
		}
		
		rootCAs := x509.NewCertPool()
		for _, certificate := range tlsConfig.RootCAs {
			if !rootCAs.AppendCertsFromPEM([]byte(certificate)) {
				logger.Panic("Unable to parse the root certificate authority certificates (Kafka.Tls.RootCAs)")
			}
		}
		brokerConfig.Net.TLS.Config = &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			RootCAs:      rootCAs,
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   0, 
		}
	}
	brokerConfig.Net.SASL.Enable = saslPlain.Enabled
	if brokerConfig.Net.SASL.Enable {
		brokerConfig.Net.SASL.User = saslPlain.User
		brokerConfig.Net.SASL.Password = saslPlain.Password
	}

	
	
	brokerConfig.Producer.MaxMessageBytes = int(sarama.MaxRequestSize) - paddingDelta

	brokerConfig.Producer.Retry.Backoff = retryOptions.Producer.RetryBackoff
	brokerConfig.Producer.Retry.Max = retryOptions.Producer.RetryMax

	
	
	brokerConfig.Producer.Partitioner = newStaticPartitioner(chosenStaticPartition)
	
	
	
	brokerConfig.Producer.RequiredAcks = sarama.WaitForAll
	
	
	brokerConfig.Producer.Return.Successes = true

	brokerConfig.Version = kafkaVersion

	return brokerConfig
}
