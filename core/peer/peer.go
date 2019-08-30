/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/common/channelconfig"
	cc "github.com/mcc-github/blockchain/common/config"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/semaphore"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	validatorv14 "github.com/mcc-github/blockchain/core/committer/txvalidator/v14"
	validatorv20 "github.com/mcc-github/blockchain/core/committer/txvalidator/v20"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	vir "github.com/mcc-github/blockchain/core/committer/txvalidator/v20/valinforetriever"
	"github.com/mcc-github/blockchain/core/common/privdata"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/api"
	gossipprivdata "github.com/mcc-github/blockchain/gossip/privdata"
	gossipservice "github.com/mcc-github/blockchain/gossip/service"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var peerLogger = flogging.MustGetLogger("peer")

type CollectionInfoShim struct {
	plugindispatcher.CollectionAndLifecycleResources
	ChannelID string
}

func (cis *CollectionInfoShim) CollectionValidationInfo(chaincodeName, collectionName string, validationState validation.State) ([]byte, error, error) {
	return cis.CollectionAndLifecycleResources.CollectionValidationInfo(cis.ChannelID, chaincodeName, collectionName, validationState)
}

type gossipSupport struct {
	channelconfig.Application
	configtx.Validator
	channelconfig.Channel
}

func ConfigBlockFromLedger(ledger ledger.PeerLedger) (*common.Block, error) {
	peerLogger.Debugf("Getting config block")

	
	blockchainInfo, err := ledger.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}
	lastBlock, err := ledger.GetBlockByNumber(blockchainInfo.Height - 1)
	if err != nil {
		return nil, err
	}

	
	configBlockIndex, err := protoutil.GetLastConfigIndexFromBlock(lastBlock)
	if err != nil {
		return nil, err
	}

	
	configBlock, err := ledger.GetBlockByNumber(configBlockIndex)
	if err != nil {
		return nil, err
	}

	peerLogger.Debugf("Got config block[%d]", configBlockIndex)
	return configBlock, nil
}


func (p *Peer) updateTrustedRoots(cm channelconfig.Resources) {
	if !p.ServerConfig.SecOpts.UseTLS {
		return
	}

	
	peerLogger.Debugf("Updating trusted root authorities for channel %s", cm.ConfigtxValidator().ChannelID())

	p.CredentialSupport.BuildTrustedRootsForChain(cm)

	
	var trustedRoots [][]byte
	for _, roots := range p.CredentialSupport.AppRootCAsByChain() {
		trustedRoots = append(trustedRoots, roots...)
	}
	trustedRoots = append(trustedRoots, p.ServerConfig.SecOpts.ClientRootCAs...)
	trustedRoots = append(trustedRoots, p.ServerConfig.SecOpts.ServerRootCAs...)

	
	err := p.Server.SetClientRootCAs(trustedRoots)
	if err != nil {
		msg := "Failed to update trusted roots from latest config block. " +
			"This peer may not be able to communicate with members of channel %s (%s)"
		peerLogger.Warningf(msg, cm.ConfigtxValidator().ChannelID(), err)
	}
}






type DeliverChainManager struct {
	Peer *Peer
}

func (d DeliverChainManager) GetChain(chainID string) deliver.Chain {
	if channel := d.Peer.Channel(chainID); channel != nil {
		return channel
	}
	return nil
}



type fileLedgerBlockStore struct {
	ledger.PeerLedger
}

func (flbs fileLedgerBlockStore) AddBlock(*common.Block) error {
	return nil
}

func (flbs fileLedgerBlockStore) RetrieveBlocks(startBlockNumber uint64) (commonledger.ResultsIterator, error) {
	return flbs.GetBlocksIterator(startBlockNumber)
}


func NewConfigSupport(peer *Peer) cc.Manager {
	return &configSupport{
		peer: peer,
	}
}

type configSupport struct {
	peer *Peer
}





func (c *configSupport) GetChannelConfig(cid string) cc.Config {
	channel := c.peer.Channel(cid)
	if channel == nil {
		peerLogger.Errorf("[channel %s] channel not associated with this peer", cid)
		return nil
	}
	return channel.Resources().ConfigtxValidator()
}


type Peer struct {
	Server            *comm.GRPCServer
	ServerConfig      comm.ServerConfig
	CredentialSupport *comm.CredentialSupport
	StoreProvider     transientstore.StoreProvider
	GossipService     *gossipservice.GossipService
	LedgerMgr         *ledgermgmt.LedgerMgr
	CryptoProvider    bccsp.BCCSP

	
	
	validationWorkersSemaphore semaphore.Semaphore

	pluginMapper       plugin.Mapper
	channelInitializer func(cid string)

	
	mutex    sync.RWMutex
	channels map[string]*Channel
}

