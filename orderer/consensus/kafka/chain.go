/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	localconfig "github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
)


const (
	indexRecvError = iota
	indexUnmarshalError
	indexRecvPass
	indexProcessConnectPass
	indexProcessTimeToCutError
	indexProcessTimeToCutPass
	indexProcessRegularError
	indexProcessRegularPass
	indexSendTimeToCutError
	indexSendTimeToCutPass
	indexExitChanPass
)

func newChain(
	consenter commonConsenter,
	support consensus.ConsenterSupport,
	lastOffsetPersisted int64,
	lastOriginalOffsetProcessed int64,
	lastResubmittedConfigOffset int64,
) (*chainImpl, error) {
	lastCutBlockNumber := getLastCutBlockNumber(support.Height())
	logger.Infof("[channel: %s] Starting chain with last persisted offset %d and last recorded block %d",
		support.ChainID(), lastOffsetPersisted, lastCutBlockNumber)

	errorChan := make(chan struct{})
	close(errorChan) 

	doneReprocessingMsgInFlight := make(chan struct{})
	
	
	
	
	
	
	
	if lastResubmittedConfigOffset == 0 || lastResubmittedConfigOffset <= lastOriginalOffsetProcessed {
		
		close(doneReprocessingMsgInFlight)
	}

	return &chainImpl{
		consenter:                   consenter,
		ConsenterSupport:            support,
		channel:                     newChannel(support.ChainID(), defaultPartition),
		lastOffsetPersisted:         lastOffsetPersisted,
		lastOriginalOffsetProcessed: lastOriginalOffsetProcessed,
		lastResubmittedConfigOffset: lastResubmittedConfigOffset,
		lastCutBlockNumber:          lastCutBlockNumber,

		errorChan:                   errorChan,
		haltChan:                    make(chan struct{}),
		startChan:                   make(chan struct{}),
		doneReprocessingMsgInFlight: doneReprocessingMsgInFlight,
	}, nil
}

type chainImpl struct {
	consenter commonConsenter
	consensus.ConsenterSupport

	channel                     channel
	lastOffsetPersisted         int64
	lastOriginalOffsetProcessed int64
	lastResubmittedConfigOffset int64
	lastCutBlockNumber          uint64

	producer        sarama.SyncProducer
	parentConsumer  sarama.Consumer
	channelConsumer sarama.PartitionConsumer

	
	doneReprocessingMutex sync.Mutex
	
	doneReprocessingMsgInFlight chan struct{}

	
	
	errorChan chan struct{}
	
	
	
	haltChan chan struct{}
	
	doneProcessingMessagesToBlocks chan struct{}
	
	startChan chan struct{}
	
	timer <-chan time.Time
}



func (chain *chainImpl) Errored() <-chan struct{} {
	select {
	case <-chain.startChan:
		return chain.errorChan
	default:
		return nil
	}
}






func (chain *chainImpl) Start() {
	go startThread(chain)
}



func (chain *chainImpl) Halt() {
	select {
	case <-chain.startChan:
		
		select {
		case <-chain.haltChan:
			
			
			
			logger.Warningf("[channel: %s] Halting of chain requested again", chain.ChainID())
		default:
			logger.Criticalf("[channel: %s] Halting of chain requested", chain.ChainID())
			
			close(chain.haltChan)
			
			<-chain.doneProcessingMessagesToBlocks
			
			chain.closeKafkaObjects()
			logger.Debugf("[channel: %s] Closed the haltChan", chain.ChainID())
		}
	default:
		logger.Warningf("[channel: %s] Waiting for chain to finish starting before halting", chain.ChainID())
		<-chain.startChan
		chain.Halt()
	}
}

func (chain *chainImpl) WaitReady() error {
	select {
	case <-chain.startChan: 
		select {
		case <-chain.haltChan: 
			return fmt.Errorf("consenter for this channel has been halted")
		case <-chain.doneReprocessing(): 
			return nil
		}
	default: 
		return fmt.Errorf("will not enqueue, consenter for this channel hasn't started yet")
	}
}

func (chain *chainImpl) doneReprocessing() <-chan struct{} {
	chain.doneReprocessingMutex.Lock()
	defer chain.doneReprocessingMutex.Unlock()
	return chain.doneReprocessingMsgInFlight
}

func (chain *chainImpl) reprocessConfigComplete() {
	chain.doneReprocessingMutex.Lock()
	defer chain.doneReprocessingMutex.Unlock()
	close(chain.doneReprocessingMsgInFlight)
}

func (chain *chainImpl) reprocessConfigPending() {
	chain.doneReprocessingMutex.Lock()
	defer chain.doneReprocessingMutex.Unlock()
	chain.doneReprocessingMsgInFlight = make(chan struct{})
}


