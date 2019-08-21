/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/consensus"
	"github.com/mcc-github/blockchain/orderer/consensus/kafka/mock"
	mockmultichannel "github.com/mcc-github/blockchain/orderer/mocks/common/multichannel"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

var mockRetryOptions = localconfig.Retry{
	ShortInterval: 50 * time.Millisecond,
	ShortTotal:    100 * time.Millisecond,
	LongInterval:  60 * time.Millisecond,
	LongTotal:     120 * time.Millisecond,
	NetworkTimeouts: localconfig.NetworkTimeouts{
		DialTimeout:  40 * time.Millisecond,
		ReadTimeout:  40 * time.Millisecond,
		WriteTimeout: 40 * time.Millisecond,
	},
	Metadata: localconfig.Metadata{
		RetryMax:     2,
		RetryBackoff: 40 * time.Millisecond,
	},
	Producer: localconfig.Producer{
		RetryMax:     2,
		RetryBackoff: 40 * time.Millisecond,
	},
	Consumer: localconfig.Consumer{
		RetryBackoff: 40 * time.Millisecond,
	},
}

func init() {
	mockLocalConfig = newMockLocalConfig(
		false,
		localconfig.SASLPlain{Enabled: false},
		mockRetryOptions,
		false)
	mockBrokerConfig = newMockBrokerConfig(
		mockLocalConfig.General.TLS,
		mockLocalConfig.Kafka.SASLPlain,
		mockLocalConfig.Kafka.Retry,
		mockLocalConfig.Kafka.Version,
		defaultPartition)
	mockConsenter = newMockConsenter(
		mockBrokerConfig,
		mockLocalConfig.General.TLS,
		mockLocalConfig.Kafka.Retry,
		mockLocalConfig.Kafka.Version)
	setupTestLogging("ERROR")
}

func TestNew(t *testing.T) {
	c, _ := New(mockLocalConfig.Kafka, &mock.MetricsProvider{}, &mock.HealthChecker{})
	_ = consensus.Consenter(c)
}

func TestHandleChain(t *testing.T) {
	consenter, _ := New(mockLocalConfig.Kafka, &disabled.Provider{}, &mock.HealthChecker{})

	oldestOffset := int64(0)
	newestOffset := int64(5)
	message := sarama.StringEncoder("messageFoo")

	mockChannel := newChannel(channelNameForTest(t), defaultPartition)

	mockBroker := sarama.NewMockBroker(t, 0)
	mockBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()).
			SetLeader(mockChannel.topic(), mockChannel.partition(), mockBroker.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t).
			SetError(mockChannel.topic(), mockChannel.partition(), sarama.ErrNoError),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(mockChannel.topic(), mockChannel.partition(), sarama.OffsetOldest, oldestOffset).
			SetOffset(mockChannel.topic(), mockChannel.partition(), sarama.OffsetNewest, newestOffset),
		"FetchRequest": sarama.NewMockFetchResponse(t, 1).
			SetMessage(mockChannel.topic(), mockChannel.partition(), newestOffset, message),
	})

	mockSupport := &mockmultichannel.ConsenterSupport{
		ChannelIDVal: mockChannel.topic(),
		SharedConfigVal: &mockconfig.Orderer{
			KafkaBrokersVal: []string{mockBroker.Addr()},
		},
	}

	mockMetadata := &cb.Metadata{Value: protoutil.MarshalOrPanic(&ab.KafkaMetadata{LastOffsetPersisted: newestOffset - 1})}
	_, err := consenter.HandleChain(mockSupport, mockMetadata)
	assert.NoError(t, err, "Expected the HandleChain call to return without errors")
}



var mockConsenter commonConsenter
var mockLocalConfig *localconfig.TopLevel
var mockBrokerConfig *sarama.Config

func extractEncodedOffset(marshalledOrdererMetadata []byte) int64 {
	omd := &cb.Metadata{}
	_ = proto.Unmarshal(marshalledOrdererMetadata, omd)
	kmd := &ab.KafkaMetadata{}
	_ = proto.Unmarshal(omd.GetValue(), kmd)
	return kmd.LastOffsetPersisted
}

func newMockBrokerConfig(
	tlsConfig localconfig.TLS,
	saslPlain localconfig.SASLPlain,
	retryOptions localconfig.Retry,
	kafkaVersion sarama.KafkaVersion,
	chosenStaticPartition int32) *sarama.Config {

	brokerConfig := newBrokerConfig(
		tlsConfig,
		saslPlain,
		retryOptions,
		kafkaVersion,
		chosenStaticPartition)
	brokerConfig.ClientID = "test"
	return brokerConfig
}

func newMockConsenter(brokerConfig *sarama.Config, tlsConfig localconfig.TLS, retryOptions localconfig.Retry, kafkaVersion sarama.KafkaVersion) *consenterImpl {
	return &consenterImpl{
		brokerConfigVal: brokerConfig,
		tlsConfigVal:    tlsConfig,
		retryOptionsVal: retryOptions,
		kafkaVersionVal: kafkaVersion,
		metrics:         NewMetrics(&disabled.Provider{}, nil),
	}
}

func newMockConsumerMessage(wrappedMessage *ab.KafkaMessage) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{
		Value: sarama.ByteEncoder(protoutil.MarshalOrPanic(wrappedMessage)),
	}
}

func newMockEnvelope(content string) *cb.Envelope {
	return &cb.Envelope{Payload: protoutil.MarshalOrPanic(&cb.Payload{
		Header: &cb.Header{ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{ChannelId: "foo"})},
		Data:   []byte(content),
	})}
}

func newMockLocalConfig(
	enableTLS bool,
	saslPlain localconfig.SASLPlain,
	retryOptions localconfig.Retry,
	verboseLog bool) *localconfig.TopLevel {

	return &localconfig.TopLevel{
		General: localconfig.General{
			TLS: localconfig.TLS{
				Enabled: enableTLS,
			},
		},
		Kafka: localconfig.Kafka{
			TLS: localconfig.TLS{
				Enabled: enableTLS,
			},
			SASLPlain: saslPlain,
			Retry:     retryOptions,
			Verbose:   verboseLog,
			Version:   sarama.V0_9_0_1, 
		},
	}
}

func setupTestLogging(logLevel string) {
	
	
	
	spec := fmt.Sprintf("orderer.consensus.kafka=%s", logLevel)
	flogging.ActivateSpec(spec)
}

func tamperBytes(original []byte) []byte {
	byteCount := len(original)
	return original[:byteCount-1]
}

func channelNameForTest(t *testing.T) string {
	return fmt.Sprintf("%s.channel", strings.Replace(strings.ToLower(t.Name()), "/", ".", -1))
}