func (p *Peer) openStore(cid string) (transientstore.Store, error) {
	store, err := p.StoreProvider.OpenStore(cid)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (p *Peer) CreateChannel(
	cid string,
	cb *common.Block,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	legacyLifecycleValidation plugindispatcher.LifecycleResources,
	newLifecycleValidation plugindispatcher.CollectionAndLifecycleResources,
) error {
	l, err := p.LedgerMgr.CreateLedger(cid, cb)
	if err != nil {
		return errors.WithMessage(err, "cannot create ledger from genesis block")
	}

	if err := p.createChannel(cid, l, cb, p.pluginMapper, deployedCCInfoProvider, legacyLifecycleValidation, newLifecycleValidation); err != nil {
		return err
	}

	p.initChannel(cid)
	return nil
}


func retrievePersistedChannelConfig(ledger ledger.PeerLedger) (*common.Config, error) {
	qe, err := ledger.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	defer qe.Done()
	return retrievePersistedConf(qe, channelConfigKey)
}


func (p *Peer) createChannel(
	cid string,
	l ledger.PeerLedger,
	cb *common.Block,
	pluginMapper plugin.Mapper,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	legacyLifecycleValidation plugindispatcher.LifecycleResources,
	newLifecycleValidation plugindispatcher.CollectionAndLifecycleResources,
) error {
	chanConf, err := retrievePersistedChannelConfig(l)
	if err != nil {
		return err
	}

	var bundle *channelconfig.Bundle
	if chanConf != nil {
		bundle, err = channelconfig.NewBundle(cid, chanConf, p.CryptoProvider)
		if err != nil {
			return err
		}
	} else {
		
		
		envelopeConfig, err := protoutil.ExtractEnvelope(cb, 0)
		if err != nil {
			return err
		}

		bundle, err = channelconfig.NewBundleFromEnvelope(envelopeConfig, p.CryptoProvider)
		if err != nil {
			return err
		}
	}

	capabilitiesSupportedOrPanic(bundle)

	channelconfig.LogSanityChecks(bundle)

	gossipEventer := p.GossipService.NewConfigEventer()

	gossipCallbackWrapper := func(bundle *channelconfig.Bundle) {
		ac, ok := bundle.ApplicationConfig()
		if !ok {
			
			ac = nil
		}
		gossipEventer.ProcessConfigUpdate(&gossipSupport{
			Validator:   bundle.ConfigtxValidator(),
			Application: ac,
			Channel:     bundle.ChannelConfig(),
		})
		p.GossipService.SuspectPeers(func(identity api.PeerIdentityType) bool {
			
			
			
			
			return true
		})
	}

	trustedRootsCallbackWrapper := func(bundle *channelconfig.Bundle) {
		p.updateTrustedRoots(bundle)
	}

	mspCallback := func(bundle *channelconfig.Bundle) {
		
		mspmgmt.XXXSetMSPManager(cid, bundle.MSPManager())
	}

	channel := &Channel{
		ledger:         l,
		resources:      bundle,
		cryptoProvider: p.CryptoProvider,
	}

	channel.bundleSource = channelconfig.NewBundleSource(
		bundle,
		gossipCallbackWrapper,
		trustedRootsCallbackWrapper,
		mspCallback,
		channel.bundleUpdate,
	)

	committer := committer.NewLedgerCommitter(l)
	validator := &txvalidator.ValidationRouter{
		CapabilityProvider: channel,
		V14Validator: validatorv14.NewTxValidator(
			cid,
			p.validationWorkersSemaphore,
			channel,
			p.pluginMapper,
		),
		V20Validator: validatorv20.NewTxValidator(
			cid,
			p.validationWorkersSemaphore,
			channel,
			channel.Ledger(),
			&vir.ValidationInfoRetrieveShim{
				New:    newLifecycleValidation,
				Legacy: legacyLifecycleValidation,
			},
			&CollectionInfoShim{
				CollectionAndLifecycleResources: newLifecycleValidation,
				ChannelID:                       bundle.ConfigtxValidator().ChannelID(),
			},
			p.pluginMapper,
			policies.PolicyManagerGetterFunc(p.GetPolicyManager),
		),
	}

	ordererAddresses := bundle.ChannelConfig().OrdererAddresses()
	if len(ordererAddresses) == 0 {
		return errors.New("no ordering service endpoint provided in configuration block")
	}

	
	store, err := p.openStore(bundle.ConfigtxValidator().ChannelID())
	if err != nil {
		return errors.Wrapf(err, "[channel %s] failed opening transient store", bundle.ConfigtxValidator().ChannelID())
	}
	channel.store = store

	simpleCollectionStore := privdata.NewSimpleCollectionStore(l, deployedCCInfoProvider)
	p.GossipService.InitializeChannel(bundle.ConfigtxValidator().ChannelID(), ordererAddresses, gossipservice.Support{
		Validator:       validator,
		Committer:       committer,
		Store:           store,
		CollectionStore: simpleCollectionStore,
		IdDeserializeFactory: gossipprivdata.IdentityDeserializerFactoryFunc(func(chainID string) msp.IdentityDeserializer {
			return mspmgmt.GetManagerForChain(chainID)
		}),
		CapabilityProvider: channel,
	})

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.channels == nil {
		p.channels = map[string]*Channel{}
	}
	p.channels[cid] = channel

	return nil
}

func (p *Peer) Channel(cid string) *Channel {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if c, ok := p.channels[cid]; ok {
		return c
	}
	return nil
}

func (p *Peer) StoreForChannel(cid string) transientstore.Store {
	if c := p.Channel(cid); c != nil {
		return c.Store()
	}
	return nil
}



func (p *Peer) GetChannelsInfo() []*pb.ChannelInfo {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var channelInfos []*pb.ChannelInfo
	for key := range p.channels {
		ci := &pb.ChannelInfo{ChannelId: key}
		channelInfos = append(channelInfos, ci)
	}
	return channelInfos
}



func (p *Peer) GetChannelConfig(cid string) channelconfig.Resources {
	if c := p.Channel(cid); c != nil {
		return c.Resources()
	}
	return nil
}



func (p *Peer) GetStableChannelConfig(cid string) channelconfig.Resources {
	if c := p.Channel(cid); c != nil {
		return c.Resources()
	}
	return nil
}



func (p *Peer) GetLedger(cid string) ledger.PeerLedger {
	if c := p.Channel(cid); c != nil {
		return c.Ledger()
	}
	return nil
}


func (p *Peer) GetMSPIDs(cid string) []string {
	if c := p.Channel(cid); c != nil {
		return c.GetMSPIDs()
	}
	return nil
}



func (p *Peer) GetPolicyManager(cid string) policies.Manager {
	if c := p.Channel(cid); c != nil {
		return c.Resources().PolicyManager()
	}
	return nil
}


func (p *Peer) initChannel(cid string) {
	if p.channelInitializer != nil {
		
		peerLogger.Debugf("Initializing channel %s", cid)
		p.channelInitializer(cid)
	}
}

func (p *Peer) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	cc := p.GetChannelConfig(cid)
	if cc == nil {
		return nil, false
	}

	return cc.ApplicationConfig()
}




