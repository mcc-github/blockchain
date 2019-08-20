/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/api"
	gcomm "github.com/mcc-github/blockchain/gossip/comm"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/gossip/gossip/algo"
	"github.com/mcc-github/blockchain/gossip/gossip/channel"
	gossipmetrics "github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/state"
	"github.com/mcc-github/blockchain/gossip/util"
	peergossip "github.com/mcc-github/blockchain/internal/peer/gossip"
	"github.com/mcc-github/blockchain/internal/peer/gossip/mocks"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/peer"
	transientstore2 "github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func init() {
	util.SetupTestLogging()
}



type signerSerializer interface {
	identity.SignerSerializer
}

type mockTransientStore struct {
}

func (*mockTransientStore) PurgeByHeight(maxBlockNumToRetain uint64) error {
	return nil
}

func (*mockTransientStore) Persist(txid string, blockHeight uint64, privateSimulationResults *rwset.TxPvtReadWriteSet) error {
	panic("implement me")
}

func (*mockTransientStore) PersistWithConfig(txid string, blockHeight uint64, privateSimulationResultsWithConfig *transientstore2.TxPvtReadWriteSetWithConfigInfo) error {
	panic("implement me")
}

func (*mockTransientStore) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (transientstore.RWSetScanner, error) {
	panic("implement me")
}

func (*mockTransientStore) PurgeByTxids(txids []string) error {
	panic("implement me")
}

func TestInitGossipService(t *testing.T) {
	grpcServer := grpc.NewServer()
	endpoint, socket := getAvailablePort(t)

	msptesttools.LoadMSPSetupForTesting()
	signer := mgmt.GetLocalSigningIdentityOrPanic()

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)

	messageCryptoService := peergossip.NewMCS(&mocks.ChannelPolicyManagerGetter{}, signer, mgmt.NewDeserializersManager(), cryptoProvider)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	dialOpts := defaultDeliverClientDialOpts()
	gossipConfig, err := gossip.GlobalConfig(endpoint, nil)
	assert.NoError(t, err)

	gossipService, err := New(
		signer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		endpoint,
		grpcServer,
		messageCryptoService,
		secAdv,
		nil,
		comm.NewCredentialSupport(),
		dialOpts,
		gossipConfig,
		&ServiceConfig{},
		&deliverservice.DeliverServiceConfig{
			ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
			ReconnectTotalTimeThreshold: deliverservice.DefaultReConnectTotalTimeThreshold,
		},
	)
	assert.NoError(t, err)

	go grpcServer.Serve(socket)
	defer grpcServer.Stop()

	defer gossipService.Stop()
}


func TestJCMInterface(t *testing.T) {
	_ = api.JoinChannelMessage(&joinChannelMessage{})
	t.Parallel()
}

func TestLeaderElectionWithDeliverClient(t *testing.T) {
	t.Parallel()
	
	
	
	

	n := 10
	serviceConfig := &ServiceConfig{
		UseLeaderElection:                true,
		OrgLeader:                        false,
		ElectionStartupGracePeriod:       election.DefStartupGracePeriod,
		ElectionMembershipSampleInterval: election.DefMembershipSampleInterval,
		ElectionLeaderAliveThreshold:     election.DefLeaderAliveThreshold,
		ElectionLeaderElectionDuration:   election.DefLeaderElectionDuration,
	}
	gossips := startPeers(t, serviceConfig, n, 0, 1, 2, 3, 4)

	channelName := "chanA"
	peerIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		peerIndexes[i] = i
	}
	addPeersToChannel(t, n, channelName, gossips, peerIndexes)

	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*20, time.Second*2)

	services := make([]*electionService, n)

	for i := 0; i < n; i++ {
		deliverServiceFactory := &mockDeliverServiceFactory{
			service: &mockDeliverService{
				running: make(map[string]bool),
			},
		}
		gossips[i].deliveryFactory = deliverServiceFactory
		deliverServiceFactory.service.running[channelName] = false

		gossips[i].InitializeChannel(channelName, []string{"endpoint"}, Support{
			Store:     &mockTransientStore{},
			Committer: &mockLedgerInfo{1},
		})
		service, exist := gossips[i].leaderElection[channelName]
		assert.True(t, exist, "Leader election service should be created for peer %d and channel %s", i, channelName)
		services[i] = &electionService{nil, false, 0}
		services[i].LeaderElectionService = service
	}

	
	assert.True(t, waitForLeaderElection(t, services, time.Second*30, time.Second*2), "One leader should be selected")

	startsNum := 0
	for i := 0; i < n; i++ {
		
		if gossips[i].deliveryService[channelName].(*mockDeliverService).running[channelName] {
			startsNum++
		}
	}

	assert.Equal(t, 1, startsNum, "Only for one peer delivery client should start")

	stopPeers(gossips)
}