func (chain *chainImpl) Order(env *cb.Envelope, configSeq uint64) error {
	return chain.order(env, configSeq, int64(0))
}

func (chain *chainImpl) order(env *cb.Envelope, configSeq uint64, originalOffset int64) error {
	marshaledEnv, err := utils.Marshal(env)
	if err != nil {
		return fmt.Errorf("cannot enqueue, unable to marshal envelope because = %s", err)
	}
	if !chain.enqueue(newNormalMessage(marshaledEnv, configSeq, originalOffset)) {
		return fmt.Errorf("cannot enqueue")
	}
	return nil
}


func (chain *chainImpl) Configure(config *cb.Envelope, configSeq uint64) error {
	return chain.configure(config, configSeq, int64(0))
}

func (chain *chainImpl) configure(config *cb.Envelope, configSeq uint64, originalOffset int64) error {
	marshaledConfig, err := utils.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot enqueue, unable to marshal config because %s", err)
	}
	if !chain.enqueue(newConfigMessage(marshaledConfig, configSeq, originalOffset)) {
		return fmt.Errorf("cannot enqueue")
	}
	return nil
}


func (chain *chainImpl) enqueue(kafkaMsg *ab.KafkaMessage) bool {
	logger.Debugf("[channel: %s] Enqueueing envelope...", chain.ChainID())
	select {
	case <-chain.startChan: 
		select {
		case <-chain.haltChan: 
			logger.Warningf("[channel: %s] consenter for this channel has been halted", chain.ChainID())
			return false
		default: 
			payload, err := utils.Marshal(kafkaMsg)
			if err != nil {
				logger.Errorf("[channel: %s] unable to marshal Kafka message because = %s", chain.ChainID(), err)
				return false
			}
			message := newProducerMessage(chain.channel, payload)
			if _, _, err = chain.producer.SendMessage(message); err != nil {
				logger.Errorf("[channel: %s] cannot enqueue envelope because = %s", chain.ChainID(), err)
				return false
			}
			logger.Debugf("[channel: %s] Envelope enqueued successfully", chain.ChainID())
			return true
		}
	default: 
		logger.Warningf("[channel: %s] Will not enqueue, consenter for this channel hasn't started yet", chain.ChainID())
		return false
	}
}


func startThread(chain *chainImpl) {
	var err error

	
	err = setupTopicForChannel(chain.consenter.retryOptions(), chain.haltChan, chain.SharedConfig().KafkaBrokers(), chain.consenter.brokerConfig(), chain.consenter.topicDetail(), chain.channel)
	if err != nil {
		
		logger.Infof("[channel: %s]: failed to create Kafka topic = %s", chain.channel.topic(), err)
	}

	
	chain.producer, err = setupProducerForChannel(chain.consenter.retryOptions(), chain.haltChan, chain.SharedConfig().KafkaBrokers(), chain.consenter.brokerConfig(), chain.channel)
	if err != nil {
		logger.Panicf("[channel: %s] Cannot set up producer = %s", chain.channel.topic(), err)
	}
	logger.Infof("[channel: %s] Producer set up successfully", chain.ChainID())

	
	if err = sendConnectMessage(chain.consenter.retryOptions(), chain.haltChan, chain.producer, chain.channel); err != nil {
		logger.Panicf("[channel: %s] Cannot post CONNECT message = %s", chain.channel.topic(), err)
	}
	logger.Infof("[channel: %s] CONNECT message posted successfully", chain.channel.topic())

	
	chain.parentConsumer, err = setupParentConsumerForChannel(chain.consenter.retryOptions(), chain.haltChan, chain.SharedConfig().KafkaBrokers(), chain.consenter.brokerConfig(), chain.channel)
	if err != nil {
		logger.Panicf("[channel: %s] Cannot set up parent consumer = %s", chain.channel.topic(), err)
	}
	logger.Infof("[channel: %s] Parent consumer set up successfully", chain.channel.topic())

	
	chain.channelConsumer, err = setupChannelConsumerForChannel(chain.consenter.retryOptions(), chain.haltChan, chain.parentConsumer, chain.channel, chain.lastOffsetPersisted+1)
	if err != nil {
		logger.Panicf("[channel: %s] Cannot set up channel consumer = %s", chain.channel.topic(), err)
	}
	logger.Infof("[channel: %s] Channel consumer set up successfully", chain.channel.topic())

	chain.doneProcessingMessagesToBlocks = make(chan struct{})

	chain.errorChan = make(chan struct{}) 
	close(chain.startChan)                

	logger.Infof("[channel: %s] Start phase completed successfully", chain.channel.topic())

	chain.processMessagesToBlocks() 
}