func (p *Peer) Initialize(
	init func(string),
	pm plugin.Mapper,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	legacyLifecycleValidation plugindispatcher.LifecycleResources,
	newLifecycleValidation plugindispatcher.CollectionAndLifecycleResources,
	nWorkers int,
) {
	
	p.validationWorkersSemaphore = semaphore.New(nWorkers)
	p.pluginMapper = pm
	p.channelInitializer = init

	ledgerIds, err := p.LedgerMgr.GetLedgerIDs()
	if err != nil {
		panic(fmt.Errorf("error in initializing ledgermgmt: %s", err))
	}

	for _, cid := range ledgerIds {
		peerLogger.Infof("Loading chain %s", cid)
		ledger, err := p.LedgerMgr.OpenLedger(cid)
		if err != nil {
			peerLogger.Errorf("Failed to load ledger %s(%s)", cid, err)
			peerLogger.Debugf("Error while loading ledger %s with message %s. We continue to the next ledger rather than abort.", cid, err)
			continue
		}
		cb, err := ConfigBlockFromLedger(ledger)
		if err != nil {
			peerLogger.Errorf("Failed to find config block on ledger %s(%s)", cid, err)
			peerLogger.Debugf("Error while looking for config block on ledger %s with message %s. We continue to the next ledger rather than abort.", cid, err)
			continue
		}
		
		err = p.createChannel(cid, ledger, cb, pm, deployedCCInfoProvider, legacyLifecycleValidation, newLifecycleValidation)
		if err != nil {
			peerLogger.Errorf("Failed to load chain %s(%s)", cid, err)
			peerLogger.Debugf("Error reloading chain %s with message %s. We continue to the next chain rather than abort.", cid, err)
			continue
		}

		p.initChannel(cid)
	}
}
