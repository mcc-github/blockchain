/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	cc "github.com/mcc-github/blockchain/common/config"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	fileledger "github.com/mcc-github/blockchain/common/ledger/blockledger/file"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/service"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
)

var peerLogger = flogging.MustGetLogger("peer")

var peerServer *comm.GRPCServer

var configTxProcessor = newConfigTxProcessor()
var ConfigTxProcessors = customtx.Processors{
	common.HeaderType_CONFIG: configTxProcessor,
}


var credSupport = comm.GetCredentialSupport()

type gossipSupport struct {
	channelconfig.Application
	configtx.Validator
	channelconfig.Channel
}

type chainSupport struct {
	bundleSource *channelconfig.BundleSource
	channelconfig.Resources
	channelconfig.Application
	ledger ledger.PeerLedger
}

var TransientStoreFactory = &storeProvider{stores: make(map[string]transientstore.Store)}

type storeProvider struct {
	stores map[string]transientstore.Store
	transientstore.StoreProvider
	sync.RWMutex
}

func (sp *storeProvider) StoreForChannel(channel string) transientstore.Store {
	sp.RLock()
	defer sp.RUnlock()
	return sp.stores[channel]
}

func (sp *storeProvider) OpenStore(ledgerID string) (transientstore.Store, error) {
	sp.Lock()
	defer sp.Unlock()
	if sp.StoreProvider == nil {
		sp.StoreProvider = transientstore.NewStoreProvider()
	}
	store, err := sp.StoreProvider.OpenStore(ledgerID)
	if err == nil {
		sp.stores[ledgerID] = store
	}
	return store, err
}

func (cs *chainSupport) Apply(configtx *common.ConfigEnvelope) error {
	err := cs.ConfigtxValidator().Validate(configtx)
	if err != nil {
		return err
	}

	
	if cs.bundleSource != nil {
		bundle, err := channelconfig.NewBundle(cs.ConfigtxValidator().ChainID(), configtx.Config)
		if err != nil {
			return err
		}

		channelconfig.LogSanityChecks(bundle)

		err = cs.bundleSource.ValidateNew(bundle)
		if err != nil {
			return err
		}

		capabilitiesSupportedOrPanic(bundle)

		cs.bundleSource.Update(bundle)
	}
	return nil
}

func capabilitiesSupportedOrPanic(res channelconfig.Resources) {
	ac, ok := res.ApplicationConfig()
	if !ok {
		peerLogger.Panicf("[channel %s] does not have application config so is incompatible", res.ConfigtxValidator().ChainID())
	}

	if err := ac.Capabilities().Supported(); err != nil {
		peerLogger.Panicf("[channel %s] incompatible: %s", res.ConfigtxValidator().ChainID(), err)
	}

	if err := res.ChannelConfig().Capabilities().Supported(); err != nil {
		peerLogger.Panicf("[channel %s] incompatible: %s", res.ConfigtxValidator().ChainID(), err)
	}
}

func (cs *chainSupport) Ledger() ledger.PeerLedger {
	return cs.ledger
}

func (cs *chainSupport) GetMSPIDs(cid string) []string {
	return GetMSPIDs(cid)
}


func (cs *chainSupport) Sequence() uint64 {
	sb := cs.bundleSource.StableBundle()
	return sb.ConfigtxValidator().Sequence()
}


func (cs *chainSupport) Reader() blockledger.Reader {
	return fileledger.NewFileLedger(fileLedgerBlockStore{cs.ledger})
}





func (cs *chainSupport) Errored() <-chan struct{} {
	return nil
}


type chain struct {
	cs        *chainSupport
	cb        *common.Block
	committer committer.Committer
}


var chains = struct {
	sync.RWMutex
	list map[string]*chain
}{list: make(map[string]*chain)}

var chainInitializer func(string)

var pluginMapper txvalidator.PluginMapper

var mockMSPIDGetter func(string) []string

func MockSetMSPIDGetter(mspIDGetter func(string) []string) {
	mockMSPIDGetter = mspIDGetter
}



var validationWorkersSemaphore *semaphore.Weighted




