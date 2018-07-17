package sarama

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)






type Client interface {
	
	
	Config() *Config

	
	Brokers() []*Broker

	
	Topics() ([]string, error)

	
	Partitions(topic string) ([]int32, error)

	
	
	
	WritablePartitions(topic string) ([]int32, error)

	
	
	Leader(topic string, partitionID int32) (*Broker, error)

	
	Replicas(topic string, partitionID int32) ([]int32, error)

	
	
	
	InSyncReplicas(topic string, partitionID int32) ([]int32, error)

	
	
	
	RefreshMetadata(topics ...string) error

	
	
	
	
	GetOffset(topic string, partitionID int32, time int64) (int64, error)

	
	
	
	
	Coordinator(consumerGroup string) (*Broker, error)

	
	
	RefreshCoordinator(consumerGroup string) error

	
	
	
	
	Close() error

	
	Closed() bool
}

const (
	
	
	
	
	OffsetNewest int64 = -1
	
	
	
	
	OffsetOldest int64 = -2
)

type client struct {
	conf           *Config
	closer, closed chan none 

	
	
	
	seedBrokers []*Broker
	deadSeeds   []*Broker

	brokers      map[int32]*Broker                       
	metadata     map[string]map[int32]*PartitionMetadata 
	coordinators map[string]int32                        

	
	
	cachedPartitionsResults map[string][maxPartitionIndex][]int32

	lock sync.RWMutex 
}