func TestWithStaticDeliverClientLeader(t *testing.T) {
	
	
	
	
	

	serviceConfig := &ServiceConfig{
		UseLeaderElection:                false,
		OrgLeader:                        true,
		ElectionStartupGracePeriod:       election.DefStartupGracePeriod,
		ElectionMembershipSampleInterval: election.DefMembershipSampleInterval,
		ElectionLeaderAliveThreshold:     election.DefLeaderAliveThreshold,
		ElectionLeaderElectionDuration:   election.DefLeaderElectionDuration,
	}
	n := 2
	gossips := startPeers(t, serviceConfig, n, 0, 1)
	channelName := "chanA"
	peerIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		peerIndexes[i] = i
	}

	addPeersToChannel(t, n, channelName, gossips, peerIndexes)

	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Second*2)

	deliverServiceFactory := &mockDeliverServiceFactory{
		service: &mockDeliverService{
			running: make(map[string]bool),
		},
	}

	for i := 0; i < n; i++ {
		gossips[i].deliveryFactory = deliverServiceFactory
		deliverServiceFactory.service.running[channelName] = false
		gossips[i].InitializeChannel(channelName, []string{"endpoint"}, Support{
			Committer: &mockLedgerInfo{1},
			Store:     &mockTransientStore{},
		})
	}

	for i := 0; i < n; i++ {
		assert.NotNil(t, gossips[i].deliveryService[channelName], "Delivery service for channel %s not initiated in peer %d", channelName, i)
		assert.True(t, gossips[i].deliveryService[channelName].(*mockDeliverService).running[channelName], "Block deliverer not started for peer %d", i)
	}

	channelName = "chanB"
	for i := 0; i < n; i++ {
		deliverServiceFactory.service.running[channelName] = false
		gossips[i].InitializeChannel(channelName, []string{"endpoint"}, Support{
			Committer: &mockLedgerInfo{1},
			Store:     &mockTransientStore{},
		})
	}

	for i := 0; i < n; i++ {
		assert.NotNil(t, gossips[i].deliveryService[channelName], "Delivery service for channel %s not initiated in peer %d", channelName, i)
		assert.True(t, gossips[i].deliveryService[channelName].(*mockDeliverService).running[channelName], "Block deliverer not started for peer %d", i)
	}

	stopPeers(gossips)
}

