/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"sync/atomic"

	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
)






type BundleSource struct {
	bundle    atomic.Value
	callbacks []BundleActor
}


type BundleActor func(bundle *Bundle)





func NewBundleSource(bundle *Bundle, callbacks ...BundleActor) *BundleSource {
	bs := &BundleSource{
		callbacks: callbacks,
	}
	bs.Update(bundle)
	return bs
}


func (bs *BundleSource) Update(newBundle *Bundle) {
	bs.bundle.Store(newBundle)
	for _, callback := range bs.callbacks {
		callback(newBundle)
	}
}









func (bs *BundleSource) StableBundle() *Bundle {
	return bs.bundle.Load().(*Bundle)
}


func (bs *BundleSource) PolicyManager() policies.Manager {
	return bs.StableBundle().policyManager
}


func (bs *BundleSource) MSPManager() msp.MSPManager {
	return bs.StableBundle().mspManager
}


func (bs *BundleSource) ChannelConfig() Channel {
	return bs.StableBundle().ChannelConfig()
}



func (bs *BundleSource) OrdererConfig() (Orderer, bool) {
	return bs.StableBundle().OrdererConfig()
}



func (bs *BundleSource) ConsortiumsConfig() (Consortiums, bool) {
	return bs.StableBundle().ConsortiumsConfig()
}



func (bs *BundleSource) ApplicationConfig() (Application, bool) {
	return bs.StableBundle().ApplicationConfig()
}


func (bs *BundleSource) ConfigtxValidator() configtx.Validator {
	return bs.StableBundle().ConfigtxValidator()
}


func (bs *BundleSource) ValidateNew(resources Resources) error {
	return bs.StableBundle().ValidateNew(resources)
}
