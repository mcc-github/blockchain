package sarama

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)


var ErrClosedConsumerGroup = errors.New("kafka: tried to use a consumer group that was closed")



type ConsumerGroup interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Consume(ctx context.Context, topics []string, handler ConsumerGroupHandler) error

	
	
	
	
	Errors() <-chan error

	
	
	Close() error
}

type consumerGroup struct {
	client    Client
	ownClient bool

	config   *Config
	consumer Consumer
	groupID  string
	memberID string
	errors   chan error

	lock      sync.Mutex
	closed    chan none
	closeOnce sync.Once
}


func NewConsumerGroup(addrs []string, groupID string, config *Config) (ConsumerGroup, error) {
	client, err := NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		_ = client.Close()
		return nil, err
	}

	c.(*consumerGroup).ownClient = true
	return c, nil
}




func NewConsumerGroupFromClient(groupID string, client Client) (ConsumerGroup, error) {
	config := client.Config()
	if !config.Version.IsAtLeast(V0_10_2_0) {
		return nil, ConfigurationError("consumer groups require Version to be >= V0_10_2_0")
	}

	consumer, err := NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		client:   client,
		consumer: consumer,
		config:   config,
		groupID:  groupID,
		errors:   make(chan error, config.ChannelBufferSize),
		closed:   make(chan none),
	}, nil
}


func (c *consumerGroup) Errors() <-chan error { return c.errors }


func (c *consumerGroup) Close() (err error) {
	c.closeOnce.Do(func() {
		close(c.closed)

		c.lock.Lock()
		defer c.lock.Unlock()

		
		if e := c.leave(); e != nil {
			err = e
		}

		
		go func() {
			close(c.errors)
		}()
		for e := range c.errors {
			err = e
		}

		if c.ownClient {
			if e := c.client.Close(); e != nil {
				err = e
			}
		}
	})
	return
}