func (chain *chainImpl) processMessagesToBlocks() ([]uint64, error) {
	counts := make([]uint64, 11) 
	msg := new(ab.KafkaMessage)

	defer func() {
		
		close(chain.doneProcessingMessagesToBlocks)
	}()

	defer func() { 
		select {
		case <-chain.errorChan: 
		default:
			close(chain.errorChan)
		}
	}()

	subscription := fmt.Sprintf("added subscription to %s/%d", chain.channel.topic(), chain.channel.partition())
	var topicPartitionSubscriptionResumed <-chan string
	var deliverSessionTimer *time.Timer
	var deliverSessionTimedOut <-chan time.Time

	for {
		select {
		case <-chain.haltChan:
			logger.Warningf("[channel: %s] Consenter for channel exiting", chain.ChainID())
			counts[indexExitChanPass]++
			return counts, nil
		case kafkaErr := <-chain.channelConsumer.Errors():
			logger.Errorf("[channel: %s] Error during consumption: %s", chain.ChainID(), kafkaErr)
			counts[indexRecvError]++
			select {
			case <-chain.errorChan: 
			default:

				switch kafkaErr.Err {
				case sarama.ErrOffsetOutOfRange:
					
					logger.Errorf("[channel: %s] Unrecoverable error during consumption: %s", chain.ChainID(), kafkaErr)
					close(chain.errorChan)
				default:
					if topicPartitionSubscriptionResumed == nil {
						
						topicPartitionSubscriptionResumed = saramaLogger.NewListener(subscription)
						
						deliverSessionTimer = time.NewTimer(chain.consenter.retryOptions().NetworkTimeouts.ReadTimeout)
						deliverSessionTimedOut = deliverSessionTimer.C
					}
				}
			}
			select {
			case <-chain.errorChan: 
				logger.Warningf("[channel: %s] Closed the errorChan", chain.ChainID())
				
				
				
				
				
				
				
				go sendConnectMessage(chain.consenter.retryOptions(), chain.haltChan, chain.producer, chain.channel)
			default: 
				logger.Warningf("[channel: %s] Deliver sessions will be dropped if consumption errors continue.", chain.ChainID())
			}
		case <-topicPartitionSubscriptionResumed:
			
			saramaLogger.RemoveListener(subscription, topicPartitionSubscriptionResumed)
			
			topicPartitionSubscriptionResumed = nil

			
			if !deliverSessionTimer.Stop() {
				<-deliverSessionTimer.C
			}
			logger.Warningf("[channel: %s] Consumption will resume.", chain.ChainID())

		case <-deliverSessionTimedOut:
			
			saramaLogger.RemoveListener(subscription, topicPartitionSubscriptionResumed)
			
			topicPartitionSubscriptionResumed = nil

			close(chain.errorChan)
			logger.Warningf("[channel: %s] Closed the errorChan", chain.ChainID())

			
			go sendConnectMessage(chain.consenter.retryOptions(), chain.haltChan, chain.producer, chain.channel)

		case in, ok := <-chain.channelConsumer.Messages():
			if !ok {
				logger.Criticalf("[channel: %s] Kafka consumer closed.", chain.ChainID())
				return counts, nil
			}

			
			
			if topicPartitionSubscriptionResumed != nil {
				
				saramaLogger.RemoveListener(subscription, topicPartitionSubscriptionResumed)
				
				topicPartitionSubscriptionResumed = nil
				
				if !deliverSessionTimer.Stop() {
					<-deliverSessionTimer.C
				}
			}

			select {
			case <-chain.errorChan: 
				chain.errorChan = make(chan struct{}) 
				logger.Infof("[channel: %s] Marked consenter as available again", chain.ChainID())
			default:
			}
			if err := proto.Unmarshal(in.Value, msg); err != nil {
				
				logger.Criticalf("[channel: %s] Unable to unmarshal consumed message = %s", chain.ChainID(), err)
				counts[indexUnmarshalError]++
				continue
			} else {
				logger.Debugf("[channel: %s] Successfully unmarshalled consumed message, offset is %d. Inspecting type...", chain.ChainID(), in.Offset)
				counts[indexRecvPass]++
			}
			switch msg.Type.(type) {
			case *ab.KafkaMessage_Connect:
				_ = chain.processConnect(chain.ChainID())
				counts[indexProcessConnectPass]++
			case *ab.KafkaMessage_TimeToCut:
				if err := chain.processTimeToCut(msg.GetTimeToCut(), in.Offset); err != nil {
					logger.Warningf("[channel: %s] %s", chain.ChainID(), err)
					logger.Criticalf("[channel: %s] Consenter for channel exiting", chain.ChainID())
					counts[indexProcessTimeToCutError]++
					return counts, err 
				}
				counts[indexProcessTimeToCutPass]++
			case *ab.KafkaMessage_Regular:
				if err := chain.processRegular(msg.GetRegular(), in.Offset); err != nil {
					logger.Warningf("[channel: %s] Error when processing incoming message of type REGULAR = %s", chain.ChainID(), err)
					counts[indexProcessRegularError]++
				} else {
					counts[indexProcessRegularPass]++
				}
			}
		case <-chain.timer:
			if err := sendTimeToCut(chain.producer, chain.channel, chain.lastCutBlockNumber+1, &chain.timer); err != nil {
				logger.Errorf("[channel: %s] cannot post time-to-cut message = %s", chain.ChainID(), err)
				
				counts[indexSendTimeToCutError]++
			} else {
				counts[indexSendTimeToCutPass]++
			}
		}
	}
}

