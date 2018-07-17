package sarama

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)


type ConsumerMessage struct {
	Key, Value     []byte
	Topic          string
	Partition      int32
	Offset         int64
	Timestamp      time.Time       
	BlockTimestamp time.Time       
	Headers        []*RecordHeader 
}



type ConsumerError struct {
	Topic     string
	Partition int32
	Err       error
}

func (ce ConsumerError) Error() string {
	return fmt.Sprintf("kafka: error while consuming %s/%d: %s", ce.Topic, ce.Partition, ce.Err)
}




type ConsumerErrors []*ConsumerError

func (ce ConsumerErrors) Error() string {
	return fmt.Sprintf("kafka: %d errors while consuming", len(ce))
}









type Consumer interface {

	
	
	
	Topics() ([]string, error)

	
	
	Partitions(topic string) ([]int32, error)

	
	
	
	
	ConsumePartition(topic string, partition int32, offset int64) (PartitionConsumer, error)

	
	
	HighWaterMarks() map[string]map[int32]int64

	
	
	Close() error
}

type consumer struct {
	client    Client
	conf      *Config
	ownClient bool

	lock            sync.Mutex
	children        map[string]map[int32]*partitionConsumer
	brokerConsumers map[*Broker]*brokerConsumer
}


func NewConsumer(addrs []string, config *Config) (Consumer, error) {
	client, err := NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	c.(*consumer).ownClient = true
	return c, nil
}



func NewConsumerFromClient(client Client) (Consumer, error) {
	
	if client.Closed() {
		return nil, ErrClosedClient
	}

	c := &consumer{
		client:          client,
		conf:            client.Config(),
		children:        make(map[string]map[int32]*partitionConsumer),
		brokerConsumers: make(map[*Broker]*brokerConsumer),
	}

	return c, nil
}

func (c *consumer) Close() error {
	if c.ownClient {
		return c.client.Close()
	}
	return nil
}

func (c *consumer) Topics() ([]string, error) {
	return c.client.Topics()
}

func (c *consumer) Partitions(topic string) ([]int32, error) {
	return c.client.Partitions(topic)
}

func (c *consumer) ConsumePartition(topic string, partition int32, offset int64) (PartitionConsumer, error) {
	child := &partitionConsumer{
		consumer:  c,
		conf:      c.conf,
		topic:     topic,
		partition: partition,
		messages:  make(chan *ConsumerMessage, c.conf.ChannelBufferSize),
		errors:    make(chan *ConsumerError, c.conf.ChannelBufferSize),
		feeder:    make(chan *FetchResponse, 1),
		trigger:   make(chan none, 1),
		dying:     make(chan none),
		fetchSize: c.conf.Consumer.Fetch.Default,
	}

	if err := child.chooseStartingOffset(offset); err != nil {
		return nil, err
	}

	var leader *Broker
	var err error
	if leader, err = c.client.Leader(child.topic, child.partition); err != nil {
		return nil, err
	}

	if err := c.addChild(child); err != nil {
		return nil, err
	}

	go withRecover(child.dispatcher)
	go withRecover(child.responseFeeder)

	child.broker = c.refBrokerConsumer(leader)
	child.broker.input <- child

	return child, nil
}

func (c *consumer) HighWaterMarks() map[string]map[int32]int64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	hwms := make(map[string]map[int32]int64)
	for topic, p := range c.children {
		hwm := make(map[int32]int64, len(p))
		for partition, pc := range p {
			hwm[partition] = pc.HighWaterMarkOffset()
		}
		hwms[topic] = hwm
	}

	return hwms
}

func (c *consumer) addChild(child *partitionConsumer) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	topicChildren := c.children[child.topic]
	if topicChildren == nil {
		topicChildren = make(map[int32]*partitionConsumer)
		c.children[child.topic] = topicChildren
	}

	if topicChildren[child.partition] != nil {
		return ConfigurationError("That topic/partition is already being consumed")
	}

	topicChildren[child.partition] = child
	return nil
}

func (c *consumer) removeChild(child *partitionConsumer) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.children[child.topic], child.partition)
}

func (c *consumer) refBrokerConsumer(broker *Broker) *brokerConsumer {
	c.lock.Lock()
	defer c.lock.Unlock()

	bc := c.brokerConsumers[broker]
	if bc == nil {
		bc = c.newBrokerConsumer(broker)
		c.brokerConsumers[broker] = bc
	}

	bc.refs++

	return bc
}

func (c *consumer) unrefBrokerConsumer(brokerWorker *brokerConsumer) {
	c.lock.Lock()
	defer c.lock.Unlock()

	brokerWorker.refs--

	if brokerWorker.refs == 0 {
		close(brokerWorker.input)
		if c.brokerConsumers[brokerWorker.broker] == brokerWorker {
			delete(c.brokerConsumers, brokerWorker.broker)
		}
	}
}

