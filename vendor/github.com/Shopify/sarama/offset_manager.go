package sarama

import (
	"sync"
	"time"
)




type OffsetManager interface {
	
	
	
	ManagePartition(topic string, partition int32) (PartitionOffsetManager, error)

	
	
	
	
	Close() error
}

type offsetManager struct {
	client Client
	conf   *Config
	group  string
	ticker *time.Ticker

	memberID   string
	generation int32

	broker     *Broker
	brokerLock sync.RWMutex

	poms     map[string]map[int32]*partitionOffsetManager
	pomsLock sync.RWMutex

	closeOnce sync.Once
	closing   chan none
	closed    chan none
}



func NewOffsetManagerFromClient(group string, client Client) (OffsetManager, error) {
	return newOffsetManagerFromClient(group, "", GroupGenerationUndefined, client)
}

func newOffsetManagerFromClient(group, memberID string, generation int32, client Client) (*offsetManager, error) {
	
	if client.Closed() {
		return nil, ErrClosedClient
	}

	conf := client.Config()
	om := &offsetManager{
		client: client,
		conf:   conf,
		group:  group,
		ticker: time.NewTicker(conf.Consumer.Offsets.CommitInterval),
		poms:   make(map[string]map[int32]*partitionOffsetManager),

		memberID:   memberID,
		generation: generation,

		closing: make(chan none),
		closed:  make(chan none),
	}
	go withRecover(om.mainLoop)

	return om, nil
}

func (om *offsetManager) ManagePartition(topic string, partition int32) (PartitionOffsetManager, error) {
	pom, err := om.newPartitionOffsetManager(topic, partition)
	if err != nil {
		return nil, err
	}

	om.pomsLock.Lock()
	defer om.pomsLock.Unlock()

	topicManagers := om.poms[topic]
	if topicManagers == nil {
		topicManagers = make(map[int32]*partitionOffsetManager)
		om.poms[topic] = topicManagers
	}

	if topicManagers[partition] != nil {
		return nil, ConfigurationError("That topic/partition is already being managed")
	}

	topicManagers[partition] = pom
	return pom, nil
}

func (om *offsetManager) Close() error {
	om.closeOnce.Do(func() {
		
		close(om.closing)
		<-om.closed

		
		om.asyncClosePOMs()

		
		for attempt := 0; attempt <= om.conf.Consumer.Offsets.Retry.Max; attempt++ {
			om.flushToBroker()
			if om.releasePOMs(false) == 0 {
				break
			}
		}

		om.releasePOMs(true)
		om.brokerLock.Lock()
		om.broker = nil
		om.brokerLock.Unlock()
	})
	return nil
}

func (om *offsetManager) fetchInitialOffset(topic string, partition int32, retries int) (int64, string, error) {
	broker, err := om.coordinator()
	if err != nil {
		if retries <= 0 {
			return 0, "", err
		}
		return om.fetchInitialOffset(topic, partition, retries-1)
	}

	req := new(OffsetFetchRequest)
	req.Version = 1
	req.ConsumerGroup = om.group
	req.AddPartition(topic, partition)

	resp, err := broker.FetchOffset(req)
	if err != nil {
		if retries <= 0 {
			return 0, "", err
		}
		om.releaseCoordinator(broker)
		return om.fetchInitialOffset(topic, partition, retries-1)
	}

	block := resp.GetBlock(topic, partition)
	if block == nil {
		return 0, "", ErrIncompleteResponse
	}

	switch block.Err {
	case ErrNoError:
		return block.Offset, block.Metadata, nil
	case ErrNotCoordinatorForConsumer:
		if retries <= 0 {
			return 0, "", block.Err
		}
		om.releaseCoordinator(broker)
		return om.fetchInitialOffset(topic, partition, retries-1)
	case ErrOffsetsLoadInProgress:
		if retries <= 0 {
			return 0, "", block.Err
		}
		select {
		case <-om.closing:
			return 0, "", block.Err
		case <-time.After(om.conf.Metadata.Retry.Backoff):
		}
		return om.fetchInitialOffset(topic, partition, retries-1)
	default:
		return 0, "", block.Err
	}
}

