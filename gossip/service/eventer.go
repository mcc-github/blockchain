/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package service

import (
	"reflect"

	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/channelconfig"
)


type Config interface {
	
	ChannelID() string

	
	Organizations() map[string]channelconfig.ApplicationOrg

	
	Sequence() uint64

	
	OrdererAddresses() []string
}


type ConfigProcessor interface {
	
	ProcessConfigUpdate(config Config)
}

type configStore struct {
	anchorPeers []*peer.AnchorPeer
	orgMap      map[string]channelconfig.ApplicationOrg
}

type configEventReceiver interface {
	updateAnchors(config Config)
	updateEndpoints(chainID string, endpoints []string)
}

type configEventer struct {
	lastConfig *configStore
	receiver   configEventReceiver
}

func newConfigEventer(receiver configEventReceiver) *configEventer {
	return &configEventer{
		receiver: receiver,
	}
}





func (ce *configEventer) ProcessConfigUpdate(config Config) {
	logger.Debugf("Processing new config for channel %s", config.ChannelID())
	orgMap := cloneOrgConfig(config.Organizations())
	if ce.lastConfig != nil && reflect.DeepEqual(ce.lastConfig.orgMap, orgMap) {
		logger.Debugf("Ignoring new config for channel %s because it contained no anchor peer updates", config.ChannelID())
	} else {

		var newAnchorPeers []*peer.AnchorPeer
		for _, group := range config.Organizations() {
			newAnchorPeers = append(newAnchorPeers, group.AnchorPeers()...)
		}

		newConfig := &configStore{
			orgMap:      orgMap,
			anchorPeers: newAnchorPeers,
		}
		ce.lastConfig = newConfig

		logger.Debugf("Calling out because config was updated for channel %s", config.ChannelID())
		ce.receiver.updateAnchors(config)
	}
	ce.receiver.updateEndpoints(config.ChannelID(), config.OrdererAddresses())
}

func cloneOrgConfig(src map[string]channelconfig.ApplicationOrg) map[string]channelconfig.ApplicationOrg {
	clone := make(map[string]channelconfig.ApplicationOrg)
	for k, v := range src {
		clone[k] = &appGrp{
			name:        v.Name(),
			mspID:       v.MSPID(),
			anchorPeers: v.AnchorPeers(),
		}
	}
	return clone
}

type appGrp struct {
	name        string
	mspID       string
	anchorPeers []*peer.AnchorPeer
}

func (ag *appGrp) Name() string {
	return ag.name
}

func (ag *appGrp) MSPID() string {
	return ag.mspID
}

func (ag *appGrp) AnchorPeers() []*peer.AnchorPeer {
	return ag.anchorPeers
}
