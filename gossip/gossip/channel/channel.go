/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	common_utils "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/filter"
	"github.com/mcc-github/blockchain/gossip/gossip/algo"
	"github.com/mcc-github/blockchain/gossip/gossip/msgstore"
	"github.com/mcc-github/blockchain/gossip/gossip/pull"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)

const DefMsgExpirationTimeout = election.DefLeaderAliveThreshold * 10



type Config struct {
	ID                          string
	PublishStateInfoInterval    time.Duration
	MaxBlockCountToStore        int
	PullPeerNum                 int
	PullInterval                time.Duration
	RequestStateInfoInterval    time.Duration
	BlockExpirationInterval     time.Duration
	StateInfoCacheSweepInterval time.Duration
	TimeForMembershipTracker    time.Duration
	DigestWaitTime              time.Duration
	RequestWaitTime             time.Duration
	ResponseWaitTime            time.Duration
	MsgExpirationTimeout        time.Duration
}


type GossipChannel interface {
	
	Self() *protoext.SignedGossipMessage

	
	GetPeers() []discovery.NetworkMember

	
	
	PeerFilter(api.SubChannelSelectionCriteria) filter.RoutingFilter

	
	IsMemberInChan(member discovery.NetworkMember) bool

	
	
	UpdateLedgerHeight(height uint64)

	
	
	UpdateChaincodes(chaincode []*proto.Chaincode)

	
	IsOrgInChannel(membersOrg api.OrgIdentityType) bool

	
	
	EligibleForChannel(member discovery.NetworkMember) bool

	
	HandleMessage(protoext.ReceivedMessage)

	
	AddToMsgStore(msg *protoext.SignedGossipMessage)

	
	
	ConfigureChannel(joinMsg api.JoinChannelMessage)

	
	LeaveChannel()

	
	Stop()
}



type Adapter interface {
	Sign(msg *proto.GossipMessage) (*protoext.SignedGossipMessage, error)

	
	GetConf() Config

	
	Gossip(message *protoext.SignedGossipMessage)

	
	Forward(message protoext.ReceivedMessage)

	
	DeMultiplex(interface{})

	
	GetMembership() []discovery.NetworkMember

	
	Lookup(PKIID common.PKIidType) *discovery.NetworkMember

	
	Send(msg *protoext.SignedGossipMessage, peers ...*comm.RemotePeer)

	
	
	ValidateStateInfoMessage(message *protoext.SignedGossipMessage) error

	
	GetOrgOfPeer(pkiID common.PKIidType) api.OrgIdentityType

	
	
	GetIdentityByPKIID(pkiID common.PKIidType) api.PeerIdentityType
}

type gossipChannel struct {
	Adapter
	sync.RWMutex
	shouldGossipStateInfo     int32
	mcs                       api.MessageCryptoService
	pkiID                     common.PKIidType
	selfOrg                   api.OrgIdentityType
	stopChan                  chan struct{}
	stateInfoMsg              *protoext.SignedGossipMessage
	orgs                      []api.OrgIdentityType
	joinMsg                   api.JoinChannelMessage
	blockMsgStore             msgstore.MessageStore
	stateInfoMsgStore         *stateInfoCache
	leaderMsgStore            msgstore.MessageStore
	chainID                   common.ChainID
	blocksPuller              pull.Mediator
	logger                    util.Logger
	stateInfoPublishScheduler *time.Ticker
	stateInfoRequestScheduler *time.Ticker
	memFilter                 *membershipFilter
	ledgerHeight              uint64
	incTime                   uint64
	leftChannel               int32
	membershipTracker         *membershipTracker
}

type membershipFilter struct {
	adapter Adapter
	*gossipChannel
}


func (mf *membershipFilter) GetMembership() []discovery.NetworkMember {
	if mf.hasLeftChannel() {
		return nil
	}

	var members []discovery.NetworkMember
	for _, mem := range mf.adapter.GetMembership() {
		if mf.eligibleForChannelAndSameOrg(mem) {
			members = append(members, mem)
		}
	}
	return members
}


