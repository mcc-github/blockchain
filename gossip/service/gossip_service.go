/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"sync"

	corecomm "github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/gossip"
	gossipmetrics "github.com/mcc-github/blockchain/gossip/metrics"
	gossipprivdata "github.com/mcc-github/blockchain/gossip/privdata"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/state"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	gproto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)


type gossipSvc interface {
	
	SelfMembershipInfo() discovery.NetworkMember

	
	SelfChannelInfo(common.ChannelID) *protoext.SignedGossipMessage

	
	Send(msg *gproto.GossipMessage, peers ...*comm.RemotePeer)

	
	SendByCriteria(*protoext.SignedGossipMessage, gossip.SendCriteria) error

	
	Peers() []discovery.NetworkMember

	
	
	PeersOfChannel(common.ChannelID) []discovery.NetworkMember

	
	
	UpdateMetadata(metadata []byte)

	
	
	UpdateLedgerHeight(height uint64, channelID common.ChannelID)

	
	
	UpdateChaincodes(chaincode []*gproto.Chaincode, channelID common.ChannelID)

	
	Gossip(msg *gproto.GossipMessage)

	
	
	PeerFilter(channel common.ChannelID, messagePredicate api.SubChannelSelectionCriteria) (filter.RoutingFilter, error)

	
	
	
	
	Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *gproto.GossipMessage, <-chan protoext.ReceivedMessage)

	
	JoinChan(joinMsg api.JoinChannelMessage, channelID common.ChannelID)

	
	
	
	
	LeaveChan(channelID common.ChannelID)

	
	
	SuspectPeers(s api.PeerSuspector)

	
	IdentityInfo() api.PeerIdentitySet

	
	IsInMyOrg(member discovery.NetworkMember) bool

	
	Stop()
}



type GossipServiceAdapter interface {
	
	PeersOfChannel(gossipcommon.ChannelID) []discovery.NetworkMember

	
	AddPayload(channelID string, payload *gproto.Payload) error

	
	Gossip(msg *gproto.GossipMessage)
}


type DeliveryServiceFactory interface {
	
	Service(g GossipServiceAdapter, endpoints []string, msc api.MessageCryptoService) (deliverservice.DeliverService, error)
}

type deliveryFactoryImpl struct {
	signer                identity.SignerSerializer
	credentialSupport     *corecomm.CredentialSupport
	deliverClientDialOpts []grpc.DialOption
	deliverServiceConfig  *deliverservice.DeliverServiceConfig
}


func (df *deliveryFactoryImpl) Service(g GossipServiceAdapter, endpoints []string, mcs api.MessageCryptoService) (deliverservice.DeliverService, error) {
	return deliverservice.NewDeliverService(&deliverservice.Config{
		ABCFactory:        deliverservice.DefaultABCFactory,
		CryptoSvc:         mcs,
		CredentialSupport: df.credentialSupport,
		Endpoints:         endpoints,
		ConnFactory: (&deliverservice.CredSupportDialerFactory{
			CredentialSupport: df.credentialSupport,
			KeepaliveOptions:  df.deliverServiceConfig.KeepaliveOptions,
			TLSEnabled:        df.deliverServiceConfig.PeerTLSEnabled,
		}).Dialer,
		Gossip:                g,
		Signer:                df.signer,
		DeliverClientDialOpts: df.deliverClientDialOpts,
		DeliverServiceConfig:  df.deliverServiceConfig,
	})
}

type privateHandler struct {
	support     Support
	coordinator gossipprivdata.Coordinator
	distributor gossipprivdata.PvtDataDistributor
	reconciler  gossipprivdata.PvtDataReconciler
}

func (p privateHandler) close() {
	p.coordinator.Close()
	p.reconciler.Stop()
}

type GossipService struct {
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
	metrics         *gossipmetrics.GossipMetrics
	serviceConfig   *ServiceConfig
}


type joinChannelMessage struct {
	seqNum              uint64
	members2AnchorPeers map[string][]api.AnchorPeer
}

func (jcm *joinChannelMessage) SequenceNumber() uint64 {
	return jcm.seqNum
}


func (jcm *joinChannelMessage) Members() []api.OrgIdentityType {
	members := make([]api.OrgIdentityType, 0, len(jcm.members2AnchorPeers))
	for org := range jcm.members2AnchorPeers {
		members = append(members, api.OrgIdentityType(org))
	}
	return members
}


func (jcm *joinChannelMessage) AnchorPeersOf(org api.OrgIdentityType) []api.AnchorPeer {
	return jcm.members2AnchorPeers[string(org)]
}

var logger = util.GetLogger(util.ServiceLogger, "")