func (chain *chainImpl) closeKafkaObjects() []error {
	var errs []error

	err := chain.channelConsumer.Close()
	if err != nil {
		logger.Errorf("[channel: %s] could not close channelConsumer cleanly = %s", chain.ChainID(), err)
		errs = append(errs, err)
	} else {
		logger.Debugf("[channel: %s] Closed the channel consumer", chain.ChainID())
	}

	err = chain.parentConsumer.Close()
	if err != nil {
		logger.Errorf("[channel: %s] could not close parentConsumer cleanly = %s", chain.ChainID(), err)
		errs = append(errs, err)
	} else {
		logger.Debugf("[channel: %s] Closed the parent consumer", chain.ChainID())
	}

	err = chain.producer.Close()
	if err != nil {
		logger.Errorf("[channel: %s] could not close producer cleanly = %s", chain.ChainID(), err)
		errs = append(errs, err)
	} else {
		logger.Debugf("[channel: %s] Closed the producer", chain.ChainID())
	}

	return errs
}



func getLastCutBlockNumber(blockchainHeight uint64) uint64 {
	return blockchainHeight - 1
}

func getOffsets(metadataValue []byte, chainID string) (persisted int64, processed int64, resubmitted int64) {
	if metadataValue != nil {
		
		kafkaMetadata := &ab.KafkaMetadata{}
		if err := proto.Unmarshal(metadataValue, kafkaMetadata); err != nil {
			logger.Panicf("[channel: %s] Ledger may be corrupted:"+
				"cannot unmarshal orderer metadata in most recent block", chainID)
		}
		return kafkaMetadata.LastOffsetPersisted,
			kafkaMetadata.LastOriginalOffsetProcessed,
			kafkaMetadata.LastResubmittedConfigOffset
	}
	return sarama.OffsetOldest - 1, int64(0), int64(0) 
}

func newConnectMessage() *ab.KafkaMessage {
	return &ab.KafkaMessage{
		Type: &ab.KafkaMessage_Connect{
			Connect: &ab.KafkaMessageConnect{
				Payload: nil,
			},
		},
	}
}

func newNormalMessage(payload []byte, configSeq uint64, originalOffset int64) *ab.KafkaMessage {
	return &ab.KafkaMessage{
		Type: &ab.KafkaMessage_Regular{
			Regular: &ab.KafkaMessageRegular{
				Payload:        payload,
				ConfigSeq:      configSeq,
				Class:          ab.KafkaMessageRegular_NORMAL,
				OriginalOffset: originalOffset,
			},
		},
	}
}

func newConfigMessage(config []byte, configSeq uint64, originalOffset int64) *ab.KafkaMessage {
	return &ab.KafkaMessage{
		Type: &ab.KafkaMessage_Regular{
			Regular: &ab.KafkaMessageRegular{
				Payload:        config,
				ConfigSeq:      configSeq,
				Class:          ab.KafkaMessageRegular_CONFIG,
				OriginalOffset: originalOffset,
			},
		},
	}
}

func newTimeToCutMessage(blockNumber uint64) *ab.KafkaMessage {
	return &ab.KafkaMessage{
		Type: &ab.KafkaMessage_TimeToCut{
			TimeToCut: &ab.KafkaMessageTimeToCut{
				BlockNumber: blockNumber,
			},
		},
	}
}

func newProducerMessage(channel channel, pld []byte) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: channel.topic(),
		Key:   sarama.StringEncoder(strconv.Itoa(int(channel.partition()))), 
		Value: sarama.ByteEncoder(pld),
	}
}

