/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"context"
	"encoding/pem"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/clock"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/consensus"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)




const DefaultSnapshotCatchUpEntries = uint64(500)





type Configurator interface {
	Configure(channel string, newNodes []cluster.RemoteNode)
}




type RPC interface {
	Step(dest uint64, msg *orderer.StepRequest) (*orderer.StepResponse, error)
	SendSubmit(dest uint64, request *orderer.SubmitRequest) error
}




type BlockPuller interface {
	PullBlock(seq uint64) *common.Block
	Close()
}

type block struct {
	b *common.Block

	
	
	
	i uint64
}


type Options struct {
	RaftID uint64

	Clock clock.Clock

	WALDir       string
	SnapDir      string
	SnapInterval uint64

	
	
	SnapshotCatchUpEntries uint64

	MemoryStorage MemoryStorage
	Logger        *flogging.FabricLogger

	TickInterval    time.Duration
	ElectionTick    int
	HeartbeatTick   int
	MaxSizePerMsg   uint64
	MaxInflightMsgs int

	RaftMetadata *etcdraft.RaftMetadata
}


type Chain struct {
	configurator Configurator
	rpc          RPC

	raftID    uint64
	channelID string

	submitC  chan *orderer.SubmitRequest
	commitC  chan block
	observeC chan<- uint64         
	haltC    chan struct{}         
	doneC    chan struct{}         
	resignC  chan struct{}         
	startC   chan struct{}         
	snapC    chan *raftpb.Snapshot 

	configChangeApplyC     chan struct{} 
	configChangeInProgress uint32        
	raftMetadataLock       sync.RWMutex

	clock clock.Clock 

	support consensus.ConsenterSupport

	leader       uint64
	appliedIndex uint64

	
	lastSnapBlockNum uint64
	syncLock         sync.Mutex       
	syncC            chan struct{}    
	confState        raftpb.ConfState 
	puller           BlockPuller      

	fresh bool 

	node    raft.Node
	storage *RaftStorage
	opts    Options

	logger *flogging.FabricLogger
}


func NewChain(
	support consensus.ConsenterSupport,
	opts Options,
	conf Configurator,
	rpc RPC,
	puller BlockPuller,
	observeC chan<- uint64) (*Chain, error) {

	lg := opts.Logger.With("channel", support.ChainID(), "node", opts.RaftID)

	fresh := !wal.Exist(opts.WALDir)

	appliedi := opts.RaftMetadata.RaftIndex
	storage, err := CreateStorage(lg, appliedi, opts.WALDir, opts.SnapDir, opts.MemoryStorage)
	if err != nil {
		return nil, errors.Errorf("failed to restore persisted raft data: %s", err)
	}

	if opts.SnapshotCatchUpEntries == 0 {
		storage.SnapshotCatchUpEntries = DefaultSnapshotCatchUpEntries
	} else {
		storage.SnapshotCatchUpEntries = opts.SnapshotCatchUpEntries
	}

	
	var snapBlkNum uint64
	if s := storage.Snapshot(); !raft.IsEmptySnap(s) {
		b := utils.UnmarshalBlockOrPanic(s.Data)
		snapBlkNum = b.Header.Number
	}

	return &Chain{
		configurator:       conf,
		rpc:                rpc,
		channelID:          support.ChainID(),
		raftID:             opts.RaftID,
		submitC:            make(chan *orderer.SubmitRequest),
		commitC:            make(chan block),
		haltC:              make(chan struct{}),
		doneC:              make(chan struct{}),
		resignC:            make(chan struct{}),
		startC:             make(chan struct{}),
		syncC:              make(chan struct{}),
		snapC:              make(chan *raftpb.Snapshot),
		configChangeApplyC: make(chan struct{}),
		observeC:           observeC,
		support:            support,
		fresh:              fresh,
		appliedIndex:       appliedi,
		lastSnapBlockNum:   snapBlkNum,
		puller:             puller,
		clock:              opts.Clock,
		logger:             lg,
		storage:            storage,
		opts:               opts,
	}, nil
}


func (c *Chain) Start() {
	c.logger.Infof("Starting Raft node")

	
	
	config := &raft.Config{
		ID:              c.raftID,
		ElectionTick:    c.opts.ElectionTick,
		HeartbeatTick:   c.opts.HeartbeatTick,
		MaxSizePerMsg:   c.opts.MaxSizePerMsg,
		MaxInflightMsgs: c.opts.MaxInflightMsgs,
		Logger:          c.logger,
		Storage:         c.opts.MemoryStorage,
		
		
		PreVote:                   true,
		DisableProposalForwarding: true, 
	}

	if err := c.configureComm(); err != nil {
		c.logger.Errorf("Failed to start chain, aborting: +%v", err)
		close(c.doneC)
		return
	}

	raftPeers := RaftPeers(c.opts.RaftMetadata.Consenters)

	if c.fresh {
		c.logger.Info("starting new raft node")
		c.node = raft.StartNode(config, raftPeers)
	} else {
		c.logger.Info("restarting raft node")
		c.node = raft.RestartNode(config)
	}

	close(c.startC)

	go c.serveRaft()
	go c.serveRequest()
}


