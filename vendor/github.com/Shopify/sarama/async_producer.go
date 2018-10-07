package sarama

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/queue"
)







type AsyncProducer interface {

	
	
	
	
	AsyncClose()

	
	
	
	
	Close() error

	
	
	Input() chan<- *ProducerMessage

	
	
	
	
	Successes() <-chan *ProducerMessage

	
	
	
	
	Errors() <-chan *ProducerError
}

type asyncProducer struct {
	client    Client
	conf      *Config
	ownClient bool

	errors                    chan *ProducerError
	input, successes, retries chan *ProducerMessage
	inFlight                  sync.WaitGroup

	brokers    map[*Broker]chan<- *ProducerMessage
	brokerRefs map[chan<- *ProducerMessage]int
	brokerLock sync.Mutex
}


func NewAsyncProducer(addrs []string, conf *Config) (AsyncProducer, error) {
	client, err := NewClient(addrs, conf)
	if err != nil {
		return nil, err
	}

	p, err := NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	p.(*asyncProducer).ownClient = true
	return p, nil
}



func NewAsyncProducerFromClient(client Client) (AsyncProducer, error) {
	
	if client.Closed() {
		return nil, ErrClosedClient
	}

	p := &asyncProducer{
		client:     client,
		conf:       client.Config(),
		errors:     make(chan *ProducerError),
		input:      make(chan *ProducerMessage),
		successes:  make(chan *ProducerMessage),
		retries:    make(chan *ProducerMessage),
		brokers:    make(map[*Broker]chan<- *ProducerMessage),
		brokerRefs: make(map[chan<- *ProducerMessage]int),
	}

	
	go withRecover(p.dispatcher)
	go withRecover(p.retryHandler)

	return p, nil
}

type flagSet int8

const (
	syn      flagSet = 1 << iota 
	fin                          
	shutdown                     
)


type ProducerMessage struct {
	Topic string 
	
	
	Key Encoder
	
	
	Value Encoder

	
	
	Headers []RecordHeader

	
	
	
	
	Metadata interface{}

	

	
	
	
	Offset int64
	
	
	Partition int32
	
	
	
	
	Timestamp time.Time

	retries     int
	flags       flagSet
	expectation chan *ProducerError
}

const producerMessageOverhead = 26 

func (m *ProducerMessage) byteSize(version int) int {
	var size int
	if version >= 2 {
		size = maximumRecordOverhead
		for _, h := range m.Headers {
			size += len(h.Key) + len(h.Value) + 2*binary.MaxVarintLen32
		}
	} else {
		size = producerMessageOverhead
	}
	if m.Key != nil {
		size += m.Key.Length()
	}
	if m.Value != nil {
		size += m.Value.Length()
	}
	return size
}

func (m *ProducerMessage) clear() {
	m.flags = 0
	m.retries = 0
}



type ProducerError struct {
	Msg *ProducerMessage
	Err error
}

func (pe ProducerError) Error() string {
	return fmt.Sprintf("kafka: Failed to produce message to topic %s: %s", pe.Msg.Topic, pe.Err)
}




type ProducerErrors []*ProducerError

func (pe ProducerErrors) Error() string {
	return fmt.Sprintf("kafka: Failed to deliver %d messages.", len(pe))
}

func (p *asyncProducer) Errors() <-chan *ProducerError {
	return p.errors
}

func (p *asyncProducer) Successes() <-chan *ProducerMessage {
	return p.successes
}

func (p *asyncProducer) Input() chan<- *ProducerMessage {
	return p.input
}