func (chain *chainImpl) processConnect(channelName string) error {
	logger.Debugf("[channel: %s] It's a connect message - ignoring", channelName)
	return nil
}

func (chain *chainImpl) processRegular(regularMessage *ab.KafkaMessageRegular, receivedOffset int64) error {
	
	
	
	
	
	
	
	commitNormalMsg := func(message *cb.Envelope, newOffset int64) {
		batches, pending := chain.BlockCutter().Ordered(message)
		logger.Debugf("[channel: %s] Ordering results: items in batch = %d, pending = %v", chain.ChainID(), len(batches), pending)
		if len(batches) == 0 {
			
			chain.lastOriginalOffsetProcessed = newOffset
			if chain.timer == nil {
				chain.timer = time.After(chain.SharedConfig().BatchTimeout())
				logger.Debugf("[channel: %s] Just began %s batch timer", chain.ChainID(), chain.SharedConfig().BatchTimeout().String())
			}
			return
		}

		chain.timer = nil

		offset := receivedOffset
		if pending || len(batches) == 2 {
			
			
			offset--
		} else {
			
			
			
			
			
			
			chain.lastOriginalOffsetProcessed = newOffset
		}

		
		block := chain.CreateNextBlock(batches[0])
		metadata := utils.MarshalOrPanic(&ab.KafkaMetadata{
			LastOffsetPersisted:         offset,
			LastOriginalOffsetProcessed: chain.lastOriginalOffsetProcessed,
			LastResubmittedConfigOffset: chain.lastResubmittedConfigOffset,
		})
		chain.WriteBlock(block, metadata)
		chain.lastCutBlockNumber++
		logger.Debugf("[channel: %s] Batch filled, just cut block %d - last persisted offset is now %d", chain.ChainID(), chain.lastCutBlockNumber, offset)

		
		if len(batches) == 2 {
			chain.lastOriginalOffsetProcessed = newOffset
			offset++

			block := chain.CreateNextBlock(batches[1])
			metadata := utils.MarshalOrPanic(&ab.KafkaMetadata{
				LastOffsetPersisted:         offset,
				LastOriginalOffsetProcessed: newOffset,
				LastResubmittedConfigOffset: chain.lastResubmittedConfigOffset,
			})
			chain.WriteBlock(block, metadata)
			chain.lastCutBlockNumber++
			logger.Debugf("[channel: %s] Batch filled, just cut block %d - last persisted offset is now %d", chain.ChainID(), chain.lastCutBlockNumber, offset)
		}
	}

	
	
	
	
	
	
	
	commitConfigMsg := func(message *cb.Envelope, newOffset int64) {
		logger.Debugf("[channel: %s] Received config message", chain.ChainID())
		batch := chain.BlockCutter().Cut()

		if batch != nil {
			logger.Debugf("[channel: %s] Cut pending messages into block", chain.ChainID())
			block := chain.CreateNextBlock(batch)
			metadata := utils.MarshalOrPanic(&ab.KafkaMetadata{
				LastOffsetPersisted:         receivedOffset - 1,
				LastOriginalOffsetProcessed: chain.lastOriginalOffsetProcessed,
				LastResubmittedConfigOffset: chain.lastResubmittedConfigOffset,
			})
			chain.WriteBlock(block, metadata)
			chain.lastCutBlockNumber++
		}

		logger.Debugf("[channel: %s] Creating isolated block for config message", chain.ChainID())
		chain.lastOriginalOffsetProcessed = newOffset
		block := chain.CreateNextBlock([]*cb.Envelope{message})
		metadata := utils.MarshalOrPanic(&ab.KafkaMetadata{
			LastOffsetPersisted:         receivedOffset,
			LastOriginalOffsetProcessed: chain.lastOriginalOffsetProcessed,
			LastResubmittedConfigOffset: chain.lastResubmittedConfigOffset,
		})
		chain.WriteConfigBlock(block, metadata)
		chain.lastCutBlockNumber++
		chain.timer = nil
	}

	seq := chain.Sequence()

	env := &cb.Envelope{}
	if err := proto.Unmarshal(regularMessage.Payload, env); err != nil {
		
		return fmt.Errorf("failed to unmarshal payload of regular message because = %s", err)
	}

	logger.Debugf("[channel: %s] Processing regular Kafka message of type %s", chain.ChainID(), regularMessage.Class.String())

	
	
	
	
	
	
	
	if regularMessage.Class == ab.KafkaMessageRegular_UNKNOWN || !chain.SharedConfig().Capabilities().Resubmission() {
		
		logger.Warningf("[channel: %s] This orderer is running in compatibility mode", chain.ChainID())

		chdr, err := utils.ChannelHeader(env)
		if err != nil {
			return fmt.Errorf("discarding bad config message because of channel header unmarshalling error = %s", err)
		}

		class := chain.ClassifyMsg(chdr)
		switch class {
		case msgprocessor.ConfigMsg:
			if _, _, err := chain.ProcessConfigMsg(env); err != nil {
				return fmt.Errorf("discarding bad config message because = %s", err)
			}

			commitConfigMsg(env, chain.lastOriginalOffsetProcessed)

		case msgprocessor.NormalMsg:
			if _, err := chain.ProcessNormalMsg(env); err != nil {
				return fmt.Errorf("discarding bad normal message because = %s", err)
			}

			commitNormalMsg(env, chain.lastOriginalOffsetProcessed)

		case msgprocessor.ConfigUpdateMsg:
			return fmt.Errorf("not expecting message of type ConfigUpdate")

		default:
			logger.Panicf("[channel: %s] Unsupported message classification: %v", chain.ChainID(), class)
		}

		return nil
	}

	switch regularMessage.Class {
	case ab.KafkaMessageRegular_UNKNOWN:
		logger.Panicf("[channel: %s] Kafka message of type UNKNOWN should have been processed already", chain.ChainID())

	case ab.KafkaMessageRegular_NORMAL:
		
		if regularMessage.OriginalOffset != 0 {
			logger.Debugf("[channel: %s] Received re-submitted normal message with original offset %d", chain.ChainID(), regularMessage.OriginalOffset)

			
			if regularMessage.OriginalOffset <= chain.lastOriginalOffsetProcessed {
				logger.Debugf(
					"[channel: %s] OriginalOffset(%d) <= LastOriginalOffsetProcessed(%d), message has been consumed already, discard",
					chain.ChainID(), regularMessage.OriginalOffset, chain.lastOriginalOffsetProcessed)
				return nil
			}

			logger.Debugf(
				"[channel: %s] OriginalOffset(%d) > LastOriginalOffsetProcessed(%d), "+
					"this is the first time we receive this re-submitted normal message",
				chain.ChainID(), regularMessage.OriginalOffset, chain.lastOriginalOffsetProcessed)

			
			
		}

		
		if regularMessage.ConfigSeq < seq {
			logger.Debugf("[channel: %s] Config sequence has advanced since this normal message got validated, re-validating", chain.ChainID())
			configSeq, err := chain.ProcessNormalMsg(env)
			if err != nil {
				return fmt.Errorf("discarding bad normal message because = %s", err)
			}

			logger.Debugf("[channel: %s] Normal message is still valid, re-submit", chain.ChainID())

			
			
			if err := chain.order(env, configSeq, receivedOffset); err != nil {
				return fmt.Errorf("error re-submitting normal message because = %s", err)
			}

			return nil
		}

		
		

		
		offset := regularMessage.OriginalOffset
		if offset == 0 {
			offset = chain.lastOriginalOffsetProcessed
		}

		commitNormalMsg(env, offset)

	case ab.KafkaMessageRegular_CONFIG:
		
		if regularMessage.OriginalOffset != 0 {
			logger.Debugf("[channel: %s] Received re-submitted config message with original offset %d", chain.ChainID(), regularMessage.OriginalOffset)

			
			if regularMessage.OriginalOffset <= chain.lastOriginalOffsetProcessed {
				logger.Debugf(
					"[channel: %s] OriginalOffset(%d) <= LastOriginalOffsetProcessed(%d), message has been consumed already, discard",
					chain.ChainID(), regularMessage.OriginalOffset, chain.lastOriginalOffsetProcessed)
				return nil
			}

			logger.Debugf(
				"[channel: %s] OriginalOffset(%d) > LastOriginalOffsetProcessed(%d), "+
					"this is the first time we receive this re-submitted config message",
				chain.ChainID(), regularMessage.OriginalOffset, chain.lastOriginalOffsetProcessed)

			if regularMessage.OriginalOffset == chain.lastResubmittedConfigOffset && 
				regularMessage.ConfigSeq == seq { 
				logger.Debugf("[channel: %s] Config message with original offset %d is the last in-flight resubmitted message"+
					"and it does not require revalidation, unblock ingress messages now", chain.ChainID(), regularMessage.OriginalOffset)
				chain.reprocessConfigComplete() 
			}

			
			
			
			
			if chain.lastResubmittedConfigOffset < regularMessage.OriginalOffset {
				chain.lastResubmittedConfigOffset = regularMessage.OriginalOffset
			}
		}

		
		if regularMessage.ConfigSeq < seq {
			logger.Debugf("[channel: %s] Config sequence has advanced since this config message got validated, re-validating", chain.ChainID())
			configEnv, configSeq, err := chain.ProcessConfigMsg(env)
			if err != nil {
				return fmt.Errorf("rejecting config message because = %s", err)
			}

			
			
			if err := chain.configure(configEnv, configSeq, receivedOffset); err != nil {
				return fmt.Errorf("error re-submitting config message because = %s", err)
			}

			logger.Debugf("[channel: %s] Resubmitted config message with offset %d, block ingress messages", chain.ChainID(), receivedOffset)
			chain.lastResubmittedConfigOffset = receivedOffset 
			chain.reprocessConfigPending()                     

			return nil
		}

		
		

		
		offset := regularMessage.OriginalOffset
		if offset == 0 {
			offset = chain.lastOriginalOffsetProcessed
		}

		commitConfigMsg(env, offset)

	default:
		return fmt.Errorf("unsupported regular kafka message type: %v", regularMessage.Class.String())
	}

	return nil
}

