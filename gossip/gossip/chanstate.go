/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"bytes"
	"sync"
	"sync/atomic"

	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/gossip/channel"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
)

type channelState struct {
	stopping int32
	sync.RWMutex
	channels map[string]channel.GossipChannel
	g        *GossipImpl
}

func (cs *channelState) stop() {
	if cs.isStopping() {
		return
	}
	atomic.StoreInt32(&cs.stopping, int32(1))
	cs.Lock()
	defer cs.Unlock()
	for _, gc := range cs.channels {
		gc.Stop()
	}
}

func (cs *channelState) isStopping() bool {
	return atomic.LoadInt32(&cs.stopping) == int32(1)
}

func (cs *channelState) lookupChannelForMsg(msg protoext.ReceivedMessage) channel.GossipChannel {
	if protoext.IsStateInfoPullRequestMsg(msg.GetGossipMessage().GossipMessage) {
		sipr := msg.GetGossipMessage().GetStateInfoPullReq()
		mac := sipr.Channel_MAC
		pkiID := msg.GetConnectionInfo().ID
		return cs.getGossipChannelByMAC(mac, pkiID)
	}
	return cs.lookupChannelForGossipMsg(msg.GetGossipMessage().GossipMessage)
}

func (cs *channelState) lookupChannelForGossipMsg(msg *proto.GossipMessage) channel.GossipChannel {
	if !protoext.IsStateInfoMsg(msg) {
		
		
		
		
		
		
		return cs.getGossipChannelByChainID(msg.Channel)
	}

	
	stateInfMsg := msg.GetStateInfo()
	return cs.getGossipChannelByMAC(stateInfMsg.Channel_MAC, stateInfMsg.PkiId)
}

func (cs *channelState) getGossipChannelByMAC(receivedMAC []byte, pkiID common.PKIidType) channel.GossipChannel {
	
	
	
	
	cs.RLock()
	defer cs.RUnlock()
	for chanName, gc := range cs.channels {
		mac := channel.GenerateMAC(pkiID, common.ChannelID(chanName))
		if bytes.Equal(mac, receivedMAC) {
			return gc
		}
	}
	return nil
}

func (cs *channelState) getGossipChannelByChainID(channelID common.ChannelID) channel.GossipChannel {
	if cs.isStopping() {
		return nil
	}
	cs.RLock()
	defer cs.RUnlock()
	return cs.channels[string(channelID)]
}

func (cs *channelState) joinChannel(joinMsg api.JoinChannelMessage, channelID common.ChannelID,
	metrics *metrics.MembershipMetrics) {
	if cs.isStopping() {
		return
	}
	cs.Lock()
	defer cs.Unlock()
	if gc, exists := cs.channels[string(channelID)]; !exists {
		pkiID := cs.g.comm.GetPKIid()
		ga := &gossipAdapterImpl{GossipImpl: cs.g, Discovery: cs.g.disc}
		gc := channel.NewGossipChannel(pkiID, cs.g.selfOrg, cs.g.mcs, channelID, ga, joinMsg, metrics, nil)
		cs.channels[string(channelID)] = gc
	} else {
		gc.ConfigureChannel(joinMsg)
	}
}

type gossipAdapterImpl struct {
	*GossipImpl
	discovery.Discovery
}

func (ga *gossipAdapterImpl) GetConf() channel.Config {
	return channel.Config{
		ID:                          ga.conf.ID,
		MaxBlockCountToStore:        ga.conf.MaxBlockCountToStore,
		PublishStateInfoInterval:    ga.conf.PublishStateInfoInterval,
		PullInterval:                ga.conf.PullInterval,
		PullPeerNum:                 ga.conf.PullPeerNum,
		RequestStateInfoInterval:    ga.conf.RequestStateInfoInterval,
		BlockExpirationInterval:     ga.conf.PullInterval * 100,
		StateInfoCacheSweepInterval: ga.conf.PullInterval * 5,
		TimeForMembershipTracker:    ga.conf.TimeForMembershipTracker,
		DigestWaitTime:              ga.conf.DigestWaitTime,
		RequestWaitTime:             ga.conf.RequestWaitTime,
		ResponseWaitTime:            ga.conf.ResponseWaitTime,
		MsgExpirationTimeout:        ga.conf.MsgExpirationTimeout,
	}
}

func (ga *gossipAdapterImpl) Sign(msg *proto.GossipMessage) (*protoext.SignedGossipMessage, error) {
	signer := func(msg []byte) ([]byte, error) {
		return ga.mcs.Sign(msg)
	}
	sMsg := &protoext.SignedGossipMessage{
		GossipMessage: msg,
	}
	e, err := sMsg.Sign(signer)
	if err != nil {
		return nil, err
	}
	return &protoext.SignedGossipMessage{
		Envelope:      e,
		GossipMessage: msg,
	}, nil
}


func (ga *gossipAdapterImpl) Gossip(msg *protoext.SignedGossipMessage) {
	ga.GossipImpl.emitter.Add(&emittedGossipMessage{
		SignedGossipMessage: msg,
		filter: func(_ common.PKIidType) bool {
			return true
		},
	})
}


func (ga *gossipAdapterImpl) Forward(msg protoext.ReceivedMessage) {
	ga.GossipImpl.emitter.Add(&emittedGossipMessage{
		SignedGossipMessage: msg.GetGossipMessage(),
		filter:              msg.GetConnectionInfo().ID.IsNotSameFilter,
	})
}

func (ga *gossipAdapterImpl) Send(msg *protoext.SignedGossipMessage, peers ...*comm.RemotePeer) {
	ga.GossipImpl.comm.Send(msg, peers...)
}



func (ga *gossipAdapterImpl) ValidateStateInfoMessage(msg *protoext.SignedGossipMessage) error {
	return ga.GossipImpl.validateStateInfoMsg(msg)
}


func (ga *gossipAdapterImpl) GetOrgOfPeer(PKIID common.PKIidType) api.OrgIdentityType {
	return ga.GossipImpl.getOrgOfPeer(PKIID)
}



func (ga *gossipAdapterImpl) GetIdentityByPKIID(pkiID common.PKIidType) api.PeerIdentityType {
	identity, err := ga.idMapper.Get(pkiID)
	if err != nil {
		return nil
	}
	return identity
}
