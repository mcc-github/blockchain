/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"sync"

	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/gossip/api"
	gossipCommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/gossip/integration"
	gossipMetrics "github.com/mcc-github/blockchain/gossip/metrics"
	privdata2 "github.com/mcc-github/blockchain/gossip/privdata"
	"github.com/mcc-github/blockchain/gossip/state"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	gproto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	gossipServiceInstance *gossipServiceImpl
	once                  sync.Once
)

type gossipSvc gossip.Gossip


type GossipService interface {
	gossip.Gossip

	
	
	DistributePrivateData(chainID string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error
	
	NewConfigEventer() ConfigProcessor
	
	InitializeChannel(chainID string, endpoints []string, support Support)
	
	AddPayload(chainID string, payload *gproto.Payload) error
}


type DeliveryServiceFactory interface {
	
	Service(g GossipService, endpoints []string, msc api.MessageCryptoService) (deliverservice.DeliverService, error)
}

type deliveryFactoryImpl struct {
	signer identity.SignerSerializer
}


func (df *deliveryFactoryImpl) Service(g GossipService, endpoints []string, mcs api.MessageCryptoService) (deliverservice.DeliverService, error) {
	return deliverservice.NewDeliverService(&deliverservice.Config{
		CryptoSvc:   mcs,
		Gossip:      g,
		Endpoints:   endpoints,
		ConnFactory: deliverservice.DefaultConnectionFactory,
		ABCFactory:  deliverservice.DefaultABCFactory,
		Signer:      df.signer,
	})
}

type privateHandler struct {
	support     Support
	coordinator privdata2.Coordinator
	distributor privdata2.PvtDataDistributor
	reconciler  privdata2.PvtDataReconciler
}

func (p privateHandler) close() {
	p.coordinator.Close()
	p.reconciler.Stop()
}

type gossipServiceImpl struct {
	gossipSvc
	privateHandlers map[string]privateHandler
	chains          map[string]state.GossipStateProvider
	leaderElection  map[string]election.LeaderElectionService
	deliveryService map[string]deliverservice.DeliverService
	deliveryFactory DeliveryServiceFactory
	lock            sync.RWMutex
	mcs             api.MessageCryptoService
	peerIdentity    []byte
	secAdv          api.SecurityAdvisor
	metrics         *gossipMetrics.GossipMetrics
}


type joinChannelMessage struct {
	seqNum              uint64
	members2AnchorPeers map[string][]api.AnchorPeer
}

func (jcm *joinChannelMessage) SequenceNumber() uint64 {
	return jcm.seqNum
}


func (jcm *joinChannelMessage) Members() []api.OrgIdentityType {
	members := make([]api.OrgIdentityType, len(jcm.members2AnchorPeers))
	i := 0
	for org := range jcm.members2AnchorPeers {
		members[i] = api.OrgIdentityType(org)
		i++
	}
	return members
}


func (jcm *joinChannelMessage) AnchorPeersOf(org api.OrgIdentityType) []api.AnchorPeer {
	return jcm.members2AnchorPeers[string(org)]
}

var logger = util.GetLogger(util.ServiceLogger, "")


func InitGossipService(peerIdentity identity.SignerSerializer, metricsProvider metrics.Provider, endpoint string, s *grpc.Server,
	certs *gossipCommon.TLSCertificates, mcs api.MessageCryptoService, secAdv api.SecurityAdvisor,
	secureDialOpts api.PeerSecureDialOpts, bootPeers ...string) error {
	
	
	
	util.GetLogger(util.ElectionLogger, "")

	return InitGossipServiceCustomDeliveryFactory(
		peerIdentity,
		metricsProvider,
		endpoint,
		s,
		certs,
		&deliveryFactoryImpl{
			signer: peerIdentity,
		},
		mcs,
		secAdv,
		secureDialOpts,
		bootPeers...,
	)
}



func InitGossipServiceCustomDeliveryFactory(
	peerIdentity identity.SignerSerializer,
	metricsProvider metrics.Provider,
	endpoint string,
	s *grpc.Server,
	certs *gossipCommon.TLSCertificates,
	factory DeliveryServiceFactory,
	mcs api.MessageCryptoService,
	secAdv api.SecurityAdvisor,
	secureDialOpts api.PeerSecureDialOpts,
	bootPeers ...string,
) error {
	var err error
	var gossip gossip.Gossip
	serializedIdentity, err := peerIdentity.Serialize()
	if err != nil {
		return err
	}
	once.Do(func() {
		if overrideEndpoint := viper.GetString("peer.gossip.endpoint"); overrideEndpoint != "" {
			endpoint = overrideEndpoint
		}

		logger.Info("Initialize gossip with endpoint", endpoint, "and bootstrap set", bootPeers)

		gossipMetrics := gossipMetrics.NewGossipMetrics(metricsProvider)

		gossip, err = integration.NewGossipComponent(serializedIdentity, endpoint, s, secAdv,
			mcs, secureDialOpts, certs, gossipMetrics, bootPeers...)
		gossipServiceInstance = &gossipServiceImpl{
			mcs:             mcs,
			gossipSvc:       gossip,
			privateHandlers: make(map[string]privateHandler),
			chains:          make(map[string]state.GossipStateProvider),
			leaderElection:  make(map[string]election.LeaderElectionService),
			deliveryService: make(map[string]deliverservice.DeliverService),
			deliveryFactory: factory,
			peerIdentity:    serializedIdentity,
			secAdv:          secAdv,
			metrics:         gossipMetrics,
		}
	})
	return errors.WithStack(err)
}


func GetGossipService() GossipService {
	return gossipServiceInstance
}


func (g *gossipServiceImpl) DistributePrivateData(chainID string, txID string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error {
	g.lock.RLock()
	handler, exists := g.privateHandlers[chainID]
	g.lock.RUnlock()
	if !exists {
		return errors.Errorf("No private data handler for %s", chainID)
	}

	if err := handler.distributor.Distribute(txID, privData, blkHt); err != nil {
		logger.Error("Failed to distributed private collection, txID", txID, "channel", chainID, "due to", err)
		return err
	}

	if err := handler.coordinator.StorePvtData(txID, privData, blkHt); err != nil {
		logger.Error("Failed to store private data into transient store, txID",
			txID, "channel", chainID, "due to", err)
		return err
	}
	return nil
}


func (g *gossipServiceImpl) NewConfigEventer() ConfigProcessor {
	return newConfigEventer(g)
}



type Support struct {
	Validator            txvalidator.Validator
	Committer            committer.Committer
	Store                privdata2.TransientStore
	Cs                   privdata.CollectionStore
	IdDeserializeFactory privdata2.IdentityDeserializerFactory
}



type DataStoreSupport struct {
	committer.Committer
	privdata2.TransientStore
}


func (g *gossipServiceImpl) InitializeChannel(chainID string, endpoints []string, support Support) {
	g.lock.Lock()
	defer g.lock.Unlock()
	
	logger.Debug("Creating state provider for chainID", chainID)
	servicesAdapter := &state.ServicesMediator{GossipAdapter: g, MCSAdapter: g.mcs}

	
	
	
	storeSupport := &DataStoreSupport{
		TransientStore: support.Store,
		Committer:      support.Committer,
	}
	
	dataRetriever := privdata2.NewDataRetriever(storeSupport)
	collectionAccessFactory := privdata2.NewCollectionAccessFactory(support.IdDeserializeFactory)
	fetcher := privdata2.NewPuller(g.metrics.PrivdataMetrics, support.Cs, g.gossipSvc, dataRetriever,
		collectionAccessFactory, chainID, privdata2.GetBtlPullMargin())

	coordinatorConfig := privdata2.CoordinatorConfig{
		TransientBlockRetention: privdata2.GetTransientBlockRetention(),
		PullRetryThreshold:      viper.GetDuration("peer.gossip.pvtData.pullRetryThreshold"),
	}
	coordinator := privdata2.NewCoordinator(privdata2.Support{
		ChainID:         chainID,
		CollectionStore: support.Cs,
		Validator:       support.Validator,
		TransientStore:  support.Store,
		Committer:       support.Committer,
		Fetcher:         fetcher,
	}, g.createSelfSignedData(), g.metrics.PrivdataMetrics, coordinatorConfig)

	reconcilerConfig := privdata2.GetReconcilerConfig()
	var reconciler privdata2.PvtDataReconciler

	if reconcilerConfig.IsEnabled {
		reconciler = privdata2.NewReconciler(chainID, g.metrics.PrivdataMetrics,
			support.Committer, fetcher, reconcilerConfig)
	} else {
		reconciler = &privdata2.NoOpReconciler{}
	}

	pushAckTimeout := viper.GetDuration("peer.gossip.pvtData.pushAckTimeout")
	g.privateHandlers[chainID] = privateHandler{
		support:     support,
		coordinator: coordinator,
		distributor: privdata2.NewDistributor(chainID, g, collectionAccessFactory, g.metrics.PrivdataMetrics, pushAckTimeout),
		reconciler:  reconciler,
	}
	g.privateHandlers[chainID].reconciler.Start()

	blockingMode := !viper.GetBool("peer.gossip.nonBlockingCommitMode")
	g.chains[chainID] = state.NewGossipStateProvider(chainID, servicesAdapter, coordinator,
		g.metrics.StateMetrics, blockingMode)
	if g.deliveryService[chainID] == nil {
		var err error
		g.deliveryService[chainID], err = g.deliveryFactory.Service(g, endpoints, g.mcs)
		if err != nil {
			logger.Warningf("Cannot create delivery client, due to %+v", errors.WithStack(err))
		}
	}

	
	
	if g.deliveryService[chainID] != nil {
		
		
		
		
		
		
		leaderElection := viper.GetBool("peer.gossip.useLeaderElection")
		isStaticOrgLeader := viper.GetBool("peer.gossip.orgLeader")

		if leaderElection && isStaticOrgLeader {
			logger.Panic("Setting both orgLeader and useLeaderElection to true isn't supported, aborting execution")
		}

		if leaderElection {
			logger.Debug("Delivery uses dynamic leader election mechanism, channel", chainID)
			g.leaderElection[chainID] = g.newLeaderElectionComponent(chainID, g.onStatusChangeFactory(chainID,
				support.Committer), g.metrics.ElectionMetrics)
		} else if isStaticOrgLeader {
			logger.Debug("This peer is configured to connect to ordering service for blocks delivery, channel", chainID)
			g.deliveryService[chainID].StartDeliverForChannel(chainID, support.Committer, func() {})
		} else {
			logger.Debug("This peer is not configured to connect to ordering service for blocks delivery, channel", chainID)
		}
	} else {
		logger.Warning("Delivery client is down won't be able to pull blocks for chain", chainID)
	}
}

func (g *gossipServiceImpl) createSelfSignedData() protoutil.SignedData {
	msg := make([]byte, 32)
	sig, err := g.mcs.Sign(msg)
	if err != nil {
		logger.Panicf("Failed creating self signed data because message signing failed: %v", err)
	}
	return protoutil.SignedData{
		Data:      msg,
		Signature: sig,
		Identity:  g.peerIdentity,
	}
}


func (g *gossipServiceImpl) updateAnchors(config Config) {
	myOrg := string(g.secAdv.OrgByPeerIdentity(api.PeerIdentityType(g.peerIdentity)))
	if !g.amIinChannel(myOrg, config) {
		logger.Error("Tried joining channel", config.ChainID(), "but our org(", myOrg, "), isn't "+
			"among the orgs of the channel:", orgListFromConfig(config), ", aborting.")
		return
	}
	jcm := &joinChannelMessage{seqNum: config.Sequence(), members2AnchorPeers: map[string][]api.AnchorPeer{}}
	for _, appOrg := range config.Organizations() {
		logger.Debug(appOrg.MSPID(), "anchor peers:", appOrg.AnchorPeers())
		jcm.members2AnchorPeers[appOrg.MSPID()] = []api.AnchorPeer{}
		for _, ap := range appOrg.AnchorPeers() {
			anchorPeer := api.AnchorPeer{
				Host: ap.Host,
				Port: int(ap.Port),
			}
			jcm.members2AnchorPeers[appOrg.MSPID()] = append(jcm.members2AnchorPeers[appOrg.MSPID()], anchorPeer)
		}
	}

	
	logger.Debug("Creating state provider for chainID", config.ChainID())
	g.JoinChan(jcm, gossipCommon.ChainID(config.ChainID()))
}

func (g *gossipServiceImpl) updateEndpoints(chainID string, endpoints []string) {
	if ds, ok := g.deliveryService[chainID]; ok {
		logger.Debugf("Updating endpoints for chainID %s", chainID)
		if err := ds.UpdateEndpoints(chainID, endpoints); err != nil {
			
			
			logger.Warningf("Failed to update ordering service endpoints, due to %s", err)
		}
	}
}


func (g *gossipServiceImpl) AddPayload(chainID string, payload *gproto.Payload) error {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.chains[chainID].AddPayload(payload)
}


func (g *gossipServiceImpl) Stop() {
	g.lock.Lock()
	defer g.lock.Unlock()

	for chainID := range g.chains {
		logger.Info("Stopping chain", chainID)
		if le, exists := g.leaderElection[chainID]; exists {
			logger.Infof("Stopping leader election for %s", chainID)
			le.Stop()
		}
		g.chains[chainID].Stop()
		g.privateHandlers[chainID].close()

		if g.deliveryService[chainID] != nil {
			g.deliveryService[chainID].Stop()
		}
	}
	g.gossipSvc.Stop()
}

func (g *gossipServiceImpl) newLeaderElectionComponent(chainID string, callback func(bool),
	electionMetrics *gossipMetrics.ElectionMetrics) election.LeaderElectionService {
	PKIid := g.mcs.GetPKIidOfCert(g.peerIdentity)
	adapter := election.NewAdapter(g, PKIid, gossipCommon.ChainID(chainID), electionMetrics)
	config := election.ElectionConfig{
		StartupGracePeriod:       util.GetDurationOrDefault("peer.gossip.election.startupGracePeriod", election.DefStartupGracePeriod),
		MembershipSampleInterval: util.GetDurationOrDefault("peer.gossip.election.membershipSampleInterval", election.DefMembershipSampleInterval),
		LeaderAliveThreshold:     util.GetDurationOrDefault("peer.gossip.election.leaderAliveThreshold", election.DefLeaderAliveThreshold),
		LeaderElectionDuration:   util.GetDurationOrDefault("peer.gossip.election.leaderElectionDuration", election.DefLeaderElectionDuration),
	}
	return election.NewLeaderElectionService(adapter, string(PKIid), callback, config)
}

func (g *gossipServiceImpl) amIinChannel(myOrg string, config Config) bool {
	for _, orgName := range orgListFromConfig(config) {
		if orgName == myOrg {
			return true
		}
	}
	return false
}

func (g *gossipServiceImpl) onStatusChangeFactory(chainID string, committer blocksprovider.LedgerInfo) func(bool) {
	return func(isLeader bool) {
		if isLeader {
			yield := func() {
				g.lock.RLock()
				le := g.leaderElection[chainID]
				g.lock.RUnlock()
				le.Yield()
			}
			logger.Info("Elected as a leader, starting delivery service for channel", chainID)
			if err := g.deliveryService[chainID].StartDeliverForChannel(chainID, committer, yield); err != nil {
				logger.Errorf("Delivery service is not able to start blocks delivery for chain, due to %+v", err)
			}
		} else {
			logger.Info("Renounced leadership, stopping delivery service for channel", chainID)
			if err := g.deliveryService[chainID].StopDeliverForChannel(chainID); err != nil {
				logger.Errorf("Delivery service is not able to stop blocks delivery for chain, due to %+v", err)
			}

		}

	}
}

func orgListFromConfig(config Config) []string {
	var orgList []string
	for _, appOrg := range config.Organizations() {
		orgList = append(orgList, appOrg.MSPID())
	}
	return orgList
}