func (chain *chainImpl) processTimeToCut(ttcMessage *ab.KafkaMessageTimeToCut, receivedOffset int64) error {
	ttcNumber := ttcMessage.GetBlockNumber()
	logger.Debugf("[channel: %s] It's a time-to-cut message for block %d", chain.ChainID(), ttcNumber)
	if ttcNumber == chain.lastCutBlockNumber+1 {
		chain.timer = nil
		logger.Debugf("[channel: %s] Nil'd the timer", chain.ChainID())
		batch := chain.BlockCutter().Cut()
		if len(batch) == 0 {
			return fmt.Errorf("got right time-to-cut message (for block %d),"+
				" no pending requests though; this might indicate a bug", chain.lastCutBlockNumber+1)
		}
		block := chain.CreateNextBlock(batch)
		metadata := utils.MarshalOrPanic(&ab.KafkaMetadata{
			LastOffsetPersisted:         receivedOffset,
			LastOriginalOffsetProcessed: chain.lastOriginalOffsetProcessed,
		})
		chain.WriteBlock(block, metadata)
		chain.lastCutBlockNumber++
		logger.Debugf("[channel: %s] Proper time-to-cut received, just cut block %d", chain.ChainID(), chain.lastCutBlockNumber)
		return nil
	} else if ttcNumber > chain.lastCutBlockNumber+1 {
		return fmt.Errorf("got larger time-to-cut message (%d) than allowed/expected (%d)"+
			" - this might indicate a bug", ttcNumber, chain.lastCutBlockNumber+1)
	}
	logger.Debugf("[channel: %s] Ignoring stale time-to-cut-message for block %d", chain.ChainID(), ttcNumber)
	return nil
}