func (c *consumerGroup) Consume(ctx context.Context, topics []string, handler ConsumerGroupHandler) error {
	
	select {
	case <-c.closed:
		return ErrClosedConsumerGroup
	default:
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	
	if len(topics) == 0 {
		return fmt.Errorf("no topics provided")
	}

	
	if err := c.client.RefreshMetadata(topics...); err != nil {
		return err
	}

	
	coordinator, err := c.client.Coordinator(c.groupID)
	if err != nil {
		return err
	}

	
	sess, err := c.newSession(ctx, coordinator, topics, handler, c.config.Consumer.Group.Rebalance.Retry.Max)
	if err == ErrClosedClient {
		return ErrClosedConsumerGroup
	} else if err != nil {
		return err
	}

	
	<-sess.ctx.Done()

	
	return sess.release(true)
}

func (c *consumerGroup) newSession(ctx context.Context, coordinator *Broker, topics []string, handler ConsumerGroupHandler, retries int) (*consumerGroupSession, error) {
	
	join, err := c.joinGroupRequest(coordinator, topics)
	if err != nil {
		_ = coordinator.Close()
		return nil, err
	}
	switch join.Err {
	case ErrNoError:
		c.memberID = join.MemberId
	case ErrUnknownMemberId, ErrIllegalGeneration: 
		c.memberID = ""
		return c.newSession(ctx, coordinator, topics, handler, retries)
	case ErrRebalanceInProgress: 
		if retries <= 0 {
			return nil, join.Err
		}

		select {
		case <-c.closed:
			return nil, ErrClosedConsumerGroup
		case <-time.After(c.config.Consumer.Group.Rebalance.Retry.Backoff):
		}

		return c.newSession(ctx, coordinator, topics, handler, retries-1)
	default:
		return nil, join.Err
	}

	
	var plan BalanceStrategyPlan
	if join.LeaderId == join.MemberId {
		members, err := join.GetMembers()
		if err != nil {
			return nil, err
		}

		plan, err = c.balance(members)
		if err != nil {
			return nil, err
		}
	}

	
	sync, err := c.syncGroupRequest(coordinator, plan, join.GenerationId)
	if err != nil {
		_ = coordinator.Close()
		return nil, err
	}
	switch sync.Err {
	case ErrNoError:
	case ErrUnknownMemberId, ErrIllegalGeneration: 
		c.memberID = ""
		return c.newSession(ctx, coordinator, topics, handler, retries)
	case ErrRebalanceInProgress: 
		if retries <= 0 {
			return nil, sync.Err
		}

		select {
		case <-c.closed:
			return nil, ErrClosedConsumerGroup
		case <-time.After(c.config.Consumer.Group.Rebalance.Retry.Backoff):
		}

		return c.newSession(ctx, coordinator, topics, handler, retries-1)
	default:
		return nil, sync.Err
	}

	
	var claims map[string][]int32
	if len(sync.MemberAssignment) > 0 {
		members, err := sync.GetMemberAssignment()
		if err != nil {
			return nil, err
		}
		claims = members.Topics

		for _, partitions := range claims {
			sort.Sort(int32Slice(partitions))
		}
	}

	return newConsumerGroupSession(c, ctx, claims, join.MemberId, join.GenerationId, handler)
}

func (c *consumerGroup) joinGroupRequest(coordinator *Broker, topics []string) (*JoinGroupResponse, error) {
	req := &JoinGroupRequest{
		GroupId:        c.groupID,
		MemberId:       c.memberID,
		SessionTimeout: int32(c.config.Consumer.Group.Session.Timeout / time.Millisecond),
		ProtocolType:   "consumer",
	}
	if c.config.Version.IsAtLeast(V0_10_1_0) {
		req.Version = 1
		req.RebalanceTimeout = int32(c.config.Consumer.Group.Rebalance.Timeout / time.Millisecond)
	}

	meta := &ConsumerGroupMemberMetadata{
		Topics:   topics,
		UserData: c.config.Consumer.Group.Member.UserData,
	}
	strategy := c.config.Consumer.Group.Rebalance.Strategy
	if err := req.AddGroupProtocolMetadata(strategy.Name(), meta); err != nil {
		return nil, err
	}

	return coordinator.JoinGroup(req)
}

func (c *consumerGroup) syncGroupRequest(coordinator *Broker, plan BalanceStrategyPlan, generationID int32) (*SyncGroupResponse, error) {
	req := &SyncGroupRequest{
		GroupId:      c.groupID,
		MemberId:     c.memberID,
		GenerationId: generationID,
	}
	for memberID, topics := range plan {
		err := req.AddGroupAssignmentMember(memberID, &ConsumerGroupMemberAssignment{
			Topics: topics,
		})
		if err != nil {
			return nil, err
		}
	}
	return coordinator.SyncGroup(req)
}

func (c *consumerGroup) heartbeatRequest(coordinator *Broker, memberID string, generationID int32) (*HeartbeatResponse, error) {
	req := &HeartbeatRequest{
		GroupId:      c.groupID,
		MemberId:     memberID,
		GenerationId: generationID,
	}

	return coordinator.Heartbeat(req)
}

func (c *consumerGroup) balance(members map[string]ConsumerGroupMemberMetadata) (BalanceStrategyPlan, error) {
	topics := make(map[string][]int32)
	for _, meta := range members {
		for _, topic := range meta.Topics {
			topics[topic] = nil
		}
	}

	for topic := range topics {
		partitions, err := c.client.Partitions(topic)
		if err != nil {
			return nil, err
		}
		topics[topic] = partitions
	}

	strategy := c.config.Consumer.Group.Rebalance.Strategy
	return strategy.Plan(members, topics)
}


func (c *consumerGroup) leave() error {
	if c.memberID == "" {
		return nil
	}

	coordinator, err := c.client.Coordinator(c.groupID)
	if err != nil {
		return err
	}

	resp, err := coordinator.LeaveGroup(&LeaveGroupRequest{
		GroupId:  c.groupID,
		MemberId: c.memberID,
	})
	if err != nil {
		_ = coordinator.Close()
		return err
	}

	
	c.memberID = ""

	
	switch resp.Err {
	case ErrRebalanceInProgress, ErrUnknownMemberId, ErrNoError:
		return nil
	default:
		return resp.Err
	}
}

func (c *consumerGroup) handleError(err error, topic string, partition int32) {
	select {
	case <-c.closed:
		return
	default:
	}

	if _, ok := err.(*ConsumerError); !ok && topic != "" && partition > -1 {
		err = &ConsumerError{
			Topic:     topic,
			Partition: partition,
			Err:       err,
		}
	}

	if c.config.Consumer.Return.Errors {
		select {
		case c.errors <- err:
		default:
		}
	} else {
		Logger.Println(err)
	}
}




type ConsumerGroupSession interface {
	
	Claims() map[string][]int32

	
	MemberID() string

	
	GenerationID() int32

	
	
	
	
	
	
	
	
	
	
	
	
	
	MarkOffset(topic string, partition int32, offset int64, metadata string)

	
	
	
	
	
	ResetOffset(topic string, partition int32, offset int64, metadata string)

	
	MarkMessage(msg *ConsumerMessage, metadata string)

	
	Context() context.Context
}

type consumerGroupSession struct {
	parent       *consumerGroup
	memberID     string
	generationID int32
	handler      ConsumerGroupHandler

	claims  map[string][]int32
	offsets *offsetManager
	ctx     context.Context
	cancel  func()

	waitGroup       sync.WaitGroup
	releaseOnce     sync.Once
	hbDying, hbDead chan none
}

func newConsumerGroupSession(parent *consumerGroup, ctx context.Context, claims map[string][]int32, memberID string, generationID int32, handler ConsumerGroupHandler) (*consumerGroupSession, error) {
	
	offsets, err := newOffsetManagerFromClient(parent.groupID, memberID, generationID, parent.client)
	if err != nil {
		return nil, err
	}

	
	ctx, cancel := context.WithCancel(ctx)

	
	sess := &consumerGroupSession{
		parent:       parent,
		memberID:     memberID,
		generationID: generationID,
		handler:      handler,
		offsets:      offsets,
		claims:       claims,
		ctx:          ctx,
		cancel:       cancel,
		hbDying:      make(chan none),
		hbDead:       make(chan none),
	}

	
	go sess.heartbeatLoop()

	
	for topic, partitions := range claims {
		for _, partition := range partitions {
			pom, err := offsets.ManagePartition(topic, partition)
			if err != nil {
				_ = sess.release(false)
				return nil, err
			}

			
			go func(topic string, partition int32) {
				for err := range pom.Errors() {
					sess.parent.handleError(err, topic, partition)
				}
			}(topic, partition)
		}
	}

	
	if err := handler.Setup(sess); err != nil {
		_ = sess.release(true)
		return nil, err
	}

	
	for topic, partitions := range claims {
		for _, partition := range partitions {
			sess.waitGroup.Add(1)

			go func(topic string, partition int32) {
				defer sess.waitGroup.Done()

				
				
				defer sess.cancel()

				
				sess.consume(topic, partition)
			}(topic, partition)
		}
	}
	return sess, nil
}

func (s *consumerGroupSession) Claims() map[string][]int32 { return s.claims }
func (s *consumerGroupSession) MemberID() string           { return s.memberID }
func (s *consumerGroupSession) GenerationID() int32        { return s.generationID }

func (s *consumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	if pom := s.offsets.findPOM(topic, partition); pom != nil {
		pom.MarkOffset(offset, metadata)
	}
}

func (s *consumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	if pom := s.offsets.findPOM(topic, partition); pom != nil {
		pom.ResetOffset(offset, metadata)
	}
}

func (s *consumerGroupSession) MarkMessage(msg *ConsumerMessage, metadata string) {
	s.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, metadata)
}