func (c *consumer) abandonBrokerConsumer(brokerWorker *brokerConsumer) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.brokerConsumers, brokerWorker.broker)
}




















type PartitionConsumer interface {

	
	
	
	
	AsyncClose()

	
	
	
	
	
	Close() error

	
	
	Messages() <-chan *ConsumerMessage

	
	
	
	
	Errors() <-chan *ConsumerError

	
	
	
	HighWaterMarkOffset() int64
}

type partitionConsumer struct {
	highWaterMarkOffset int64 
	consumer            *consumer
	conf                *Config
	topic               string
	partition           int32

	broker   *brokerConsumer
	messages chan *ConsumerMessage
	errors   chan *ConsumerError
	feeder   chan *FetchResponse

	trigger, dying chan none
	responseResult error

	fetchSize int32
	offset    int64
}

var errTimedOut = errors.New("timed out feeding messages to the user") 

func (child *partitionConsumer) sendError(err error) {
	cErr := &ConsumerError{
		Topic:     child.topic,
		Partition: child.partition,
		Err:       err,
	}

	if child.conf.Consumer.Return.Errors {
		child.errors <- cErr
	} else {
		Logger.Println(cErr)
	}
}

func (child *partitionConsumer) dispatcher() {
	for range child.trigger {
		select {
		case <-child.dying:
			close(child.trigger)
		case <-time.After(child.conf.Consumer.Retry.Backoff):
			if child.broker != nil {
				child.consumer.unrefBrokerConsumer(child.broker)
				child.broker = nil
			}

			Logger.Printf("consumer/%s/%d finding new broker\n", child.topic, child.partition)
			if err := child.dispatch(); err != nil {
				child.sendError(err)
				child.trigger <- none{}
			}
		}
	}

	if child.broker != nil {
		child.consumer.unrefBrokerConsumer(child.broker)
	}
	child.consumer.removeChild(child)
	close(child.feeder)
}

func (child *partitionConsumer) dispatch() error {
	if err := child.consumer.client.RefreshMetadata(child.topic); err != nil {
		return err
	}

	var leader *Broker
	var err error
	if leader, err = child.consumer.client.Leader(child.topic, child.partition); err != nil {
		return err
	}

	child.broker = child.consumer.refBrokerConsumer(leader)

	child.broker.input <- child

	return nil
}

func (child *partitionConsumer) chooseStartingOffset(offset int64) error {
	newestOffset, err := child.consumer.client.GetOffset(child.topic, child.partition, OffsetNewest)
	if err != nil {
		return err
	}
	oldestOffset, err := child.consumer.client.GetOffset(child.topic, child.partition, OffsetOldest)
	if err != nil {
		return err
	}

	switch {
	case offset == OffsetNewest:
		child.offset = newestOffset
	case offset == OffsetOldest:
		child.offset = oldestOffset
	case offset >= oldestOffset && offset <= newestOffset:
		child.offset = offset
	default:
		return ErrOffsetOutOfRange
	}

	return nil
}

func (child *partitionConsumer) Messages() <-chan *ConsumerMessage {
	return child.messages
}

func (child *partitionConsumer) Errors() <-chan *ConsumerError {
	return child.errors
}

func (child *partitionConsumer) AsyncClose() {
	
	
	
	
	close(child.dying)
}

