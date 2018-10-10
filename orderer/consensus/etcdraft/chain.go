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





type Storage interface {
	raft.Storage
	Append(entries []raftpb.Entry) error
}


type Options struct {
	RaftID uint64

	Clock clock.Clock

	Storage Storage
	Logger  *flogging.FabricLogger

	TickInterval    time.Duration
	ElectionTick    int
	HeartbeatTick   int
	MaxSizePerMsg   uint64
	MaxInflightMsgs int
	Peers           []raft.Peer
}


type Chain struct {
	communication cluster.Communicator
	raftID        uint64

	submitC  chan *orderer.SubmitRequest
	commitC  chan *common.Block
	observeC chan<- uint64 
	haltC    chan struct{}
	doneC    chan struct{}

	clock clock.Clock

	support consensus.ConsenterSupport

	leaderLock   sync.RWMutex
	leader       uint64
	appliedIndex uint64

	node    raft.Node
	storage Storage
	opts    Options

	logger *flogging.FabricLogger
}


func NewChain(support consensus.ConsenterSupport, opts Options, observe chan<- uint64, comm cluster.Communicator) (*Chain, error) {
	return &Chain{
		communication: comm,
		raftID:        opts.RaftID,
		submitC:       make(chan *orderer.SubmitRequest),
		commitC:       make(chan *common.Block),
		haltC:         make(chan struct{}),
		doneC:         make(chan struct{}),
		observeC:      observe,
		support:       support,
		clock:         opts.Clock,
		logger:        opts.Logger.With("channel", support.ChainID()),
		storage:       opts.Storage,
		opts:          opts,
	}, nil
}


func (c *Chain) Start() {
	config := &raft.Config{
		ID:              c.raftID,
		ElectionTick:    c.opts.ElectionTick,
		HeartbeatTick:   c.opts.HeartbeatTick,
		MaxSizePerMsg:   c.opts.MaxSizePerMsg,
		MaxInflightMsgs: c.opts.MaxInflightMsgs,
		Logger:          c.logger,
		Storage:         c.opts.Storage,
	}

	if err := c.configureComm(); err != nil {
		c.logger.Errorf("Failed starting chain, aborting: +%v", err)
		close(c.doneC)
		return
	}

	c.node = raft.StartNode(config, c.opts.Peers)

	go c.serveRaft()
	go c.serveRequest()
}


func (c *Chain) Order(env *common.Envelope, configSeq uint64) error {
	return c.Submit(&orderer.SubmitRequest{LastValidationSeq: configSeq, Content: env}, 0)
}


func (c *Chain) Configure(env *common.Envelope, configSeq uint64) error {
	if err := c.checkConfigUpdateValidity(env); err != nil {
		return err
	}
	return c.Submit(&orderer.SubmitRequest{LastValidationSeq: configSeq, Content: env}, 0)
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
		return errors.Errorf("channel creation requests not supported yet")
	case int32(common.HeaderType_CONFIG):
		configEnv, err := configtx.UnmarshalConfigEnvelope(payload.Data)
		if err != nil {
			return err
		}
		configUpdateEnv, err := utils.EnvelopeToConfigUpdate(configEnv.LastUpdate)
		if err != nil {
			return err
		}
		configUpdate, err := configtx.UnmarshalConfigUpdate(configUpdateEnv.ConfigUpdate)
		if err != nil {
			return err
		}

		
		
		if ordererConfigGroup, ok := configUpdate.WriteSet.Groups["Orderer"]; ok {
			if _, ok := ordererConfigGroup.Values["ConsensusType"]; ok {
				return errors.Errorf("updates to ConsensusType not supported currently")
			}
		}
		return nil

	default:
		return errors.Errorf("config transaction has unknown header type")
	}
}


func (c *Chain) WaitReady() error {
	return nil
}


func (c *Chain) Errored() <-chan struct{} {
	return c.doneC
}


func (c *Chain) Halt() {
	select {
	case c.haltC <- struct{}{}:
	case <-c.doneC:
		return
	}
	<-c.doneC
}


func (c *Chain) Step(req *orderer.StepRequest, sender uint64) error {
	panic("not implemented")
}