func (c *Chain) Order(env *common.Envelope, configSeq uint64) error {
	return c.Submit(&orderer.SubmitRequest{LastValidationSeq: configSeq, Content: env, Channel: c.channelID}, 0)
}


func (c *Chain) Configure(env *common.Envelope, configSeq uint64) error {
	if err := c.checkConfigUpdateValidity(env); err != nil {
		return err
	}
	return c.Submit(&orderer.SubmitRequest{LastValidationSeq: configSeq, Content: env, Channel: c.channelID}, 0)
}


func (c *Chain) checkConfigUpdateValidity(ctx *common.Envelope) error {
	var err error
	payload, err := utils.UnmarshalPayload(ctx.Payload)
	if err != nil {
		return err
	}
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return err
	}

	switch chdr.Type {
	case int32(common.HeaderType_ORDERER_TRANSACTION):
		return nil
	case int32(common.HeaderType_CONFIG):
		configUpdate, err := configtx.UnmarshalConfigUpdateFromPayload(payload)
		if err != nil {
			return err
		}

		
		if ordererConfigGroup, ok := configUpdate.WriteSet.Groups["Orderer"]; ok {
			if val, ok := ordererConfigGroup.Values["ConsensusType"]; ok {
				return c.checkConsentersSet(val)
			}
		}
		return nil

	default:
		return errors.Errorf("config transaction has unknown header type")
	}
}





func (c *Chain) WaitReady() error {
	if err := c.isRunning(); err != nil {
		return err
	}

	c.syncLock.Lock()
	ch := c.syncC
	c.syncLock.Unlock()

	select {
	case <-ch:
	case <-c.doneC:
		return errors.Errorf("chain is stopped")
	}

	return nil
}


func (c *Chain) Errored() <-chan struct{} {
	return c.doneC
}


func (c *Chain) Halt() {
	select {
	case <-c.startC:
	default:
		c.logger.Warnf("Attempted to halt a chain that has not started")
		return
	}

	select {
	case c.haltC <- struct{}{}:
	case <-c.doneC:
		return
	}
	<-c.doneC
}

func (c *Chain) isRunning() error {
	select {
	case <-c.startC:
	default:
		return errors.Errorf("chain is not started")
	}

	select {
	case <-c.doneC:
		return errors.Errorf("chain is stopped")
	default:
	}

	return nil
}


func (c *Chain) Step(req *orderer.StepRequest, sender uint64) error {
	if err := c.isRunning(); err != nil {
		return err
	}

	stepMsg := &raftpb.Message{}
	if err := proto.Unmarshal(req.Payload, stepMsg); err != nil {
		return fmt.Errorf("failed to unmarshal StepRequest payload to Raft Message: %s", err)
	}

	if err := c.node.Step(context.TODO(), *stepMsg); err != nil {
		return fmt.Errorf("failed to process Raft Step message: %s", err)
	}

	return nil
}





func (c *Chain) Submit(req *orderer.SubmitRequest, sender uint64) error {
	if err := c.isRunning(); err != nil {
		return err
	}

	lead := atomic.LoadUint64(&c.leader)

	if lead == raft.None {
		return errors.Errorf("no Raft leader")
	}

	if lead == c.raftID {
		select {
		case c.submitC <- req:
			return nil
		case <-c.doneC:
			return errors.Errorf("chain is stopped")
		}
	}

	c.logger.Debugf("Forwarding submit request to Raft leader %d", lead)
	return c.rpc.SendSubmit(lead, req)
}