func (om *offsetManager) coordinator() (*Broker, error) {
	om.brokerLock.RLock()
	broker := om.broker
	om.brokerLock.RUnlock()

	if broker != nil {
		return broker, nil
	}

	om.brokerLock.Lock()
	defer om.brokerLock.Unlock()

	if broker := om.broker; broker != nil {
		return broker, nil
	}

	if err := om.client.RefreshCoordinator(om.group); err != nil {
		return nil, err
	}

	broker, err := om.client.Coordinator(om.group)
	if err != nil {
		return nil, err
	}

	om.broker = broker
	return broker, nil
}

func (om *offsetManager) releaseCoordinator(b *Broker) {
	om.brokerLock.Lock()
	if om.broker == b {
		om.broker = nil
	}
	om.brokerLock.Unlock()
}

func (om *offsetManager) mainLoop() {
	defer om.ticker.Stop()
	defer close(om.closed)

	for {
		select {
		case <-om.ticker.C:
			om.flushToBroker()
			om.releasePOMs(false)
		case <-om.closing:
			return
		}
	}
}

func (om *offsetManager) flushToBroker() {
	req := om.constructRequest()
	if req == nil {
		return
	}

	broker, err := om.coordinator()
	if err != nil {
		om.handleError(err)
		return
	}

	resp, err := broker.CommitOffset(req)
	if err != nil {
		om.handleError(err)
		om.releaseCoordinator(broker)
		_ = broker.Close()
		return
	}

	om.handleResponse(broker, req, resp)
}

func (om *offsetManager) constructRequest() *OffsetCommitRequest {
	var r *OffsetCommitRequest
	var perPartitionTimestamp int64
	if om.conf.Consumer.Offsets.Retention == 0 {
		perPartitionTimestamp = ReceiveTime
		r = &OffsetCommitRequest{
			Version:                 1,
			ConsumerGroup:           om.group,
			ConsumerID:              om.memberID,
			ConsumerGroupGeneration: om.generation,
		}
	} else {
		r = &OffsetCommitRequest{
			Version:                 2,
			RetentionTime:           int64(om.conf.Consumer.Offsets.Retention / time.Millisecond),
			ConsumerGroup:           om.group,
			ConsumerID:              om.memberID,
			ConsumerGroupGeneration: om.generation,
		}

	}

	om.pomsLock.RLock()
	defer om.pomsLock.RUnlock()

	for _, topicManagers := range om.poms {
		for _, pom := range topicManagers {
			pom.lock.Lock()
			if pom.dirty {
				r.AddBlock(pom.topic, pom.partition, pom.offset, perPartitionTimestamp, pom.metadata)
			}
			pom.lock.Unlock()
		}
	}

	if len(r.blocks) > 0 {
		return r
	}

	return nil
}

func (om *offsetManager) handleResponse(broker *Broker, req *OffsetCommitRequest, resp *OffsetCommitResponse) {
	om.pomsLock.RLock()
	defer om.pomsLock.RUnlock()

	for _, topicManagers := range om.poms {
		for _, pom := range topicManagers {
			if req.blocks[pom.topic] == nil || req.blocks[pom.topic][pom.partition] == nil {
				continue
			}

			var err KError
			var ok bool

			if resp.Errors[pom.topic] == nil {
				pom.handleError(ErrIncompleteResponse)
				continue
			}
			if err, ok = resp.Errors[pom.topic][pom.partition]; !ok {
				pom.handleError(ErrIncompleteResponse)
				continue
			}

			switch err {
			case ErrNoError:
				block := req.blocks[pom.topic][pom.partition]
				pom.updateCommitted(block.offset, block.metadata)
			case ErrNotLeaderForPartition, ErrLeaderNotAvailable,
				ErrConsumerCoordinatorNotAvailable, ErrNotCoordinatorForConsumer:
				
				om.releaseCoordinator(broker)
			case ErrOffsetMetadataTooLarge, ErrInvalidCommitOffsetSize:
				
				pom.handleError(err)
			case ErrOffsetsLoadInProgress:
				
				break
			case ErrUnknownTopicOrPartition:
				
				
				
				
				fallthrough
			default:
				
				pom.handleError(err)
				om.releaseCoordinator(broker)
			}
		}
	}
}