func sendConnectMessage(retryOptions localconfig.Retry, exitChan chan struct{}, producer sarama.SyncProducer, channel channel) error {
	logger.Infof("[channel: %s] About to post the CONNECT message...", channel.topic())

	payload := utils.MarshalOrPanic(newConnectMessage())
	message := newProducerMessage(channel, payload)

	retryMsg := "Attempting to post the CONNECT message..."
	postConnect := newRetryProcess(retryOptions, exitChan, channel, retryMsg, func() error {
		select {
		case <-exitChan:
			logger.Debugf("[channel: %s] Consenter for channel exiting, aborting retry", channel)
			return nil
		default:
			_, _, err := producer.SendMessage(message)
			return err
		}
	})

	return postConnect.retry()
}

func sendTimeToCut(producer sarama.SyncProducer, channel channel, timeToCutBlockNumber uint64, timer *<-chan time.Time) error {
	logger.Debugf("[channel: %s] Time-to-cut block %d timer expired", channel.topic(), timeToCutBlockNumber)
	*timer = nil
	payload := utils.MarshalOrPanic(newTimeToCutMessage(timeToCutBlockNumber))
	message := newProducerMessage(channel, payload)
	_, _, err := producer.SendMessage(message)
	return err
}


func setupChannelConsumerForChannel(retryOptions localconfig.Retry, haltChan chan struct{}, parentConsumer sarama.Consumer, channel channel, startFrom int64) (sarama.PartitionConsumer, error) {
	var err error
	var channelConsumer sarama.PartitionConsumer

	logger.Infof("[channel: %s] Setting up the channel consumer for this channel (start offset: %d)...", channel.topic(), startFrom)

	retryMsg := "Connecting to the Kafka cluster"
	setupChannelConsumer := newRetryProcess(retryOptions, haltChan, channel, retryMsg, func() error {
		channelConsumer, err = parentConsumer.ConsumePartition(channel.topic(), channel.partition(), startFrom)
		return err
	})

	return channelConsumer, setupChannelConsumer.retry()
}