func (child *partitionConsumer) Close() error {
	child.AsyncClose()

	go withRecover(func() {
		for range child.messages {
			
		}
	})

	var errors ConsumerErrors
	for err := range child.errors {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (child *partitionConsumer) HighWaterMarkOffset() int64 {
	return atomic.LoadInt64(&child.highWaterMarkOffset)
}

func (child *partitionConsumer) responseFeeder() {
	var msgs []*ConsumerMessage
	expiryTicker := time.NewTicker(child.conf.Consumer.MaxProcessingTime)
	firstAttempt := true

feederLoop:
	for response := range child.feeder {
		msgs, child.responseResult = child.parseResponse(response)

		for i, msg := range msgs {
		messageSelect:
			select {
			case child.messages <- msg:
				firstAttempt = true
			case <-expiryTicker.C:
				if !firstAttempt {
					child.responseResult = errTimedOut
					child.broker.acks.Done()
					for _, msg = range msgs[i:] {
						child.messages <- msg
					}
					child.broker.input <- child
					expiryTicker.Stop()
					continue feederLoop
				} else {
					
					
					firstAttempt = false
					goto messageSelect
				}
			}
		}

		child.broker.acks.Done()
	}

	expiryTicker.Stop()
	close(child.messages)
	close(child.errors)
}

func (child *partitionConsumer) parseMessages(msgSet *MessageSet) ([]*ConsumerMessage, error) {
	var messages []*ConsumerMessage
	var incomplete bool
	prelude := true

	for _, msgBlock := range msgSet.Messages {
		for _, msg := range msgBlock.Messages() {
			offset := msg.Offset
			if msg.Msg.Version >= 1 {
				baseOffset := msgBlock.Offset - msgBlock.Messages()[len(msgBlock.Messages())-1].Offset
				offset += baseOffset
			}
			if prelude && offset < child.offset {
				continue
			}
			prelude = false

			if offset >= child.offset {
				messages = append(messages, &ConsumerMessage{
					Topic:          child.topic,
					Partition:      child.partition,
					Key:            msg.Msg.Key,
					Value:          msg.Msg.Value,
					Offset:         offset,
					Timestamp:      msg.Msg.Timestamp,
					BlockTimestamp: msgBlock.Msg.Timestamp,
				})
				child.offset = offset + 1
			} else {
				incomplete = true
			}
		}
	}

	if incomplete || len(messages) == 0 {
		return nil, ErrIncompleteResponse
	}
	return messages, nil
}

func (child *partitionConsumer) parseRecords(batch *RecordBatch) ([]*ConsumerMessage, error) {
	var messages []*ConsumerMessage
	var incomplete bool
	prelude := true
	originalOffset := child.offset

	for _, rec := range batch.Records {
		offset := batch.FirstOffset + rec.OffsetDelta
		if prelude && offset < child.offset {
			continue
		}
		prelude = false

		if offset >= child.offset {
			messages = append(messages, &ConsumerMessage{
				Topic:     child.topic,
				Partition: child.partition,
				Key:       rec.Key,
				Value:     rec.Value,
				Offset:    offset,
				Timestamp: batch.FirstTimestamp.Add(rec.TimestampDelta),
				Headers:   rec.Headers,
			})
			child.offset = offset + 1
		} else {
			incomplete = true
		}
	}

	if incomplete {
		return nil, ErrIncompleteResponse
	}

	child.offset = batch.FirstOffset + int64(batch.LastOffsetDelta) + 1
	if child.offset <= originalOffset {
		return nil, ErrConsumerOffsetNotAdvanced
	}

	return messages, nil
}

func (child *partitionConsumer) parseResponse(response *FetchResponse) ([]*ConsumerMessage, error) {
	block := response.GetBlock(child.topic, child.partition)
	if block == nil {
		return nil, ErrIncompleteResponse
	}

	if block.Err != ErrNoError {
		return nil, block.Err
	}

	nRecs, err := block.numRecords()
	if err != nil {
		return nil, err
	}
	if nRecs == 0 {
		partialTrailingMessage, err := block.isPartial()
		if err != nil {
			return nil, err
		}
		
		
		if partialTrailingMessage {
			if child.conf.Consumer.Fetch.Max > 0 && child.fetchSize == child.conf.Consumer.Fetch.Max {
				
				child.sendError(ErrMessageTooLarge)
				child.offset++ 
			} else {
				child.fetchSize *= 2
				if child.conf.Consumer.Fetch.Max > 0 && child.fetchSize > child.conf.Consumer.Fetch.Max {
					child.fetchSize = child.conf.Consumer.Fetch.Max
				}
			}
		}

		return nil, nil
	}

	
	child.fetchSize = child.conf.Consumer.Fetch.Default
	atomic.StoreInt64(&child.highWaterMarkOffset, block.HighWaterMarkOffset)

	messages := []*ConsumerMessage{}
	for _, records := range block.RecordsSet {
		if control, err := records.isControl(); err != nil || control {
			continue
		}

		switch records.recordsType {
		case legacyRecords:
			messageSetMessages, err := child.parseMessages(records.msgSet)
			if err != nil {
				return nil, err
			}

			messages = append(messages, messageSetMessages...)
		case defaultRecords:
			recordBatchMessages, err := child.parseRecords(records.recordBatch)
			if err != nil {
				return nil, err
			}

			messages = append(messages, recordBatchMessages...)
		default:
			return nil, fmt.Errorf("unknown records type: %v", records.recordsType)
		}
	}

	return messages, nil
}



type brokerConsumer struct {
	consumer         *consumer
	broker           *Broker
	input            chan *partitionConsumer
	newSubscriptions chan []*partitionConsumer
	wait             chan none
	subscriptions    map[*partitionConsumer]none
	acks             sync.WaitGroup
	refs             int
}

func (c *consumer) newBrokerConsumer(broker *Broker) *brokerConsumer {
	bc := &brokerConsumer{
		consumer:         c,
		broker:           broker,
		input:            make(chan *partitionConsumer),
		newSubscriptions: make(chan []*partitionConsumer),
		wait:             make(chan none),
		subscriptions:    make(map[*partitionConsumer]none),
		refs:             0,
	}

	go withRecover(bc.subscriptionManager)
	go withRecover(bc.subscriptionConsumer)

	return bc
}

func (bc *brokerConsumer) subscriptionManager() {
	var buffer []*partitionConsumer

	
	
	
	
	
	for {
		if len(buffer) > 0 {
			select {
			case event, ok := <-bc.input:
				if !ok {
					goto done
				}
				buffer = append(buffer, event)
			case bc.newSubscriptions <- buffer:
				buffer = nil
			case bc.wait <- none{}:
			}
		} else {
			select {
			case event, ok := <-bc.input:
				if !ok {
					goto done
				}
				buffer = append(buffer, event)
			case bc.newSubscriptions <- nil:
			}
		}
	}

done:
	close(bc.wait)
	if len(buffer) > 0 {
		bc.newSubscriptions <- buffer
	}
	close(bc.newSubscriptions)
}

func (bc *brokerConsumer) subscriptionConsumer() {
	<-bc.wait 

	
	for newSubscriptions := range bc.newSubscriptions {
		bc.updateSubscriptions(newSubscriptions)

		if len(bc.subscriptions) == 0 {
			
			
			<-bc.wait
			continue
		}

		response, err := bc.fetchNewMessages()

		if err != nil {
			Logger.Printf("consumer/broker/%d disconnecting due to error processing FetchRequest: %s\n", bc.broker.ID(), err)
			bc.abort(err)
			return
		}

		bc.acks.Add(len(bc.subscriptions))
		for child := range bc.subscriptions {
			child.feeder <- response
		}
		bc.acks.Wait()
		bc.handleResponses()
	}
}

func (bc *brokerConsumer) updateSubscriptions(newSubscriptions []*partitionConsumer) {
	for _, child := range newSubscriptions {
		bc.subscriptions[child] = none{}
		Logger.Printf("consumer/broker/%d added subscription to %s/%d\n", bc.broker.ID(), child.topic, child.partition)
	}

	for child := range bc.subscriptions {
		select {
		case <-child.dying:
			Logger.Printf("consumer/broker/%d closed dead subscription to %s/%d\n", bc.broker.ID(), child.topic, child.partition)
			close(child.trigger)
			delete(bc.subscriptions, child)
		default:
			break
		}
	}
}

func (bc *brokerConsumer) handleResponses() {
	
	for child := range bc.subscriptions {
		result := child.responseResult
		child.responseResult = nil

		switch result {
		case nil:
			break
		case errTimedOut:
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because consuming was taking too long\n",
				bc.broker.ID(), child.topic, child.partition)
			delete(bc.subscriptions, child)
		case ErrOffsetOutOfRange:
			
			
			child.sendError(result)
			Logger.Printf("consumer/%s/%d shutting down because %s\n", child.topic, child.partition, result)
			close(child.trigger)
			delete(bc.subscriptions, child)
		case ErrUnknownTopicOrPartition, ErrNotLeaderForPartition, ErrLeaderNotAvailable, ErrReplicaNotAvailable:
			
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because %s\n",
				bc.broker.ID(), child.topic, child.partition, result)
			child.trigger <- none{}
			delete(bc.subscriptions, child)
		default:
			
			child.sendError(result)
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because %s\n",
				bc.broker.ID(), child.topic, child.partition, result)
			child.trigger <- none{}
			delete(bc.subscriptions, child)
		}
	}
}