func NewClient(addrs []string, conf *Config) (Client, error) {
	Logger.Println("Initializing new client")

	if conf == nil {
		conf = NewConfig()
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	if len(addrs) < 1 {
		return nil, ConfigurationError("You must provide at least one broker address")
	}

	client := &client{
		conf:                    conf,
		closer:                  make(chan none),
		closed:                  make(chan none),
		brokers:                 make(map[int32]*Broker),
		metadata:                make(map[string]map[int32]*PartitionMetadata),
		cachedPartitionsResults: make(map[string][maxPartitionIndex][]int32),
		coordinators:            make(map[string]int32),
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, index := range random.Perm(len(addrs)) {
		client.seedBrokers = append(client.seedBrokers, NewBroker(addrs[index]))
	}

	if conf.Metadata.Full {
		
		err := client.RefreshMetadata()
		switch err {
		case nil:
			break
		case ErrLeaderNotAvailable, ErrReplicaNotAvailable, ErrTopicAuthorizationFailed, ErrClusterAuthorizationFailed:
			
			Logger.Println(err)
		default:
			close(client.closed) 
			_ = client.Close()
			return nil, err
		}
	}
	go withRecover(client.backgroundMetadataUpdater)

	Logger.Println("Successfully initialized new client")

	return client, nil
}

func (client *client) Config() *Config {
	return client.conf
}

func (client *client) Brokers() []*Broker {
	client.lock.RLock()
	defer client.lock.RUnlock()
	brokers := make([]*Broker, 0)
	for _, broker := range client.brokers {
		brokers = append(brokers, broker)
	}
	return brokers
}

func (client *client) Close() error {
	if client.Closed() {
		
		
		Logger.Printf("Close() called on already closed client")
		return ErrClosedClient
	}

	
	close(client.closer)
	<-client.closed

	client.lock.Lock()
	defer client.lock.Unlock()
	Logger.Println("Closing Client")

	for _, broker := range client.brokers {
		safeAsyncClose(broker)
	}

	for _, broker := range client.seedBrokers {
		safeAsyncClose(broker)
	}

	client.brokers = nil
	client.metadata = nil

	return nil
}

func (client *client) Closed() bool {
	return client.brokers == nil
}

func (client *client) Topics() ([]string, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	client.lock.RLock()
	defer client.lock.RUnlock()

	ret := make([]string, 0, len(client.metadata))
	for topic := range client.metadata {
		ret = append(ret, topic)
	}

	return ret, nil
}

func (client *client) Partitions(topic string) ([]int32, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	partitions := client.cachedPartitions(topic, allPartitions)

	if len(partitions) == 0 {
		err := client.RefreshMetadata(topic)
		if err != nil {
			return nil, err
		}
		partitions = client.cachedPartitions(topic, allPartitions)
	}

	if partitions == nil {
		return nil, ErrUnknownTopicOrPartition
	}

	return partitions, nil
}

func (client *client) WritablePartitions(topic string) ([]int32, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	partitions := client.cachedPartitions(topic, writablePartitions)

	
	
	
	
	
	
	if len(partitions) == 0 {
		err := client.RefreshMetadata(topic)
		if err != nil {
			return nil, err
		}
		partitions = client.cachedPartitions(topic, writablePartitions)
	}

	if partitions == nil {
		return nil, ErrUnknownTopicOrPartition
	}

	return partitions, nil
}

func (client *client) Replicas(topic string, partitionID int32) ([]int32, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	metadata := client.cachedMetadata(topic, partitionID)

	if metadata == nil {
		err := client.RefreshMetadata(topic)
		if err != nil {
			return nil, err
		}
		metadata = client.cachedMetadata(topic, partitionID)
	}

	if metadata == nil {
		return nil, ErrUnknownTopicOrPartition
	}

	if metadata.Err == ErrReplicaNotAvailable {
		return dupInt32Slice(metadata.Replicas), metadata.Err
	}
	return dupInt32Slice(metadata.Replicas), nil
}

func (client *client) InSyncReplicas(topic string, partitionID int32) ([]int32, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	metadata := client.cachedMetadata(topic, partitionID)

	if metadata == nil {
		err := client.RefreshMetadata(topic)
		if err != nil {
			return nil, err
		}
		metadata = client.cachedMetadata(topic, partitionID)
	}

	if metadata == nil {
		return nil, ErrUnknownTopicOrPartition
	}

	if metadata.Err == ErrReplicaNotAvailable {
		return dupInt32Slice(metadata.Isr), metadata.Err
	}
	return dupInt32Slice(metadata.Isr), nil
}

func (client *client) Leader(topic string, partitionID int32) (*Broker, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	leader, err := client.cachedLeader(topic, partitionID)

	if leader == nil {
		err = client.RefreshMetadata(topic)
		if err != nil {
			return nil, err
		}
		leader, err = client.cachedLeader(topic, partitionID)
	}

	return leader, err
}

func (client *client) RefreshMetadata(topics ...string) error {
	if client.Closed() {
		return ErrClosedClient
	}

	
	
	
	for _, topic := range topics {
		if len(topic) == 0 {
			return ErrInvalidTopic 
		}
	}

	return client.tryRefreshMetadata(topics, client.conf.Metadata.Retry.Max)
}

func (client *client) GetOffset(topic string, partitionID int32, time int64) (int64, error) {
	if client.Closed() {
		return -1, ErrClosedClient
	}

	offset, err := client.getOffset(topic, partitionID, time)

	if err != nil {
		if err := client.RefreshMetadata(topic); err != nil {
			return -1, err
		}
		return client.getOffset(topic, partitionID, time)
	}

	return offset, err
}

func (client *client) Coordinator(consumerGroup string) (*Broker, error) {
	if client.Closed() {
		return nil, ErrClosedClient
	}

	coordinator := client.cachedCoordinator(consumerGroup)

	if coordinator == nil {
		if err := client.RefreshCoordinator(consumerGroup); err != nil {
			return nil, err
		}
		coordinator = client.cachedCoordinator(consumerGroup)
	}

	if coordinator == nil {
		return nil, ErrConsumerCoordinatorNotAvailable
	}

	_ = coordinator.Open(client.conf)
	return coordinator, nil
}

func (client *client) RefreshCoordinator(consumerGroup string) error {
	if client.Closed() {
		return ErrClosedClient
	}

	response, err := client.getConsumerMetadata(consumerGroup, client.conf.Metadata.Retry.Max)
	if err != nil {
		return err
	}

	client.lock.Lock()
	defer client.lock.Unlock()
	client.registerBroker(response.Coordinator)
	client.coordinators[consumerGroup] = response.Coordinator.ID()
	return nil
}






func (client *client) registerBroker(broker *Broker) {
	if client.brokers[broker.ID()] == nil {
		client.brokers[broker.ID()] = broker
		Logger.Printf("client/brokers registered new broker #%d at %s", broker.ID(), broker.Addr())
	} else if broker.Addr() != client.brokers[broker.ID()].Addr() {
		safeAsyncClose(client.brokers[broker.ID()])
		client.brokers[broker.ID()] = broker
		Logger.Printf("client/brokers replaced registered broker #%d with %s", broker.ID(), broker.Addr())
	}
}



func (client *client) deregisterBroker(broker *Broker) {
	client.lock.Lock()
	defer client.lock.Unlock()

	if len(client.seedBrokers) > 0 && broker == client.seedBrokers[0] {
		client.deadSeeds = append(client.deadSeeds, broker)
		client.seedBrokers = client.seedBrokers[1:]
	} else {
		
		
		
		
		Logger.Printf("client/brokers deregistered broker #%d at %s", broker.ID(), broker.Addr())
		delete(client.brokers, broker.ID())
	}
}

func (client *client) resurrectDeadBrokers() {
	client.lock.Lock()
	defer client.lock.Unlock()

	Logger.Printf("client/brokers resurrecting %d dead seed brokers", len(client.deadSeeds))
	client.seedBrokers = append(client.seedBrokers, client.deadSeeds...)
	client.deadSeeds = nil
}

func (client *client) any() *Broker {
	client.lock.RLock()
	defer client.lock.RUnlock()

	if len(client.seedBrokers) > 0 {
		_ = client.seedBrokers[0].Open(client.conf)
		return client.seedBrokers[0]
	}

	
	for _, broker := range client.brokers {
		_ = broker.Open(client.conf)
		return broker
	}

	return nil
}



type partitionType int

const (
	allPartitions partitionType = iota
	writablePartitions
	

	
	maxPartitionIndex
)

func (client *client) cachedMetadata(topic string, partitionID int32) *PartitionMetadata {
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions := client.metadata[topic]
	if partitions != nil {
		return partitions[partitionID]
	}

	return nil
}

func (client *client) cachedPartitions(topic string, partitionSet partitionType) []int32 {
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions, exists := client.cachedPartitionsResults[topic]

	if !exists {
		return nil
	}
	return partitions[partitionSet]
}

func (client *client) setPartitionCache(topic string, partitionSet partitionType) []int32 {
	partitions := client.metadata[topic]

	if partitions == nil {
		return nil
	}

	ret := make([]int32, 0, len(partitions))
	for _, partition := range partitions {
		if partitionSet == writablePartitions && partition.Err == ErrLeaderNotAvailable {
			continue
		}
		ret = append(ret, partition.ID)
	}

	sort.Sort(int32Slice(ret))
	return ret
}

func (client *client) cachedLeader(topic string, partitionID int32) (*Broker, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions := client.metadata[topic]
	if partitions != nil {
		metadata, ok := partitions[partitionID]
		if ok {
			if metadata.Err == ErrLeaderNotAvailable {
				return nil, ErrLeaderNotAvailable
			}
			b := client.brokers[metadata.Leader]
			if b == nil {
				return nil, ErrLeaderNotAvailable
			}
			_ = b.Open(client.conf)
			return b, nil
		}
	}

	return nil, ErrUnknownTopicOrPartition
}

func (client *client) getOffset(topic string, partitionID int32, time int64) (int64, error) {
	broker, err := client.Leader(topic, partitionID)
	if err != nil {
		return -1, err
	}

	request := &OffsetRequest{}
	if client.conf.Version.IsAtLeast(V0_10_1_0) {
		request.Version = 1
	}
	request.AddBlock(topic, partitionID, time, 1)

	response, err := broker.GetAvailableOffsets(request)
	if err != nil {
		_ = broker.Close()
		return -1, err
	}

	block := response.GetBlock(topic, partitionID)
	if block == nil {
		_ = broker.Close()
		return -1, ErrIncompleteResponse
	}
	if block.Err != ErrNoError {
		return -1, block.Err
	}
	if len(block.Offsets) != 1 {
		return -1, ErrOffsetOutOfRange
	}

	return block.Offsets[0], nil
}



func (client *client) backgroundMetadataUpdater() {
	defer close(client.closed)

	if client.conf.Metadata.RefreshFrequency == time.Duration(0) {
		return
	}

	ticker := time.NewTicker(client.conf.Metadata.RefreshFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			topics := []string{}
			if !client.conf.Metadata.Full {
				if specificTopics, err := client.Topics(); err != nil {
					Logger.Println("Client background metadata topic load:", err)
					break
				} else if len(specificTopics) == 0 {
					Logger.Println("Client background metadata update: no specific topics to update")
					break
				} else {
					topics = specificTopics
				}
			}

			if err := client.RefreshMetadata(topics...); err != nil {
				Logger.Println("Client background metadata update:", err)
			}
		case <-client.closer:
			return
		}
	}
}

func (client *client) tryRefreshMetadata(topics []string, attemptsRemaining int) error {
	retry := func(err error) error {
		if attemptsRemaining > 0 {
			Logger.Printf("client/metadata retrying after %dms... (%d attempts remaining)\n", client.conf.Metadata.Retry.Backoff/time.Millisecond, attemptsRemaining)
			time.Sleep(client.conf.Metadata.Retry.Backoff)
			return client.tryRefreshMetadata(topics, attemptsRemaining-1)
		}
		return err
	}

	for broker := client.any(); broker != nil; broker = client.any() {
		if len(topics) > 0 {
			Logger.Printf("client/metadata fetching metadata for %v from broker %s\n", topics, broker.addr)
		} else {
			Logger.Printf("client/metadata fetching metadata for all topics from broker %s\n", broker.addr)
		}
		response, err := broker.GetMetadata(&MetadataRequest{Topics: topics})

		switch err.(type) {
		case nil:
			
			shouldRetry, err := client.updateMetadata(response)
			if shouldRetry {
				Logger.Println("client/metadata found some partitions to be leaderless")
				return retry(err) 
			}
			return err

		case PacketEncodingError:
			
			return err
		default:
			
			Logger.Println("client/metadata got error from broker while fetching metadata:", err)
			_ = broker.Close()
			client.deregisterBroker(broker)
		}
	}

	Logger.Println("client/metadata no available broker to send metadata request to")
	client.resurrectDeadBrokers()
	return retry(ErrOutOfBrokers)
}


func (client *client) updateMetadata(data *MetadataResponse) (retry bool, err error) {
	client.lock.Lock()
	defer client.lock.Unlock()

	
	
	
	
	for _, broker := range data.Brokers {
		client.registerBroker(broker)
	}

	for _, topic := range data.Topics {
		delete(client.metadata, topic.Name)
		delete(client.cachedPartitionsResults, topic.Name)

		switch topic.Err {
		case ErrNoError:
			break
		case ErrInvalidTopic, ErrTopicAuthorizationFailed: 
			err = topic.Err
			continue
		case ErrUnknownTopicOrPartition: 
			err = topic.Err
			retry = true
			continue
		case ErrLeaderNotAvailable: 
			retry = true
			break
		default: 
			Logger.Printf("Unexpected topic-level metadata error: %s", topic.Err)
			err = topic.Err
			continue
		}

		client.metadata[topic.Name] = make(map[int32]*PartitionMetadata, len(topic.Partitions))
		for _, partition := range topic.Partitions {
			client.metadata[topic.Name][partition.ID] = partition
			if partition.Err == ErrLeaderNotAvailable {
				retry = true
			}
		}

		var partitionCache [maxPartitionIndex][]int32
		partitionCache[allPartitions] = client.setPartitionCache(topic.Name, allPartitions)
		partitionCache[writablePartitions] = client.setPartitionCache(topic.Name, writablePartitions)
		client.cachedPartitionsResults[topic.Name] = partitionCache
	}

	return
}

func (client *client) cachedCoordinator(consumerGroup string) *Broker {
	client.lock.RLock()
	defer client.lock.RUnlock()
	if coordinatorID, ok := client.coordinators[consumerGroup]; ok {
		return client.brokers[coordinatorID]
	}
	return nil
}

func (client *client) getConsumerMetadata(consumerGroup string, attemptsRemaining int) (*ConsumerMetadataResponse, error) {
	retry := func(err error) (*ConsumerMetadataResponse, error) {
		if attemptsRemaining > 0 {
			Logger.Printf("client/coordinator retrying after %dms... (%d attempts remaining)\n", client.conf.Metadata.Retry.Backoff/time.Millisecond, attemptsRemaining)
			time.Sleep(client.conf.Metadata.Retry.Backoff)
			return client.getConsumerMetadata(consumerGroup, attemptsRemaining-1)
		}
		return nil, err
	}

	for broker := client.any(); broker != nil; broker = client.any() {
		Logger.Printf("client/coordinator requesting coordinator for consumergroup %s from %s\n", consumerGroup, broker.Addr())

		request := new(ConsumerMetadataRequest)
		request.ConsumerGroup = consumerGroup

		response, err := broker.GetConsumerMetadata(request)

		if err != nil {
			Logger.Printf("client/coordinator request to broker %s failed: %s\n", broker.Addr(), err)

			switch err.(type) {
			case PacketEncodingError:
				return nil, err
			default:
				_ = broker.Close()
				client.deregisterBroker(broker)
				continue
			}
		}

		switch response.Err {
		case ErrNoError:
			Logger.Printf("client/coordinator coordinator for consumergroup %s is #%d (%s)\n", consumerGroup, response.Coordinator.ID(), response.Coordinator.Addr())
			return response, nil

		case ErrConsumerCoordinatorNotAvailable:
			Logger.Printf("client/coordinator coordinator for consumer group %s is not available\n", consumerGroup)

			
			
			
			if _, err := client.Leader("__consumer_offsets", 0); err != nil {
				Logger.Printf("client/coordinator the __consumer_offsets topic is not initialized completely yet. Waiting 2 seconds...\n")
				time.Sleep(2 * time.Second)
			}

			return retry(ErrConsumerCoordinatorNotAvailable)
		default:
			return nil, response.Err
		}
	}

	Logger.Println("client/coordinator no available broker to send consumer metadata request to")
	client.resurrectDeadBrokers()
	return retry(ErrOutOfBrokers)
}