func TestWithStaticDeliverClientNotLeader(t *testing.T) {

	serviceConfig := &ServiceConfig{
		UseLeaderElection:                false,
		OrgLeader:                        false,
		ElectionStartupGracePeriod:       election.DefStartupGracePeriod,
		ElectionMembershipSampleInterval: election.DefMembershipSampleInterval,
		ElectionLeaderAliveThreshold:     election.DefLeaderAliveThreshold,
		ElectionLeaderElectionDuration:   election.DefLeaderElectionDuration,
	}
	n := 2
	gossips := startPeers(t, serviceConfig, n, 0, 1)

	channelName := "chanA"
	peerIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		peerIndexes[i] = i
	}

	addPeersToChannel(t, n, channelName, gossips, peerIndexes)

	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Second*2)

	deliverServiceFactory := &mockDeliverServiceFactory{
		service: &mockDeliverService{
			running: make(map[string]bool),
		},
	}

	for i := 0; i < n; i++ {
		gossips[i].deliveryFactory = deliverServiceFactory
		deliverServiceFactory.service.running[channelName] = false
		gossips[i].InitializeChannel(channelName, []string{"endpoint"}, Support{
			Committer: &mockLedgerInfo{1},
			Store:     &mockTransientStore{},
		})
	}

	for i := 0; i < n; i++ {
		assert.NotNil(t, gossips[i].deliveryService[channelName], "Delivery service for channel %s not initiated in peer %d", channelName, i)
		assert.False(t, gossips[i].deliveryService[channelName].(*mockDeliverService).running[channelName], "Block deliverer should not be started for peer %d", i)
	}

	stopPeers(gossips)
}

func TestWithStaticDeliverClientBothStaticAndLeaderElection(t *testing.T) {

	serviceConfig := &ServiceConfig{
		UseLeaderElection:                true,
		OrgLeader:                        true,
		ElectionStartupGracePeriod:       election.DefStartupGracePeriod,
		ElectionMembershipSampleInterval: election.DefMembershipSampleInterval,
		ElectionLeaderAliveThreshold:     election.DefLeaderAliveThreshold,
		ElectionLeaderElectionDuration:   election.DefLeaderElectionDuration,
	}
	n := 2
	gossips := startPeers(t, serviceConfig, n, 0, 1)

	channelName := "chanA"
	peerIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		peerIndexes[i] = i
	}

	addPeersToChannel(t, n, channelName, gossips, peerIndexes)

	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Second*2)

	deliverServiceFactory := &mockDeliverServiceFactory{
		service: &mockDeliverService{
			running: make(map[string]bool),
		},
	}

	for i := 0; i < n; i++ {
		gossips[i].deliveryFactory = deliverServiceFactory
		assert.Panics(t, func() {
			gossips[i].InitializeChannel(channelName, []string{"endpoint"}, Support{
				Committer: &mockLedgerInfo{1},
				Store:     &mockTransientStore{},
			})
		}, "Dynamic leader election based and static connection to ordering service can't exist simultaneously")
	}

	stopPeers(gossips)
}

type mockDeliverServiceFactory struct {
	service *mockDeliverService
}

func (mf *mockDeliverServiceFactory) Service(g GossipServiceAdapter, endpoints []string, mcs api.MessageCryptoService) (deliverservice.DeliverService, error) {
	return mf.service, nil
}

type mockDeliverService struct {
	running map[string]bool
}

func (ds *mockDeliverService) UpdateEndpoints(chainID string, endpoints []string) error {
	panic("implement me")
}

func (ds *mockDeliverService) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error {
	ds.running[chainID] = true
	return nil
}

func (ds *mockDeliverService) StopDeliverForChannel(chainID string) error {
	ds.running[chainID] = false
	return nil
}

func (ds *mockDeliverService) Stop() {
}

type mockLedgerInfo struct {
	Height uint64
}

func (li *mockLedgerInfo) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	panic("implement me")
}

func (li *mockLedgerInfo) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	panic("implement me")
}

func (li *mockLedgerInfo) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	panic("implement me")
}

func (li *mockLedgerInfo) CommitLegacy(blockAndPvtData *ledger.BlockAndPvtData, commitOpts *ledger.CommitOptions) error {
	panic("implement me")
}

func (li *mockLedgerInfo) CommitPvtDataOfOldBlocks(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error) {
	panic("implement me")
}

func (li *mockLedgerInfo) GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error) {
	panic("implement me")
}


func (li *mockLedgerInfo) LedgerHeight() (uint64, error) {
	return li.Height, nil
}

