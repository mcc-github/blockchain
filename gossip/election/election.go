/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package election

import (
	"bytes"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcc-github/blockchain/gossip/util"
)






















































type LeaderElectionAdapter interface {
	
	Gossip(Msg)

	
	Accept() <-chan Msg

	
	CreateMessage(isDeclaration bool) Msg

	
	Peers() []Peer

	
	ReportMetrics(isLeader bool)
}

type leadershipCallback func(isLeader bool)


type LeaderElectionService interface {
	
	IsLeader() bool

	
	Stop()

	
	
	Yield()
}

type peerID []byte

func (p peerID) String() string {
	if p == nil {
		return "<nil>"
	}
	return hex.EncodeToString(p)
}


type Peer interface {
	
	ID() peerID
}


type Msg interface {
	
	SenderID() peerID
	
	IsProposal() bool
	
	IsDeclaration() bool
}

func noopCallback(_ bool) {
}

const (
	DefStartupGracePeriod       = time.Second * 15
	DefMembershipSampleInterval = time.Second
	DefLeaderAliveThreshold     = time.Second * 10
	DefLeaderElectionDuration   = time.Second * 5
)

type ElectionConfig struct {
	StartupGracePeriod       time.Duration
	MembershipSampleInterval time.Duration
	LeaderAliveThreshold     time.Duration
	LeaderElectionDuration   time.Duration
}


func NewLeaderElectionService(adapter LeaderElectionAdapter, id string, callback leadershipCallback, config ElectionConfig) LeaderElectionService {
	if len(id) == 0 {
		panic("Empty id")
	}
	le := &leaderElectionSvcImpl{
		id:            peerID(id),
		proposals:     util.NewSet(),
		adapter:       adapter,
		stopChan:      make(chan struct{}),
		interruptChan: make(chan struct{}, 1),
		logger:        util.GetLogger(util.ElectionLogger, ""),
		callback:      noopCallback,
		config:        config,
	}

	if callback != nil {
		le.callback = callback
	}

	go le.start()
	return le
}


type leaderElectionSvcImpl struct {
	id        peerID
	proposals *util.Set
	sync.Mutex
	stopChan      chan struct{}
	interruptChan chan struct{}
	stopWG        sync.WaitGroup
	isLeader      int32
	leaderExists  int32
	yield         int32
	sleeping      bool
	adapter       LeaderElectionAdapter
	logger        util.Logger
	callback      leadershipCallback
	yieldTimer    *time.Timer
	config        ElectionConfig
}

func (le *leaderElectionSvcImpl) start() {
	le.stopWG.Add(2)
	go le.handleMessages()
	le.waitForMembershipStabilization(le.config.StartupGracePeriod)
	go le.run()
}

func (le *leaderElectionSvcImpl) handleMessages() {
	le.logger.Debug(le.id, ": Entering")
	defer le.logger.Debug(le.id, ": Exiting")
	defer le.stopWG.Done()
	msgChan := le.adapter.Accept()
	for {
		select {
		case <-le.stopChan:
			return
		case msg := <-msgChan:
			if !le.isAlive(msg.SenderID()) {
				le.logger.Debug(le.id, ": Got message from", msg.SenderID(), "but it is not in the view")
				break
			}
			le.handleMessage(msg)
		}
	}
}

func (le *leaderElectionSvcImpl) handleMessage(msg Msg) {
	msgType := "proposal"
	if msg.IsDeclaration() {
		msgType = "declaration"
	}
	le.logger.Debug(le.id, ":", msg.SenderID(), "sent us", msgType)
	le.Lock()
	defer le.Unlock()

	if msg.IsProposal() {
		le.proposals.Add(string(msg.SenderID()))
	} else if msg.IsDeclaration() {
		atomic.StoreInt32(&le.leaderExists, int32(1))
		if le.sleeping && len(le.interruptChan) == 0 {
			le.interruptChan <- struct{}{}
		}
		if bytes.Compare(msg.SenderID(), le.id) < 0 && le.IsLeader() {
			le.stopBeingLeader()
		}
	} else {
		
		le.logger.Error("Got a message that's not a proposal and not a declaration")
	}
}



func (le *leaderElectionSvcImpl) waitForInterrupt(timeout time.Duration) {
	le.logger.Debug(le.id, ": Entering")
	defer le.logger.Debug(le.id, ": Exiting")
	le.Lock()
	le.sleeping = true
	le.Unlock()

	select {
	case <-le.interruptChan:
	case <-le.stopChan:
	case <-time.After(timeout):
	}

	le.Lock()
	le.sleeping = false
	
	
	
	le.drainInterruptChannel()
	le.Unlock()
}

func (le *leaderElectionSvcImpl) run() {
	defer le.stopWG.Done()
	for !le.shouldStop() {
		if !le.isLeaderExists() {
			le.leaderElection()
		}
		
		
		if le.isLeaderExists() && le.isYielding() {
			le.stopYielding()
		}
		if le.shouldStop() {
			return
		}
		if le.IsLeader() {
			le.leader()
		} else {
			le.follower()
		}
	}
}

