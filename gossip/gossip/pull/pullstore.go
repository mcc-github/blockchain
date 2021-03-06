/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pull

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/mcc-github/blockchain-protos-go/gossip"
	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/gossip/algo"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)


const (
	HelloMsgType MsgType = iota
	DigestMsgType
	RequestMsgType
	ResponseMsgType
)


type MsgType int


type MessageHook func(itemIDs []string, items []*protoext.SignedGossipMessage, msg protoext.ReceivedMessage)


type Sender interface {
	
	Send(msg *protoext.SignedGossipMessage, peers ...*comm.RemotePeer)
}


type MembershipService interface {
	
	GetMembership() []discovery.NetworkMember
}


type Config struct {
	ID                string
	PullInterval      time.Duration 
	Channel           common.ChannelID
	PeerCountToSelect int 
	Tag               proto.GossipMessage_Tag
	MsgType           proto.PullMsgType
	PullEngineConfig  algo.PullEngineConfig
}


type IngressDigestFilter func(digestMsg *proto.DataDigest) *proto.DataDigest



type EgressDigestFilter func(helloMsg protoext.ReceivedMessage) func(digestItem string) bool


func (df EgressDigestFilter) byContext() algo.DigestFilter {
	return func(context interface{}) func(digestItem string) bool {
		return func(digestItem string) bool {
			return df(context.(protoext.ReceivedMessage))(digestItem)
		}
	}
}


type MsgConsumer func(message *protoext.SignedGossipMessage)


type IdentifierExtractor func(*protoext.SignedGossipMessage) string



type PullAdapter struct {
	Sndr             Sender
	MemSvc           MembershipService
	IdExtractor      IdentifierExtractor
	MsgCons          MsgConsumer
	EgressDigFilter  EgressDigestFilter
	IngressDigFilter IngressDigestFilter
}







type Mediator interface {
	
	Stop()

	
	RegisterMsgHook(MsgType, MessageHook)

	
	Add(*protoext.SignedGossipMessage)

	
	
	Remove(digest string)

	
	HandleMessage(msg protoext.ReceivedMessage)
}


type pullMediatorImpl struct {
	sync.RWMutex
	*PullAdapter
	msgType2Hook map[MsgType][]MessageHook
	config       Config
	logger       util.Logger
	itemID2Msg   map[string]*protoext.SignedGossipMessage
	engine       *algo.PullEngine
}


func NewPullMediator(config Config, adapter *PullAdapter) Mediator {
	egressDigFilter := adapter.EgressDigFilter

	acceptAllFilter := func(_ protoext.ReceivedMessage) func(string) bool {
		return func(_ string) bool {
			return true
		}
	}

	if egressDigFilter == nil {
		egressDigFilter = acceptAllFilter
	}

	p := &pullMediatorImpl{
		PullAdapter:  adapter,
		msgType2Hook: make(map[MsgType][]MessageHook),
		config:       config,
		logger:       util.GetLogger(util.PullLogger, config.ID),
		itemID2Msg:   make(map[string]*protoext.SignedGossipMessage),
	}

	p.engine = algo.NewPullEngineWithFilter(p, config.PullInterval, egressDigFilter.byContext(), config.PullEngineConfig)

	if adapter.IngressDigFilter == nil {
		
		adapter.IngressDigFilter = func(digestMsg *proto.DataDigest) *proto.DataDigest {
			return digestMsg
		}
	}
	return p

}