func (li *mockLedgerInfo) DoesPvtDataInfoExistInLedger(blkNum uint64) (bool, error) {
	return false, nil
}


func (li *mockLedgerInfo) Commit(block *common.Block) error {
	return nil
}


func (li *mockLedgerInfo) GetBlocks(blockSeqs []uint64) []*common.Block {
	return make([]*common.Block, 0)
}


func (li *mockLedgerInfo) Close() {
}

func TestLeaderElectionWithRealGossip(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	

	
	serviceConfig := &ServiceConfig{
		UseLeaderElection:                false,
		OrgLeader:                        false,
		ElectionStartupGracePeriod:       election.DefStartupGracePeriod,
		ElectionMembershipSampleInterval: election.DefMembershipSampleInterval,
		ElectionLeaderAliveThreshold:     election.DefLeaderAliveThreshold,
		ElectionLeaderElectionDuration:   election.DefLeaderElectionDuration,
	}

	n := 10
	gossips := startPeers(t, serviceConfig, n, 0, 1, 2, 3, 4)
	
	channelName := "chanA"
	peerIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		peerIndexes[i] = i
	}
	addPeersToChannel(t, n, channelName, gossips, peerIndexes)

	waitForFullMembershipOrFailNow(t, channelName, gossips, n, time.Second*30, time.Second*2)

	logger.Warning("Starting leader election services")

	
	services := make([]*electionService, n)

	electionMetrics := gossipmetrics.NewGossipMetrics(&disabled.Provider{}).ElectionMetrics

	for i := 0; i < n; i++ {
		services[i] = &electionService{nil, false, 0}
		services[i].LeaderElectionService = gossips[i].newLeaderElectionComponent(channelName, services[i].callback, electionMetrics)
	}

	logger.Warning("Waiting for leader election")

	assert.True(t, waitForLeaderElection(t, services, time.Second*30, time.Second*2), "One leader should be selected")

	startsNum := 0
	for i := 0; i < n; i++ {
		
		if services[i].callbackInvokeRes {
			startsNum++
		}
	}
	
	assert.Equal(t, 1, startsNum, "Only for one peer callback function should be called - chanA")

	
	
	secondChannelPeerIndexes := []int{1, 3, 5, 7}
	secondChannelName := "chanB"
	secondChannelServices := make([]*electionService, len(secondChannelPeerIndexes))
	addPeersToChannel(t, n, secondChannelName, gossips, secondChannelPeerIndexes)

	secondChannelGossips := make([]*gossipGRPC, 0)
	for _, i := range secondChannelPeerIndexes {
		secondChannelGossips = append(secondChannelGossips, gossips[i])
	}
	waitForFullMembershipOrFailNow(t, secondChannelName, secondChannelGossips, len(secondChannelGossips), time.Second*30, time.Millisecond*100)

	for idx, i := range secondChannelPeerIndexes {
		secondChannelServices[idx] = &electionService{nil, false, 0}
		secondChannelServices[idx].LeaderElectionService =
			gossips[i].newLeaderElectionComponent(secondChannelName, secondChannelServices[idx].callback, electionMetrics)
	}

	assert.True(t, waitForLeaderElection(t, secondChannelServices, time.Second*30, time.Second*2), "One leader should be selected for chanB")
	assert.True(t, waitForLeaderElection(t, services, time.Second*30, time.Second*2), "One leader should be selected for chanA")

	startsNum = 0
	for i := 0; i < n; i++ {
		if services[i].callbackInvokeRes {
			startsNum++
		}
	}
	assert.Equal(t, 1, startsNum, "Only for one peer callback function should be called - chanA")

	startsNum = 0
	for i := 0; i < len(secondChannelServices); i++ {
		if secondChannelServices[i].callbackInvokeRes {
			startsNum++
		}
	}
	assert.Equal(t, 1, startsNum, "Only for one peer callback function should be called - chanB")

	
	

	logger.Warning("Killing 2 peers, initiation new leader election")

	stopPeers(gossips[:2])

	waitForFullMembershipOrFailNow(t, channelName, gossips[2:], n-2, time.Second*30, time.Millisecond*100)
	waitForFullMembershipOrFailNow(t, secondChannelName, secondChannelGossips[1:], len(secondChannelGossips)-1, time.Second*30, time.Millisecond*100)

	assert.True(t, waitForLeaderElection(t, services[2:], time.Second*30, time.Second*2), "One leader should be selected after re-election - chanA")
	assert.True(t, waitForLeaderElection(t, secondChannelServices[1:], time.Second*30, time.Second*2), "One leader should be selected after re-election - chanB")

	startsNum = 0
	for i := 2; i < n; i++ {
		if services[i].callbackInvokeRes {
			startsNum++
		}
	}
	assert.Equal(t, 1, startsNum, "Only for one peer callback function should be called after re-election - chanA")

	startsNum = 0
	for i := 1; i < len(secondChannelServices); i++ {
		if secondChannelServices[i].callbackInvokeRes {
			startsNum++
		}
	}
	assert.Equal(t, 1, startsNum, "Only for one peer callback function should be called after re-election - chanB")

	stopServices(secondChannelServices)
	stopServices(services)
	stopPeers(gossips[2:])
}

