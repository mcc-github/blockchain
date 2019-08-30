/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/fileledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/msp"
)


type Channel struct {
	ledger         ledger.PeerLedger
	store          transientstore.Store
	cryptoProvider bccsp.BCCSP

	
	applyLock sync.Mutex
	
	
	bundleSource *channelconfig.BundleSource

	
	lock sync.RWMutex
	
	
	resources channelconfig.Resources
}


func (c *Channel) Apply(configtx *common.ConfigEnvelope) error {
	c.applyLock.Lock()
	defer c.applyLock.Unlock()

	configTxValidator := c.Resources().ConfigtxValidator()
	err := configTxValidator.Validate(configtx)
	if err != nil {
		return err
	}

	bundle, err := channelconfig.NewBundle(configTxValidator.ChannelID(), configtx.Config, c.cryptoProvider)
	if err != nil {
		return err
	}

	channelconfig.LogSanityChecks(bundle)
	err = c.bundleSource.ValidateNew(bundle)
	if err != nil {
		return err
	}

	capabilitiesSupportedOrPanic(bundle)

	c.bundleSource.Update(bundle)
	return nil
}



func (c *Channel) bundleUpdate(b *channelconfig.Bundle) {
	c.lock.Lock()
	c.resources = b
	c.lock.Unlock()
}


func (c *Channel) Resources() channelconfig.Resources {
	c.lock.RLock()
	res := c.resources
	c.lock.RUnlock()
	return res
}


func (c *Channel) Sequence() uint64 {
	return c.Resources().ConfigtxValidator().Sequence()
}



func (c *Channel) PolicyManager() policies.Manager {
	return c.Resources().PolicyManager()
}



func (c *Channel) Capabilities() channelconfig.ApplicationCapabilities {
	ac, ok := c.Resources().ApplicationConfig()
	if !ok {
		return nil
	}
	return ac.Capabilities()
}



func (c *Channel) GetMSPIDs() []string {
	ac, ok := c.Resources().ApplicationConfig()
	if !ok || ac.Organizations() == nil {
		return nil
	}

	var mspIDs []string
	for _, org := range ac.Organizations() {
		mspIDs = append(mspIDs, org.MSPID())
	}

	return mspIDs
}



func (c *Channel) MSPManager() msp.MSPManager {
	return c.Resources().MSPManager()
}


func (c *Channel) Ledger() ledger.PeerLedger {
	return c.ledger
}


func (c *Channel) Store() transientstore.Store {
	return c.store
}



func (c *Channel) Reader() blockledger.Reader {
	return fileledger.NewFileLedger(fileLedgerBlockStore{c.ledger})
}





func (c *Channel) Errored() <-chan struct{} {
	
	
	return nil
}

func capabilitiesSupportedOrPanic(res channelconfig.Resources) {
	ac, ok := res.ApplicationConfig()
	if !ok {
		peerLogger.Panicf("[channel %s] does not have application config so is incompatible", res.ConfigtxValidator().ChannelID())
	}

	if err := ac.Capabilities().Supported(); err != nil {
		peerLogger.Panicf("[channel %s] incompatible: %s", res.ConfigtxValidator().ChannelID(), err)
	}

	if err := res.ChannelConfig().Capabilities().Supported(); err != nil {
		peerLogger.Panicf("[channel %s] incompatible: %s", res.ConfigtxValidator().ChannelID(), err)
	}
}