func (p *pullMediatorImpl) HandleMessage(m protoext.ReceivedMessage) {
	if m.GetGossipMessage() == nil || !protoext.IsPullMsg(m.GetGossipMessage().GossipMessage) {
		return
	}
	msg := m.GetGossipMessage()
	msgType := protoext.GetPullMsgType(msg.GossipMessage)
	if msgType != p.config.MsgType {
		return
	}

	p.logger.Debug(msg)

	itemIDs := []string{}
	items := []*protoext.SignedGossipMessage{}
	var pullMsgType MsgType

	if helloMsg := msg.GetHello(); helloMsg != nil {
		pullMsgType = HelloMsgType
		p.engine.OnHello(helloMsg.Nonce, m)
	} else if digest := msg.GetDataDig(); digest != nil {
		d := p.PullAdapter.IngressDigFilter(digest)
		itemIDs = util.BytesToStrings(d.Digests)
		pullMsgType = DigestMsgType
		p.engine.OnDigest(itemIDs, d.Nonce, m)
	} else if req := msg.GetDataReq(); req != nil {
		itemIDs = util.BytesToStrings(req.Digests)
		pullMsgType = RequestMsgType
		p.engine.OnReq(itemIDs, req.Nonce, m)
	} else if res := msg.GetDataUpdate(); res != nil {
		itemIDs = make([]string, len(res.Data))
		items = make([]*protoext.SignedGossipMessage, len(res.Data))
		pullMsgType = ResponseMsgType
		for i, pulledMsg := range res.Data {
			msg, err := protoext.EnvelopeToGossipMessage(pulledMsg)
			if err != nil {
				p.logger.Warningf("Data update contains an invalid message: %+v", errors.WithStack(err))
				return
			}
			p.MsgCons(msg)
			itemIDs[i] = p.IdExtractor(msg)
			items[i] = msg
			p.Lock()
			p.itemID2Msg[itemIDs[i]] = msg
			p.logger.Debugf("Added %s to the in memory item map, total items: %d", itemIDs[i], len(p.itemID2Msg))
			p.Unlock()
		}
		p.engine.OnRes(itemIDs, res.Nonce)
	}

	
	for _, h := range p.hooksByMsgType(pullMsgType) {
		h(itemIDs, items, m)
	}
}

func (p *pullMediatorImpl) Stop() {
	p.engine.Stop()
}


func (p *pullMediatorImpl) RegisterMsgHook(pullMsgType MsgType, hook MessageHook) {
	p.Lock()
	defer p.Unlock()
	p.msgType2Hook[pullMsgType] = append(p.msgType2Hook[pullMsgType], hook)

}


func (p *pullMediatorImpl) Add(msg *protoext.SignedGossipMessage) {
	p.Lock()
	defer p.Unlock()
	itemID := p.IdExtractor(msg)
	p.itemID2Msg[itemID] = msg
	p.engine.Add(itemID)
	p.logger.Debugf("Added %s, total items: %d", itemID, len(p.itemID2Msg))
}



func (p *pullMediatorImpl) Remove(digest string) {
	p.Lock()
	defer p.Unlock()
	delete(p.itemID2Msg, digest)
	p.engine.Remove(digest)
	p.logger.Debugf("Removed %s, total items: %d", digest, len(p.itemID2Msg))
}


func (p *pullMediatorImpl) SelectPeers() []string {
	remotePeers := SelectEndpoints(p.config.PeerCountToSelect, p.MemSvc.GetMembership())
	endpoints := make([]string, len(remotePeers))
	for i, peer := range remotePeers {
		endpoints[i] = peer.Endpoint
	}
	return endpoints
}




func (p *pullMediatorImpl) Hello(dest string, nonce uint64) {
	helloMsg := &proto.GossipMessage{
		Channel: p.config.Channel,
		Tag:     p.config.Tag,
		Content: &proto.GossipMessage_Hello{
			Hello: &proto.GossipHello{
				Nonce:    nonce,
				Metadata: nil,
				MsgType:  p.config.MsgType,
			},
		},
	}

	p.logger.Debug("Sending", p.config.MsgType, "hello to", dest)
	sMsg, err := protoext.NoopSign(helloMsg)
	if err != nil {
		p.logger.Errorf("Failed creating SignedGossipMessage: %+v", errors.WithStack(err))
		return
	}
	p.Sndr.Send(sMsg, p.peersWithEndpoints(dest)...)
}



func (p *pullMediatorImpl) SendDigest(digest []string, nonce uint64, context interface{}) {
	digMsg := &proto.GossipMessage{
		Channel: p.config.Channel,
		Tag:     p.config.Tag,
		Nonce:   0,
		Content: &proto.GossipMessage_DataDig{
			DataDig: &proto.DataDigest{
				MsgType: p.config.MsgType,
				Nonce:   nonce,
				Digests: util.StringsToBytes(digest),
			},
		},
	}
	remotePeer := context.(protoext.ReceivedMessage).GetConnectionInfo()
	if p.logger.IsEnabledFor(zapcore.DebugLevel) {
		p.logger.Debug("Sending", p.config.MsgType, "digest:", formattedDigests(digMsg.GetDataDig()), "to", remotePeer)
	}

	context.(protoext.ReceivedMessage).Respond(digMsg)
}