func NewGossipChannel(pkiID common.PKIidType, org api.OrgIdentityType, mcs api.MessageCryptoService,
	chainID common.ChainID, adapter Adapter, joinMsg api.JoinChannelMessage,
	metrics *metrics.MembershipMetrics, logger util.Logger) GossipChannel {
	gc := &gossipChannel{
		incTime:                   uint64(time.Now().UnixNano()),
		selfOrg:                   org,
		pkiID:                     pkiID,
		mcs:                       mcs,
		Adapter:                   adapter,
		stopChan:                  make(chan struct{}, 1),
		shouldGossipStateInfo:     int32(0),
		stateInfoPublishScheduler: time.NewTicker(adapter.GetConf().PublishStateInfoInterval),
		stateInfoRequestScheduler: time.NewTicker(adapter.GetConf().RequestStateInfoInterval),
		orgs:                      []api.OrgIdentityType{},
		chainID:                   chainID,
	}

	if logger == nil {
		gc.logger = util.GetLogger(util.ChannelLogger, adapter.GetConf().ID)
	} else {
		gc.logger = logger
	}

	gc.memFilter = &membershipFilter{adapter: gc.Adapter, gossipChannel: gc}

	comparator := protoext.NewGossipMessageComparator(adapter.GetConf().MaxBlockCountToStore)

	gc.blocksPuller = gc.createBlockPuller()

	seqNumFromMsg := func(m interface{}) string {
		return fmt.Sprintf("%d", m.(*protoext.SignedGossipMessage).GetDataMsg().Payload.SeqNum)
	}
	gc.blockMsgStore = msgstore.NewMessageStoreExpirable(comparator, func(m interface{}) {
		gc.logger.Debugf("Removing %s from the message store", seqNumFromMsg(m))
		gc.blocksPuller.Remove(seqNumFromMsg(m))
	}, gc.GetConf().BlockExpirationInterval, nil, nil, func(m interface{}) {
		gc.logger.Debugf("Removing %s from the message store", seqNumFromMsg(m))
		gc.blocksPuller.Remove(seqNumFromMsg(m))
	})

	hashPeerExpiredInMembership := func(o interface{}) bool {
		pkiID := o.(*protoext.SignedGossipMessage).GetStateInfo().PkiId
		return gc.Lookup(pkiID) == nil
	}
	verifyStateInfoMsg := func(msg *protoext.SignedGossipMessage, orgs ...api.OrgIdentityType) bool {
		si := msg.GetStateInfo()
		
		if bytes.Equal(gc.pkiID, si.PkiId) {
			return true
		}
		peerIdentity := adapter.GetIdentityByPKIID(si.PkiId)
		if len(peerIdentity) == 0 {
			gc.logger.Warning("Identity for peer", si.PkiId, "doesn't exist")
			return false
		}
		isOrgInChan := func(org api.OrgIdentityType) bool {
			if len(orgs) == 0 {
				if !gc.IsOrgInChannel(org) {
					return false
				}
			} else {
				found := false
				for _, chanMember := range orgs {
					if bytes.Equal(chanMember, org) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			return true
		}

		org := gc.GetOrgOfPeer(si.PkiId)
		if !isOrgInChan(org) {
			gc.logger.Warning("peer", peerIdentity, "'s organization(", string(org), ") isn't in the channel", string(chainID))
			return false
		}
		if err := gc.mcs.VerifyByChannel(chainID, peerIdentity, msg.Signature, msg.Payload); err != nil {
			gc.logger.Warningf("Peer %v isn't eligible for channel %s : %+v", peerIdentity, string(chainID), err)
			return false
		}
		return true
	}
	gc.stateInfoMsgStore = newStateInfoCache(gc.GetConf().StateInfoCacheSweepInterval, hashPeerExpiredInMembership, verifyStateInfoMsg)

	ttl := adapter.GetConf().MsgExpirationTimeout
	pol := protoext.NewGossipMessageComparator(0)

	gc.leaderMsgStore = msgstore.NewMessageStoreExpirable(pol, msgstore.Noop, ttl, nil, nil, nil)

	gc.ConfigureChannel(joinMsg)

	
	go gc.periodicalInvocation(gc.publishStateInfo, gc.stateInfoPublishScheduler.C)
	
	go gc.periodicalInvocation(gc.requestStateInfo, gc.stateInfoRequestScheduler.C)

	ticker := time.NewTicker(gc.GetConf().TimeForMembershipTracker)
	gc.membershipTracker = &membershipTracker{
		getPeersToTrack: gc.GetPeers,
		report:          gc.reportMembershipChanges,
		stopChan:        make(chan struct{}, 1),
		tickerChannel:   ticker.C,
		metrics:         metrics,
		chainID:         chainID,
	}

	go gc.membershipTracker.trackMembershipChanges()
	return gc
}

func (gc *gossipChannel) reportMembershipChanges(input ...interface{}) {
	gc.logger.Info(input...)
}


func (gc *gossipChannel) Stop() {
	close(gc.stopChan)
	close(gc.membershipTracker.stopChan)
	gc.blocksPuller.Stop()
	gc.stateInfoPublishScheduler.Stop()
	gc.stateInfoRequestScheduler.Stop()
	gc.leaderMsgStore.Stop()
	gc.stateInfoMsgStore.Stop()
	gc.blockMsgStore.Stop()
}

func (gc *gossipChannel) periodicalInvocation(fn func(), c <-chan time.Time) {
	for {
		select {
		case <-c:
			fn()
		case <-gc.stopChan:
			return
		}
	}
}


func (gc *gossipChannel) Self() *protoext.SignedGossipMessage {
	gc.RLock()
	defer gc.RUnlock()
	return gc.stateInfoMsg
}


func (gc *gossipChannel) LeaveChannel() {
	gc.Lock()
	defer gc.Unlock()

	atomic.StoreInt32(&gc.leftChannel, 1)

	var chaincodes []*proto.Chaincode
	var height uint64
	if prevMsg := gc.stateInfoMsg; prevMsg != nil {
		chaincodes = prevMsg.GetStateInfo().Properties.Chaincodes
		height = prevMsg.GetStateInfo().Properties.LedgerHeight
	}
	gc.updateProperties(height, chaincodes, true)
}

func (gc *gossipChannel) hasLeftChannel() bool {
	return atomic.LoadInt32(&gc.leftChannel) == 1
}


func (gc *gossipChannel) GetPeers() []discovery.NetworkMember {
	var members []discovery.NetworkMember
	if gc.hasLeftChannel() {
		return members
	}

	for _, member := range gc.GetMembership() {
		if !gc.EligibleForChannel(member) {
			continue
		}
		stateInf := gc.stateInfoMsgStore.MsgByID(member.PKIid)
		if stateInf == nil {
			continue
		}
		props := stateInf.GetStateInfo().Properties
		if props != nil && props.LeftChannel {
			continue
		}
		member.Properties = stateInf.GetStateInfo().Properties
		member.Envelope = stateInf.Envelope
		members = append(members, member)
	}
	return members
}

func (gc *gossipChannel) requestStateInfo() {
	req, err := gc.createStateInfoRequest()
	if err != nil {
		gc.logger.Warningf("Failed creating SignedGossipMessage: %+v", err)
		return
	}
	endpoints := filter.SelectPeers(gc.GetConf().PullPeerNum, gc.GetMembership(), gc.IsMemberInChan)
	gc.Send(req, endpoints...)
}

func (gc *gossipChannel) eligibleForChannelAndSameOrg(member discovery.NetworkMember) bool {
	sameOrg := func(networkMember discovery.NetworkMember) bool {
		return bytes.Equal(gc.GetOrgOfPeer(networkMember.PKIid), gc.selfOrg)
	}
	return filter.CombineRoutingFilters(gc.EligibleForChannel, sameOrg)(member)
}

func (gc *gossipChannel) publishStateInfo() {
	if atomic.LoadInt32(&gc.shouldGossipStateInfo) == int32(0) {
		return
	}
	gc.RLock()
	stateInfoMsg := gc.stateInfoMsg
	gc.RUnlock()
	gc.Gossip(stateInfoMsg)
	if len(gc.GetMembership()) > 0 {
		atomic.StoreInt32(&gc.shouldGossipStateInfo, int32(0))
	}
}

func (gc *gossipChannel) createBlockPuller() pull.Mediator {
	conf := pull.Config{
		MsgType:           proto.PullMsgType_BLOCK_MSG,
		Channel:           []byte(gc.chainID),
		ID:                gc.GetConf().ID,
		PeerCountToSelect: gc.GetConf().PullPeerNum,
		PullInterval:      gc.GetConf().PullInterval,
		Tag:               proto.GossipMessage_CHAN_AND_ORG,
		PullEngineConfig: algo.PullEngineConfig{
			DigestWaitTime:   gc.GetConf().DigestWaitTime,
			RequestWaitTime:  gc.GetConf().RequestWaitTime,
			ResponseWaitTime: gc.GetConf().ResponseWaitTime,
		},
	}
	seqNumFromMsg := func(msg *protoext.SignedGossipMessage) string {
		dataMsg := msg.GetDataMsg()
		if dataMsg == nil || dataMsg.Payload == nil {
			gc.logger.Warning("Non-data block or with no payload")
			return ""
		}
		return fmt.Sprintf("%d", dataMsg.Payload.SeqNum)
	}
	adapter := &pull.PullAdapter{
		Sndr:        gc,
		MemSvc:      gc.memFilter,
		IdExtractor: seqNumFromMsg,
		MsgCons: func(msg *protoext.SignedGossipMessage) {
			gc.DeMultiplex(msg)
		},
	}

	adapter.IngressDigFilter = func(digestMsg *proto.DataDigest) *proto.DataDigest {
		gc.RLock()
		height := gc.ledgerHeight
		gc.RUnlock()
		digests := digestMsg.Digests
		digestMsg.Digests = nil
		for i := range digests {
			seqNum, err := strconv.ParseUint(string(digests[i]), 10, 64)
			if err != nil {
				gc.logger.Warningf("Can't parse digest %s : %+v", digests[i], err)
				continue
			}
			if seqNum >= height {
				digestMsg.Digests = append(digestMsg.Digests, digests[i])
			}

		}
		return digestMsg
	}

	return pull.NewPullMediator(conf, adapter)
}


func (gc *gossipChannel) IsMemberInChan(member discovery.NetworkMember) bool {
	org := gc.GetOrgOfPeer(member.PKIid)
	if org == nil {
		return false
	}

	return gc.IsOrgInChannel(org)
}



func (gc *gossipChannel) PeerFilter(messagePredicate api.SubChannelSelectionCriteria) filter.RoutingFilter {
	return func(member discovery.NetworkMember) bool {
		peerIdentity := gc.GetIdentityByPKIID(member.PKIid)
		if len(peerIdentity) == 0 {
			return false
		}
		msg := gc.stateInfoMsgStore.MembershipStore.MsgByID(member.PKIid)
		if msg == nil {
			return false
		}

		return messagePredicate(api.PeerSignature{
			Message:      msg.Payload,
			Signature:    msg.Signature,
			PeerIdentity: peerIdentity,
		})
	}
}


func (gc *gossipChannel) IsOrgInChannel(membersOrg api.OrgIdentityType) bool {
	gc.RLock()
	defer gc.RUnlock()
	for _, orgOfChan := range gc.orgs {
		if bytes.Equal(orgOfChan, membersOrg) {
			return true
		}
	}
	return false
}



func (gc *gossipChannel) EligibleForChannel(member discovery.NetworkMember) bool {
	peerIdentity := gc.GetIdentityByPKIID(member.PKIid)
	if len(peerIdentity) == 0 {
		gc.logger.Warning("Identity for peer", member.PKIid, "doesn't exist")
		return false
	}
	msg := gc.stateInfoMsgStore.MsgByID(member.PKIid)
	if msg == nil {
		return false
	}
	return true
}


func (gc *gossipChannel) AddToMsgStore(msg *protoext.SignedGossipMessage) {
	if protoext.IsDataMsg(msg.GossipMessage) {
		gc.Lock()
		defer gc.Unlock()
		added := gc.blockMsgStore.Add(msg)
		if added {
			gc.logger.Debugf("Adding %v to the block puller", msg)
			gc.blocksPuller.Add(msg)
		}
	}

	if protoext.IsStateInfoMsg(msg.GossipMessage) {
		gc.stateInfoMsgStore.Add(msg)
	}
}



func (gc *gossipChannel) ConfigureChannel(joinMsg api.JoinChannelMessage) {
	gc.Lock()
	defer gc.Unlock()

	if len(joinMsg.Members()) == 0 {
		gc.logger.Warning("Received join channel message with empty set of members")
		return
	}

	if gc.joinMsg == nil {
		gc.joinMsg = joinMsg
	}

	if gc.joinMsg.SequenceNumber() > (joinMsg.SequenceNumber()) {
		gc.logger.Warning("Already have a more updated JoinChannel message(", gc.joinMsg.SequenceNumber(), ") than", joinMsg.SequenceNumber())
		return
	}

	gc.orgs = joinMsg.Members()
	gc.joinMsg = joinMsg
	gc.stateInfoMsgStore.validate(joinMsg.Members())
}


func (gc *gossipChannel) HandleMessage(msg protoext.ReceivedMessage) {
	if !gc.verifyMsg(msg) {
		gc.logger.Warning("Failed verifying message:", msg.GetGossipMessage().GossipMessage)
		return
	}
	m := msg.GetGossipMessage()
	if !protoext.IsChannelRestricted(m.GossipMessage) {
		gc.logger.Warning("Got message", msg.GetGossipMessage(), "but it's not a per-channel message, discarding it")
		return
	}
	orgID := gc.GetOrgOfPeer(msg.GetConnectionInfo().ID)
	if len(orgID) == 0 {
		gc.logger.Debug("Couldn't find org identity of peer", msg.GetConnectionInfo())
		return
	}
	if !gc.IsOrgInChannel(orgID) {
		gc.logger.Warning("Point to point message came from", msg.GetConnectionInfo(),
			", org(", string(orgID), ") but it's not eligible for the channel", string(gc.chainID))
		return
	}

	if protoext.IsStateInfoPullRequestMsg(m.GossipMessage) {
		msg.Respond(gc.createStateInfoSnapshot(orgID))
		return
	}

	if protoext.IsStateInfoSnapshot(m.GossipMessage) {
		gc.handleStateInfSnapshot(m.GossipMessage, msg.GetConnectionInfo().ID)
		return
	}

	if protoext.IsDataMsg(m.GossipMessage) || protoext.IsStateInfoMsg(m.GossipMessage) {
		added := false

		if protoext.IsDataMsg(m.GossipMessage) {
			if m.GetDataMsg().Payload == nil {
				gc.logger.Warning("Payload is empty, got it from", msg.GetConnectionInfo().ID)
				return
			}
			
			if !gc.blockMsgStore.CheckValid(msg.GetGossipMessage()) {
				return
			}
			if !gc.verifyBlock(m.GossipMessage, msg.GetConnectionInfo().ID) {
				gc.logger.Warning("Failed verifying block", m.GetDataMsg().Payload.SeqNum)
				return
			}
			gc.Lock()
			added = gc.blockMsgStore.Add(msg.GetGossipMessage())
			if added {
				gc.logger.Debugf("Adding %v to the block puller", msg.GetGossipMessage())
				gc.blocksPuller.Add(msg.GetGossipMessage())
			}
			gc.Unlock()
		} else { 
			
			added = gc.stateInfoMsgStore.Add(msg.GetGossipMessage())
		}

		if added {
			
			gc.Forward(msg)
			
			gc.DeMultiplex(m)
		}
		return
	}

	if protoext.IsPullMsg(m.GossipMessage) && protoext.GetPullMsgType(m.GossipMessage) == proto.PullMsgType_BLOCK_MSG {
		if gc.hasLeftChannel() {
			gc.logger.Info("Received Pull message from", msg.GetConnectionInfo().Endpoint, "but left the channel", string(gc.chainID))
			return
		}
		
		
		if gc.stateInfoMsgStore.MsgByID(msg.GetConnectionInfo().ID) == nil {
			gc.logger.Debug("Don't have StateInfo message of peer", msg.GetConnectionInfo())
			return
		}
		if !gc.eligibleForChannelAndSameOrg(discovery.NetworkMember{PKIid: msg.GetConnectionInfo().ID}) {
			gc.logger.Warning(msg.GetConnectionInfo(), "isn't eligible for pulling blocks of", string(gc.chainID))
			return
		}
		if protoext.IsDataUpdate(m.GossipMessage) {
			
			
			
			var msgs []*protoext.SignedGossipMessage
			var items []*proto.Envelope
			filteredEnvelopes := []*proto.Envelope{}
			for _, item := range m.GetDataUpdate().Data {
				gMsg, err := protoext.EnvelopeToGossipMessage(item)
				if err != nil {
					gc.logger.Warningf("Data update contains an invalid message: %+v", err)
					return
				}
				if !bytes.Equal(gMsg.Channel, []byte(gc.chainID)) {
					gc.logger.Warning("DataUpdate message contains item with channel", gMsg.Channel, "but should be", gc.chainID)
					return
				}
				
				if !gc.blockMsgStore.CheckValid(gMsg) {
					continue
				}
				if !gc.verifyBlock(gMsg.GossipMessage, msg.GetConnectionInfo().ID) {
					return
				}
				msgs = append(msgs, gMsg)
				items = append(items, item)
			}

			gc.Lock()
			defer gc.Unlock()

			for i, gMsg := range msgs {
				item := items[i]
				added := gc.blockMsgStore.Add(gMsg)
				if !added {
					
					
					continue
				}
				filteredEnvelopes = append(filteredEnvelopes, item)
			}

			
			m.GetDataUpdate().Data = filteredEnvelopes
		}
		gc.blocksPuller.HandleMessage(msg)
	}

	if protoext.IsLeadershipMsg(m.GossipMessage) {
		
		added := gc.leaderMsgStore.Add(m)
		if added {
			gc.DeMultiplex(m)
		}
	}
}

func (gc *gossipChannel) handleStateInfSnapshot(m *proto.GossipMessage, sender common.PKIidType) {
	chanName := string(gc.chainID)
	for _, envelope := range m.GetStateSnapshot().Elements {
		stateInf, err := protoext.EnvelopeToGossipMessage(envelope)
		if err != nil {
			gc.logger.Warningf("Channel %s : StateInfo snapshot contains an invalid message: %+v", chanName, err)
			return
		}
		if !protoext.IsStateInfoMsg(stateInf.GossipMessage) {
			gc.logger.Warning("Channel", chanName, ": Element of StateInfoSnapshot isn't a StateInfoMessage:",
				stateInf, "message sent from", sender)
			return
		}
		si := stateInf.GetStateInfo()
		orgID := gc.GetOrgOfPeer(si.PkiId)
		if orgID == nil {
			gc.logger.Debug("Channel", chanName, ": Couldn't find org identity of peer",
				string(si.PkiId), "message sent from", string(sender))
			return
		}

		if !gc.IsOrgInChannel(orgID) {
			gc.logger.Warning("Channel", chanName, ": Peer", stateInf.GetStateInfo().PkiId,
				"is not in an eligible org, can't process a stateInfo from it, sent from", sender)
			return
		}

		expectedMAC := GenerateMAC(si.PkiId, gc.chainID)
		if !bytes.Equal(si.Channel_MAC, expectedMAC) {
			gc.logger.Warning("Channel", chanName, ": StateInfo message", stateInf,
				", has an invalid MAC. Expected", expectedMAC, ", got", si.Channel_MAC, ", sent from", sender)
			return
		}
		err = gc.ValidateStateInfoMessage(stateInf)
		if err != nil {
			gc.logger.Warningf("Channel %s: Failed validating state info message: %v sent from %v : %+v", chanName, stateInf, sender, err)
			return
		}

		if gc.Lookup(si.PkiId) == nil {
			
			
			continue
		}

		gc.stateInfoMsgStore.Add(stateInf)
	}
}

func (gc *gossipChannel) verifyBlock(msg *proto.GossipMessage, sender common.PKIidType) bool {
	if !protoext.IsDataMsg(msg) {
		gc.logger.Warning("Received from ", sender, "a DataUpdate message that contains a non-block GossipMessage:", msg)
		return false
	}
	payload := msg.GetDataMsg().Payload
	if payload == nil {
		gc.logger.Warning("Received empty payload from", sender)
		return false
	}
	seqNum := payload.SeqNum
	rawBlock := payload.Data
	err := gc.mcs.VerifyBlock(msg.Channel, seqNum, rawBlock)
	if err != nil {
		gc.logger.Warningf("Received blockchainated block from %v in DataUpdate: %+v", sender, err)
		return false
	}
	return true
}

func (gc *gossipChannel) createStateInfoSnapshot(requestersOrg api.OrgIdentityType) *proto.GossipMessage {
	sameOrg := bytes.Equal(gc.selfOrg, requestersOrg)
	rawElements := gc.stateInfoMsgStore.Get()
	elements := []*proto.Envelope{}
	for _, rawEl := range rawElements {
		msg := rawEl.(*protoext.SignedGossipMessage)
		orgOfCurrentMsg := gc.GetOrgOfPeer(msg.GetStateInfo().PkiId)
		
		
		if sameOrg || !bytes.Equal(orgOfCurrentMsg, gc.selfOrg) {
			elements = append(elements, msg.Envelope)
			continue
		}
		
		
		if netMember := gc.Lookup(msg.GetStateInfo().PkiId); netMember == nil || netMember.Endpoint == "" {
			continue
		}
		elements = append(elements, msg.Envelope)
	}

	return &proto.GossipMessage{
		Channel: gc.chainID,
		Tag:     proto.GossipMessage_CHAN_OR_ORG,
		Nonce:   0,
		Content: &proto.GossipMessage_StateSnapshot{
			StateSnapshot: &proto.StateInfoSnapshot{
				Elements: elements,
			},
		},
	}
}

func (gc *gossipChannel) verifyMsg(msg protoext.ReceivedMessage) bool {
	if msg == nil {
		gc.logger.Warning("Messsage is nil")
		return false
	}
	m := msg.GetGossipMessage()
	if m == nil {
		gc.logger.Warning("Message content is empty")
		return false
	}

	if msg.GetConnectionInfo().ID == nil {
		gc.logger.Warning("Message has nil PKI-ID")
		return false
	}

	if protoext.IsStateInfoMsg(m.GossipMessage) {
		si := m.GetStateInfo()
		expectedMAC := GenerateMAC(si.PkiId, gc.chainID)
		if !bytes.Equal(expectedMAC, si.Channel_MAC) {
			gc.logger.Warning("Message contains wrong channel MAC(", si.Channel_MAC, "), expected", expectedMAC)
			return false
		}
		return true
	}

	if protoext.IsStateInfoPullRequestMsg(m.GossipMessage) {
		sipr := m.GetStateInfoPullReq()
		expectedMAC := GenerateMAC(msg.GetConnectionInfo().ID, gc.chainID)
		if !bytes.Equal(expectedMAC, sipr.Channel_MAC) {
			gc.logger.Warning("Message contains wrong channel MAC(", sipr.Channel_MAC, "), expected", expectedMAC)
			return false
		}
		return true
	}

	if !bytes.Equal(m.Channel, []byte(gc.chainID)) {
		gc.logger.Warning("Message contains wrong channel(", m.Channel, "), expected", gc.chainID)
		return false
	}
	return true
}

func (gc *gossipChannel) createStateInfoRequest() (*protoext.SignedGossipMessage, error) {
	return protoext.NoopSign(&proto.GossipMessage{
		Tag:   proto.GossipMessage_CHAN_OR_ORG,
		Nonce: 0,
		Content: &proto.GossipMessage_StateInfoPullReq{
			StateInfoPullReq: &proto.StateInfoPullRequest{
				Channel_MAC: GenerateMAC(gc.pkiID, gc.chainID),
			},
		},
	})
}



func (gc *gossipChannel) UpdateLedgerHeight(height uint64) {
	gc.Lock()
	defer gc.Unlock()

	var chaincodes []*proto.Chaincode
	var leftChannel bool
	if prevMsg := gc.stateInfoMsg; prevMsg != nil {
		leftChannel = prevMsg.GetStateInfo().Properties.LeftChannel
		chaincodes = prevMsg.GetStateInfo().Properties.Chaincodes
	}
	gc.updateProperties(height, chaincodes, leftChannel)
}



func (gc *gossipChannel) UpdateChaincodes(chaincodes []*proto.Chaincode) {
	gc.Lock()
	defer gc.Unlock()

	var ledgerHeight uint64 = 1
	var leftChannel bool
	if prevMsg := gc.stateInfoMsg; prevMsg != nil {
		ledgerHeight = prevMsg.GetStateInfo().Properties.LedgerHeight
		leftChannel = prevMsg.GetStateInfo().Properties.LeftChannel
	}
	gc.updateProperties(ledgerHeight, chaincodes, leftChannel)
}



func (gc *gossipChannel) updateStateInfo(msg *protoext.SignedGossipMessage) {
	gc.stateInfoMsgStore.Add(msg)
	gc.ledgerHeight = msg.GetStateInfo().Properties.LedgerHeight
	gc.stateInfoMsg = msg
	atomic.StoreInt32(&gc.shouldGossipStateInfo, int32(1))
}

func (gc *gossipChannel) updateProperties(ledgerHeight uint64, chaincodes []*proto.Chaincode, leftChannel bool) {
	stateInfMsg := &proto.StateInfo{
		Channel_MAC: GenerateMAC(gc.pkiID, gc.chainID),
		PkiId:       gc.pkiID,
		Timestamp: &proto.PeerTime{
			IncNum: gc.incTime,
			SeqNum: uint64(time.Now().UnixNano()),
		},
		Properties: &proto.Properties{
			LeftChannel:  leftChannel,
			LedgerHeight: ledgerHeight,
			Chaincodes:   chaincodes,
		},
	}
	m := &proto.GossipMessage{
		Nonce: 0,
		Tag:   proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfo{
			StateInfo: stateInfMsg,
		},
	}

	msg, err := gc.Sign(m)
	if err != nil {
		gc.logger.Error("Failed signing message:", err)
		return
	}
	gc.updateStateInfo(msg)
}

func newStateInfoCache(sweepInterval time.Duration, hasExpired func(interface{}) bool, verifyFunc membershipPredicate) *stateInfoCache {
	membershipStore := util.NewMembershipStore()
	pol := protoext.NewGossipMessageComparator(0)

	s := &stateInfoCache{
		verify:          verifyFunc,
		MembershipStore: membershipStore,
		stopChan:        make(chan struct{}),
	}
	invalidationTrigger := func(m interface{}) {
		pkiID := m.(*protoext.SignedGossipMessage).GetStateInfo().PkiId
		membershipStore.Remove(pkiID)
	}
	s.MessageStore = msgstore.NewMessageStore(pol, invalidationTrigger)

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			case <-time.After(sweepInterval):
				s.Purge(hasExpired)
			}
		}
	}()
	return s
}