func setupParentConsumerForChannel(retryOptions localconfig.Retry, haltChan chan struct{}, brokers []string, brokerConfig *sarama.Config, channel channel) (sarama.Consumer, error) {
	var err error
	var parentConsumer sarama.Consumer

	logger.Infof("[channel: %s] Setting up the parent consumer for this channel...", channel.topic())

	retryMsg := "Connecting to the Kafka cluster"
	setupParentConsumer := newRetryProcess(retryOptions, haltChan, channel, retryMsg, func() error {
		parentConsumer, err = sarama.NewConsumer(brokers, brokerConfig)
		return err
	})

	return parentConsumer, setupParentConsumer.retry()
}


func setupProducerForChannel(retryOptions localconfig.Retry, haltChan chan struct{}, brokers []string, brokerConfig *sarama.Config, channel channel) (sarama.SyncProducer, error) {
	var err error
	var producer sarama.SyncProducer

	logger.Infof("[channel: %s] Setting up the producer for this channel...", channel.topic())

	retryMsg := "Connecting to the Kafka cluster"
	setupProducer := newRetryProcess(retryOptions, haltChan, channel, retryMsg, func() error {
		producer, err = sarama.NewSyncProducer(brokers, brokerConfig)
		return err
	})

	return producer, setupProducer.retry()
}


func setupTopicForChannel(retryOptions localconfig.Retry, haltChan chan struct{}, brokers []string, brokerConfig *sarama.Config, topicDetail *sarama.TopicDetail, channel channel) error {

	
	if !brokerConfig.Version.IsAtLeast(sarama.V0_10_1_0) {
		return nil
	}

	logger.Infof("[channel: %s] Setting up the topic for this channel...",
		channel.topic())

	retryMsg := fmt.Sprintf("Creating Kafka topic [%s] for channel [%s]",
		channel.topic(), channel.String())

	setupTopic := newRetryProcess(
		retryOptions,
		haltChan,
		channel,
		retryMsg,
		func() error {

			var err error
			clusterMembers := map[int32]*sarama.Broker{}
			var controllerId int32

			
			for _, address := range brokers {
				broker := sarama.NewBroker(address)
				err = broker.Open(brokerConfig)

				if err != nil {
					continue
				}

				var ok bool
				ok, err = broker.Connected()
				if !ok {
					continue
				}
				defer broker.Close()

				
				var apiVersion int16
				if brokerConfig.Version.IsAtLeast(sarama.V0_11_0_0) {
					
					
					apiVersion = 4
				} else {
					apiVersion = 1
				}
				metadata, err := broker.GetMetadata(&sarama.MetadataRequest{
					Version:                apiVersion,
					Topics:                 []string{channel.topic()},
					AllowAutoTopicCreation: false})

				if err != nil {
					continue
				}

				controllerId = metadata.ControllerID
				for _, broker := range metadata.Brokers {
					clusterMembers[broker.ID()] = broker
				}

				for _, topic := range metadata.Topics {
					if topic.Name == channel.topic() {
						if topic.Err != sarama.ErrUnknownTopicOrPartition {
							
							return nil
						}
					}
				}
				break
			}

			
			if len(clusterMembers) == 0 {
				return fmt.Errorf(
					"error creating topic [%s]; failed to retrieve metadata for the cluster",
					channel.topic())
			}

			
			controller := clusterMembers[controllerId]
			err = controller.Open(brokerConfig)

			if err != nil {
				return err
			}

			var ok bool
			ok, err = controller.Connected()
			if !ok {
				return err
			}
			defer controller.Close()

			
			req := &sarama.CreateTopicsRequest{
				Version: 0,
				TopicDetails: map[string]*sarama.TopicDetail{
					channel.topic(): topicDetail},
				Timeout: 3 * time.Second}
			resp := &sarama.CreateTopicsResponse{}
			resp, err = controller.CreateTopics(req)
			if err != nil {
				return err
			}

			
			if topicErr, ok := resp.TopicErrors[channel.topic()]; ok {
				
				if topicErr.Err == sarama.ErrNoError ||
					topicErr.Err == sarama.ErrTopicAlreadyExists {
					return nil
				}
				if topicErr.Err == sarama.ErrInvalidTopic {
					
					logger.Warningf("[channel: %s] Failed to set up topic = %s",
						channel.topic(), topicErr.Err.Error())
					go func() {
						haltChan <- struct{}{}
					}()
				}
				return fmt.Errorf("error creating topic: [%s]",
					topicErr.Err.Error())
			}

			return nil
		})

	return setupTopic.retry()
}