func (p *pullMediatorImpl) SendReq(dest string, items []string, nonce uint64) {
	req := &proto.GossipMessage{
		Channel: p.config.Channel,
		Tag:     p.config.Tag,
		Nonce:   0,
		Content: &proto.GossipMessage_DataReq{
			DataReq: &proto.DataRequest{
				MsgType: p.config.MsgType,
				Nonce:   nonce,
				Digests: util.StringsToBytes(items),
			},
		},
	}
	if p.logger.IsEnabledFor(zapcore.DebugLevel) {
		p.logger.Debug("Sending", formattedDigests(req.GetDataReq()), "to", dest)
	}
	sMsg, err := protoext.NoopSign(req)
	if err != nil {
		p.logger.Warningf("Failed creating SignedGossipMessage: %+v", errors.WithStack(err))
		return
	}
	p.Sndr.Send(sMsg, p.peersWithEndpoints(dest)...)
}



type typedDigester interface {
	GetMsgType() gossip.PullMsgType
	GetDigests() [][]byte
}

func formattedDigests(dd typedDigester) []string {
	switch dd.GetMsgType() {
	case gossip.PullMsgType_IDENTITY_MSG:
		return digestsToHex(dd.GetDigests())
	default:
		return digestsAsStrings(dd.GetDigests())
	}
}

func digestsAsStrings(digests [][]byte) []string {
	a := make([]string, len(digests))
	for i, dig := range digests {
		a[i] = string(dig)
	}
	return a
}

func digestsToHex(digests [][]byte) []string {
	a := make([]string, len(digests))
	for i, dig := range digests {
		a[i] = hex.EncodeToString(dig)
	}
	return a
}


func (p *pullMediatorImpl) SendRes(items []string, context interface{}, nonce uint64) {
	items2return := []*proto.Envelope{}
	p.RLock()
	defer p.RUnlock()
	for _, item := range items {
		if msg, exists := p.itemID2Msg[item]; exists {
			items2return = append(items2return, msg.Envelope)
		}
	}
	returnedUpdate := &proto.GossipMessage{
		Channel: p.config.Channel,
		Tag:     p.config.Tag,
		Nonce:   0,
		Content: &proto.GossipMessage_DataUpdate{
			DataUpdate: &proto.DataUpdate{
				MsgType: p.config.MsgType,
				Nonce:   nonce,
				Data:    items2return,
			},
		},
	}
	remotePeer := context.(protoext.ReceivedMessage).GetConnectionInfo()
	p.logger.Debug("Sending", len(returnedUpdate.GetDataUpdate().Data), p.config.MsgType, "items to", remotePeer)
	context.(protoext.ReceivedMessage).Respond(returnedUpdate)
}

func (p *pullMediatorImpl) peersWithEndpoints(endpoints ...string) []*comm.RemotePeer {
	peers := []*comm.RemotePeer{}
	for _, member := range p.MemSvc.GetMembership() {
		for _, endpoint := range endpoints {
			if member.PreferredEndpoint() == endpoint {
				peers = append(peers, &comm.RemotePeer{Endpoint: member.PreferredEndpoint(), PKIID: member.PKIid})
			}
		}
	}
	return peers
}

func (p *pullMediatorImpl) hooksByMsgType(msgType MsgType) []MessageHook {
	p.RLock()
	defer p.RUnlock()
	returnedHooks := []MessageHook{}
	for _, h := range p.msgType2Hook[msgType] {
		returnedHooks = append(returnedHooks, h)
	}
	return returnedHooks
}


func SelectEndpoints(k int, peerPool []discovery.NetworkMember) []*comm.RemotePeer {
	if len(peerPool) < k {
		k = len(peerPool)
	}

	indices := util.GetRandomIndices(k, len(peerPool)-1)
	endpoints := make([]*comm.RemotePeer, len(indices))
	for i, j := range indices {
		endpoints[i] = &comm.RemotePeer{Endpoint: peerPool[j].PreferredEndpoint(), PKIID: peerPool[j].PKIid}
	}
	return endpoints
}
