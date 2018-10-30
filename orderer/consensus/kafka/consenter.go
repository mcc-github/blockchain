/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/Shopify/sarama"
	localconfig "github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
	logging "github.com/op/go-logging"
)


func New(config localconfig.Kafka) consensus.Consenter {
	if config.Verbose {
		logging.SetLevel(logging.DEBUG, "orderer.consensus.kafka.sarama")
	}
	brokerConfig := newBrokerConfig(
		config.TLS,
		config.SASLPlain,
		config.Retry,
		config.Version,
		defaultPartition)
	return &consenterImpl{
		brokerConfigVal: brokerConfig,
		tlsConfigVal:    config.TLS,
		retryOptionsVal: config.Retry,
		kafkaVersionVal: config.Version,
		topicDetailVal: &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: config.Topic.ReplicationFactor,
		},
	}
}




type consenterImpl struct {
	brokerConfigVal *sarama.Config
	tlsConfigVal    localconfig.TLS
	retryOptionsVal localconfig.Retry
	kafkaVersionVal sarama.KafkaVersion
	topicDetailVal  *sarama.TopicDetail
}






func (consenter *consenterImpl) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	lastOffsetPersisted, lastOriginalOffsetProcessed, lastResubmittedConfigOffset := getOffsets(metadata.Value, support.ChainID())
	return newChain(consenter, support, lastOffsetPersisted, lastOriginalOffsetProcessed, lastResubmittedConfigOffset)
}





type commonConsenter interface {
	brokerConfig() *sarama.Config
	retryOptions() localconfig.Retry
	topicDetail() *sarama.TopicDetail
}

func (consenter *consenterImpl) brokerConfig() *sarama.Config {
	return consenter.brokerConfigVal
}

func (consenter *consenterImpl) retryOptions() localconfig.Retry {
	return consenter.retryOptionsVal
}

func (consenter *consenterImpl) topicDetail() *sarama.TopicDetail {
	return consenter.topicDetailVal
}


type closeable interface {
	close() error
}