type electionService struct {
	election.LeaderElectionService
	callbackInvokeRes   bool
	callbackInvokeCount int
}

func (es *electionService) callback(isLeader bool) {
	es.callbackInvokeRes = isLeader
	es.callbackInvokeCount = es.callbackInvokeCount + 1
}

type joinChanMsg struct {
}



func (jmc *joinChanMsg) SequenceNumber() uint64 {
	return uint64(time.Now().UnixNano())
}


func (jmc *joinChanMsg) Members() []api.OrgIdentityType {
	return []api.OrgIdentityType{orgInChannelA}
}


func (jmc *joinChanMsg) AnchorPeersOf(org api.OrgIdentityType) []api.AnchorPeer {
	return []api.AnchorPeer{}
}

func waitForFullMembershipOrFailNow(t *testing.T, channel string, gossips []*gossipGRPC, peersNum int, timeout time.Duration, testPollInterval time.Duration) {
	end := time.Now().Add(timeout)
	var correctPeers int
	for time.Now().Before(end) {
		correctPeers = 0
		for _, g := range gossips {
			if len(g.PeersOfChannel(gossipcommon.ChannelID(channel))) == (peersNum - 1) {
				correctPeers++
			}
		}
		if correctPeers == peersNum {
			return
		}
		time.Sleep(testPollInterval)
	}
	t.Fatalf("Failed to establish full channel membership. Only %d out of %d peers have full membership", correctPeers, peersNum)
}

func waitForMultipleLeadersElection(t *testing.T, services []*electionService, leadersNum int, timeout time.Duration, testPollInterval time.Duration) bool {
	logger.Warning("Waiting for", leadersNum, "leaders")
	end := time.Now().Add(timeout)
	correctNumberOfLeadersFound := false
	leaders := 0
	for time.Now().Before(end) {
		leaders = 0
		for _, s := range services {
			if s.IsLeader() {
				leaders++
			}
		}
		if leaders == leadersNum {
			if correctNumberOfLeadersFound {
				return true
			}
			correctNumberOfLeadersFound = true
		} else {
			correctNumberOfLeadersFound = false
		}
		time.Sleep(testPollInterval)
	}
	logger.Warning("Incorrect number of leaders", leaders)
	for i, s := range services {
		logger.Warning("Peer at index", i, "is leader", s.IsLeader())
	}
	return false
}

func waitForLeaderElection(t *testing.T, services []*electionService, timeout time.Duration, testPollInterval time.Duration) bool {
	return waitForMultipleLeadersElection(t, services, 1, timeout, testPollInterval)
}