func (c *Chain) serveRequest() {
	ticking := false
	timer := c.clock.NewTimer(time.Second)
	
	
	if !timer.Stop() {
		<-timer.C()
	}

	
	start := func() {
		if !ticking {
			ticking = true
			timer.Reset(c.support.SharedConfig().BatchTimeout())
		}
	}

	stop := func() {
		if !timer.Stop() && ticking {
			
			<-timer.C()
		}
		ticking = false
	}

	if s := c.storage.Snapshot(); !raft.IsEmptySnap(s) {
		if err := c.catchUp(&s); err != nil {
			c.logger.Errorf("Failed to recover from snapshot taken at Term %d and Index %d: %s",
				s.Metadata.Term, s.Metadata.Index, err)
		}
	} else {
		close(c.syncC)
	}

	for {
		select {
		case msg := <-c.submitC:
			batches, pending, err := c.ordered(msg)
			if err != nil {
				c.logger.Errorf("Failed to order message: %s", err)
			}
			if pending {
				start() 
			} else {
				stop()
			}

			if err := c.commitBatches(batches...); err != nil {
				c.logger.Errorf("Failed to commit block: %s", err)
			}

		case b := <-c.commitC:
			c.writeBlock(b)

		case <-c.resignC:
			_ = c.support.BlockCutter().Cut()
			stop()

		case <-timer.C():
			ticking = false

			batch := c.support.BlockCutter().Cut()
			if len(batch) == 0 {
				c.logger.Warningf("Batch timer expired with no pending requests, this might indicate a bug")
				continue
			}

			c.logger.Debugf("Batch timer expired, creating block")
			if err := c.commitBatches(batch); err != nil {
				c.logger.Errorf("Failed to commit block: %s", err)
			}

		case sn := <-c.snapC:
			if err := c.catchUp(sn); err != nil {
				c.logger.Errorf("Failed to recover from snapshot taken at Term %d and Index %d: %s",
					sn.Metadata.Term, sn.Metadata.Index, err)
			}

		case <-c.doneC:
			c.logger.Infof("Stop serving requests")
			return
		}
	}
}

func (c *Chain) writeBlock(b block) {
	if utils.IsConfigBlock(b.b) {
		if err := c.writeConfigBlock(b); err != nil {
			c.logger.Panicf("failed to write configuration block, %+v", err)
		}
		return
	}

	c.raftMetadataLock.Lock()
	c.opts.RaftMetadata.RaftIndex = b.i
	m := utils.MarshalOrPanic(c.opts.RaftMetadata)
	c.raftMetadataLock.Unlock()

	c.support.WriteBlock(b.b, m)
}







func (c *Chain) ordered(msg *orderer.SubmitRequest) (batches [][]*common.Envelope, pending bool, err error) {
	seq := c.support.Sequence()

	if c.isConfig(msg.Content) {
		
		if msg.LastValidationSeq < seq {
			msg.Content, _, err = c.support.ProcessConfigMsg(msg.Content)
			if err != nil {
				return nil, true, errors.Errorf("bad config message: %s", err)
			}
		}
		batch := c.support.BlockCutter().Cut()
		batches = [][]*common.Envelope{}
		if len(batch) != 0 {
			batches = append(batches, batch)
		}
		batches = append(batches, []*common.Envelope{msg.Content})
		return batches, false, nil
	}
	
	if msg.LastValidationSeq < seq {
		if _, err := c.support.ProcessNormalMsg(msg.Content); err != nil {
			return nil, true, errors.Errorf("bad normal message: %s", err)
		}
	}
	batches, pending = c.support.BlockCutter().Ordered(msg.Content)
	return batches, pending, nil

}

func (c *Chain) commitBatches(batches ...[]*common.Envelope) error {
	for _, batch := range batches {
		b := c.support.CreateNextBlock(batch)
		data := utils.MarshalOrPanic(b)
		if err := c.node.Propose(context.TODO(), data); err != nil {
			return errors.Errorf("failed to propose data to Raft node: %s", err)
		}

		select {
		case block := <-c.commitC:
			c.writeBlock(block)
		case <-c.resignC:
			return errors.Errorf("aborted block committing: lost leadership")

		case <-c.doneC:
			return nil
		}
	}

	return nil
}

func (c *Chain) catchUp(snap *raftpb.Snapshot) error {
	b, err := utils.UnmarshalBlock(snap.Data)
	if err != nil {
		return errors.Errorf("failed to unmarshal snapshot data to block: %s", err)
	}

	c.logger.Infof("Catching up with snapshot taken at block %d", b.Header.Number)

	next := c.support.Height()
	if next > b.Header.Number {
		c.logger.Warnf("Snapshot is at block %d, local block number is %d, no sync needed", b.Header.Number, next-1)
		return nil
	}

	c.syncLock.Lock()
	c.syncC = make(chan struct{})
	c.syncLock.Unlock()
	defer func() {
		close(c.syncC)
		c.puller.Close()
	}()

	for next <= b.Header.Number {
		block := c.puller.PullBlock(next)
		if block == nil {
			return errors.Errorf("failed to fetch block %d from cluster", next)
		}

		if utils.IsConfigBlock(block) {
			c.support.WriteConfigBlock(block, nil)
		} else {
			c.support.WriteBlock(block, nil)
		}

		next++
	}

	c.logger.Infof("Finished syncing with cluster up to block %d (incl.)", b.Header.Number)
	return nil
}