func New(
	peerIdentity identity.SignerSerializer,
	gossipMetrics *gossipmetrics.GossipMetrics,
	endpoint string,
	s *grpc.Server,
	mcs api.MessageCryptoService,
	secAdv api.SecurityAdvisor,
	secureDialOpts api.PeerSecureDialOpts,
	credSupport *corecomm.CredentialSupport,
	deliverClientDialOpts []grpc.DialOption,
	gossipConfig *gossip.Config,
	serviceConfig *ServiceConfig,
	deliverServiceConfig *deliverservice.DeliverServiceConfig,
) (*GossipService, error) {
	serializedIdentity, err := peerIdentity.Serialize()
	if err != nil {
		return nil, err
	}

	logger.Infof("Initialize gossip with endpoint %s", endpoint)

	gossipComponent := gossip.New(
		gossipConfig,
		s,
		secAdv,
		mcs,
		serializedIdentity,
		secureDialOpts,
		gossipMetrics,
	)

	return &GossipService{
		gossipSvc:       gossipComponent,
		mcs:             mcs,
		privateHandlers: make(map[string]privateHandler),
		chains:          make(map[string]state.GossipStateProvider),
		leaderElection:  make(map[string]election.LeaderElectionService),
		deliveryService: make(map[string]deliverservice.DeliverService),
		deliveryFactory: &deliveryFactoryImpl{
			signer:                peerIdentity,
			credentialSupport:     credSupport,
			deliverClientDialOpts: deliverClientDialOpts,
			deliverServiceConfig:  deliverServiceConfig,
		},
		peerIdentity:  serializedIdentity,
		secAdv:        secAdv,
		metrics:       gossipMetrics,
		serviceConfig: serviceConfig,
	}, nil
}