func waitUntilOrFailBlocking(t *testing.T, f func(), timeout time.Duration) {
	successChan := make(chan struct{}, 1)
	go func() {
		f()
		successChan <- struct{}{}
	}()
	select {
	case <-time.NewTimer(timeout).C:
		break
	case <-successChan:
		return
	}
	util.PrintStackTrace()
	assert.Fail(t, "Timeout expired!")
}

func stopServices(services []*electionService) {
	stoppingWg := sync.WaitGroup{}
	stoppingWg.Add(len(services))
	for i, sI := range services {
		go func(i int, s_i election.LeaderElectionService) {
			defer stoppingWg.Done()
			s_i.Stop()
		}(i, sI)
	}
	stoppingWg.Wait()
	time.Sleep(time.Second * time.Duration(2))
}

func stopPeers(peers []*gossipGRPC) {
	stoppingWg := sync.WaitGroup{}
	stoppingWg.Add(len(peers))
	for i, pI := range peers {
		go func(i int, p_i *GossipService) {
			defer stoppingWg.Done()
			p_i.Stop()
		}(i, pI.GossipService)
	}
	stoppingWg.Wait()
	time.Sleep(time.Second * time.Duration(2))
}

func addPeersToChannel(t *testing.T, n int, channel string, peers []*gossipGRPC, peerIndexes []int) {
	jcm := &joinChanMsg{}

	wg := sync.WaitGroup{}
	for _, i := range peerIndexes {
		wg.Add(1)
		go func(i int) {
			peers[i].JoinChan(jcm, gossipcommon.ChannelID(channel))
			peers[i].UpdateLedgerHeight(0, gossipcommon.ChannelID(channel))
			wg.Done()
		}(i)
	}
	waitUntilOrFailBlocking(t, wg.Wait, time.Second*10)
}

func startPeers(t *testing.T, serviceConfig *ServiceConfig, n int, boot ...int) []*gossipGRPC {
	var ports []int
	var grpcs []*comm.GRPCServer
	var certs []*gossipcommon.TLSCertificates
	var secDialOpts []api.PeerSecureDialOpts

	for i := 0; i < n; i++ {
		port, grpc, cert, secDialOpt, _ := util.CreateGRPCLayer()
		ports = append(ports, port)
		grpcs = append(grpcs, grpc)
		certs = append(certs, cert)
		secDialOpts = append(secDialOpts, secDialOpt)
	}

	var bootPorts []int
	for _, index := range boot {
		bootPorts = append(bootPorts, ports[index])
	}

	peers := make([]*gossipGRPC, n)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			peers[i] = newGossipInstance(serviceConfig, ports[i], i, grpcs[i], certs[i], secDialOpts[i], 100, bootPorts...)
			wg.Done()
		}(i)
	}
	waitUntilOrFailBlocking(t, wg.Wait, time.Second*10)

	return peers
}