func (s *consumerGroupSession) Context() context.Context {
	return s.ctx
}

func (s *consumerGroupSession) consume(topic string, partition int32) {
	
	select {
	case <-s.ctx.Done():
		return
	case <-s.parent.closed:
		return
	default:
	}

	
	offset := s.parent.config.Consumer.Offsets.Initial
	if pom := s.offsets.findPOM(topic, partition); pom != nil {
		offset, _ = pom.NextOffset()
	}

	
	claim, err := newConsumerGroupClaim(s, topic, partition, offset)
	if err != nil {
		s.parent.handleError(err, topic, partition)
		return
	}

	
	go func() {
		for err := range claim.Errors() {
			s.parent.handleError(err, topic, partition)
		}
	}()

	
	go func() {
		select {
		case <-s.ctx.Done():
		case <-s.parent.closed:
		}
		claim.AsyncClose()
	}()

	
	if err := s.handler.ConsumeClaim(s, claim); err != nil {
		s.parent.handleError(err, topic, partition)
	}

	
	claim.AsyncClose()
	for _, err := range claim.waitClosed() {
		s.parent.handleError(err, topic, partition)
	}
}

func (s *consumerGroupSession) release(withCleanup bool) (err error) {
	
	s.cancel()

	
	s.waitGroup.Wait()

	
	s.releaseOnce.Do(func() {
		if withCleanup {
			if e := s.handler.Cleanup(s); e != nil {
				s.parent.handleError(err, "", -1)
				err = e
			}
		}

		if e := s.offsets.Close(); e != nil {
			err = e
		}

		close(s.hbDying)
		<-s.hbDead
	})

	return
}

