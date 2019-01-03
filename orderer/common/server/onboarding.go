/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
)

const (
	defaultReplicationBackgroundRefreshInterval = time.Minute * 5
	replicationBackgroundInitialRefreshInterval = time.Second * 10
)

type replicationInitiator struct {
	channelLister cluster.ChannelLister
	logger        *flogging.FabricLogger
	secOpts       *comm.SecureOptions
	conf          *localconfig.TopLevel
	lf            cluster.LedgerFactory
	signer        crypto.LocalSigner
}

func (ri *replicationInitiator) replicateIfNeeded(bootstrapBlock *common.Block) {
	if bootstrapBlock.Header.Number == 0 {
		ri.logger.Debug("Booted with a genesis block, replication isn't an option")
		return
	}
	ri.replicateNeededChannels(bootstrapBlock)
}

func (ri *replicationInitiator) createReplicator(bootstrapBlock *common.Block, filter func(string) bool) *cluster.Replicator {
	consenterCert := etcdraft.ConsenterCertificate(ri.secOpts.Certificate)
	systemChannelName, err := utils.GetChainIDFromBlock(bootstrapBlock)
	if err != nil {
		ri.logger.Panicf("Failed extracting system channel name from bootstrap block: %v", err)
	}
	pullerConfig := cluster.PullerConfigFromTopLevelConfig(systemChannelName, ri.conf, ri.secOpts.Key, ri.secOpts.Certificate, ri.signer)
	puller, err := cluster.BlockPullerFromConfigBlock(pullerConfig, bootstrapBlock)
	if err != nil {
		ri.logger.Panicf("Failed creating puller config from bootstrap block: %v", err)
	}

	pullerLogger := flogging.MustGetLogger("orderer.common.cluster")

	replicator := &cluster.Replicator{
		Filter:           filter,
		LedgerFactory:    ri.lf,
		SystemChannel:    systemChannelName,
		BootBlock:        bootstrapBlock,
		Logger:           pullerLogger,
		AmIPartOfChannel: consenterCert.IsConsenterOfChannel,
		Puller:           puller,
		ChannelLister: &cluster.ChainInspector{
			Logger:          pullerLogger,
			Puller:          puller,
			LastConfigBlock: bootstrapBlock,
		},
	}

	
	if ri.channelLister != nil {
		replicator.ChannelLister = ri.channelLister
	}

	return replicator
}

func (ri *replicationInitiator) replicateNeededChannels(bootstrapBlock *common.Block) {
	replicator := ri.createReplicator(bootstrapBlock, cluster.AnyChannel)
	defer replicator.Puller.Close()
	replicationNeeded, err := replicator.IsReplicationNeeded()
	if err != nil {
		ri.logger.Panicf("Failed determining whether replication is needed: %v", err)
	}

	if !replicationNeeded {
		ri.logger.Info("Replication isn't needed")
		return
	}

	ri.logger.Info("Will now replicate chains")
	replicator.ReplicateChains()
}



func (ri *replicationInitiator) ReplicateChains(lastConfigBlock *common.Block, chains []string) []string {
	ri.logger.Info("Will now replicate chains", chains)
	wantedChannels := make(map[string]struct{})
	for _, chain := range chains {
		wantedChannels[chain] = struct{}{}
	}
	filter := func(channelName string) bool {
		_, exists := wantedChannels[channelName]
		return exists
	}
	replicator := ri.createReplicator(lastConfigBlock, filter)
	replicator.DoNotPanicIfClusterNotReachable = true
	defer replicator.Puller.Close()
	return replicator.ReplicateChains()
}

type ledgerFactory struct {
	blockledger.Factory
}

func (lf *ledgerFactory) GetOrCreate(chainID string) (cluster.LedgerWriter, error) {
	return lf.Factory.GetOrCreate(chainID)
}




type ChainReplicator interface {
	
	
	ReplicateChains(lastConfigBlock *common.Block, chains []string) []string
}


type inactiveChainReplicator struct {
	logger                            *flogging.FabricLogger
	retrieveLastSysChannelConfigBlock func() *common.Block
	replicator                        ChainReplicator
	scheduleChan                      <-chan time.Time
	quitChan                          chan struct{}
	lock                              sync.RWMutex
	chains2CreationCallbacks          map[string]chainCreation
}

func (dc *inactiveChainReplicator) Channels() []cluster.ChannelGenesisBlock {
	dc.lock.RLock()
	defer dc.lock.RUnlock()

	var res []cluster.ChannelGenesisBlock
	for name, chain := range dc.chains2CreationCallbacks {
		res = append(res, cluster.ChannelGenesisBlock{
			ChannelName:  name,
			GenesisBlock: chain.genesisBlock,
		})
	}
	return res
}

func (dc *inactiveChainReplicator) Close() {}

type chainCreation struct {
	create       func()
	genesisBlock *common.Block
}



func (dc *inactiveChainReplicator) TrackChain(chain string, genesisBlock *common.Block, createChainCallback func()) {
	if genesisBlock == nil {
		dc.logger.Panicf("Called with a nil genesis block")
	}
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.logger.Infof("Adding %s to the set of chains to track", chain)
	dc.chains2CreationCallbacks[chain] = chainCreation{
		genesisBlock: genesisBlock,
		create:       createChainCallback,
	}
}

func (dc *inactiveChainReplicator) run() {
	for {
		select {
		case <-dc.scheduleChan:
			dc.replicateDisabledChains()
		case <-dc.quitChan:
			return
		}
	}
}

func (dc *inactiveChainReplicator) replicateDisabledChains() {
	chains := dc.listInactiveChains()
	if len(chains) == 0 {
		dc.logger.Debugf("No inactive chains to try to replicate")
		return
	}
	dc.logger.Infof("Found %d inactive chains: %v", len(chains), chains)
	lastSystemChannelConfigBlock := dc.retrieveLastSysChannelConfigBlock()
	replicatedChains := dc.replicator.ReplicateChains(lastSystemChannelConfigBlock, chains)
	dc.logger.Infof("Successfully replicated %d chains: %v", len(replicatedChains), replicatedChains)
	dc.lock.Lock()
	defer dc.lock.Unlock()
	for _, chainName := range replicatedChains {
		chain := dc.chains2CreationCallbacks[chainName]
		delete(dc.chains2CreationCallbacks, chainName)
		chain.create()
	}
}

func (dc *inactiveChainReplicator) stop() {
	close(dc.quitChan)
}

func (dc *inactiveChainReplicator) listInactiveChains() []string {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	var chains []string
	for chain := range dc.chains2CreationCallbacks {
		chains = append(chains, chain)
	}
	return chains
}

func exponentialDurationSeries(initialDuration, maxDuration time.Duration) func() time.Duration {
	exp := &exponentialDuration{
		n:   initialDuration,
		max: maxDuration,
	}
	return exp.next
}

type exponentialDuration struct {
	n   time.Duration
	max time.Duration
}

func (exp *exponentialDuration) next() time.Duration {
	n := exp.n
	exp.n *= 2
	if exp.n > exp.max {
		exp.n = exp.max
	}
	return n
}

func makeTickChannel(computeSleepDuration func() time.Duration, sleep func(time.Duration)) <-chan time.Time {
	c := make(chan time.Time)
	go func() {
		for {
			sleep(computeSleepDuration())
			c <- time.Now()
		}
	}()

	return c
}