func (p *asyncProducer) Close() error {
	p.AsyncClose()

	if p.conf.Producer.Return.Successes {
		go withRecover(func() {
			for range p.successes {
			}
		})
	}

	var errors ProducerErrors
	if p.conf.Producer.Return.Errors {
		for event := range p.errors {
			errors = append(errors, event)
		}
	} else {
		<-p.errors
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (p *asyncProducer) AsyncClose() {
	go withRecover(p.shutdown)
}



func (p *asyncProducer) dispatcher() {
	handlers := make(map[string]chan<- *ProducerMessage)
	shuttingDown := false

	for msg := range p.input {
		if msg == nil {
			Logger.Println("Something tried to send a nil message, it was ignored.")
			continue
		}

		if msg.flags&shutdown != 0 {
			shuttingDown = true
			p.inFlight.Done()
			continue
		} else if msg.retries == 0 {
			if shuttingDown {
				
				
				pErr := &ProducerError{Msg: msg, Err: ErrShuttingDown}
				if p.conf.Producer.Return.Errors {
					p.errors <- pErr
				} else {
					Logger.Println(pErr)
				}
				continue
			}
			p.inFlight.Add(1)
		}

		version := 1
		if p.conf.Version.IsAtLeast(V0_11_0_0) {
			version = 2
		} else if msg.Headers != nil {
			p.returnError(msg, ConfigurationError("Producing headers requires Kafka at least v0.11"))
			continue
		}
		if msg.byteSize(version) > p.conf.Producer.MaxMessageBytes {
			p.returnError(msg, ErrMessageSizeTooLarge)
			continue
		}

		handler := handlers[msg.Topic]
		if handler == nil {
			handler = p.newTopicProducer(msg.Topic)
			handlers[msg.Topic] = handler
		}

		handler <- msg
	}

	for _, handler := range handlers {
		close(handler)
	}
}



type topicProducer struct {
	parent *asyncProducer
	topic  string
	input  <-chan *ProducerMessage

	breaker     *breaker.Breaker
	handlers    map[int32]chan<- *ProducerMessage
	partitioner Partitioner
}

func (p *asyncProducer) newTopicProducer(topic string) chan<- *ProducerMessage {
	input := make(chan *ProducerMessage, p.conf.ChannelBufferSize)
	tp := &topicProducer{
		parent:      p,
		topic:       topic,
		input:       input,
		breaker:     breaker.New(3, 1, 10*time.Second),
		handlers:    make(map[int32]chan<- *ProducerMessage),
		partitioner: p.conf.Producer.Partitioner(topic),
	}
	go withRecover(tp.dispatch)
	return input
}

func (tp *topicProducer) dispatch() {
	for msg := range tp.input {
		if msg.retries == 0 {
			if err := tp.partitionMessage(msg); err != nil {
				tp.parent.returnError(msg, err)
				continue
			}
		}

		handler := tp.handlers[msg.Partition]
		if handler == nil {
			handler = tp.parent.newPartitionProducer(msg.Topic, msg.Partition)
			tp.handlers[msg.Partition] = handler
		}

		handler <- msg
	}

	for _, handler := range tp.handlers {
		close(handler)
	}
}

func (tp *topicProducer) partitionMessage(msg *ProducerMessage) error {
	var partitions []int32

	err := tp.breaker.Run(func() (err error) {
		var requiresConsistency = false
		if ep, ok := tp.partitioner.(DynamicConsistencyPartitioner); ok {
			requiresConsistency = ep.MessageRequiresConsistency(msg)
		} else {
			requiresConsistency = tp.partitioner.RequiresConsistency()
		}

		if requiresConsistency {
			partitions, err = tp.parent.client.Partitions(msg.Topic)
		} else {
			partitions, err = tp.parent.client.WritablePartitions(msg.Topic)
		}
		return
	})

	if err != nil {
		return err
	}

	numPartitions := int32(len(partitions))

	if numPartitions == 0 {
		return ErrLeaderNotAvailable
	}

	choice, err := tp.partitioner.Partition(msg, numPartitions)

	if err != nil {
		return err
	} else if choice < 0 || choice >= numPartitions {
		return ErrInvalidPartition
	}

	msg.Partition = partitions[choice]

	return nil
}




type partitionProducer struct {
	parent    *asyncProducer
	topic     string
	partition int32
	input     <-chan *ProducerMessage

	leader  *Broker
	breaker *breaker.Breaker
	output  chan<- *ProducerMessage

	
	
	
	
	highWatermark int
	retryState    []partitionRetryState
}

type partitionRetryState struct {
	buf          []*ProducerMessage
	expectChaser bool
}

func (p *asyncProducer) newPartitionProducer(topic string, partition int32) chan<- *ProducerMessage {
	input := make(chan *ProducerMessage, p.conf.ChannelBufferSize)
	pp := &partitionProducer{
		parent:    p,
		topic:     topic,
		partition: partition,
		input:     input,

		breaker:    breaker.New(3, 1, 10*time.Second),
		retryState: make([]partitionRetryState, p.conf.Producer.Retry.Max+1),
	}
	go withRecover(pp.dispatch)
	return input
}

func (pp *partitionProducer) dispatch() {
	
	
	pp.leader, _ = pp.parent.client.Leader(pp.topic, pp.partition)
	if pp.leader != nil {
		pp.output = pp.parent.getBrokerProducer(pp.leader)
		pp.parent.inFlight.Add(1) 
		pp.output <- &ProducerMessage{Topic: pp.topic, Partition: pp.partition, flags: syn}
	}

	for msg := range pp.input {
		if msg.retries > pp.highWatermark {
			
			pp.newHighWatermark(msg.retries)
			time.Sleep(pp.parent.conf.Producer.Retry.Backoff)
		} else if pp.highWatermark > 0 {
			
			if msg.retries < pp.highWatermark {
				
				if msg.flags&fin == fin {
					pp.retryState[msg.retries].expectChaser = false
					pp.parent.inFlight.Done() 
				} else {
					pp.retryState[msg.retries].buf = append(pp.retryState[msg.retries].buf, msg)
				}
				continue
			} else if msg.flags&fin == fin {
				
				
				pp.retryState[pp.highWatermark].expectChaser = false
				pp.flushRetryBuffers()
				pp.parent.inFlight.Done() 
				continue
			}
		}

		
		

		if pp.output == nil {
			if err := pp.updateLeader(); err != nil {
				pp.parent.returnError(msg, err)
				time.Sleep(pp.parent.conf.Producer.Retry.Backoff)
				continue
			}
			Logger.Printf("producer/leader/%s/%d selected broker %d\n", pp.topic, pp.partition, pp.leader.ID())
		}

		pp.output <- msg
	}

	if pp.output != nil {
		pp.parent.unrefBrokerProducer(pp.leader, pp.output)
	}
}

func (pp *partitionProducer) newHighWatermark(hwm int) {
	Logger.Printf("producer/leader/%s/%d state change to [retrying-%d]\n", pp.topic, pp.partition, hwm)
	pp.highWatermark = hwm

	
	
	pp.retryState[pp.highWatermark].expectChaser = true
	pp.parent.inFlight.Add(1) 
	pp.output <- &ProducerMessage{Topic: pp.topic, Partition: pp.partition, flags: fin, retries: pp.highWatermark - 1}

	
	Logger.Printf("producer/leader/%s/%d abandoning broker %d\n", pp.topic, pp.partition, pp.leader.ID())
	pp.parent.unrefBrokerProducer(pp.leader, pp.output)
	pp.output = nil
}

func (pp *partitionProducer) flushRetryBuffers() {
	Logger.Printf("producer/leader/%s/%d state change to [flushing-%d]\n", pp.topic, pp.partition, pp.highWatermark)
	for {
		pp.highWatermark--

		if pp.output == nil {
			if err := pp.updateLeader(); err != nil {
				pp.parent.returnErrors(pp.retryState[pp.highWatermark].buf, err)
				goto flushDone
			}
			Logger.Printf("producer/leader/%s/%d selected broker %d\n", pp.topic, pp.partition, pp.leader.ID())
		}

		for _, msg := range pp.retryState[pp.highWatermark].buf {
			pp.output <- msg
		}

	flushDone:
		pp.retryState[pp.highWatermark].buf = nil
		if pp.retryState[pp.highWatermark].expectChaser {
			Logger.Printf("producer/leader/%s/%d state change to [retrying-%d]\n", pp.topic, pp.partition, pp.highWatermark)
			break
		} else if pp.highWatermark == 0 {
			Logger.Printf("producer/leader/%s/%d state change to [normal]\n", pp.topic, pp.partition)
			break
		}
	}
}

func (pp *partitionProducer) updateLeader() error {
	return pp.breaker.Run(func() (err error) {
		if err = pp.parent.client.RefreshMetadata(pp.topic); err != nil {
			return err
		}

		if pp.leader, err = pp.parent.client.Leader(pp.topic, pp.partition); err != nil {
			return err
		}

		pp.output = pp.parent.getBrokerProducer(pp.leader)
		pp.parent.inFlight.Add(1) 
		pp.output <- &ProducerMessage{Topic: pp.topic, Partition: pp.partition, flags: syn}

		return nil
	})
}


func (p *asyncProducer) newBrokerProducer(broker *Broker) chan<- *ProducerMessage {
	var (
		input     = make(chan *ProducerMessage)
		bridge    = make(chan *produceSet)
		responses = make(chan *brokerProducerResponse)
	)

	bp := &brokerProducer{
		parent:         p,
		broker:         broker,
		input:          input,
		output:         bridge,
		responses:      responses,
		buffer:         newProduceSet(p),
		currentRetries: make(map[string]map[int32]error),
	}
	go withRecover(bp.run)

	
	go withRecover(func() {
		for set := range bridge {
			request := set.buildRequest()

			response, err := broker.Produce(request)

			responses <- &brokerProducerResponse{
				set: set,
				err: err,
				res: response,
			}
		}
		close(responses)
	})

	return input
}

type brokerProducerResponse struct {
	set *produceSet
	err error
	res *ProduceResponse
}



type brokerProducer struct {
	parent *asyncProducer
	broker *Broker

	input     <-chan *ProducerMessage
	output    chan<- *produceSet
	responses <-chan *brokerProducerResponse

	buffer     *produceSet
	timer      <-chan time.Time
	timerFired bool

	closing        error
	currentRetries map[string]map[int32]error
}

func (bp *brokerProducer) run() {
	var output chan<- *produceSet
	Logger.Printf("producer/broker/%d starting up\n", bp.broker.ID())

	for {
		select {
		case msg := <-bp.input:
			if msg == nil {
				bp.shutdown()
				return
			}

			if msg.flags&syn == syn {
				Logger.Printf("producer/broker/%d state change to [open] on %s/%d\n",
					bp.broker.ID(), msg.Topic, msg.Partition)
				if bp.currentRetries[msg.Topic] == nil {
					bp.currentRetries[msg.Topic] = make(map[int32]error)
				}
				bp.currentRetries[msg.Topic][msg.Partition] = nil
				bp.parent.inFlight.Done()
				continue
			}

			if reason := bp.needsRetry(msg); reason != nil {
				bp.parent.retryMessage(msg, reason)

				if bp.closing == nil && msg.flags&fin == fin {
					
					delete(bp.currentRetries[msg.Topic], msg.Partition)
					Logger.Printf("producer/broker/%d state change to [closed] on %s/%d\n",
						bp.broker.ID(), msg.Topic, msg.Partition)
				}

				continue
			}

			if bp.buffer.wouldOverflow(msg) {
				if err := bp.waitForSpace(msg); err != nil {
					bp.parent.retryMessage(msg, err)
					continue
				}
			}

			if err := bp.buffer.add(msg); err != nil {
				bp.parent.returnError(msg, err)
				continue
			}

			if bp.parent.conf.Producer.Flush.Frequency > 0 && bp.timer == nil {
				bp.timer = time.After(bp.parent.conf.Producer.Flush.Frequency)
			}
		case <-bp.timer:
			bp.timerFired = true
		case output <- bp.buffer:
			bp.rollOver()
		case response := <-bp.responses:
			bp.handleResponse(response)
		}

		if bp.timerFired || bp.buffer.readyToFlush() {
			output = bp.output
		} else {
			output = nil
		}
	}
}

func (bp *brokerProducer) shutdown() {
	for !bp.buffer.empty() {
		select {
		case response := <-bp.responses:
			bp.handleResponse(response)
		case bp.output <- bp.buffer:
			bp.rollOver()
		}
	}
	close(bp.output)
	for response := range bp.responses {
		bp.handleResponse(response)
	}

	Logger.Printf("producer/broker/%d shut down\n", bp.broker.ID())
}

func (bp *brokerProducer) needsRetry(msg *ProducerMessage) error {
	if bp.closing != nil {
		return bp.closing
	}

	return bp.currentRetries[msg.Topic][msg.Partition]
}

func (bp *brokerProducer) waitForSpace(msg *ProducerMessage) error {
	Logger.Printf("producer/broker/%d maximum request accumulated, waiting for space\n", bp.broker.ID())

	for {
		select {
		case response := <-bp.responses:
			bp.handleResponse(response)
			
			if reason := bp.needsRetry(msg); reason != nil {
				return reason
			} else if !bp.buffer.wouldOverflow(msg) {
				return nil
			}
		case bp.output <- bp.buffer:
			bp.rollOver()
			return nil
		}
	}
}

func (bp *brokerProducer) rollOver() {
	bp.timer = nil
	bp.timerFired = false
	bp.buffer = newProduceSet(bp.parent)
}

func (bp *brokerProducer) handleResponse(response *brokerProducerResponse) {
	if response.err != nil {
		bp.handleError(response.set, response.err)
	} else {
		bp.handleSuccess(response.set, response.res)
	}

	if bp.buffer.empty() {
		bp.rollOver() 
	}
}

func (bp *brokerProducer) handleSuccess(sent *produceSet, response *ProduceResponse) {
	
	
	sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) {
		if response == nil {
			
			bp.parent.returnSuccesses(msgs)
			return
		}

		block := response.GetBlock(topic, partition)
		if block == nil {
			bp.parent.returnErrors(msgs, ErrIncompleteResponse)
			return
		}

		switch block.Err {
		
		case ErrNoError:
			if bp.parent.conf.Version.IsAtLeast(V0_10_0_0) && !block.Timestamp.IsZero() {
				for _, msg := range msgs {
					msg.Timestamp = block.Timestamp
				}
			}
			for i, msg := range msgs {
				msg.Offset = block.Offset + int64(i)
			}
			bp.parent.returnSuccesses(msgs)
		
		case ErrInvalidMessage, ErrUnknownTopicOrPartition, ErrLeaderNotAvailable, ErrNotLeaderForPartition,
			ErrRequestTimedOut, ErrNotEnoughReplicas, ErrNotEnoughReplicasAfterAppend:
			Logger.Printf("producer/broker/%d state change to [retrying] on %s/%d because %v\n",
				bp.broker.ID(), topic, partition, block.Err)
			bp.currentRetries[topic][partition] = block.Err
			bp.parent.retryMessages(msgs, block.Err)
			bp.parent.retryMessages(bp.buffer.dropPartition(topic, partition), block.Err)
		
		default:
			bp.parent.returnErrors(msgs, block.Err)
		}
	})
}

