/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/mcc-github/blockchain-lib-go/healthz"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/consensus"
)




type healthChecker interface {
	RegisterChecker(component string, checker healthz.HealthChecker) error
}


func New(config localconfig.Kafka, metricsProvider metrics.Provider, healthChecker healthChecker) (consensus.Consenter, *Metrics) {
	if config.Verbose {
		flogging.ActivateSpec(flogging.Global.Spec() + ":orderer.consensus.kafka.sarama=debug")
	}

	brokerConfig := newBrokerConfig(
		config.TLS,
		config.SASLPlain,
		config.Retry,
		config.Version,
		defaultPartition)

	metrics := NewMetrics(metricsProvider, brokerConfig.MetricRegistry)

	return &consenterImpl{
		brokerConfigVal: brokerConfig,
		tlsConfigVal:    config.TLS,
		retryOptionsVal: config.Retry,
		kafkaVersionVal: config.Version,
		topicDetailVal: &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: config.Topic.ReplicationFactor,
		},
		healthChecker: healthChecker,
		metrics:       metrics,
	}, metrics
}




type consenterImpl struct {
	brokerConfigVal *sarama.Config
	tlsConfigVal    localconfig.TLS
	retryOptionsVal localconfig.Retry
	kafkaVersionVal sarama.KafkaVersion
	topicDetailVal  *sarama.TopicDetail
	healthChecker   healthChecker
	metrics         *Metrics
}






func (consenter *consenterImpl) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	lastOffsetPersisted, lastOriginalOffsetProcessed, lastResubmittedConfigOffset := getOffsets(metadata.Value, support.ChannelID())
	ch, err := newChain(consenter, support, lastOffsetPersisted, lastOriginalOffsetProcessed, lastResubmittedConfigOffset)
	if err != nil {
		return nil, err
	}
	consenter.healthChecker.RegisterChecker(ch.channel.String(), ch)
	return ch, nil
}





type commonConsenter interface {
	brokerConfig() *sarama.Config
	retryOptions() localconfig.Retry
	topicDetail() *sarama.TopicDetail
	Metrics() *Metrics
}

func (consenter *consenterImpl) Metrics() *Metrics {
	return consenter.metrics
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