func (le *leaderElectionSvcImpl) leaderElection() {
	le.logger.Debug(le.id, ": Entering")
	defer le.logger.Debug(le.id, ": Exiting")
	
	
	if le.isYielding() {
		return
	}
	
	le.propose()
	
	le.waitForInterrupt(le.config.LeaderElectionDuration)
	
	
	if le.isLeaderExists() {
		le.logger.Info(le.id, ": Some peer is already a leader")
		return
	}

	if le.isYielding() {
		le.logger.Debug(le.id, ": Aborting leader election because yielding")
		return
	}
	
	
	for _, o := range le.proposals.ToArray() {
		id := o.(string)
		if bytes.Compare(peerID(id), le.id) < 0 {
			return
		}
	}
	
	
	le.beLeader()
	atomic.StoreInt32(&le.leaderExists, int32(1))
}


func (le *leaderElectionSvcImpl) propose() {
	le.logger.Debug(le.id, ": Entering")
	le.logger.Debug(le.id, ": Exiting")
	leadershipProposal := le.adapter.CreateMessage(false)
	le.adapter.Gossip(leadershipProposal)
}

func (le *leaderElectionSvcImpl) follower() {
	le.logger.Debug(le.id, ": Entering")
	defer le.logger.Debug(le.id, ": Exiting")

	le.proposals.Clear()
	atomic.StoreInt32(&le.leaderExists, int32(0))
	le.adapter.ReportMetrics(false)
	select {
	case <-time.After(le.config.LeaderAliveThreshold):
	case <-le.stopChan:
	}
}

func (le *leaderElectionSvcImpl) leader() {
	leaderDeclaration := le.adapter.CreateMessage(true)
	le.adapter.Gossip(leaderDeclaration)
	le.adapter.ReportMetrics(true)
	le.waitForInterrupt(le.config.LeaderAliveThreshold / 2)
}



func (le *leaderElectionSvcImpl) waitForMembershipStabilization(timeLimit time.Duration) {
	le.logger.Debug(le.id, ": Entering")
	defer le.logger.Debug(le.id, ": Exiting, peers found", len(le.adapter.Peers()))
	endTime := time.Now().Add(timeLimit)
	viewSize := len(le.adapter.Peers())
	for !le.shouldStop() {
		time.Sleep(le.config.MembershipSampleInterval)
		newSize := len(le.adapter.Peers())
		if newSize == viewSize || time.Now().After(endTime) || le.isLeaderExists() {
			return
		}
		viewSize = newSize
	}
}



func (le *leaderElectionSvcImpl) drainInterruptChannel() {
	if len(le.interruptChan) == 1 {
		<-le.interruptChan
	}
}


func (le *leaderElectionSvcImpl) isAlive(id peerID) bool {
	for _, p := range le.adapter.Peers() {
		if bytes.Equal(p.ID(), id) {
			return true
		}
	}
	return false
}

func (le *leaderElectionSvcImpl) isLeaderExists() bool {
	return atomic.LoadInt32(&le.leaderExists) == int32(1)
}


func (le *leaderElectionSvcImpl) IsLeader() bool {
	isLeader := atomic.LoadInt32(&le.isLeader) == int32(1)
	le.logger.Debug(le.id, ": Returning", isLeader)
	return isLeader
}

func (le *leaderElectionSvcImpl) beLeader() {
	le.logger.Info(le.id, ": Becoming a leader")
	atomic.StoreInt32(&le.isLeader, int32(1))
	le.callback(true)
}

func (le *leaderElectionSvcImpl) stopBeingLeader() {
	le.logger.Info(le.id, "Stopped being a leader")
	atomic.StoreInt32(&le.isLeader, int32(0))
	le.callback(false)
}

func (le *leaderElectionSvcImpl) shouldStop() bool {
	select {
	case <-le.stopChan:
		return true
	default:
		return false
	}
}

func (le *leaderElectionSvcImpl) isYielding() bool {
	return atomic.LoadInt32(&le.yield) == int32(1)
}

func (le *leaderElectionSvcImpl) stopYielding() {
	le.logger.Debug("Stopped yielding")
	le.Lock()
	defer le.Unlock()
	atomic.StoreInt32(&le.yield, int32(0))
	le.yieldTimer.Stop()
}



func (le *leaderElectionSvcImpl) Yield() {
	le.Lock()
	defer le.Unlock()
	if !le.IsLeader() || le.isYielding() {
		return
	}
	
	atomic.StoreInt32(&le.yield, int32(1))
	
	le.stopBeingLeader()
	
	atomic.StoreInt32(&le.leaderExists, int32(0))
	
	le.yieldTimer = time.AfterFunc(le.config.LeaderAliveThreshold*6, func() {
		atomic.StoreInt32(&le.yield, int32(0))
	})
}


func (le *leaderElectionSvcImpl) Stop() {
	select {
	case <-le.stopChan:
	default:
		close(le.stopChan)
		le.logger.Debug(le.id, ": Entering")
		defer le.logger.Debug(le.id, ": Exiting")
		le.stopWG.Wait()
	}
}