func newGossipInstance(serviceConfig *ServiceConfig, port int, id int, gRPCServer *comm.GRPCServer, certs *gossipcommon.TLSCertificates,
	secureDialOpts api.PeerSecureDialOpts, maxMsgCount int, bootPorts ...int) *gossipGRPC {
	conf := &gossip.Config{
		BindPort:                     port,
		BootstrapPeers:               bootPeers(bootPorts...),
		ID:                           fmt.Sprintf("p%d", id),
		MaxBlockCountToStore:         maxMsgCount,
		MaxPropagationBurstLatency:   time.Duration(500) * time.Millisecond,
		MaxPropagationBurstSize:      20,
		PropagateIterations:          1,
		PropagatePeerNum:             3,
		PullInterval:                 time.Duration(2) * time.Second,
		PullPeerNum:                  5,
		InternalEndpoint:             fmt.Sprintf("127.0.0.1:%d", port),
		ExternalEndpoint:             fmt.Sprintf("1.2.3.4:%d", port),
		PublishCertPeriod:            time.Duration(4) * time.Second,
		PublishStateInfoInterval:     time.Duration(1) * time.Second,
		RequestStateInfoInterval:     time.Duration(1) * time.Second,
		TimeForMembershipTracker:     time.Second * 5,
		TLSCerts:                     certs,
		DigestWaitTime:               algo.DefDigestWaitTime,
		RequestWaitTime:              algo.DefRequestWaitTime,
		ResponseWaitTime:             algo.DefResponseWaitTime,
		DialTimeout:                  gcomm.DefDialTimeout,
		ConnTimeout:                  gcomm.DefConnTimeout,
		RecvBuffSize:                 gcomm.DefRecvBuffSize,
		SendBuffSize:                 gcomm.DefSendBuffSize,
		MsgExpirationTimeout:         channel.DefMsgExpirationTimeout,
		AliveTimeInterval:            discovery.DefAliveTimeInterval,
		AliveExpirationTimeout:       discovery.DefAliveExpirationTimeout,
		AliveExpirationCheckInterval: discovery.DefAliveExpirationCheckInterval,
		ReconnectInterval:            time.Duration(1) * time.Second,
	}
	selfID := api.PeerIdentityType(conf.InternalEndpoint)
	cryptoService := &naiveCryptoService{}
	metrics := gossipmetrics.NewGossipMetrics(&disabled.Provider{})
	gossip := gossip.New(
		conf,
		gRPCServer.Server(),
		&orgCryptoService{},
		cryptoService,
		selfID,
		secureDialOpts,
		metrics,
	)
	go gRPCServer.Start()

	gossipService := &GossipService{
		mcs:             cryptoService,
		gossipSvc:       gossip,
		chains:          make(map[string]state.GossipStateProvider),
		leaderElection:  make(map[string]election.LeaderElectionService),
		privateHandlers: make(map[string]privateHandler),
		deliveryService: make(map[string]deliverservice.DeliverService),
		deliveryFactory: &deliveryFactoryImpl{
			credentialSupport: comm.NewCredentialSupport(),
		},
		peerIdentity:  api.PeerIdentityType(conf.InternalEndpoint),
		metrics:       metrics,
		serviceConfig: serviceConfig,
	}

	return &gossipGRPC{GossipService: gossipService, grpc: gRPCServer}
}

type gossipGRPC struct {
	*GossipService
	grpc *comm.GRPCServer
}

func (g *gossipGRPC) Stop() {
	g.GossipService.Stop()
	g.grpc.Stop()
}

func bootPeers(ports ...int) []string {
	var peers []string
	for _, port := range ports {
		peers = append(peers, fmt.Sprintf("127.0.0.1:%d", port))
	}
	return peers
}

func getAvailablePort(t *testing.T) (endpoint string, ll net.Listener) {
	ll, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	endpoint = ll.Addr().String()
	return endpoint, ll
}

type naiveCryptoService struct {
}

type orgCryptoService struct {
}



func (*orgCryptoService) OrgByPeerIdentity(identity api.PeerIdentityType) api.OrgIdentityType {
	return orgInChannelA
}



func (*orgCryptoService) Verify(joinChanMsg api.JoinChannelMessage) error {
	return nil
}

func (naiveCryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	return time.Now().Add(time.Hour), nil
}



func (*naiveCryptoService) VerifyByChannel(_ gossipcommon.ChannelID, _ api.PeerIdentityType, _, _ []byte) error {
	return nil
}

func (*naiveCryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	return nil
}


func (*naiveCryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) gossipcommon.PKIidType {
	return gossipcommon.PKIidType(peerIdentity)
}



func (*naiveCryptoService) VerifyBlock(chainID gossipcommon.ChannelID, seqNum uint64, signedBlock []byte) error {
	return nil
}



func (*naiveCryptoService) Sign(msg []byte) ([]byte, error) {
	return msg, nil
}