func Initialize(init func(string), ccp ccprovider.ChaincodeProvider, sccp sysccprovider.SystemChaincodeProvider, pm txvalidator.PluginMapper, pr *platforms.Registry, deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider) {
	nWorkers := viper.GetInt("peer.validatorPoolSize")
	if nWorkers <= 0 {
		nWorkers = runtime.NumCPU()
	}
	validationWorkersSemaphore = semaphore.NewWeighted(int64(nWorkers))

	pluginMapper = pm
	chainInitializer = init

	var cb *common.Block
	var ledger ledger.PeerLedger
	ledgermgmt.Initialize(&ledgermgmt.Initializer{
		CustomTxProcessors:            ConfigTxProcessors,
		PlatformRegistry:              pr,
		DeployedChaincodeInfoProvider: deployedCCInfoProvider,
	})
	ledgerIds, err := ledgermgmt.GetLedgerIDs()
	if err != nil {
		panic(fmt.Errorf("Error in initializing ledgermgmt: %s", err))
	}
	for _, cid := range ledgerIds {
		peerLogger.Infof("Loading chain %s", cid)
		if ledger, err = ledgermgmt.OpenLedger(cid); err != nil {
			peerLogger.Warningf("Failed to load ledger %s(%s)", cid, err)
			peerLogger.Debugf("Error while loading ledger %s with message %s. We continue to the next ledger rather than abort.", cid, err)
			continue
		}
		if cb, err = getCurrConfigBlockFromLedger(ledger); err != nil {
			peerLogger.Warningf("Failed to find config block on ledger %s(%s)", cid, err)
			peerLogger.Debugf("Error while looking for config block on ledger %s with message %s. We continue to the next ledger rather than abort.", cid, err)
			continue
		}
		
		if err = createChain(cid, ledger, cb, ccp, sccp, pm); err != nil {
			peerLogger.Warningf("Failed to load chain %s(%s)", cid, err)
			peerLogger.Debugf("Error reloading chain %s with message %s. We continue to the next chain rather than abort.", cid, err)
			continue
		}

		InitChain(cid)
	}
}


func InitChain(cid string) {
	if chainInitializer != nil {
		
		peerLogger.Debugf("Initializing channel %s", cid)
		chainInitializer(cid)
	}
}

func getCurrConfigBlockFromLedger(ledger ledger.PeerLedger) (*common.Block, error) {
	peerLogger.Debugf("Getting config block")

	
	blockchainInfo, err := ledger.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}
	lastBlock, err := ledger.GetBlockByNumber(blockchainInfo.Height - 1)
	if err != nil {
		return nil, err
	}

	
	configBlockIndex, err := utils.GetLastConfigIndexFromBlock(lastBlock)
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


func createChain(cid string, ledger ledger.PeerLedger, cb *common.Block, ccp ccprovider.ChaincodeProvider, sccp sysccprovider.SystemChaincodeProvider, pm txvalidator.PluginMapper) error {
	chanConf, err := retrievePersistedChannelConfig(ledger)
	if err != nil {
		return err
	}

	var bundle *channelconfig.Bundle

	if chanConf != nil {
		bundle, err = channelconfig.NewBundle(cid, chanConf)
		if err != nil {
			return err
		}
	} else {
		
		
		envelopeConfig, err := utils.ExtractEnvelope(cb, 0)
		if err != nil {
			return err
		}

		bundle, err = channelconfig.NewBundleFromEnvelope(envelopeConfig)
		if err != nil {
			return err
		}
	}

	capabilitiesSupportedOrPanic(bundle)

	channelconfig.LogSanityChecks(bundle)

	gossipEventer := service.GetGossipService().NewConfigEventer()

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
		service.GetGossipService().SuspectPeers(func(identity api.PeerIdentityType) bool {
			
			
			
			
			return true
		})
	}

	trustedRootsCallbackWrapper := func(bundle *channelconfig.Bundle) {
		updateTrustedRoots(bundle)
	}

	mspCallback := func(bundle *channelconfig.Bundle) {
		
		mspmgmt.XXXSetMSPManager(cid, bundle.MSPManager())
	}

	ac, ok := bundle.ApplicationConfig()
	if !ok {
		ac = nil
	}

	cs := &chainSupport{
		Application: ac, 
		ledger:      ledger,
	}

	peerSingletonCallback := func(bundle *channelconfig.Bundle) {
		ac, ok := bundle.ApplicationConfig()
		if !ok {
			ac = nil
		}
		cs.Application = ac
		cs.Resources = bundle
	}

	cs.bundleSource = channelconfig.NewBundleSource(
		bundle,
		gossipCallbackWrapper,
		trustedRootsCallbackWrapper,
		mspCallback,
		peerSingletonCallback,
	)

	vcs := struct {
		*chainSupport
		*semaphore.Weighted
	}{cs, validationWorkersSemaphore}
	validator := txvalidator.NewTxValidator(vcs, sccp, pm)
	c := committer.NewLedgerCommitterReactive(ledger, func(block *common.Block) error {
		chainID, err := utils.GetChainIDFromBlock(block)
		if err != nil {
			return err
		}
		return SetCurrConfigBlock(block, chainID)
	})

	ordererAddresses := bundle.ChannelConfig().OrdererAddresses()
	if len(ordererAddresses) == 0 {
		return errors.New("no ordering service endpoint provided in configuration block")
	}

	
	store, err := TransientStoreFactory.OpenStore(bundle.ConfigtxValidator().ChainID())
	if err != nil {
		return errors.Wrapf(err, "[channel %s] failed opening transient store", bundle.ConfigtxValidator().ChainID())
	}
	csStoreSupport := &collectionSupport{
		PeerLedger: ledger,
	}
	simpleCollectionStore := privdata.NewSimpleCollectionStore(csStoreSupport)

	service.GetGossipService().InitializeChannel(bundle.ConfigtxValidator().ChainID(), ordererAddresses, service.Support{
		Validator:            validator,
		Committer:            c,
		Store:                store,
		Cs:                   simpleCollectionStore,
		IdDeserializeFactory: csStoreSupport,
	})

	chains.Lock()
	defer chains.Unlock()
	chains.list[cid] = &chain{
		cs:        cs,
		cb:        cb,
		committer: c,
	}

	return nil
}