func (om *offsetManager) handleError(err error) {
	om.pomsLock.RLock()
	defer om.pomsLock.RUnlock()

	for _, topicManagers := range om.poms {
		for _, pom := range topicManagers {
			pom.handleError(err)
		}
	}
}

func (om *offsetManager) asyncClosePOMs() {
	om.pomsLock.RLock()
	defer om.pomsLock.RUnlock()

	for _, topicManagers := range om.poms {
		for _, pom := range topicManagers {
			pom.AsyncClose()
		}
	}
}


func (om *offsetManager) releasePOMs(force bool) (remaining int) {
	om.pomsLock.Lock()
	defer om.pomsLock.Unlock()

	for topic, topicManagers := range om.poms {
		for partition, pom := range topicManagers {
			pom.lock.Lock()
			releaseDue := pom.done && (force || !pom.dirty)
			pom.lock.Unlock()

			if releaseDue {
				pom.release()

				delete(om.poms[topic], partition)
				if len(om.poms[topic]) == 0 {
					delete(om.poms, topic)
				}
			}
		}
		remaining += len(om.poms[topic])
	}
	return
}

func (om *offsetManager) findPOM(topic string, partition int32) *partitionOffsetManager {
	om.pomsLock.RLock()
	defer om.pomsLock.RUnlock()

	if partitions, ok := om.poms[topic]; ok {
		if pom, ok := partitions[partition]; ok {
			return pom
		}
	}
	return nil
}






type PartitionOffsetManager interface {
	
	
	
	
	
	NextOffset() (int64, string)

	
	
	
	
	
	
	
	
	
	
	
	
	
	MarkOffset(offset int64, metadata string)

	
	
	
	
	
	ResetOffset(offset int64, metadata string)

	
	
	
	
	Errors() <-chan *ConsumerError

	
	
	
	
	
	AsyncClose()

	
	
	
	
	Close() error
}

type partitionOffsetManager struct {
	parent    *offsetManager
	topic     string
	partition int32

	lock     sync.Mutex
	offset   int64
	metadata string
	dirty    bool
	done     bool

	releaseOnce sync.Once
	errors      chan *ConsumerError
}

func (om *offsetManager) newPartitionOffsetManager(topic string, partition int32) (*partitionOffsetManager, error) {
	offset, metadata, err := om.fetchInitialOffset(topic, partition, om.conf.Metadata.Retry.Max)
	if err != nil {
		return nil, err
	}

	return &partitionOffsetManager{
		parent:    om,
		topic:     topic,
		partition: partition,
		errors:    make(chan *ConsumerError, om.conf.ChannelBufferSize),
		offset:    offset,
		metadata:  metadata,
	}, nil
}

func (pom *partitionOffsetManager) Errors() <-chan *ConsumerError {
	return pom.errors
}

func (pom *partitionOffsetManager) MarkOffset(offset int64, metadata string) {
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if offset > pom.offset {
		pom.offset = offset
		pom.metadata = metadata
		pom.dirty = true
	}
}

func (pom *partitionOffsetManager) ResetOffset(offset int64, metadata string) {
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if offset <= pom.offset {
		pom.offset = offset
		pom.metadata = metadata
		pom.dirty = true
	}
}

func (pom *partitionOffsetManager) updateCommitted(offset int64, metadata string) {
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if pom.offset == offset && pom.metadata == metadata {
		pom.dirty = false
	}
}

func (pom *partitionOffsetManager) NextOffset() (int64, string) {
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if pom.offset >= 0 {
		return pom.offset, pom.metadata
	}

	return pom.parent.conf.Consumer.Offsets.Initial, ""
}

func (pom *partitionOffsetManager) AsyncClose() {
	pom.lock.Lock()
	pom.done = true
	pom.lock.Unlock()
}

func (pom *partitionOffsetManager) Close() error {
	pom.AsyncClose()

	var errors ConsumerErrors
	for err := range pom.errors {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (pom *partitionOffsetManager) handleError(err error) {
	cErr := &ConsumerError{
		Topic:     pom.topic,
		Partition: pom.partition,
		Err:       err,
	}

	if pom.parent.conf.Consumer.Return.Errors {
		pom.errors <- cErr
	} else {
		Logger.Println(cErr)
	}
}

func (pom *partitionOffsetManager) release() {
	pom.releaseOnce.Do(func() {
		go close(pom.errors)
	})
}