func (s *consumerGroupSession) heartbeatLoop() {
	defer close(s.hbDead)
	defer s.cancel() 

	pause := time.NewTicker(s.parent.config.Consumer.Group.Heartbeat.Interval)
	defer pause.Stop()

	retries := s.parent.config.Metadata.Retry.Max
	for {
		coordinator, err := s.parent.client.Coordinator(s.parent.groupID)
		if err != nil {
			if retries <= 0 {
				s.parent.handleError(err, "", -1)
				return
			}

			select {
			case <-s.hbDying:
				return
			case <-time.After(s.parent.config.Metadata.Retry.Backoff):
				retries--
			}
			continue
		}

		resp, err := s.parent.heartbeatRequest(coordinator, s.memberID, s.generationID)
		if err != nil {
			_ = coordinator.Close()
			retries--
			continue
		}

		switch resp.Err {
		case ErrNoError:
			retries = s.parent.config.Metadata.Retry.Max
		case ErrRebalanceInProgress, ErrUnknownMemberId, ErrIllegalGeneration:
			return
		default:
			s.parent.handleError(err, "", -1)
			return
		}

		select {
		case <-pause.C:
		case <-s.hbDying:
			return
		}
	}
}









type ConsumerGroupHandler interface {
	
	Setup(ConsumerGroupSession) error

	
	
	Cleanup(ConsumerGroupSession) error

	
	
	
	ConsumeClaim(ConsumerGroupSession, ConsumerGroupClaim) error
}


type ConsumerGroupClaim interface {
	
	Topic() string

	
	Partition() int32

	
	InitialOffset() int64

	
	
	
	HighWaterMarkOffset() int64

	
	
	
	
	
	Messages() <-chan *ConsumerMessage
}

type consumerGroupClaim struct {
	topic     string
	partition int32
	offset    int64
	PartitionConsumer
}

func newConsumerGroupClaim(sess *consumerGroupSession, topic string, partition int32, offset int64) (*consumerGroupClaim, error) {
	pcm, err := sess.parent.consumer.ConsumePartition(topic, partition, offset)
	if err == ErrOffsetOutOfRange {
		offset = sess.parent.config.Consumer.Offsets.Initial
		pcm, err = sess.parent.consumer.ConsumePartition(topic, partition, offset)
	}
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range pcm.Errors() {
			sess.parent.handleError(err, topic, partition)
		}
	}()

	return &consumerGroupClaim{
		topic:             topic,
		partition:         partition,
		offset:            offset,
		PartitionConsumer: pcm,
	}, nil
}

func (c *consumerGroupClaim) Topic() string        { return c.topic }
func (c *consumerGroupClaim) Partition() int32     { return c.partition }
func (c *consumerGroupClaim) InitialOffset() int64 { return c.offset }


func (c *consumerGroupClaim) waitClosed() (errs ConsumerErrors) {
	go func() {
		for range c.Messages() {
		}
	}()

	for err := range c.Errors() {
		errs = append(errs, err)
	}
	return
}
