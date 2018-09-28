/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/




package multichannel

import (
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

const (
	pkgLogID = "orderer/commmon/multichannel"

	msgVersion = int32(0)
	epoch      = 0
)

var logger = flogging.MustGetLogger(pkgLogID)


func checkResources(res channelconfig.Resources) error {
	channelconfig.LogSanityChecks(res)
	oc, ok := res.OrdererConfig()
	if !ok {
		return errors.New("config does not contain orderer config")
	}
	if err := oc.Capabilities().Supported(); err != nil {
		return errors.Wrapf(err, "config requires unsupported orderer capabilities: %s", err)
	}
	if err := res.ChannelConfig().Capabilities().Supported(); err != nil {
		return errors.Wrapf(err, "config requires unsupported channel capabilities: %s", err)
	}
	return nil
}


func checkResourcesOrPanic(res channelconfig.Resources) {
	if err := checkResources(res); err != nil {
		logger.Panicf("[channel %s] %s", res.ConfigtxValidator().ChainID(), err)
	}
}

type mutableResources interface {
	channelconfig.Resources
	Update(*channelconfig.Bundle)
}

type configResources struct {
	mutableResources
}

func (cr *configResources) CreateBundle(channelID string, config *cb.Config) (*channelconfig.Bundle, error) {
	return channelconfig.NewBundle(channelID, config)
}

func (cr *configResources) Update(bndl *channelconfig.Bundle) {
	checkResourcesOrPanic(bndl)
	cr.mutableResources.Update(bndl)
}

func (cr *configResources) SharedConfig() channelconfig.Orderer {
	oc, ok := cr.OrdererConfig()
	if !ok {
		logger.Panicf("[channel %s] has no orderer configuration", cr.ConfigtxValidator().ChainID())
	}
	return oc
}

type ledgerResources struct {
	*configResources
	blockledger.ReadWriter
}


type Registrar struct {
	lock   sync.RWMutex
	chains map[string]*ChainSupport

	consenters      map[string]consensus.Consenter
	ledgerFactory   blockledger.Factory
	signer          crypto.LocalSigner
	systemChannelID string
	systemChannel   *ChainSupport
	templator       msgprocessor.ChannelConfigTemplator
	callbacks       []func(bundle *channelconfig.Bundle)
}

func getConfigTx(reader blockledger.Reader) *cb.Envelope {
	lastBlock := blockledger.GetBlock(reader, reader.Height()-1)
	index, err := utils.GetLastConfigIndexFromBlock(lastBlock)
	if err != nil {
		logger.Panicf("Chain did not have appropriately encoded last config in its latest block: %s", err)
	}
	configBlock := blockledger.GetBlock(reader, index)
	if configBlock == nil {
		logger.Panicf("Config block does not exist")
	}

	return utils.ExtractEnvelopeOrPanic(configBlock, 0)
}


func NewRegistrar(ledgerFactory blockledger.Factory,
	signer crypto.LocalSigner, callbacks ...func(bundle *channelconfig.Bundle)) *Registrar {
	r := &Registrar{
		chains:        make(map[string]*ChainSupport),
		ledgerFactory: ledgerFactory,
		signer:        signer,
		callbacks:     callbacks,
	}

	return r
}

func (r *Registrar) Initialize(consenters map[string]consensus.Consenter) {
	r.consenters = consenters
	existingChains := r.ledgerFactory.ChainIDs()
	for _, chainID := range existingChains {
		rl, err := r.ledgerFactory.GetOrCreate(chainID)
		if err != nil {
			logger.Panicf("Ledger factory reported chainID %s but could not retrieve it: %s", chainID, err)
		}
		configTx := getConfigTx(rl)
		if configTx == nil {
			logger.Panic("Programming error, configTx should never be nil here")
		}
		ledgerResources := r.newLedgerResources(configTx)
		chainID := ledgerResources.ConfigtxValidator().ChainID()

		if _, ok := ledgerResources.ConsortiumsConfig(); ok {
			if r.systemChannelID != "" {
				logger.Panicf("There appear to be two system chains %s and %s", r.systemChannelID, chainID)
			}
			chain := newChainSupport(
				r,
				ledgerResources,
				r.consenters,
				r.signer)
			r.templator = msgprocessor.NewDefaultTemplator(chain)
			chain.Processor = msgprocessor.NewSystemChannel(chain, r.templator, msgprocessor.CreateSystemChannelFilters(r, chain))

			
			iter, pos := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Oldest{Oldest: &ab.SeekOldest{}}})
			defer iter.Close()
			if pos != uint64(0) {
				logger.Panicf("Error iterating over system channel: '%s', expected position 0, got %d", chainID, pos)
			}
			genesisBlock, status := iter.Next()
			if status != cb.Status_SUCCESS {
				logger.Panicf("Error reading genesis block of system channel '%s'", chainID)
			}
			logger.Infof("Starting system channel '%s' with genesis block hash %x and orderer type %s", chainID, genesisBlock.Header.Hash(), chain.SharedConfig().ConsensusType())

			r.chains[chainID] = chain
			r.systemChannelID = chainID
			r.systemChannel = chain
			
			defer chain.start()
		} else {
			logger.Debugf("Starting chain: %s", chainID)
			chain := newChainSupport(
				r,
				ledgerResources,
				r.consenters,
				r.signer)
			r.chains[chainID] = chain
			chain.start()
		}

	}

	if r.systemChannelID == "" {
		logger.Panicf("No system chain found.  If bootstrapping, does your system channel contain a consortiums group definition?")
	}
}


