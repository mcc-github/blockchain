package mocks

import (
	"sync"
	"sync/atomic"

	"github.com/Shopify/sarama"
)




type Consumer struct {
	l                  sync.Mutex
	t                  ErrorReporter
	config             *sarama.Config
	partitionConsumers map[string]map[int32]*PartitionConsumer
	metadata           map[string][]int32
}




func NewConsumer(t ErrorReporter, config *sarama.Config) *Consumer {
	if config == nil {
		config = sarama.NewConfig()
	}

	c := &Consumer{
		t:                  t,
		config:             config,
		partitionConsumers: make(map[string]map[int32]*PartitionConsumer),
	}
	return c
}








func (c *Consumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.partitionConsumers[topic] == nil || c.partitionConsumers[topic][partition] == nil {
		c.t.Errorf("No expectations set for %s/%d", topic, partition)
		return nil, errOutOfExpectations
	}

	pc := c.partitionConsumers[topic][partition]
	if pc.consumed {
		return nil, sarama.ConfigurationError("The topic/partition is already being consumed")
	}

	if pc.offset != AnyOffset && pc.offset != offset {
		c.t.Errorf("Unexpected offset when calling ConsumePartition for %s/%d. Expected %d, got %d.", topic, partition, pc.offset, offset)
	}

	pc.consumed = true
	return pc, nil
}


func (c *Consumer) Topics() ([]string, error) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.metadata == nil {
		c.t.Errorf("Unexpected call to Topics. Initialize the mock's topic metadata with SetMetadata.")
		return nil, sarama.ErrOutOfBrokers
	}

	var result []string
	for topic := range c.metadata {
		result = append(result, topic)
	}
	return result, nil
}


func (c *Consumer) Partitions(topic string) ([]int32, error) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.metadata == nil {
		c.t.Errorf("Unexpected call to Partitions. Initialize the mock's topic metadata with SetMetadata.")
		return nil, sarama.ErrOutOfBrokers
	}
	if c.metadata[topic] == nil {
		return nil, sarama.ErrUnknownTopicOrPartition
	}

	return c.metadata[topic], nil
}

func (c *Consumer) HighWaterMarks() map[string]map[int32]int64 {
	c.l.Lock()
	defer c.l.Unlock()

	hwms := make(map[string]map[int32]int64, len(c.partitionConsumers))
	for topic, partitionConsumers := range c.partitionConsumers {
		hwm := make(map[int32]int64, len(partitionConsumers))
		for partition, pc := range partitionConsumers {
			hwm[partition] = pc.HighWaterMarkOffset()
		}
		hwms[topic] = hwm
	}

	return hwms
}



func (c *Consumer) Close() error {
	c.l.Lock()
	defer c.l.Unlock()

	for _, partitions := range c.partitionConsumers {
		for _, partitionConsumer := range partitions {
			_ = partitionConsumer.Close()
		}
	}

	return nil
}







func (c *Consumer) SetTopicMetadata(metadata map[string][]int32) {
	c.l.Lock()
	defer c.l.Unlock()

	c.metadata = metadata
}







func (c *Consumer) ExpectConsumePartition(topic string, partition int32, offset int64) *PartitionConsumer {
	c.l.Lock()
	defer c.l.Unlock()

	if c.partitionConsumers[topic] == nil {
		c.partitionConsumers[topic] = make(map[int32]*PartitionConsumer)
	}

	if c.partitionConsumers[topic][partition] == nil {
		c.partitionConsumers[topic][partition] = &PartitionConsumer{
			t:         c.t,
			topic:     topic,
			partition: partition,
			offset:    offset,
			messages:  make(chan *sarama.ConsumerMessage, c.config.ChannelBufferSize),
			errors:    make(chan *sarama.ConsumerError, c.config.ChannelBufferSize),
		}
	}

	return c.partitionConsumers[topic][partition]
}










type PartitionConsumer struct {
	highWaterMarkOffset     int64 
	l                       sync.Mutex
	t                       ErrorReporter
	topic                   string
	partition               int32
	offset                  int64
	messages                chan *sarama.ConsumerMessage
	errors                  chan *sarama.ConsumerError
	singleClose             sync.Once
	consumed                bool
	errorsShouldBeDrained   bool
	messagesShouldBeDrained bool
}






func (pc *PartitionConsumer) AsyncClose() {
	pc.singleClose.Do(func() {
		close(pc.messages)
		close(pc.errors)
	})
}



func (pc *PartitionConsumer) Close() error {
	if !pc.consumed {
		pc.t.Errorf("Expectations set on %s/%d, but no partition consumer was started.", pc.topic, pc.partition)
		return errPartitionConsumerNotStarted
	}

	if pc.errorsShouldBeDrained && len(pc.errors) > 0 {
		pc.t.Errorf("Expected the errors channel for %s/%d to be drained on close, but found %d errors.", pc.topic, pc.partition, len(pc.errors))
	}

	if pc.messagesShouldBeDrained && len(pc.messages) > 0 {
		pc.t.Errorf("Expected the messages channel for %s/%d to be drained on close, but found %d messages.", pc.topic, pc.partition, len(pc.messages))
	}

	pc.AsyncClose()

	var (
		closeErr error
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		var errs = make(sarama.ConsumerErrors, 0)
		for err := range pc.errors {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			closeErr = errs
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range pc.messages {
			
		}
	}()

	wg.Wait()
	return closeErr
}


func (pc *PartitionConsumer) Errors() <-chan *sarama.ConsumerError {
	return pc.errors
}


func (pc *PartitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return pc.messages
}

func (pc *PartitionConsumer) HighWaterMarkOffset() int64 {
	return atomic.LoadInt64(&pc.highWaterMarkOffset) + 1
}










func (pc *PartitionConsumer) YieldMessage(msg *sarama.ConsumerMessage) {
	pc.l.Lock()
	defer pc.l.Unlock()

	msg.Topic = pc.topic
	msg.Partition = pc.partition
	msg.Offset = atomic.AddInt64(&pc.highWaterMarkOffset, 1)

	pc.messages <- msg
}






func (pc *PartitionConsumer) YieldError(err error) {
	pc.errors <- &sarama.ConsumerError{
		Topic:     pc.topic,
		Partition: pc.partition,
		Err:       err,
	}
}




func (pc *PartitionConsumer) ExpectMessagesDrainedOnClose() {
	pc.messagesShouldBeDrained = true
}




func (pc *PartitionConsumer) ExpectErrorsDrainedOnClose() {
	pc.errorsShouldBeDrained = true
}