func (c *Chain) Submit(req *orderer.SubmitRequest, sender uint64) error {
	c.leaderLock.RLock()
	defer c.leaderLock.RUnlock()

	if c.leader == raft.None {
		return errors.Errorf("no raft leader")
	}

	if c.leader == c.raftID {
		select {
		case c.submitC <- req:
			return nil
		case <-c.doneC:
			return errors.Errorf("chain is stopped")
		}
	}

	
	return errors.Errorf("only single raft node is currently supported")
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

		case <-c.doneC:
			c.logger.Infof("Stop serving requests")
			return
		}
	}
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
			return errors.Errorf("failed to propose data to raft: %s", err)
		}

		select {
		case block := <-c.commitC:
			if utils.IsConfigBlock(block) {
				c.support.WriteConfigBlock(block, nil)
			} else {
				c.support.WriteBlock(block, nil)
			}

		case <-c.doneC:
			return nil
		}
	}

	return nil
}

func (c *Chain) serveRaft() {
	ticker := c.clock.NewTicker(c.opts.TickInterval)

	for {
		select {
		case <-ticker.C():
			c.node.Tick()

		case rd := <-c.node.Ready():
			c.storage.Append(rd.Entries)
			
			c.apply(c.entriesToApply(rd.CommittedEntries))
			c.node.Advance()

			if rd.SoftState != nil {
				c.leaderLock.Lock()
				newLead := atomic.LoadUint64(&rd.SoftState.Lead)
				if newLead != c.leader {
					c.logger.Infof("Raft leader changed on node %x: %x -> %x", c.raftID, c.leader, newLead)
					c.leader = newLead

					
					select {
					case c.observeC <- newLead:
					default:
					}
				}
				c.leaderLock.Unlock()
			}

		case <-c.haltC:
			close(c.doneC)
			ticker.Stop()
			c.node.Stop()
			c.logger.Infof("Raft node %x stopped", c.raftID)
			return
		}
	}
}

func (c *Chain) apply(ents []raftpb.Entry) {
	for i := range ents {
		switch ents[i].Type {
		case raftpb.EntryNormal:
			if len(ents[i].Data) == 0 {
				break
			}

			c.commitC <- utils.UnmarshalBlockOrPanic(ents[i].Data)

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			if err := cc.Unmarshal(ents[i].Data); err != nil {
				c.logger.Warnf("Failed to unmarshal ConfChange data: %s", err)
				continue
			}

			c.node.ApplyConfChange(cc)
		}

		c.appliedIndex = ents[i].Index
	}
}


func (c *Chain) entriesToApply(ents []raftpb.Entry) (nents []raftpb.Entry) {
	if len(ents) == 0 {
		return
	}

	firstIdx := ents[0].Index
	if firstIdx > c.appliedIndex+1 {
		c.logger.Panicf("first index of committed entry[%d] should <= progress.appliedIndex[%d]+1", firstIdx, c.appliedIndex)
	}

	
	
	
	
	if c.appliedIndex-firstIdx+1 < uint64(len(ents)) {
		nents = ents[c.appliedIndex-firstIdx+1:]
	}
	return nents
}

func (c *Chain) isConfig(env *common.Envelope) bool {
	h, err := utils.ChannelHeader(env)
	if err != nil {
		c.logger.Panicf("programming error: failed to extract channel header from envelope")
	}

	return h.Type == int32(common.HeaderType_CONFIG) || h.Type == int32(common.HeaderType_ORDERER_TRANSACTION)
}

func (c *Chain) configureComm() error {
	nodes, err := c.nodeConfigFromMetadata()
	if err != nil {
		return err
	}
	c.communication.Configure(c.support.ChainID(), nodes)
	return nil
}

func (c *Chain) nodeConfigFromMetadata() ([]cluster.RemoteNode, error) {
	var nodes []cluster.RemoteNode
	m := &etcdraft.Metadata{}
	if err := proto.Unmarshal(c.support.SharedConfig().ConsensusMetadata(), m); err != nil {
		return nil, errors.Wrap(err, "failed extracting consensus metadata")
	}

	for id, consenter := range m.Consenters {
		raftID := uint64(id + 1)
		
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