func (*naiveCryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	equal := bytes.Equal(signature, message)
	if !equal {
		return fmt.Errorf("Wrong signature:%v, %v", signature, message)
	}
	return nil
}

var orgInChannelA = api.OrgIdentityType("ORG1")

func TestInvalidInitialization(t *testing.T) {
	grpcServer := grpc.NewServer()
	endpoint, socket := getAvailablePort(t)

	mockSignerSerializer := &mocks.SignerSerializer{}
	mockSignerSerializer.SerializeReturns(api.PeerIdentityType("peer-identity"), nil)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	dialOpts := defaultDeliverClientDialOpts()
	gossipConfig, err := gossip.GlobalConfig(endpoint, nil)
	assert.NoError(t, err)

	gossipService, err := New(
		mockSignerSerializer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		endpoint,
		grpcServer,
		&naiveCryptoService{},
		secAdv,
		nil,
		comm.NewCredentialSupport(),
		dialOpts,
		gossipConfig,
		&ServiceConfig{},
		&deliverservice.DeliverServiceConfig{
			PeerTLSEnabled:              false,
			ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
			ReconnectTotalTimeThreshold: deliverservice.DefaultReConnectTotalTimeThreshold,
		},
	)
	assert.NoError(t, err)
	gService := gossipService
	defer gService.Stop()

	go grpcServer.Serve(socket)
	defer grpcServer.Stop()

	dc, err := gService.deliveryFactory.Service(gService, []string{}, &naiveCryptoService{})
	assert.EqualError(t, err, "no endpoints specified")
	assert.Nil(t, dc)

	endpoint2, socket2 := getAvailablePort(t)
	defer socket2.Close()

	dc, err = gService.deliveryFactory.Service(gService, []string{endpoint2}, &naiveCryptoService{})
	assert.NotNil(t, dc)
	assert.NoError(t, err)
}

func TestChannelConfig(t *testing.T) {
	
	grpcServer := grpc.NewServer()
	endpoint, socket := getAvailablePort(t)

	mockSignerSerializer := &mocks.SignerSerializer{}
	mockSignerSerializer.SerializeReturns(api.PeerIdentityType("peer-identity"), nil)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	dialOpts := defaultDeliverClientDialOpts()
	gossipConfig, err := gossip.GlobalConfig(endpoint, nil)
	assert.NoError(t, err)

	gossipService, err := New(
		mockSignerSerializer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		endpoint,
		grpcServer,
		&naiveCryptoService{},
		secAdv,
		nil,
		nil,
		dialOpts,
		gossipConfig,
		&ServiceConfig{},
		&deliverservice.DeliverServiceConfig{
			ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
			ReconnectTotalTimeThreshold: deliverservice.DefaultReConnectTotalTimeThreshold,
		},
	)
	assert.NoError(t, err)
	gService := gossipService
	defer gService.Stop()

	go grpcServer.Serve(socket)
	defer grpcServer.Stop()

	jcm := &joinChannelMessage{seqNum: 1, members2AnchorPeers: map[string][]api.AnchorPeer{
		"A": {{Host: "host", Port: 5000}},
	}}

	assert.Equal(t, uint64(1), jcm.SequenceNumber())

	mc := &mockConfig{
		sequence: 1,
		orgs: map[string]channelconfig.ApplicationOrg{
			string(orgInChannelA): &appGrp{
				mspID:       string(orgInChannelA),
				anchorPeers: []*peer.AnchorPeer{},
			},
		},
	}
	gService.JoinChan(jcm, gossipcommon.ChannelID("A"))
	gService.updateAnchors(mc)
	assert.True(t, gService.amIinChannel(string(orgInChannelA), mc))
}

func defaultDeliverClientDialOpts() []grpc.DialOption {
	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	dialOpts = append(
		dialOpts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)))
	kaOpts := comm.DefaultKeepaliveOptions
	dialOpts = append(dialOpts, comm.ClientKeepaliveOptions(kaOpts)...)

	return dialOpts
}