func (c *Chain) serveRaft() {
	ticker := c.clock.NewTicker(c.opts.TickInterval)

	for {
		select {
		case <-ticker.C():
			c.node.Tick()

		case rd := <-c.node.Ready():
			if err := c.storage.Store(rd.Entries, rd.HardState, rd.Snapshot); err != nil {
				c.logger.Panicf("Failed to persist etcd/raft data: %s", err)
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				c.snapC <- &rd.Snapshot

				b := utils.UnmarshalBlockOrPanic(rd.Snapshot.Data)
				c.lastSnapBlockNum = b.Header.Number
				c.confState = rd.Snapshot.Metadata.ConfState
				c.appliedIndex = rd.Snapshot.Metadata.Index
			}

			c.apply(rd.CommittedEntries)
			c.node.Advance()

			
			
			c.send(rd.Messages)

			if rd.SoftState != nil {
				newLead := atomic.LoadUint64(&rd.SoftState.Lead)
				lead := atomic.LoadUint64(&c.leader)
				if newLead != lead {
					c.logger.Infof("Raft leader changed: %d -> %d", lead, newLead)
					atomic.StoreUint64(&c.leader, newLead)

					if lead == c.raftID {
						c.resignC <- struct{}{}
					}

					
					select {
					case c.observeC <- newLead:
					default:
					}
				}
			}

		case <-c.haltC:
			ticker.Stop()
			c.node.Stop()
			c.storage.Close()
			c.logger.Infof("Raft node stopped")
			close(c.doneC) 
			return
		}
	}
}

func (c *Chain) apply(ents []raftpb.Entry) {
	if len(ents) == 0 {
		return
	}

	if ents[0].Index > c.appliedIndex+1 {
		c.logger.Panicf("first index of committed entry[%d] should <= appliedIndex[%d]+1", ents[0].Index, c.appliedIndex)
	}

	var appliedb uint64
	var position int
	for i := range ents {
		switch ents[i].Type {
		case raftpb.EntryNormal:
			
			
			if len(ents[i].Data) == 0 || ents[i].Index <= c.appliedIndex {
				break
			}

			b := utils.UnmarshalBlockOrPanic(ents[i].Data)
			
			
			
			c.raftMetadataLock.RLock()
			m := c.opts.RaftMetadata
			c.raftMetadataLock.RUnlock()

			isConfigMembershipUpdate, err := IsMembershipUpdate(b, m)
			if err != nil {
				c.logger.Warnf("Error while attempting to determine membership update, due to %s", err)
			}
			
			
			if isConfigMembershipUpdate {
				
				
				atomic.StoreUint32(&c.configChangeInProgress, uint32(1))
			}

			c.commitC <- block{b, ents[i].Index}

			appliedb = b.Header.Number
			position = i

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			if err := cc.Unmarshal(ents[i].Data); err != nil {
				c.logger.Warnf("Failed to unmarshal ConfChange data: %s", err)
				continue
			}

			c.confState = *c.node.ApplyConfChange(cc)

			
			
			isConfChangeInProgress := atomic.LoadUint32(&c.configChangeInProgress)
			if isConfChangeInProgress == 1 {
				
				c.configChangeApplyC <- struct{}{}
			}
		}

		if ents[i].Index > c.appliedIndex {
			c.appliedIndex = ents[i].Index
		}
	}

	if c.opts.SnapInterval == 0 || appliedb == 0 {
		
		
		return
	}

	if appliedb-c.lastSnapBlockNum >= c.opts.SnapInterval {
		c.logger.Infof("Taking snapshot at block %d, last snapshotted block number is %d", appliedb, c.lastSnapBlockNum)
		if err := c.storage.TakeSnapshot(c.appliedIndex, &c.confState, ents[position].Data); err != nil {
			c.logger.Fatalf("Failed to create snapshot at index %d", c.appliedIndex)
		}

		c.lastSnapBlockNum = appliedb
	}
}

func (c *Chain) send(msgs []raftpb.Message) {
	for _, msg := range msgs {
		if msg.To == 0 {
			continue
		}

		status := raft.SnapshotFinish

		msgBytes := utils.MarshalOrPanic(&msg)
		_, err := c.rpc.Step(msg.To, &orderer.StepRequest{Channel: c.support.ChainID(), Payload: msgBytes})
		if err != nil {
			
			c.logger.Errorf("Failed to send StepRequest to %d, because: %s", msg.To, err)

			status = raft.SnapshotFailure
		}

		if msg.Type == raftpb.MsgSnap {
			c.node.ReportSnapshot(msg.To, status)
		}
	}
}