func CreateChainFromBlock(cb *common.Block, ccp ccprovider.ChaincodeProvider, sccp sysccprovider.SystemChaincodeProvider) error {
	cid, err := utils.GetChainIDFromBlock(cb)
	if err != nil {
		return err
	}

	var l ledger.PeerLedger
	if l, err = ledgermgmt.CreateLedger(cb); err != nil {
		return errors.WithMessage(err, "cannot create ledger from genesis block")
	}

	return createChain(cid, l, cb, ccp, sccp, pluginMapper)
}



func GetLedger(cid string) ledger.PeerLedger {
	chains.RLock()
	defer chains.RUnlock()
	if c, ok := chains.list[cid]; ok {
		return c.cs.ledger
	}
	return nil
}



func GetStableChannelConfig(cid string) channelconfig.Resources {
	chains.RLock()
	defer chains.RUnlock()
	if c, ok := chains.list[cid]; ok {
		return c.cs.bundleSource.StableBundle()
	}
	return nil
}



func GetChannelConfig(cid string) channelconfig.Resources {
	chains.RLock()
	defer chains.RUnlock()
	if c, ok := chains.list[cid]; ok {
		return c.cs
	}
	return nil
}



func GetPolicyManager(cid string) policies.Manager {
	chains.RLock()
	defer chains.RUnlock()
	if c, ok := chains.list[cid]; ok {
		return c.cs.PolicyManager()
	}
	return nil
}



func GetCurrConfigBlock(cid string) *common.Block {
	chains.RLock()
	defer chains.RUnlock()
	if c, ok := chains.list[cid]; ok {
		return c.cb
	}
	return nil
}


func updateTrustedRoots(cm channelconfig.Resources) {
	
	peerLogger.Debugf("Updating trusted root authorities for channel %s", cm.ConfigtxValidator().ChainID())
	var serverConfig comm.ServerConfig
	var err error
	
	serverConfig, err = GetServerConfig()
	if err == nil && serverConfig.SecOpts.UseTLS {
		buildTrustedRootsForChain(cm)

		
		trustedRoots := [][]byte{}
		credSupport.RLock()
		defer credSupport.RUnlock()
		for _, roots := range credSupport.AppRootCAsByChain {
			trustedRoots = append(trustedRoots, roots...)
		}
		
		if len(serverConfig.SecOpts.ClientRootCAs) > 0 {
			trustedRoots = append(trustedRoots, serverConfig.SecOpts.ClientRootCAs...)
		}
		if len(serverConfig.SecOpts.ServerRootCAs) > 0 {
			trustedRoots = append(trustedRoots, serverConfig.SecOpts.ServerRootCAs...)
		}

		server := peerServer
		
		if server != nil {
			err := server.SetClientRootCAs(trustedRoots)
			if err != nil {
				msg := "Failed to update trusted roots for peer from latest config " +
					"block.  This peer may not be able to communicate " +
					"with members of channel %s (%s)"
				peerLogger.Warningf(msg, cm.ConfigtxValidator().ChainID(), err)
			}
		}
	}
}