func (bp *brokerProducer) handleError(sent *produceSet, err error) {
	switch err.(type) {
	case PacketEncodingError:
		sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) {
			bp.parent.returnErrors(msgs, err)
		})
	default:
		Logger.Printf("producer/broker/%d state change to [closing] because %s\n", bp.broker.ID(), err)
		bp.parent.abandonBrokerConnection(bp.broker)
		_ = bp.broker.Close()
		bp.closing = err
		sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) {
			bp.parent.retryMessages(msgs, err)
		})
		bp.buffer.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) {
			bp.parent.retryMessages(msgs, err)
		})
		bp.rollOver()
	}
}




func (p *asyncProducer) retryHandler() {
	var msg *ProducerMessage
	buf := queue.New()

	for {
		if buf.Length() == 0 {
			msg = <-p.retries
		} else {
			select {
			case msg = <-p.retries:
			case p.input <- buf.Peek().(*ProducerMessage):
				buf.Remove()
				continue
			}
		}

		if msg == nil {
			return
		}

		buf.Add(msg)
	}
}



func (p *asyncProducer) shutdown() {
	Logger.Println("Producer shutting down.")
	p.inFlight.Add(1)
	p.input <- &ProducerMessage{flags: shutdown}

	p.inFlight.Wait()

	if p.ownClient {
		err := p.client.Close()
		if err != nil {
			Logger.Println("producer/shutdown failed to close the embedded client:", err)
		}
	}

	close(p.input)
	close(p.retries)
	close(p.errors)
	close(p.successes)
}