func (bc *brokerConsumer) abort(err error) {
	bc.consumer.abandonBrokerConsumer(bc)
	_ = bc.broker.Close() 

	for child := range bc.subscriptions {
		child.sendError(err)
		child.trigger <- none{}
	}

	for newSubscriptions := range bc.newSubscriptions {
		if len(newSubscriptions) == 0 {
			<-bc.wait
			continue
		}
		for _, child := range newSubscriptions {
			child.sendError(err)
			child.trigger <- none{}
		}
	}
}

func (bc *brokerConsumer) fetchNewMessages() (*FetchResponse, error) {
	request := &FetchRequest{
		MinBytes:    bc.consumer.conf.Consumer.Fetch.Min,
		MaxWaitTime: int32(bc.consumer.conf.Consumer.MaxWaitTime / time.Millisecond),
	}
	if bc.consumer.conf.Version.IsAtLeast(V0_10_0_0) {
		request.Version = 2
	}
	if bc.consumer.conf.Version.IsAtLeast(V0_10_1_0) {
		request.Version = 3
		request.MaxBytes = MaxResponseSize
	}
	if bc.consumer.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 4
		request.Isolation = ReadUncommitted 
	}

	for child := range bc.subscriptions {
		request.AddBlock(child.topic, child.partition, child.offset, child.fetchSize)
	}

	return bc.broker.Fetch(request)
}