type membershipPredicate func(msg *protoext.SignedGossipMessage, orgs ...api.OrgIdentityType) bool




type stateInfoCache struct {
	verify membershipPredicate
	*util.MembershipStore
	msgstore.MessageStore
	stopChan chan struct{}
}

func (cache *stateInfoCache) validate(orgs []api.OrgIdentityType) {
	for _, m := range cache.Get() {
		msg := m.(*protoext.SignedGossipMessage)
		if !cache.verify(msg, orgs...) {
			cache.delete(msg)
		}
	}
}




func (cache *stateInfoCache) Add(msg *protoext.SignedGossipMessage) bool {
	if !cache.MessageStore.CheckValid(msg) {
		return false
	}
	if !cache.verify(msg) {
		return false
	}
	added := cache.MessageStore.Add(msg)
	if added {
		pkiID := msg.GetStateInfo().PkiId
		cache.MembershipStore.Put(pkiID, msg)
	}
	return added
}

func (cache *stateInfoCache) delete(msg *protoext.SignedGossipMessage) {
	cache.Purge(func(o interface{}) bool {
		pkiID := o.(*protoext.SignedGossipMessage).GetStateInfo().PkiId
		return bytes.Equal(pkiID, msg.GetStateInfo().PkiId)
	})
	cache.Remove(msg.GetStateInfo().PkiId)
}