func (g *GossipService) DistributePrivateData(channelID string, txID string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error {
	g.lock.RLock()
	handler, exists := g.privateHandlers[channelID]
	g.lock.RUnlock()
	if !exists {
		return errors.Errorf("No private data handler for %s", channelID)
	}

	if err := handler.distributor.Distribute(txID, privData, blkHt); err != nil {
		logger.Error("Failed to distributed private collection, txID", txID, "channel", channelID, "due to", err)
		return err
	}

	if err := handler.coordinator.StorePvtData(txID, privData, blkHt); err != nil {
		logger.Error("Failed to store private data into transient store, txID",
			txID, "channel", channelID, "due to", err)
		return err
	}
	return nil
}


func (g *GossipService) NewConfigEventer() ConfigProcessor {
	return newConfigEventer(g)
}



type Support struct {
	Validator            txvalidator.Validator
	Committer            committer.Committer
	Store                gossipprivdata.TransientStore
	CollectionStore      privdata.CollectionStore
	IdDeserializeFactory gossipprivdata.IdentityDeserializerFactory
	CapabilityProvider   gossipprivdata.CapabilityProvider
}



type DataStoreSupport struct {
	committer.Committer
	gossipprivdata.TransientStore
}


func (g *GossipService) InitializeChannel(channelID string, endpoints []string, support Support) {
	g.lock.Lock()
	defer g.lock.Unlock()
	
	logger.Debug("Creating state provider for channelID", channelID)
	servicesAdapter := &state.ServicesMediator{GossipAdapter: g, MCSAdapter: g.mcs}

	
	
	
	storeSupport := &DataStoreSupport{
		TransientStore: support.Store,
		Committer:      support.Committer,
	}
	
	dataRetriever := gossipprivdata.NewDataRetriever(storeSupport)
	collectionAccessFactory := gossipprivdata.NewCollectionAccessFactory(support.IdDeserializeFactory)
	fetcher := gossipprivdata.NewPuller(g.metrics.PrivdataMetrics, support.CollectionStore, g.gossipSvc, dataRetriever,
		collectionAccessFactory, channelID, g.serviceConfig.BtlPullMargin)

	coordinatorConfig := gossipprivdata.CoordinatorConfig{
		TransientBlockRetention:        g.serviceConfig.TransientstoreMaxBlockRetention,
		PullRetryThreshold:             g.serviceConfig.PvtDataPullRetryThreshold,
		SkipPullingInvalidTransactions: g.serviceConfig.SkipPullingInvalidTransactionsDuringCommit,
	}
	coordinator := gossipprivdata.NewCoordinator(gossipprivdata.Support{
		ChainID:            channelID,
		CollectionStore:    support.CollectionStore,
		Validator:          support.Validator,
		TransientStore:     support.Store,
		Committer:          support.Committer,
		Fetcher:            fetcher,
		CapabilityProvider: support.CapabilityProvider,
	}, g.createSelfSignedData(), g.metrics.PrivdataMetrics, coordinatorConfig)

	privdataConfig := gossipprivdata.GlobalConfig()
	var reconciler gossipprivdata.PvtDataReconciler

	if privdataConfig.ReconciliationEnabled {
		reconciler = gossipprivdata.NewReconciler(channelID, g.metrics.PrivdataMetrics,
			support.Committer, fetcher, privdataConfig)
	} else {
		reconciler = &gossipprivdata.NoOpReconciler{}
	}

	pushAckTimeout := g.serviceConfig.PvtDataPushAckTimeout
	g.privateHandlers[channelID] = privateHandler{
		support:     support,
		coordinator: coordinator,
		distributor: gossipprivdata.NewDistributor(channelID, g, collectionAccessFactory, g.metrics.PrivdataMetrics, pushAckTimeout),
		reconciler:  reconciler,
	}
	g.privateHandlers[channelID].reconciler.Start()

	blockingMode := !g.serviceConfig.NonBlockingCommitMode
	stateConfig := state.GlobalConfig()
	g.chains[channelID] = state.NewGossipStateProvider(
		channelID,
		servicesAdapter,
		coordinator,
		g.metrics.StateMetrics,
		blockingMode,
		stateConfig)
	if g.deliveryService[channelID] == nil {
		var err error
		g.deliveryService[channelID], err = g.deliveryFactory.Service(g, endpoints, g.mcs)
		if err != nil {
			logger.Warningf("Cannot create delivery client, due to %+v", errors.WithStack(err))
		}
	}

	
	
	if g.deliveryService[channelID] != nil {
		
		
		
		
		
		
		leaderElection := g.serviceConfig.UseLeaderElection
		isStaticOrgLeader := g.serviceConfig.OrgLeader

		if leaderElection && isStaticOrgLeader {
			logger.Panic("Setting both orgLeader and useLeaderElection to true isn't supported, aborting execution")
		}

		if leaderElection {
			logger.Debug("Delivery uses dynamic leader election mechanism, channel", channelID)
			g.leaderElection[channelID] = g.newLeaderElectionComponent(channelID, g.onStatusChangeFactory(channelID,
				support.Committer), g.metrics.ElectionMetrics)
		} else if isStaticOrgLeader {
			logger.Debug("This peer is configured to connect to ordering service for blocks delivery, channel", channelID)
			g.deliveryService[channelID].StartDeliverForChannel(channelID, support.Committer, func() {})
		} else {
			logger.Debug("This peer is not configured to connect to ordering service for blocks delivery, channel", channelID)
		}
	} else {
		logger.Warning("Delivery client is down won't be able to pull blocks for chain", channelID)
	}

}

func (g *GossipService) createSelfSignedData() protoutil.SignedData {
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


func (g *GossipService) updateAnchors(config Config) {
	myOrg := string(g.secAdv.OrgByPeerIdentity(api.PeerIdentityType(g.peerIdentity)))
	if !g.amIinChannel(myOrg, config) {
		logger.Error("Tried joining channel", config.ChannelID(), "but our org(", myOrg, "), isn't "+
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

	
	logger.Debug("Creating state provider for channelID", config.ChannelID())
	g.JoinChan(jcm, gossipcommon.ChannelID(config.ChannelID()))
}

func (g *GossipService) updateEndpoints(channelID string, endpoints []string) {
	if ds, ok := g.deliveryService[channelID]; ok {
		logger.Debugf("Updating endpoints for channelID %s", channelID)
		if err := ds.UpdateEndpoints(channelID, endpoints); err != nil {
			
			
			logger.Warningf("Failed to update ordering service endpoints, due to %s", err)
		}
	}
}


func (g *GossipService) AddPayload(channelID string, payload *gproto.Payload) error {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.chains[channelID].AddPayload(payload)
}


func (g *GossipService) Stop() {
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

func (g *GossipService) newLeaderElectionComponent(channelID string, callback func(bool),
	electionMetrics *gossipmetrics.ElectionMetrics) election.LeaderElectionService {
	PKIid := g.mcs.GetPKIidOfCert(g.peerIdentity)
	adapter := election.NewAdapter(g, PKIid, gossipcommon.ChannelID(channelID), electionMetrics)
	config := election.ElectionConfig{
		StartupGracePeriod:       g.serviceConfig.ElectionStartupGracePeriod,
		MembershipSampleInterval: g.serviceConfig.ElectionMembershipSampleInterval,
		LeaderAliveThreshold:     g.serviceConfig.ElectionLeaderAliveThreshold,
		LeaderElectionDuration:   g.serviceConfig.ElectionLeaderElectionDuration,
	}
	return election.NewLeaderElectionService(adapter, string(PKIid), callback, config)
}

func (g *GossipService) amIinChannel(myOrg string, config Config) bool {
	for _, orgName := range orgListFromConfig(config) {
		if orgName == myOrg {
			return true
		}
	}
	return false
}

func (g *GossipService) onStatusChangeFactory(channelID string, committer blocksprovider.LedgerInfo) func(bool) {
	return func(isLeader bool) {
		if isLeader {
			yield := func() {
				g.lock.RLock()
				le := g.leaderElection[channelID]
				g.lock.RUnlock()
				le.Yield()
			}
			logger.Info("Elected as a leader, starting delivery service for channel", channelID)
			if err := g.deliveryService[channelID].StartDeliverForChannel(channelID, committer, yield); err != nil {
				logger.Errorf("Delivery service is not able to start blocks delivery for chain, due to %+v", err)
			}
		} else {
			logger.Info("Renounced leadership, stopping delivery service for channel", channelID)
			if err := g.deliveryService[channelID].StopDeliverForChannel(channelID); err != nil {
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
