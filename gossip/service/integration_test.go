/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/stretchr/testify/assert"
)

type transientStoreMock struct {
}

func (*transientStoreMock) PurgeByHeight(maxBlockNumToRetain uint64) error {
	return nil
}

func (*transientStoreMock) Persist(txid string, blockHeight uint64, privateSimulationResults *rwset.TxPvtReadWriteSet) error {
	panic("implement me")
}

func (*transientStoreMock) PersistWithConfig(txid string, blockHeight uint64, privateSimulationResultsWithConfig *transientstore2.TxPvtReadWriteSetWithConfigInfo) error {
	panic("implement me")
}

func (*transientStoreMock) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	panic("implement me")
}

func (*transientStoreMock) PurgeByTxids(txids []string) error {
	panic("implement me")
}

type embeddingDeliveryService struct {
	startOnce sync.Once
	stopOnce  sync.Once
	deliverservice.DeliverService
	startSignal sync.WaitGroup
	stopSignal  sync.WaitGroup
}

func newEmbeddingDeliveryService(ds deliverservice.DeliverService) *embeddingDeliveryService {
	eds := &embeddingDeliveryService{
		DeliverService: ds,
	}
	eds.startSignal.Add(1)
	eds.stopSignal.Add(1)
	return eds
}

func (eds *embeddingDeliveryService) waitForDeliveryServiceActivation() {
	eds.startSignal.Wait()
}

func (eds *embeddingDeliveryService) waitForDeliveryServiceTermination() {
	eds.stopSignal.Wait()
}

func (eds *embeddingDeliveryService) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error {
	eds.startOnce.Do(func() {
		eds.startSignal.Done()
	})
	return eds.DeliverService.StartDeliverForChannel(chainID, ledgerInfo, finalizer)
}

func (eds *embeddingDeliveryService) StopDeliverForChannel(chainID string) error {
	eds.stopOnce.Do(func() {
		eds.stopSignal.Done()
	})
	return eds.DeliverService.StopDeliverForChannel(chainID)
}

func (eds *embeddingDeliveryService) Stop() {
	eds.DeliverService.Stop()
}

type embeddingDeliveryServiceFactory struct {
	DeliveryServiceFactory
}

func (edsf *embeddingDeliveryServiceFactory) Service(g GossipServiceAdapter, endpoints []string, mcs api.MessageCryptoService) (deliverservice.DeliverService, error) {
	ds, err := edsf.DeliveryServiceFactory.Service(g, endpoints, mcs)
	if err != nil {
		panic(err)
	}
	return newEmbeddingDeliveryService(ds), nil
}

func TestLeaderYield(t *testing.T) {
	
	
	
	
	takeOverMaxTimeout := time.Minute
	
	
	
	serviceConfig := &ServiceConfig{
		UseLeaderElection:          true,
		OrgLeader:                  false,
		ElectionStartupGracePeriod: election.DefStartupGracePeriod,
		
		
		ElectionMembershipSampleInterval: time.Millisecond * 100,
		ElectionLeaderAliveThreshold:     time.Second * 5,
		
		
		ElectionLeaderElectionDuration: time.Millisecond * 500,
	}
	n := 2
	gossips := startPeers(t, serviceConfig, n, 0, 1)
	defer stopPeers(gossips)
	channelName := "channelA"
	peerIndexes := []int{0, 1}
	
	addPeersToChannel(t, n, channelName, gossips, peerIndexes)
	
	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Millisecond*100)

	endpoint, socket := getAvailablePort(t)
	socket.Close()

	
	newGossipService := func(i int) *GossipService {
		gs := gossips[i].GossipService
		gs.deliveryFactory = &embeddingDeliveryServiceFactory{&deliveryFactoryImpl{
			credentialSupport: comm.NewCredentialSupport(),
			deliverServiceConfig: &deliverservice.DeliverServiceConfig{
				PeerTLSEnabled:              false,
				ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
				ReconnectTotalTimeThreshold: time.Second,
				ConnectionTimeout:           time.Millisecond * 100,
			},
		}}
		gs.InitializeChannel(channelName, []string{endpoint}, Support{
			Committer: &mockLedgerInfo{1},
			Store:     &transientStoreMock{},
		})
		return gs
	}

	
	
	pkiID0 := gossips[0].peerIdentity
	pkiID1 := gossips[1].peerIdentity
	var firstLeaderIdx, secondLeaderIdx int
	if bytes.Compare(pkiID0, pkiID1) < 0 {
		firstLeaderIdx = 0
		secondLeaderIdx = 1
	} else {
		firstLeaderIdx = 1
		secondLeaderIdx = 0
	}
	p0 := newGossipService(firstLeaderIdx)
	p1 := newGossipService(secondLeaderIdx)

	
	getLeader := func() int {
		p0.lock.RLock()
		p1.lock.RLock()
		defer p0.lock.RUnlock()
		defer p1.lock.RUnlock()

		if p0.leaderElection[channelName].IsLeader() {
			return 0
		}
		if p1.leaderElection[channelName].IsLeader() {
			return 1
		}
		return -1
	}

	ds0 := p0.deliveryService[channelName].(*embeddingDeliveryService)

	
	ds0.waitForDeliveryServiceActivation()
	t.Log("p0 started its delivery service")
	
	assert.Equal(t, 0, getLeader())
	
	ds0.waitForDeliveryServiceTermination()
	t.Log("p0 stopped its delivery service")
	
	assert.NotEqual(t, 0, getLeader())
	
	timeLimit := time.Now().Add(takeOverMaxTimeout)
	for getLeader() != 1 && time.Now().Before(timeLimit) {
		time.Sleep(100 * time.Millisecond)
	}
	if time.Now().After(timeLimit) && getLeader() != 1 {
		util.PrintStackTrace()
		t.Fatalf("p1 hasn't taken over leadership within %v: %d", takeOverMaxTimeout, getLeader())
	}
	t.Log("p1 has taken over leadership")
	p0.chains[channelName].Stop()
	p1.chains[channelName].Stop()
	p0.deliveryService[channelName].Stop()
	p1.deliveryService[channelName].Stop()
}