func (r *Registrar) SystemChannelID() string {
	return r.systemChannelID
}




func (r *Registrar) BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, *ChainSupport, error) {
	chdr, err := utils.ChannelHeader(msg)
	if err != nil {
		return nil, false, nil, fmt.Errorf("could not determine channel ID: %s", err)
	}

	cs, ok := r.GetChain(chdr.ChannelId)
	if !ok {
		cs = r.systemChannel
	}

	isConfig := false
	switch cs.ClassifyMsg(chdr) {
	case msgprocessor.ConfigUpdateMsg:
		isConfig = true
	case msgprocessor.ConfigMsg:
		return chdr, false, nil, errors.New("message is of type that cannot be processed directly")
	default:
	}

	return chdr, isConfig, cs, nil
}


func (r *Registrar) GetChain(chainID string) (*ChainSupport, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	cs, ok := r.chains[chainID]
	return cs, ok
}

func (r *Registrar) newLedgerResources(configTx *cb.Envelope) *ledgerResources {
	payload, err := utils.UnmarshalPayload(configTx.Payload)
	if err != nil {
		logger.Panicf("Error umarshaling envelope to payload: %s", err)
	}

	if payload.Header == nil {
		logger.Panicf("Missing channel header: %s", err)
	}

	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Panicf("Error unmarshaling channel header: %s", err)
	}

	configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		logger.Panicf("Error umarshaling config envelope from payload data: %s", err)
	}

	bundle, err := channelconfig.NewBundle(chdr.ChannelId, configEnvelope.Config)
	if err != nil {
		logger.Panicf("Error creating channelconfig bundle: %s", err)
	}

	checkResourcesOrPanic(bundle)

	ledger, err := r.ledgerFactory.GetOrCreate(chdr.ChannelId)
	if err != nil {
		logger.Panicf("Error getting ledger for %s", chdr.ChannelId)
	}

	return &ledgerResources{
		configResources: &configResources{
			mutableResources: channelconfig.NewBundleSource(bundle, r.callbacks...),
		},
		ReadWriter: ledger,
	}
}

func (r *Registrar) newChain(configtx *cb.Envelope) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ledgerResources := r.newLedgerResources(configtx)
	ledgerResources.Append(blockledger.CreateNextBlock(ledgerResources, []*cb.Envelope{configtx}))

	
	newChains := make(map[string]*ChainSupport)
	for key, value := range r.chains {
		newChains[key] = value
	}

	cs := newChainSupport(r, ledgerResources, r.consenters, r.signer)
	chainID := ledgerResources.ConfigtxValidator().ChainID()

	logger.Infof("Created and starting new chain %s", chainID)

	newChains[string(chainID)] = cs
	cs.start()

	r.chains = newChains
}


func (r *Registrar) ChannelsCount() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.chains)
}


func (r *Registrar) NewChannelConfig(envConfigUpdate *cb.Envelope) (channelconfig.Resources, error) {
	return r.templator.NewChannelConfig(envConfigUpdate)
}


func (r *Registrar) CreateBundle(channelID string, config *cb.Config) (channelconfig.Resources, error) {
	return channelconfig.NewBundle(channelID, config)
}