func (p *asyncProducer) returnError(msg *ProducerMessage, err error) {
	msg.clear()
	pErr := &ProducerError{Msg: msg, Err: err}
	if p.conf.Producer.Return.Errors {
		p.errors <- pErr
	} else {
		Logger.Println(pErr)
	}
	p.inFlight.Done()
}

func (p *asyncProducer) returnErrors(batch []*ProducerMessage, err error) {
	for _, msg := range batch {
		p.returnError(msg, err)
	}
}

func (p *asyncProducer) returnSuccesses(batch []*ProducerMessage) {
	for _, msg := range batch {
		if p.conf.Producer.Return.Successes {
			msg.clear()
			p.successes <- msg
		}
		p.inFlight.Done()
	}
}

func (p *asyncProducer) retryMessage(msg *ProducerMessage, err error) {
	if msg.retries >= p.conf.Producer.Retry.Max {
		p.returnError(msg, err)
	} else {
		msg.retries++
		p.retries <- msg
	}
}

func (p *asyncProducer) retryMessages(batch []*ProducerMessage, err error) {
	for _, msg := range batch {
		p.retryMessage(msg, err)
	}
}

func (p *asyncProducer) getBrokerProducer(broker *Broker) chan<- *ProducerMessage {
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	bp := p.brokers[broker]

	if bp == nil {
		bp = p.newBrokerProducer(broker)
		p.brokers[broker] = bp
		p.brokerRefs[bp] = 0
	}

	p.brokerRefs[bp]++

	return bp
}

func (p *asyncProducer) unrefBrokerProducer(broker *Broker, bp chan<- *ProducerMessage) {
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	p.brokerRefs[bp]--
	if p.brokerRefs[bp] == 0 {
		close(bp)
		delete(p.brokerRefs, bp)

		if p.brokers[broker] == bp {
			delete(p.brokers, broker)
		}
	}
}

func (p *asyncProducer) abandonBrokerConnection(broker *Broker) {
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	delete(p.brokers, broker)
}