func buildTrustedRootsForChain(cm channelconfig.Resources) {
	credSupport.Lock()
	defer credSupport.Unlock()

	appRootCAs := [][]byte{}
	ordererRootCAs := [][]byte{}
	appOrgMSPs := make(map[string]struct{})
	ordOrgMSPs := make(map[string]struct{})

	if ac, ok := cm.ApplicationConfig(); ok {
		
		for _, appOrg := range ac.Organizations() {
			appOrgMSPs[appOrg.MSPID()] = struct{}{}
		}
	}

	if ac, ok := cm.OrdererConfig(); ok {
		
		for _, ordOrg := range ac.Organizations() {
			ordOrgMSPs[ordOrg.MSPID()] = struct{}{}
		}
	}

	cid := cm.ConfigtxValidator().ChainID()
	peerLogger.Debugf("updating root CAs for channel [%s]", cid)
	msps, err := cm.MSPManager().GetMSPs()
	if err != nil {
		peerLogger.Errorf("Error getting root CAs for channel %s (%s)", cid, err)
	}
	if err == nil {
		for k, v := range msps {
			
			if v.GetType() == msp.FABRIC {
				for _, root := range v.GetTLSRootCerts() {
					
					if _, ok := appOrgMSPs[k]; ok {
						peerLogger.Debugf("adding app root CAs for MSP [%s]", k)
						appRootCAs = append(appRootCAs, root)
					}
					
					if _, ok := ordOrgMSPs[k]; ok {
						peerLogger.Debugf("adding orderer root CAs for MSP [%s]", k)
						ordererRootCAs = append(ordererRootCAs, root)
					}
				}
				for _, intermediate := range v.GetTLSIntermediateCerts() {
					
					if _, ok := appOrgMSPs[k]; ok {
						peerLogger.Debugf("adding app root CAs for MSP [%s]", k)
						appRootCAs = append(appRootCAs, intermediate)
					}
					
					if _, ok := ordOrgMSPs[k]; ok {
						peerLogger.Debugf("adding orderer root CAs for MSP [%s]", k)
						ordererRootCAs = append(ordererRootCAs, intermediate)
					}
				}
			}
		}
		credSupport.AppRootCAsByChain[cid] = appRootCAs
		credSupport.OrdererRootCAsByChain[cid] = ordererRootCAs
	}
}


func GetMSPIDs(cid string) []string {
	chains.RLock()
	defer chains.RUnlock()

	
	
	if mockMSPIDGetter != nil {
		return mockMSPIDGetter(cid)
	}
	if c, ok := chains.list[cid]; ok {
		if c == nil || c.cs == nil {
			return nil
		}
		ac, ok := c.cs.ApplicationConfig()
		if !ok || ac.Organizations() == nil {
			return nil
		}

		orgs := ac.Organizations()
		toret := make([]string, len(orgs))
		i := 0
		for _, org := range orgs {
			toret[i] = org.MSPID()
			i++
		}

		return toret
	}
	return nil
}


func SetCurrConfigBlock(block *common.Block, cid string) error {
	chains.Lock()
	defer chains.Unlock()
	if c, ok := chains.list[cid]; ok {
		c.cb = block
		return nil
	}
	return errors.Errorf("[channel %s] channel not associated with this peer", cid)
}


func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}



func GetChannelsInfo() []*pb.ChannelInfo {
	
	var channelInfoArray []*pb.ChannelInfo

	chains.RLock()
	defer chains.RUnlock()
	for key := range chains.list {
		channelInfo := &pb.ChannelInfo{ChannelId: key}

		
		channelInfoArray = append(channelInfoArray, channelInfo)
	}

	return channelInfoArray
}


func NewChannelPolicyManagerGetter() policies.ChannelPolicyManagerGetter {
	return &channelPolicyManagerGetter{}
}

type channelPolicyManagerGetter struct{}

func (c *channelPolicyManagerGetter) Manager(channelID string) (policies.Manager, bool) {
	policyManager := GetPolicyManager(channelID)
	return policyManager, policyManager != nil
}



func NewPeerServer(listenAddress string, serverConfig comm.ServerConfig) (*comm.GRPCServer, error) {
	var err error
	peerServer, err = comm.NewGRPCServer(listenAddress, serverConfig)
	if err != nil {
		peerLogger.Errorf("Failed to create peer server (%s)", err)
		return nil, err
	}
	return peerServer, nil
}

type collectionSupport struct {
	ledger.PeerLedger
}

func (cs *collectionSupport) GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error) {
	return cs.NewQueryExecutor()
}

func (*collectionSupport) GetIdentityDeserializer(chainID string) msp.IdentityDeserializer {
	return mspmgmt.GetManagerForChain(chainID)
}






type DeliverChainManager struct {
}

func (DeliverChainManager) GetChain(chainID string) (deliver.Chain, bool) {
	channel, ok := chains.list[chainID]
	if !ok {
		return nil, ok
	}
	return channel.cs, ok
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


func NewConfigSupport() cc.Manager {
	return &configSupport{}
}

type configSupport struct {
}





func (*configSupport) GetChannelConfig(channel string) cc.Config {
	chains.RLock()
	defer chains.RUnlock()
	chain := chains.list[channel]
	if chain == nil {
		peerLogger.Errorf("[channel %s] channel not associated with this peer", channel)
		return nil
	}
	return chain.cs.bundleSource.ConfigtxValidator()
}