func (c *Chain) isConfig(env *common.Envelope) bool {
	h, err := utils.ChannelHeader(env)
	if err != nil {
		c.logger.Panicf("failed to extract channel header from envelope")
	}

	return h.Type == int32(common.HeaderType_CONFIG) || h.Type == int32(common.HeaderType_ORDERER_TRANSACTION)
}

func (c *Chain) configureComm() error {
	nodes, err := c.remotePeers()
	if err != nil {
		return err
	}

	c.configurator.Configure(c.channelID, nodes)
	return nil
}

func (c *Chain) remotePeers() ([]cluster.RemoteNode, error) {
	var nodes []cluster.RemoteNode
	for raftID, consenter := range c.opts.RaftMetadata.Consenters {
		
		if raftID == c.raftID {
			continue
		}
		serverCertAsDER, err := c.pemToDER(consenter.ServerTlsCert, raftID, "server")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		clientCertAsDER, err := c.pemToDER(consenter.ClientTlsCert, raftID, "client")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		nodes = append(nodes, cluster.RemoteNode{
			ID:            raftID,
			Endpoint:      fmt.Sprintf("%s:%d", consenter.Host, consenter.Port),
			ServerTLSCert: serverCertAsDER,
			ClientTLSCert: clientCertAsDER,
		})
	}
	return nodes, nil
}

func (c *Chain) pemToDER(pemBytes []byte, id uint64, certType string) ([]byte, error) {
	bl, _ := pem.Decode(pemBytes)
	if bl == nil {
		c.logger.Errorf("Rejecting PEM block of %s TLS cert for node %d, offending PEM is: %s", certType, id, string(pemBytes))
		return nil, errors.Errorf("invalid PEM block")
	}
	return bl.Bytes, nil
}


func (c *Chain) checkConsentersSet(configValue *common.ConfigValue) error {
	
	updatedMetadata, err := MetadataFromConfigValue(configValue)
	if err != nil {
		return err
	}

	c.raftMetadataLock.RLock()
	changes := ComputeMembershipChanges(c.opts.RaftMetadata.Consenters, updatedMetadata.Consenters)
	c.raftMetadataLock.RUnlock()

	if changes.TotalChanges > 1 {
		return errors.New("update of more than one consenters at a time is not supported")
	}

	return nil
}



func (c *Chain) updateMembership(metadata *etcdraft.RaftMetadata, change *raftpb.ConfChange) error {
	lead := atomic.LoadUint64(&c.leader)
	
	if lead == c.raftID {
		if err := c.node.ProposeConfChange(context.TODO(), *change); err != nil {
			return errors.Errorf("failed to propose configuration update to Raft node: %s", err)
		}
	}

	var err error

	select {
	case <-c.configChangeApplyC:
		
		c.raftMetadataLock.Lock()
		c.opts.RaftMetadata = metadata
		c.raftMetadataLock.Unlock()

		
		err = c.configureComm()
	case <-c.resignC:
		
		
		c.logger.Debug("Raft cluster leader has changed, new leader should re-propose config change based on last config block")
	case <-c.doneC:
		c.logger.Debug("shutting down node, aborting config change update")
	}

	
	atomic.StoreUint32(&c.configChangeInProgress, uint32(0))
	return err
}




func (c *Chain) writeConfigBlock(b block) error {
	metadata, err := ConsensusMetadataFromConfigBlock(b.b)
	if err != nil {
		c.logger.Panicf("error reading consensus metadata, because of %s", err)
	}

	c.raftMetadataLock.RLock()
	raftMetadata := proto.Clone(c.opts.RaftMetadata).(*etcdraft.RaftMetadata)
	
	
	if raftMetadata.Consenters == nil {
		raftMetadata.Consenters = map[uint64]*etcdraft.Consenter{}
	}
	c.raftMetadataLock.RUnlock()

	var changes *MembershipChanges
	if metadata != nil {
		changes = ComputeMembershipChanges(raftMetadata.Consenters, metadata.Consenters)
	}

	confChange := changes.UpdateRaftMetadataAndConfChange(raftMetadata)
	raftMetadata.RaftIndex = b.i

	raftMetadataBytes := utils.MarshalOrPanic(raftMetadata)
	
	c.support.WriteConfigBlock(b.b, raftMetadataBytes)
	if confChange != nil {
		if err := c.updateMembership(raftMetadata, confChange); err != nil {
			return errors.Wrap(err, "failed to update Raft with consenters membership changes")
		}
	}
	return nil
}
