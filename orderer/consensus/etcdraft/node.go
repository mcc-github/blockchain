/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"context"
	"crypto/sha256"
	"sync"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/clock"
	"github.com/gogo/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protoutil"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
)

type node struct {
	chainID string
	logger  *flogging.FabricLogger
	metrics *Metrics

	unreachableLock sync.RWMutex
	unreachable     map[uint64]struct{}

	storage *RaftStorage
	config  *raft.Config

	rpc RPC

	chain *Chain

	tickInterval time.Duration
	clock        clock.Clock

	metadata *etcdraft.BlockMetadata

	subscriberC chan chan uint64

	raft.Node
}

func (n *node) start(fresh, join bool) {
	raftPeers := RaftPeers(n.metadata.ConsenterIds)
	n.logger.Debugf("Starting raft node: #peers: %v", len(raftPeers))

	var campaign bool
	if fresh {
		if join {
			raftPeers = nil
			n.logger.Info("Starting raft node to join an existing channel")
		} else {
			n.logger.Info("Starting raft node as part of a new channel")

			
			
			sha := sha256.Sum256([]byte(n.chainID))
			number, _ := proto.DecodeVarint(sha[24:])
			if n.config.ID == number%uint64(len(raftPeers))+1 {
				campaign = true
			}
		}
		n.Node = raft.StartNode(n.config, raftPeers)
	} else {
		n.logger.Info("Restarting raft node")
		n.Node = raft.RestartNode(n.config)
	}

	n.subscriberC = make(chan chan uint64)

	go n.run(campaign)
}

func (n *node) run(campaign bool) {
	raftTicker := n.clock.NewTicker(n.tickInterval)

	if s := n.storage.Snapshot(); !raft.IsEmptySnap(s) {
		n.chain.snapC <- &s
	}

	elected := make(chan struct{})
	if campaign {
		n.logger.Infof("This node is picked to start campaign")
		go func() {
			
			
			
			
			
			
			
			
			campaignTicker := n.clock.NewTicker(n.tickInterval * time.Duration(n.config.HeartbeatTick) * 2)
			defer campaignTicker.Stop()

			for {
				select {
				case <-campaignTicker.C():
					n.Campaign(context.TODO())
				case <-elected:
					return
				case <-n.chain.doneC:
					return
				}
			}
		}()
	}

	var notifyLeaderChangeC chan uint64

	for {
		select {
		case <-raftTicker.C():
			n.Tick()

		case rd := <-n.Ready():
			startStoring := n.clock.Now()
			if err := n.storage.Store(rd.Entries, rd.HardState, rd.Snapshot); err != nil {
				n.logger.Panicf("Failed to persist etcd/raft data: %s", err)
			}
			duration := n.clock.Since(startStoring).Seconds()
			n.metrics.DataPersistDuration.Observe(float64(duration))

			if !raft.IsEmptySnap(rd.Snapshot) {
				n.chain.snapC <- &rd.Snapshot
			}

			if notifyLeaderChangeC != nil && rd.SoftState != nil {
				if l := atomic.LoadUint64(&rd.SoftState.Lead); l != raft.None {
					select {
					case notifyLeaderChangeC <- l:
					default:
					}

					notifyLeaderChangeC = nil
				}
			}

			
			if len(rd.CommittedEntries) != 0 || rd.SoftState != nil {
				n.chain.applyC <- apply{rd.CommittedEntries, rd.SoftState}
			}

			if campaign && rd.SoftState != nil {
				leader := atomic.LoadUint64(&rd.SoftState.Lead) 
				if leader != raft.None {
					n.logger.Infof("Leader %d is present, quit campaign", leader)
					campaign = false
					close(elected)
				}
			}

			n.Advance()

			
			
			n.send(rd.Messages)

		case notifyLeaderChangeC = <-n.subscriberC:

		case <-n.chain.haltC:
			raftTicker.Stop()
			n.Stop()
			n.storage.Close()
			n.logger.Infof("Raft node stopped")
			close(n.chain.doneC) 
			return
		}
	}
}

func (n *node) send(msgs []raftpb.Message) {
	n.unreachableLock.RLock()
	defer n.unreachableLock.RUnlock()

	for _, msg := range msgs {
		if msg.To == 0 {
			continue
		}

		status := raft.SnapshotFinish

		msgBytes := protoutil.MarshalOrPanic(&msg)
		err := n.rpc.SendConsensus(msg.To, &orderer.ConsensusRequest{Channel: n.chainID, Payload: msgBytes})
		if err != nil {
			n.ReportUnreachable(msg.To)
			n.logSendFailure(msg.To, err)

			status = raft.SnapshotFailure
		} else if _, ok := n.unreachable[msg.To]; ok {
			n.logger.Infof("Successfully sent StepRequest to %d after failed attempt(s)", msg.To)
			delete(n.unreachable, msg.To)
		}

		if msg.Type == raftpb.MsgSnap {
			n.ReportSnapshot(msg.To, status)
		}
	}
}





func (n *node) abdicateLeader(currentLead uint64) {
	status := n.Status()

	if status.Lead != raft.None && status.Lead != currentLead {
		n.logger.Warn("Leader has changed since asked to transfer leadership")
		return
	}

	
	notifyc := make(chan uint64, 1)
	select {
	case n.subscriberC <- notifyc:
	case <-n.chain.doneC:
		return
	}

	
	if status.RaftState == raft.StateLeader {
		var transferee uint64
		for id, pr := range status.Progress {
			if id == status.ID {
				continue 
			}

			if pr.RecentActive && !pr.Paused {
				transferee = id
				break
			}

			n.logger.Debugf("Node %d is not qualified as transferee because it's either paused or not active", id)
		}

		if transferee == raft.None {
			n.logger.Errorf("No follower is qualified as transferee, abort leader transfer")
			return
		}

		n.logger.Infof("Transferring leadership to %d", transferee)
		n.TransferLeadership(context.TODO(), status.ID, transferee)
	}

	timer := n.clock.NewTimer(time.Duration(n.config.ElectionTick) * n.tickInterval)
	defer timer.Stop() 

	select {
	case <-timer.C():
		n.logger.Warn("Leader transfer timeout")
	case l := <-notifyc:
		n.logger.Infof("Leader has been transferred from %d to %d", currentLead, l)
	case <-n.chain.doneC:
	}
}

func (n *node) logSendFailure(dest uint64, err error) {
	if _, ok := n.unreachable[dest]; ok {
		n.logger.Debugf("Failed to send StepRequest to %d, because: %s", dest, err)
		return
	}

	n.logger.Errorf("Failed to send StepRequest to %d, because: %s", dest, err)
	n.unreachable[dest] = struct{}{}
}

func (n *node) takeSnapshot(index uint64, cs raftpb.ConfState, data []byte) {
	if err := n.storage.TakeSnapshot(index, cs, data); err != nil {
		n.logger.Errorf("Failed to create snapshot at index %d: %s", index, err)
	}
}

func (n *node) lastIndex() uint64 {
	i, _ := n.storage.ram.LastIndex()
	return i
}