func (cache *stateInfoCache) Stop() {
	cache.stopChan <- struct{}{}
}



func GenerateMAC(pkiID common.PKIidType, channelID common.ChainID) []byte {
	
	var preImage []byte
	preImage = append(preImage, []byte(pkiID)...)
	preImage = append(preImage, []byte(channelID)...)
	return common_utils.ComputeSHA256(preImage)
}


type membershipTracker struct {
	getPeersToTrack func() []discovery.NetworkMember
	report          func(...interface{})
	stopChan        chan struct{}
	tickerChannel   <-chan time.Time
	metrics         *metrics.MembershipMetrics
	chainID         common.ChainID
}


func endpoints(members discovery.Members) [][]string {
	var currView [][]string
	for _, member := range members {
		ep := member.Endpoint
		epi := member.InternalEndpoint
		var endPoints []string
		if ep != epi {
			endPoints = append(endPoints, ep, epi)
		} else {
			endPoints = append(endPoints, ep)
		}
		currView = append(currView, endPoints)
	}
	return currView
}


func (mt *membershipTracker) checkIfPeersChanged(prevPeers discovery.Members, currPeers discovery.Members,
	prevSetPeers map[string]struct{}, currSetPeers map[string]struct{}) {
	var currView [][]string

	wereInPrev := endpoints(prevPeers.Filter(func(member discovery.NetworkMember) bool {
		_, exists := currSetPeers[string(member.PKIid)]
		return !exists
	}))
	newInCurr := endpoints(currPeers.Filter(func(member discovery.NetworkMember) bool {
		_, exists := prevSetPeers[string(member.PKIid)]
		return !exists
	}))
	currView = endpoints(currPeers)

	if !reflect.DeepEqual(wereInPrev, newInCurr) {
		if len(wereInPrev) == 0 {
			mt.report("Membership view has changed. peers went online: ", newInCurr, ", current view: ", currView)
		} else if len(newInCurr) == 0 {
			mt.report("Membership view has changed. peers went offline: ", wereInPrev, ", current view: ", currView)
		} else {
			mt.report("Membership view has changed. peers went offline: ", wereInPrev, ", peers went online: ", newInCurr, ", current view: ", currView)
		}
	}
}

func (mt *membershipTracker) createSetOfPeers(peersToMakeSet []discovery.NetworkMember) map[string]struct{} {
	setPeers := make(map[string]struct{})
	for _, prevPeer := range peersToMakeSet {
		prevPeerID := string(prevPeer.PKIid)
		setPeers[prevPeerID] = struct{}{}
	}
	return setPeers
}

func (mt *membershipTracker) trackMembershipChanges() {
	prev := mt.getPeersToTrack()
	prevSetPeers := mt.createSetOfPeers(prev)
	for {
		
		select {
		case <-mt.stopChan:
			return
		case <-mt.tickerChannel:
			currPeers := mt.getPeersToTrack()
			mt.metrics.Total.With("channel", string(mt.chainID)).Set(float64(len(currPeers)))
			currSetPeers := mt.createSetOfPeers(currPeers)
			mt.checkIfPeersChanged(prev, currPeers, prevSetPeers, currSetPeers)
			prev = currPeers
			prevSetPeers = mt.createSetOfPeers(prev)
		}
	}
}
