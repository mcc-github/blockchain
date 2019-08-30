/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gproto "github.com/golang/protobuf/proto"
	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/metrics/mocks"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type msgMutator func(message *proto.Envelope)

var shortenedWaitTime = time.Millisecond * 300
var conf = Config{
	ID:                          "test",
	PublishStateInfoInterval:    time.Millisecond * 100,
	MaxBlockCountToStore:        100,
	PullPeerNum:                 3,
	PullInterval:                time.Second,
	RequestStateInfoInterval:    time.Millisecond * 100,
	BlockExpirationInterval:     time.Second * 6,
	StateInfoCacheSweepInterval: time.Second,
	TimeForMembershipTracker:    time.Second * 5,
	DigestWaitTime:              shortenedWaitTime / 2,
	RequestWaitTime:             shortenedWaitTime,
	ResponseWaitTime:            shortenedWaitTime,
	MsgExpirationTimeout:        DefMsgExpirationTimeout,
}

var disabledMetrics = metrics.NewGossipMetrics(&disabled.Provider{}).MembershipMetrics

func init() {
	util.SetupTestLogging()
	factory.InitFactories(nil)
}

var (
	
	
	channelA                  = common.ChannelID("A")
	orgInChannelA             = api.OrgIdentityType("ORG1")
	orgNotInChannelA          = api.OrgIdentityType("ORG2")
	pkiIDInOrg1               = common.PKIidType("pkiIDInOrg1")
	pkiIDnilOrg               = common.PKIidType("pkIDnilOrg")
	pkiIDInOrg1ButNotEligible = common.PKIidType("pkiIDInOrg1ButNotEligible")
	pkiIDinOrg2               = common.PKIidType("pkiIDinOrg2")
)

type joinChanMsg struct {
	getTS               func() time.Time
	members2AnchorPeers map[string][]api.AnchorPeer
}




func (jcm *joinChanMsg) SequenceNumber() uint64 {
	if jcm.getTS != nil {
		return uint64(jcm.getTS().UnixNano())
	}
	return uint64(time.Now().UnixNano())
}


func (jcm *joinChanMsg) Members() []api.OrgIdentityType {
	if jcm.members2AnchorPeers == nil {
		return []api.OrgIdentityType{orgInChannelA}
	}
	members := make([]api.OrgIdentityType, len(jcm.members2AnchorPeers))
	i := 0
	for org := range jcm.members2AnchorPeers {
		members[i] = api.OrgIdentityType(org)
		i++
	}
	return members
}


func (jcm *joinChanMsg) AnchorPeersOf(org api.OrgIdentityType) []api.AnchorPeer {
	if jcm.members2AnchorPeers == nil {
		return []api.AnchorPeer{}
	}
	return jcm.members2AnchorPeers[string(org)]
}

type cryptoService struct {
	mocked bool
	mock.Mock
}

func (cs *cryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	return time.Now().Add(time.Hour), nil
}

func (cs *cryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	panic("Should not be called in this test")
}

func (cs *cryptoService) VerifyByChannel(channel common.ChannelID, identity api.PeerIdentityType, _, _ []byte) error {
	if !cs.mocked {
		return nil
	}
	args := cs.Called(identity)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (cs *cryptoService) VerifyBlock(channelID common.ChannelID, seqNum uint64, signedBlock []byte) error {
	args := cs.Called(signedBlock)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (cs *cryptoService) Sign(msg []byte) ([]byte, error) {
	panic("Should not be called in this test")
}

func (cs *cryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	panic("Should not be called in this test")
}

func (cs *cryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	panic("Should not be called in this test")
}

type receivedMsg struct {
	PKIID common.PKIidType
	msg   *protoext.SignedGossipMessage
	mock.Mock
}



func (m *receivedMsg) GetSourceEnvelope() *proto.Envelope {
	return m.msg.Envelope
}

func (m *receivedMsg) GetGossipMessage() *protoext.SignedGossipMessage {
	return m.msg
}

func (m *receivedMsg) Respond(msg *proto.GossipMessage) {
	m.Called(msg)
}


func (m *receivedMsg) Ack(err error) {

}

func (m *receivedMsg) GetConnectionInfo() *protoext.ConnectionInfo {
	return &protoext.ConnectionInfo{
		ID: m.PKIID,
	}
}

type gossipAdapterMock struct {
	mock.Mock
	sync.RWMutex
}

func (ga *gossipAdapterMock) On(methodName string, arguments ...interface{}) *mock.Call {
	ga.Lock()
	defer ga.Unlock()
	return ga.Mock.On(methodName, arguments...)
}

func (ga *gossipAdapterMock) Sign(msg *proto.GossipMessage) (*protoext.SignedGossipMessage, error) {
	return protoext.NoopSign(msg)
}

func (ga *gossipAdapterMock) GetConf() Config {
	args := ga.Called()
	return args.Get(0).(Config)
}

func (ga *gossipAdapterMock) Gossip(msg *protoext.SignedGossipMessage) {
	ga.Called(msg)
}

func (ga *gossipAdapterMock) Forward(msg protoext.ReceivedMessage) {
	ga.Called(msg)
}

func (ga *gossipAdapterMock) DeMultiplex(msg interface{}) {
	ga.Called(msg)
}

func (ga *gossipAdapterMock) GetMembership() []discovery.NetworkMember {
	args := ga.Called()
	arg := args.Get(0)
	if f, isFunc := arg.(func() []discovery.NetworkMember); isFunc {
		return f()
	} else {
		return arg.([]discovery.NetworkMember)
	}
}


func (ga *gossipAdapterMock) Lookup(PKIID common.PKIidType) *discovery.NetworkMember {
	
	if !ga.wasMocked("Lookup") {
		return &discovery.NetworkMember{}
	}
	args := ga.Called(PKIID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*discovery.NetworkMember)
}

func (ga *gossipAdapterMock) Send(msg *protoext.SignedGossipMessage, peers ...*comm.RemotePeer) {
	
	if !ga.wasMocked("Send") {
		return
	}
	ga.Called(msg, peers)
}

func (ga *gossipAdapterMock) ValidateStateInfoMessage(msg *protoext.SignedGossipMessage) error {
	args := ga.Called(msg)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (ga *gossipAdapterMock) GetOrgOfPeer(PKIIID common.PKIidType) api.OrgIdentityType {
	args := ga.Called(PKIIID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(api.OrgIdentityType)
}

func (ga *gossipAdapterMock) GetIdentityByPKIID(pkiID common.PKIidType) api.PeerIdentityType {
	if ga.wasMocked("GetIdentityByPKIID") {
		return ga.Called(pkiID).Get(0).(api.PeerIdentityType)
	}
	return api.PeerIdentityType(pkiID)
}

func (ga *gossipAdapterMock) wasMocked(methodName string) bool {
	ga.RLock()
	defer ga.RUnlock()
	for _, ec := range ga.ExpectedCalls {
		if ec.Method == methodName {
			return true
		}
	}
	return false
}

func configureAdapter(adapter *gossipAdapterMock, members ...discovery.NetworkMember) {
	adapter.On("GetConf").Return(conf)
	adapter.On("GetMembership").Return(members)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1ButNotEligible).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDinOrg2).Return(orgNotInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDnilOrg).Return(nil)
	adapter.On("GetOrgOfPeer", mock.Anything).Return(api.OrgIdentityType(nil))
}

func TestBadInput(t *testing.T) {
	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{},
		disabledMetrics, nil).(*gossipChannel)
	assert.False(t, gc.verifyMsg(nil))
	assert.False(t, gc.verifyMsg(&receivedMsg{msg: nil, PKIID: nil}))
}

func TestSelf(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	pkiID1 := common.PKIidType("1")
	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgInChannelA): {},
		},
	}
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	adapter.On("Gossip", mock.Anything)
	gc := NewGossipChannel(pkiID1, orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, nil)
	gc.UpdateLedgerHeight(1)
	gMsg := gc.Self().GossipMessage
	env := gc.Self().Envelope
	sMsg, _ := protoext.EnvelopeToGossipMessage(env)
	assert.True(t, gproto.Equal(gMsg, sMsg.GossipMessage))
	assert.Equal(t, gMsg.GetStateInfo().Properties.LedgerHeight, uint64(1))
	assert.Equal(t, gMsg.GetStateInfo().PkiId, []byte("1"))
}

func TestMsgStoreNotExpire(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}

	pkiID1 := common.PKIidType("1")
	pkiID2 := common.PKIidType("2")
	pkiID3 := common.PKIidType("3")

	peer1 := discovery.NetworkMember{PKIid: pkiID2, InternalEndpoint: "1", Endpoint: "1"}
	peer2 := discovery.NetworkMember{PKIid: pkiID2, InternalEndpoint: "2", Endpoint: "2"}
	peer3 := discovery.NetworkMember{PKIid: pkiID3, InternalEndpoint: "3", Endpoint: "3"}

	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgInChannelA): {},
		},
	}

	adapter := new(gossipAdapterMock)
	adapter.On("GetOrgOfPeer", pkiID1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiID2).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiID3).Return(orgInChannelA)

	adapter.On("ValidateStateInfoMessage", mock.Anything).Return(nil)
	adapter.On("GetMembership").Return([]discovery.NetworkMember{peer2, peer3})
	adapter.On("DeMultiplex", mock.Anything)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("GetConf").Return(conf)

	gc := NewGossipChannel(pkiID1, orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, nil)
	gc.UpdateLedgerHeight(1)
	
	gc.HandleMessage(&receivedMsg{PKIID: pkiID2, msg: createStateInfoMsg(1, pkiID2, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID3, msg: createStateInfoMsg(1, pkiID3, channelA)})

	simulateStateInfoRequest := func(pkiID []byte, outChan chan *protoext.SignedGossipMessage) {
		sentMessages := make(chan *proto.GossipMessage, 1)
		
		s, _ := protoext.NoopSign(&proto.GossipMessage{
			Tag: proto.GossipMessage_CHAN_OR_ORG,
			Content: &proto.GossipMessage_StateInfoPullReq{
				StateInfoPullReq: &proto.StateInfoPullRequest{
					Channel_MAC: GenerateMAC(pkiID, channelA),
				},
			},
		})
		snapshotReq := &receivedMsg{
			PKIID: pkiID,
			msg:   s,
		}
		snapshotReq.On("Respond", mock.Anything).Run(func(args mock.Arguments) {
			sentMessages <- args.Get(0).(*proto.GossipMessage)
		})

		go gc.HandleMessage(snapshotReq)
		select {
		case <-time.After(time.Second):
			t.Fatal("Haven't received a state info snapshot on time")
		case msg := <-sentMessages:
			for _, el := range msg.GetStateSnapshot().Elements {
				sMsg, err := protoext.EnvelopeToGossipMessage(el)
				assert.NoError(t, err)
				outChan <- sMsg
			}
		}
	}

	c := make(chan *protoext.SignedGossipMessage, 3)
	simulateStateInfoRequest(pkiID2, c)
	assert.Len(t, c, 3)

	c = make(chan *protoext.SignedGossipMessage, 3)
	simulateStateInfoRequest(pkiID3, c)
	assert.Len(t, c, 3)

	
	adapter.On("Lookup", pkiID1).Return(&peer1)
	adapter.On("Lookup", pkiID2).Return(&peer2)
	adapter.On("Lookup", pkiID3).Return(nil)
	
	
	time.Sleep(conf.StateInfoCacheSweepInterval * 2)

	c = make(chan *protoext.SignedGossipMessage, 3)
	simulateStateInfoRequest(pkiID2, c)
	assert.Len(t, c, 2)

	c = make(chan *protoext.SignedGossipMessage, 3)
	simulateStateInfoRequest(pkiID3, c)
	assert.Len(t, c, 2)
}

func TestLeaveChannel(t *testing.T) {
	
	
	
	
	
	
	
	t.Parallel()

	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			"ORG1": {},
			"ORG2": {},
		},
	}

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	members := []discovery.NetworkMember{
		{PKIid: pkiIDInOrg1},
		{PKIid: pkiIDinOrg2},
	}
	var helloPullWG sync.WaitGroup
	helloPullWG.Add(1)
	configureAdapter(adapter, members...)
	gc := NewGossipChannel(common.PKIidType("p0"), orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, nil)
	adapter.On("Send", mock.Anything, mock.Anything).Run(func(arguments mock.Arguments) {
		msg := arguments.Get(0).(*protoext.SignedGossipMessage)
		if protoext.IsPullMsg(msg.GossipMessage) {
			helloPullWG.Done()
			assert.False(t, gc.(*gossipChannel).hasLeftChannel())
		}
	})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDinOrg2, msg: createStateInfoMsg(1, pkiIDinOrg2, channelA)})
	
	gc.HandleMessage(&receivedMsg{msg: createDataMsg(2, channelA), PKIID: pkiIDInOrg1})
	assert.Len(t, gc.GetPeers(), 2)
	
	stateInfoMsg := &receivedMsg{PKIID: pkiIDinOrg2, msg: createStateInfoMsg(0, pkiIDinOrg2, channelA)}
	stateInfoMsg.GetGossipMessage().GetStateInfo().Properties.LeftChannel = true
	gc.HandleMessage(stateInfoMsg)
	assert.Len(t, gc.GetPeers(), 1)
	
	assert.Equal(t, pkiIDInOrg1, gc.GetPeers()[0].PKIid)
	var digestSendTime int32
	var DigestSentWg sync.WaitGroup
	DigestSentWg.Add(1)
	hello := createHelloMsg(pkiIDInOrg1)
	hello.On("Respond", mock.Anything).Run(func(arguments mock.Arguments) {
		atomic.AddInt32(&digestSendTime, 1)
		
		assert.Equal(t, int32(1), atomic.LoadInt32(&digestSendTime))
		DigestSentWg.Done()
	})
	
	helloPullWG.Wait()
	go gc.HandleMessage(hello)
	DigestSentWg.Wait()
	
	gc.LeaveChannel()
	
	go gc.HandleMessage(hello)
	
	assert.Len(t, gc.GetPeers(), 0)
	
	
	time.Sleep(conf.PullInterval * 3)

}

func TestChannelPeriodicalPublishStateInfo(t *testing.T) {
	t.Parallel()
	ledgerHeight := 5
	receivedMsg := int32(0)
	stateInfoReceptionChan := make(chan *protoext.SignedGossipMessage, 1)

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)

	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("Gossip", mock.Anything).Run(func(arg mock.Arguments) {
		if atomic.LoadInt32(&receivedMsg) == int32(1) {
			return
		}

		atomic.StoreInt32(&receivedMsg, int32(1))
		msg := arg.Get(0).(*protoext.SignedGossipMessage)
		stateInfoReceptionChan <- msg
	})

	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	gc.UpdateLedgerHeight(uint64(ledgerHeight))
	defer gc.Stop()

	var msg *protoext.SignedGossipMessage
	select {
	case <-time.After(time.Second * 5):
		t.Fatal("Haven't sent stateInfo on time")
	case m := <-stateInfoReceptionChan:
		msg = m
	}

	assert.Equal(t, ledgerHeight, int(msg.GetStateInfo().Properties.LedgerHeight))
}

func TestChannelMsgStoreEviction(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything).Run(func(arg mock.Arguments) {
	})

	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	defer gc.Stop()
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(100, pkiIDInOrg1, channelA)})

	var wg sync.WaitGroup

	msgsPerPhase := uint64(50)
	lastPullPhase := make(chan uint64, msgsPerPhase)
	totalPhases := uint64(4)
	phaseNum := uint64(0)
	wg.Add(int(totalPhases))

	adapter.On("Send", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(0).(*protoext.SignedGossipMessage)
		
		if !protoext.IsPullMsg(msg.GossipMessage) {
			return
		}
		
		if atomic.LoadUint64(&phaseNum) == totalPhases && protoext.IsHelloMsg(msg.GossipMessage) {
			return
		}

		start := atomic.LoadUint64(&phaseNum) * msgsPerPhase
		end := start + msgsPerPhase
		if protoext.IsHelloMsg(msg.GossipMessage) {
			
			atomic.AddUint64(&phaseNum, uint64(1))
		}

		
		currSeq := sequence(start, end)
		pullPhase := simulatePullPhase(gc, t, &wg, func(envelope *proto.Envelope) {}, currSeq...)
		pullPhase(args)

		
		if protoext.IsDataReq(msg.GossipMessage) && atomic.LoadUint64(&phaseNum) == totalPhases {
			for _, seq := range currSeq {
				lastPullPhase <- seq
			}
			close(lastPullPhase)
		}
	})
	
	wg.Wait()

	msgSentFromPullMediator := make(chan *proto.GossipMessage, 1)

	helloMsg := createHelloMsg(pkiIDInOrg1)
	helloMsg.On("Respond", mock.Anything).Run(func(arg mock.Arguments) {
		msg := arg.Get(0).(*proto.GossipMessage)
		if !protoext.IsDigestMsg(msg) {
			return
		}
		msgSentFromPullMediator <- msg
	})
	gc.HandleMessage(helloMsg)
	select {
	case msg := <-msgSentFromPullMediator:
		
		
		msgSentFromPullMediator <- msg
	case <-time.After(time.Second * 5):
		t.Fatal("Didn't reply with a digest on time")
	}
	
	assert.Len(t, msgSentFromPullMediator, 1)
	msg := <-msgSentFromPullMediator
	
	assert.True(t, protoext.IsDigestMsg(msg))
	assert.Len(t, msg.GetDataDig().Digests, adapter.GetConf().MaxBlockCountToStore+1)
	
	
	
	for seq := range lastPullPhase {
		assert.Contains(t, msg.GetDataDig().Digests, []byte(fmt.Sprintf("%d", seq)))
	}
}

func TestChannelPull(t *testing.T) {
	t.Parallel()
	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	receivedBlocksChan := make(chan *protoext.SignedGossipMessage, 2)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything).Run(func(arg mock.Arguments) {
		msg := arg.Get(0).(*protoext.SignedGossipMessage)
		if !protoext.IsDataMsg(msg.GossipMessage) {
			return
		}
		
		assert.True(t, protoext.IsDataMsg(msg.GossipMessage))
		receivedBlocksChan <- msg
	})
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	go gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(100, pkiIDInOrg1, channelA)})

	var wg sync.WaitGroup
	wg.Add(1)
	pullPhase := simulatePullPhase(gc, t, &wg, func(envelope *proto.Envelope) {}, 10, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase)

	wg.Wait()
	for expectedSeq := 10; expectedSeq <= 11; expectedSeq++ {
		select {
		case <-time.After(time.Second * 5):
			t.Fatal("Haven't received blocks on time")
		case msg := <-receivedBlocksChan:
			assert.Equal(t, uint64(expectedSeq), msg.GetDataMsg().Payload.SeqNum)
		}
	}
}

func TestChannelPullAccessControl(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	
	

	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	cs.Mock = mock.Mock{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)

	pkiID1 := common.PKIidType("1")
	pkiID2 := common.PKIidType("2")
	pkiID3 := common.PKIidType("3")

	peer1 := discovery.NetworkMember{PKIid: pkiID1, InternalEndpoint: "1", Endpoint: "1"}
	peer2 := discovery.NetworkMember{PKIid: pkiID2, InternalEndpoint: "2", Endpoint: "2"}
	peer3 := discovery.NetworkMember{PKIid: pkiID3, InternalEndpoint: "3", Endpoint: "3"}

	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(api.OrgIdentityType("ORG1"))
	adapter.On("GetOrgOfPeer", pkiID1).Return(api.OrgIdentityType("ORG1"))
	adapter.On("GetOrgOfPeer", pkiID2).Return(api.OrgIdentityType("ORG2"))
	adapter.On("GetOrgOfPeer", pkiID3).Return(api.OrgIdentityType("ORG2"))

	adapter.On("GetMembership").Return([]discovery.NetworkMember{peer1, peer2, peer3})
	adapter.On("DeMultiplex", mock.Anything)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("GetConf").Return(conf)

	sentHello := int32(0)
	adapter.On("Send", mock.Anything, mock.Anything).Run(func(arg mock.Arguments) {
		msg := arg.Get(0).(*protoext.SignedGossipMessage)
		if !protoext.IsHelloMsg(msg.GossipMessage) {
			return
		}
		atomic.StoreInt32(&sentHello, int32(1))
		peerID := string(arg.Get(1).([]*comm.RemotePeer)[0].PKIID)
		assert.Equal(t, "1", peerID)
		assert.NotEqual(t, "2", peerID, "Sent hello to peer 2 but it's in a different org")
		assert.NotEqual(t, "3", peerID, "Sent hello to peer 3 but it's in a different org")
	})

	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			"ORG1": {},
			"ORG2": {},
		},
	}
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(100, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID1, msg: createStateInfoMsg(100, pkiID1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID2, msg: createStateInfoMsg(100, pkiID2, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID3, msg: createStateInfoMsg(100, pkiID3, channelA)})

	respondedChan := make(chan *proto.GossipMessage, 1)
	messageRelayer := func(arg mock.Arguments) {
		msg := arg.Get(0).(*proto.GossipMessage)
		respondedChan <- msg
	}

	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(5, channelA), PKIID: pkiIDInOrg1})

	helloMsg := createHelloMsg(pkiID1)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	go gc.HandleMessage(helloMsg)
	select {
	case <-respondedChan:
	case <-time.After(time.Second):
		assert.Fail(t, "Didn't reply to a hello within a timely manner")
	}

	helloMsg = createHelloMsg(pkiID2)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	go gc.HandleMessage(helloMsg)
	select {
	case <-respondedChan:
		assert.Fail(t, "Shouldn't have replied to a hello, because the peer is from a foreign org")
	case <-time.After(time.Second):
	}

	
	time.Sleep(time.Second * 3)
	
	assert.Equal(t, int32(1), atomic.LoadInt32(&sentHello))
}

func TestChannelPeerNotInChannel(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	gossipMessagesSentFromChannel := make(chan *proto.GossipMessage, 1)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)

	
	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(5, channelA)})
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())

	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(5, channelA), PKIID: pkiIDinOrg2})
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())
	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(5, channelA), PKIID: pkiIDInOrg1})
	assert.Equal(t, 1, gc.(*gossipChannel).blockMsgStore.Size())

	
	messageRelayer := func(arg mock.Arguments) {
		msg := arg.Get(0).(*proto.GossipMessage)
		gossipMessagesSentFromChannel <- msg
	}
	
	
	gc.HandleMessage(&receivedMsg{msg: createStateInfoMsg(10, pkiIDInOrg1, channelA), PKIID: pkiIDInOrg1})
	helloMsg := createHelloMsg(pkiIDInOrg1)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(helloMsg)
	select {
	case <-gossipMessagesSentFromChannel:
	case <-time.After(time.Second * 5):
		t.Fatal("Didn't reply with a digest on time")
	}
	
	helloMsg = createHelloMsg(pkiIDinOrg2)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(helloMsg)
	select {
	case <-gossipMessagesSentFromChannel:
		t.Fatal("Responded with digest, but shouldn't have since peer is in ORG2 and its not in the channel")
	case <-time.After(time.Second * 1):
	}

	
	
	gc.HandleMessage(&receivedMsg{msg: createStateInfoMsg(10, pkiIDInOrg1ButNotEligible, channelA), PKIID: pkiIDInOrg1ButNotEligible})
	
	cs.On("VerifyByChannel", mock.Anything).Return(errors.New("Not eligible"))
	cs.mocked = true
	
	gc.ConfigureChannel(&joinChanMsg{})
	helloMsg = createHelloMsg(pkiIDInOrg1ButNotEligible)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(helloMsg)
	select {
	case <-gossipMessagesSentFromChannel:
		t.Fatal("Responded with digest, but shouldn't have since peer is not eligible for the channel")
	case <-time.After(time.Second * 1):
	}

	cs.Mock = mock.Mock{}

	
	req, _ := gc.(*gossipChannel).createStateInfoRequest()
	validReceivedMsg := &receivedMsg{
		msg:   req,
		PKIID: pkiIDInOrg1,
	}
	validReceivedMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(validReceivedMsg)
	select {
	case <-gossipMessagesSentFromChannel:
	case <-time.After(time.Second * 5):
		t.Fatal("Didn't reply with a digest on time")
	}

	
	invalidReceivedMsg := &receivedMsg{
		msg:   req,
		PKIID: pkiIDinOrg2,
	}
	invalidReceivedMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(invalidReceivedMsg)
	select {
	case <-gossipMessagesSentFromChannel:
		t.Fatal("Responded with digest, but shouldn't have since peer is in ORG2 and its not in the channel")
	case <-time.After(time.Second * 1):
	}

	
	req2, _ := gc.(*gossipChannel).createStateInfoRequest()
	req2.GetStateInfoPullReq().Channel_MAC = GenerateMAC(pkiIDInOrg1, common.ChannelID("B"))
	invalidReceivedMsg2 := &receivedMsg{
		msg:   req2,
		PKIID: pkiIDInOrg1,
	}
	invalidReceivedMsg2.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(invalidReceivedMsg2)
	select {
	case <-gossipMessagesSentFromChannel:
		t.Fatal("Responded with stateInfo request, but shouldn't have since it has the wrong MAC")
	case <-time.After(time.Second * 1):
	}
}

func TestChannelIsInChannel(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)

	assert.False(t, gc.IsOrgInChannel(nil))
	assert.True(t, gc.IsOrgInChannel(orgInChannelA))
	assert.False(t, gc.IsOrgInChannel(orgNotInChannelA))
	assert.True(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
}

func TestChannelIsSubscribed(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	gc.HandleMessage(&receivedMsg{msg: createStateInfoMsg(10, pkiIDInOrg1, channelA), PKIID: pkiIDInOrg1})
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
}

func TestChannelAddToMessageStore(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	demuxedMsgs := make(chan *protoext.SignedGossipMessage, 1)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything).Run(func(arg mock.Arguments) {
		demuxedMsgs <- arg.Get(0).(*protoext.SignedGossipMessage)
	})

	
	gc.AddToMsgStore(createHelloMsg(pkiIDInOrg1).GetGossipMessage())

	
	
	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(11, channelA), PKIID: pkiIDInOrg1})
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't detected a demultiplexing within a time period")
	case <-demuxedMsgs:
	}
	gc.AddToMsgStore(dataMsgOfChannel(12, channelA))
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(12, channelA), PKIID: pkiIDInOrg1})
	select {
	case <-time.After(time.Second):
	case <-demuxedMsgs:
		t.Fatal("Demultiplexing detected, even though it wasn't supposed to happen")
	}

	gc.AddToMsgStore(createStateInfoMsg(10, pkiIDInOrg1, channelA))
	helloMsg := createHelloMsg(pkiIDInOrg1)
	respondedChan := make(chan struct{}, 1)
	helloMsg.On("Respond", mock.Anything).Run(func(arg mock.Arguments) {
		respondedChan <- struct{}{}
	})
	gc.HandleMessage(helloMsg)
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't responded to hello message within a time period")
	case <-respondedChan:
	}

	gc.HandleMessage(&receivedMsg{msg: createStateInfoMsg(10, pkiIDInOrg1, channelA), PKIID: pkiIDInOrg1})
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
}

func TestChannelBlockExpiration(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	demuxedMsgs := make(chan *protoext.SignedGossipMessage, 1)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything).Run(func(arg mock.Arguments) {
		demuxedMsgs <- arg.Get(0).(*protoext.SignedGossipMessage)
	})
	respondedChan := make(chan *proto.GossipMessage, 1)
	messageRelayer := func(arg mock.Arguments) {
		msg := arg.Get(0).(*proto.GossipMessage)
		respondedChan <- msg
	}

	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(11, channelA), PKIID: pkiIDInOrg1})
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't detected a demultiplexing within a time period")
	case <-demuxedMsgs:
	}

	
	stateInfoMsg := createStateInfoMsg(10, pkiIDInOrg1, channelA)
	gc.AddToMsgStore(stateInfoMsg)
	helloMsg := createHelloMsg(pkiIDInOrg1)
	helloMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(helloMsg)
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't responded to hello message within a time period")
	case msg := <-respondedChan:
		if protoext.IsDigestMsg(msg) {
			assert.Equal(t, 1, len(msg.GetDataDig().Digests), "Number of digests returned by channel blockPuller incorrect")
		} else {
			t.Fatal("Not correct pull msg type in response - expect digest")
		}
	}

	time.Sleep(gc.(*gossipChannel).GetConf().BlockExpirationInterval + time.Second)

	
	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(11, channelA), PKIID: pkiIDInOrg1})
	select {
	case <-time.After(time.Second):
	case <-demuxedMsgs:
		t.Fatal("Demultiplexing detected, even though it wasn't supposed to happen")
	}

	
	gc.AddToMsgStore(stateInfoMsg)
	gc.HandleMessage(helloMsg)
	select {
	case <-time.After(time.Second):
	case <-respondedChan:
		t.Fatal("No digest should be sent")
	}

	time.Sleep(gc.(*gossipChannel).GetConf().BlockExpirationInterval + time.Second)
	
	
	gc.HandleMessage(&receivedMsg{msg: dataMsgOfChannel(11, channelA), PKIID: pkiIDInOrg1})
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't detected a demultiplexing within a time period")
	case <-demuxedMsgs:
	}

	
	gc.AddToMsgStore(stateInfoMsg)
	gc.HandleMessage(helloMsg)
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't responded to hello message within a time period")
	case msg := <-respondedChan:
		if protoext.IsDigestMsg(msg) {
			assert.Equal(t, 1, len(msg.GetDataDig().Digests), "Number of digests returned by channel blockPuller incorrect")
		} else {
			t.Fatal("Not correct pull msg type in response - expect digest")
		}
	}

	gc.Stop()
}

func TestChannelBadBlocks(t *testing.T) {
	t.Parallel()
	receivedMessages := make(chan *protoext.SignedGossipMessage, 1)
	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)

	adapter.On("DeMultiplex", mock.Anything).Run(func(args mock.Arguments) {
		receivedMessages <- args.Get(0).(*protoext.SignedGossipMessage)
	})

	
	gc.HandleMessage(&receivedMsg{msg: createDataMsg(1, channelA), PKIID: pkiIDInOrg1})
	assert.Len(t, receivedMessages, 1)
	<-receivedMessages 

	
	gc.HandleMessage(&receivedMsg{msg: createDataMsg(2, common.ChannelID("B")), PKIID: pkiIDInOrg1})
	assert.Len(t, receivedMessages, 0)

	
	dataMsg := createDataMsg(3, channelA)
	dataMsg.GetDataMsg().Payload = nil
	gc.HandleMessage(&receivedMsg{msg: dataMsg, PKIID: pkiIDInOrg1})
	assert.Len(t, receivedMessages, 0)

	
	cs.Mock = mock.Mock{}
	cs.On("VerifyBlock", mock.Anything).Return(errors.New("Bad signature"))
	gc.HandleMessage(&receivedMsg{msg: createDataMsg(4, channelA), PKIID: pkiIDInOrg1})
	assert.Len(t, receivedMessages, 0)
}

func TestChannelPulledBadBlocks(t *testing.T) {
	t.Parallel()

	
	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("DeMultiplex", mock.Anything)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})

	var wg sync.WaitGroup
	wg.Add(1)

	changeChan := func(env *proto.Envelope) {
		sMsg, _ := protoext.EnvelopeToGossipMessage(env)
		sMsg.Channel = []byte("B")
		sMsg, _ = protoext.NoopSign(sMsg.GossipMessage)
		env.Payload = sMsg.Payload
	}

	pullPhase1 := simulatePullPhase(gc, t, &wg, changeChan, 10, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase1)
	adapter.On("DeMultiplex", mock.Anything)
	wg.Wait()
	gc.Stop()
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())

	
	cs = &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(errors.New("Bad block"))
	adapter = new(gossipAdapterMock)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	gc = NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})

	var wg2 sync.WaitGroup
	wg2.Add(1)
	noop := func(env *proto.Envelope) {

	}
	pullPhase2 := simulatePullPhase(gc, t, &wg2, noop, 10, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase2)
	wg2.Wait()
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())

	
	cs = &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter = new(gossipAdapterMock)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	gc = NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})

	var wg3 sync.WaitGroup
	wg3.Add(1)
	emptyBlock := func(env *proto.Envelope) {
		sMsg, _ := protoext.EnvelopeToGossipMessage(env)
		sMsg.GossipMessage.GetDataMsg().Payload = nil
		sMsg, _ = protoext.NoopSign(sMsg.GossipMessage)
		env.Payload = sMsg.Payload
	}
	pullPhase3 := simulatePullPhase(gc, t, &wg3, emptyBlock, 10, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase3)
	wg3.Wait()
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())

	
	cs = &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)

	adapter = new(gossipAdapterMock)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	gc = NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})

	var wg4 sync.WaitGroup
	wg4.Add(1)
	nonBlockMsg := func(env *proto.Envelope) {
		sMsg, _ := protoext.EnvelopeToGossipMessage(env)
		sMsg.Content = createHelloMsg(pkiIDInOrg1).GetGossipMessage().Content
		sMsg, _ = protoext.NoopSign(sMsg.GossipMessage)
		env.Payload = sMsg.Payload
	}
	pullPhase4 := simulatePullPhase(gc, t, &wg4, nonBlockMsg, 10, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase4)
	wg4.Wait()
	assert.Equal(t, 0, gc.(*gossipChannel).blockMsgStore.Size())
}

func TestChannelStateInfoSnapshot(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	adapter.On("Lookup", mock.Anything).Return(&discovery.NetworkMember{Endpoint: "localhost:5000"})
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	sentMessages := make(chan *proto.GossipMessage, 10)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("ValidateStateInfoMessage", mock.Anything).Return(nil)

	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(common.ChannelID("B"), createStateInfoMsg(4, pkiIDInOrg1, channelA))})
	assert.Empty(t, gc.GetPeers())
	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDInOrg1, common.ChannelID("B")))})
	assert.Empty(t, gc.GetPeers())

	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDinOrg2, channelA))})
	assert.Empty(t, gc.GetPeers())

	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDinOrg2, msg: stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDInOrg1, channelA))})
	assert.Empty(t, gc.GetPeers())

	
	sim := createStateInfoMsg(4, pkiIDInOrg1, channelA)
	sim.GetStateInfo().Channel_MAC = append(sim.GetStateInfo().Channel_MAC, 1)
	sim, _ = protoext.NoopSign(sim.GossipMessage)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(channelA, sim)})
	assert.Empty(t, gc.GetPeers())

	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDInOrg1, channelA))})

	
	stateInfoMsg := &receivedMsg{PKIID: pkiIDInOrg1, msg: stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDInOrg1, channelA))}
	gc.HandleMessage(stateInfoMsg)
	assert.NotEmpty(t, gc.GetPeers())
	assert.Equal(t, 4, int(gc.GetPeers()[0].Properties.LedgerHeight))

	
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Tag: proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfoPullReq{
			StateInfoPullReq: &proto.StateInfoPullRequest{
				Channel_MAC: append(GenerateMAC(pkiIDInOrg1, channelA), 1),
			},
		},
	})
	snapshotReq := &receivedMsg{
		PKIID: pkiIDInOrg1,
		msg:   sMsg,
	}
	snapshotReq.On("Respond", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages <- args.Get(0).(*proto.GossipMessage)
	})

	go gc.HandleMessage(snapshotReq)
	select {
	case <-time.After(time.Second):
	case <-sentMessages:
		assert.Fail(t, "Shouldn't have responded to this StateInfoSnapshot request because of bad MAC")
	}

	
	sMsg, _ = protoext.NoopSign(&proto.GossipMessage{
		Tag: proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfoPullReq{
			StateInfoPullReq: &proto.StateInfoPullRequest{
				Channel_MAC: GenerateMAC(pkiIDInOrg1, channelA),
			},
		},
	})
	snapshotReq = &receivedMsg{
		PKIID: pkiIDInOrg1,
		msg:   sMsg,
	}
	snapshotReq.On("Respond", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages <- args.Get(0).(*proto.GossipMessage)
	})

	go gc.HandleMessage(snapshotReq)
	select {
	case <-time.After(time.Second):
		t.Fatal("Haven't received a state info snapshot on time")
	case msg := <-sentMessages:
		elements := msg.GetStateSnapshot().Elements
		assert.Len(t, elements, 1)
		sMsg, err := protoext.EnvelopeToGossipMessage(elements[0])
		assert.NoError(t, err)
		assert.Equal(t, 4, int(sMsg.GetStateInfo().Properties.LedgerHeight))
	}

	
	invalidStateInfoSnapshot := stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, pkiIDInOrg1, channelA))
	invalidStateInfoSnapshot.GetStateSnapshot().Elements = []*proto.Envelope{createHelloMsg(pkiIDInOrg1).GetSourceEnvelope()}
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: invalidStateInfoSnapshot})

	
	invalidStateInfoSnapshot = stateInfoSnapshotForChannel(channelA, createStateInfoMsg(4, common.PKIidType("unknown"), channelA))
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: invalidStateInfoSnapshot})
}

func TestInterOrgExternalEndpointDisclosure(t *testing.T) {
	t.Parallel()
	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	pkiID1 := common.PKIidType("withExternalEndpoint")
	pkiID2 := common.PKIidType("noExternalEndpoint")
	pkiID3 := common.PKIidType("pkiIDinOrg2")
	adapter.On("Lookup", pkiID1).Return(&discovery.NetworkMember{Endpoint: "localhost:5000"})
	adapter.On("Lookup", pkiID2).Return(&discovery.NetworkMember{})
	adapter.On("Lookup", pkiID3).Return(&discovery.NetworkMember{})
	adapter.On("GetOrgOfPeer", pkiID1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiID2).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiID3).Return(api.OrgIdentityType("ORG2"))
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	configureAdapter(adapter)
	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgInChannelA): {},
			"ORG2":                {},
		},
	}
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiID1, msg: createStateInfoMsg(0, pkiID1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID2, msg: createStateInfoMsg(0, pkiID2, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiID2, msg: createStateInfoMsg(0, pkiID3, channelA)})

	sentMessages := make(chan *proto.GossipMessage, 10)

	
	
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Tag: proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfoPullReq{
			StateInfoPullReq: &proto.StateInfoPullRequest{
				Channel_MAC: GenerateMAC(pkiID3, channelA),
			},
		},
	})
	snapshotReq := &receivedMsg{
		PKIID: pkiID3,
		msg:   sMsg,
	}
	snapshotReq.On("Respond", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages <- args.Get(0).(*proto.GossipMessage)
	})

	go gc.HandleMessage(snapshotReq)
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "Should have responded to this StateInfoSnapshot, but didn't")
	case msg := <-sentMessages:
		elements := msg.GetStateSnapshot().Elements
		assert.Len(t, elements, 2)
		m1, _ := protoext.EnvelopeToGossipMessage(elements[0])
		m2, _ := protoext.EnvelopeToGossipMessage(elements[1])
		pkiIDs := [][]byte{m1.GetStateInfo().PkiId, m2.GetStateInfo().PkiId}
		assert.Contains(t, pkiIDs, []byte(pkiID1))
		assert.Contains(t, pkiIDs, []byte(pkiID3))
	}

	
	
	sMsg, _ = protoext.NoopSign(&proto.GossipMessage{
		Tag: proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfoPullReq{
			StateInfoPullReq: &proto.StateInfoPullRequest{
				Channel_MAC: GenerateMAC(pkiID2, channelA),
			},
		},
	})
	snapshotReq = &receivedMsg{
		PKIID: pkiID2,
		msg:   sMsg,
	}
	snapshotReq.On("Respond", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages <- args.Get(0).(*proto.GossipMessage)
	})

	go gc.HandleMessage(snapshotReq)
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "Should have responded to this StateInfoSnapshot, but didn't")
	case msg := <-sentMessages:
		elements := msg.GetStateSnapshot().Elements
		assert.Len(t, elements, 3)
		m1, _ := protoext.EnvelopeToGossipMessage(elements[0])
		m2, _ := protoext.EnvelopeToGossipMessage(elements[1])
		m3, _ := protoext.EnvelopeToGossipMessage(elements[2])
		pkiIDs := [][]byte{m1.GetStateInfo().PkiId, m2.GetStateInfo().PkiId, m3.GetStateInfo().PkiId}
		assert.Contains(t, pkiIDs, []byte(pkiID1))
		assert.Contains(t, pkiIDs, []byte(pkiID2))
		assert.Contains(t, pkiIDs, []byte(pkiID3))
	}
}

func TestChannelStop(t *testing.T) {
	t.Parallel()

	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	adapter := new(gossipAdapterMock)
	var sendCount int32
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("Send", mock.Anything, mock.Anything).Run(func(mock.Arguments) {
		atomic.AddInt32(&sendCount, int32(1))
	})
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	time.Sleep(time.Second)
	gc.Stop()
	oldCount := atomic.LoadInt32(&sendCount)
	t1 := time.Now()
	for {
		if time.Since(t1).Nanoseconds() > (time.Second * 15).Nanoseconds() {
			t.Fatal("Stop failed")
		}
		time.Sleep(time.Second)
		newCount := atomic.LoadInt32(&sendCount)
		if newCount == oldCount {
			break
		}
		oldCount = newCount
	}
}

func TestChannelReconfigureChannel(t *testing.T) {
	t.Parallel()

	
	
	
	
	
	

	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})

	adapter.On("GetConf").Return(conf)
	adapter.On("GetMembership").Return([]discovery.NetworkMember{})
	adapter.On("OrgByPeerIdentity", api.PeerIdentityType(orgInChannelA)).Return(orgInChannelA)
	adapter.On("OrgByPeerIdentity", api.PeerIdentityType(orgNotInChannelA)).Return(orgNotInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDinOrg2).Return(orgNotInChannelA)

	outdatedJoinChanMsg := &joinChanMsg{
		getTS: func() time.Time {
			return time.Now()
		},
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgNotInChannelA): {},
		},
	}

	newJoinChanMsg := &joinChanMsg{
		getTS: func() time.Time {
			return time.Now().Add(time.Millisecond * 100)
		},
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgInChannelA): {},
		},
	}

	updatedJoinChanMsg := &joinChanMsg{
		getTS: func() time.Time {
			return time.Now().Add(time.Millisecond * 200)
		},
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgNotInChannelA): {},
		},
	}

	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, api.JoinChannelMessage(newJoinChanMsg),
		disabledMetrics, nil)

	
	gc.ConfigureChannel(api.JoinChannelMessage(newJoinChanMsg))

	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)

	assert.True(t, gc.IsOrgInChannel(orgInChannelA))
	assert.False(t, gc.IsOrgInChannel(orgNotInChannelA))
	assert.True(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDinOrg2}))

	gc.ConfigureChannel(outdatedJoinChanMsg)
	assert.True(t, gc.IsOrgInChannel(orgInChannelA))
	assert.False(t, gc.IsOrgInChannel(orgNotInChannelA))
	assert.True(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDinOrg2}))

	gc.ConfigureChannel(updatedJoinChanMsg)
	gc.ConfigureChannel(updatedJoinChanMsg)
	assert.False(t, gc.IsOrgInChannel(orgInChannelA))
	assert.True(t, gc.IsOrgInChannel(orgNotInChannelA))
	assert.False(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.True(t, gc.IsMemberInChan(discovery.NetworkMember{PKIid: pkiIDinOrg2}))

	
	sMsg, _ := gc.(*gossipChannel).createStateInfoRequest()
	invalidReceivedMsg := &receivedMsg{
		msg:   sMsg,
		PKIID: pkiIDInOrg1,
	}
	gossipMessagesSentFromChannel := make(chan *proto.GossipMessage, 1)
	messageRelayer := func(arg mock.Arguments) {
		msg := arg.Get(0).(*proto.GossipMessage)
		gossipMessagesSentFromChannel <- msg
	}
	invalidReceivedMsg.On("Respond", mock.Anything).Run(messageRelayer)
	gc.HandleMessage(invalidReceivedMsg)
	select {
	case <-gossipMessagesSentFromChannel:
		t.Fatal("Responded with digest, but shouldn't have since peer is in ORG2 and its not in the channel")
	case <-time.After(time.Second * 1):
	}
}

func TestChannelNoAnchorPeers(t *testing.T) {
	t.Parallel()

	
	

	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})

	adapter.On("GetConf").Return(conf)
	adapter.On("GetMembership").Return([]discovery.NetworkMember{})
	adapter.On("OrgByPeerIdentity", api.PeerIdentityType(orgInChannelA)).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDinOrg2).Return(orgNotInChannelA)

	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(orgInChannelA): {},
		},
	}

	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, api.JoinChannelMessage(jcm),
		disabledMetrics, nil)
	assert.True(t, gc.IsOrgInChannel(orgInChannelA))
}

func TestGossipChannelEligibility(t *testing.T) {
	t.Parallel()

	
	
	
	
	

	cs := &cryptoService{}
	selfPKIID := common.PKIidType("p")
	adapter := new(gossipAdapterMock)
	pkiIDinOrg3 := common.PKIidType("pkiIDinOrg3")
	members := []discovery.NetworkMember{
		{PKIid: pkiIDInOrg1},
		{PKIid: pkiIDInOrg1ButNotEligible},
		{PKIid: pkiIDinOrg2},
		{PKIid: pkiIDinOrg3},
	}
	adapter.On("GetMembership").Return(members)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	adapter.On("GetConf").Return(conf)

	
	org1 := api.OrgIdentityType("ORG1")
	org2 := api.OrgIdentityType("ORG2")
	org3 := api.OrgIdentityType("ORG3")

	adapter.On("GetOrgOfPeer", selfPKIID).Return(org1)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(org1)
	adapter.On("GetOrgOfPeer", pkiIDinOrg2).Return(org2)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1ButNotEligible).Return(org1)
	adapter.On("GetOrgOfPeer", pkiIDinOrg3).Return(org3)

	gc := NewGossipChannel(selfPKIID, orgInChannelA, cs, channelA, adapter, &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(org1): {},
			string(org2): {},
		},
	}, disabledMetrics, nil)
	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDinOrg2, msg: createStateInfoMsg(1, pkiIDinOrg2, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1ButNotEligible, msg: createStateInfoMsg(1, pkiIDInOrg1ButNotEligible, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDinOrg3, msg: createStateInfoMsg(1, pkiIDinOrg3, channelA)})

	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1ButNotEligible}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg3}))

	
	assert.True(t, gc.PeerFilter(func(signature api.PeerSignature) bool {
		return true
	})(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.True(t, gc.PeerFilter(func(signature api.PeerSignature) bool {
		return true
	})(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	
	assert.False(t, gc.PeerFilter(func(signature api.PeerSignature) bool {
		return true
	})(discovery.NetworkMember{PKIid: pkiIDinOrg3}))

	
	assert.True(t, gc.PeerFilter(func(signature api.PeerSignature) bool {
		return bytes.Equal(signature.PeerIdentity, []byte("pkiIDinOrg2"))
	})(discovery.NetworkMember{PKIid: pkiIDinOrg2}))

	assert.False(t, gc.PeerFilter(func(signature api.PeerSignature) bool {
		return bytes.Equal(signature.PeerIdentity, []byte("pkiIDinOrg2"))
	})(discovery.NetworkMember{PKIid: pkiIDInOrg1}))

	
	gc.ConfigureChannel(&joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(org1): {},
		},
	})

	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1ButNotEligible}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg3}))

	
	cs.mocked = true
	cs.On("VerifyByChannel", api.PeerIdentityType(pkiIDInOrg1ButNotEligible)).Return(errors.New("Not a channel reader"))
	cs.On("VerifyByChannel", mock.Anything).Return(nil)
	gc.ConfigureChannel(&joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(org1): {},
		},
	})
	assert.True(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1ButNotEligible}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg3}))

	
	
	
	adapter.On("GetIdentityByPKIID", pkiIDInOrg1).Return(api.PeerIdentityType(nil))
	adapter.On("GetIdentityByPKIID", pkiIDinOrg2).Return(api.PeerIdentityType(pkiIDinOrg2))
	adapter.On("GetIdentityByPKIID", pkiIDInOrg1ButNotEligible).Return(api.PeerIdentityType(pkiIDInOrg1ButNotEligible))
	adapter.On("GetIdentityByPKIID", pkiIDinOrg3).Return(api.PeerIdentityType(pkiIDinOrg3))

	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1ButNotEligible}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg3}))

	
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(2, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(2, pkiIDinOrg2, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(2, pkiIDInOrg1ButNotEligible, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(2, pkiIDinOrg3, channelA)})

	
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg2}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDInOrg1ButNotEligible}))
	assert.False(t, gc.EligibleForChannel(discovery.NetworkMember{PKIid: pkiIDinOrg3}))
}

func TestChannelGetPeers(t *testing.T) {
	t.Parallel()

	
	
	
	
	
	
	cs := &cryptoService{}
	adapter := new(gossipAdapterMock)
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	members := []discovery.NetworkMember{
		{PKIid: pkiIDInOrg1},
		{PKIid: pkiIDInOrg1ButNotEligible},
		{PKIid: pkiIDinOrg2},
	}
	configureAdapter(adapter, members...)
	gc := NewGossipChannel(common.PKIidType("p0"), orgInChannelA, cs, channelA, adapter, &joinChanMsg{},
		disabledMetrics, nil)
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDinOrg2, channelA)})
	assert.Len(t, gc.GetPeers(), 1)
	assert.Equal(t, pkiIDInOrg1, gc.GetPeers()[0].PKIid)

	
	gMsg, _ := protoext.EnvelopeToGossipMessage(gc.GetPeers()[0].Envelope)
	assert.Equal(t, []byte(pkiIDInOrg1), gMsg.GetStateInfo().PkiId)

	gc.HandleMessage(&receivedMsg{msg: createStateInfoMsg(10, pkiIDInOrg1ButNotEligible, channelA), PKIID: pkiIDInOrg1ButNotEligible})
	cs.On("VerifyByChannel", mock.Anything).Return(errors.New("Not eligible"))
	cs.mocked = true
	
	gc.ConfigureChannel(&joinChanMsg{})
	assert.Len(t, gc.GetPeers(), 0)

	
	
	gc = NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	msg := &receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)}
	msg.GetGossipMessage().GetStateInfo().Channel_MAC = GenerateMAC(pkiIDinOrg2, channelA)
	gc.HandleMessage(msg)
	assert.Len(t, gc.GetPeers(), 0)
}

func TestOnDemandGossip(t *testing.T) {
	t.Parallel()

	
	

	cs := &cryptoService{}

	adapter := new(gossipAdapterMock)
	configureAdapter(adapter)

	adapter.ExpectedCalls = append(adapter.ExpectedCalls[:1], adapter.ExpectedCalls[2:]...)
	var lock sync.RWMutex
	var membershipKnown bool
	var signal bool
	adapter.On("GetMembership").Return(func() []discovery.NetworkMember {
		lock.RLock()
		defer lock.RUnlock()
		if !membershipKnown {
			return []discovery.NetworkMember{}
		}
		return []discovery.NetworkMember{{}}
	})

	gossipedEvents := make(chan struct{})
	adapter.On("Gossip", mock.Anything).Run(func(mock.Arguments) {
		lock.Lock()
		defer lock.Unlock()
		if signal {
			membershipKnown = true
		}
		gossipedEvents <- struct{}{}
	})
	adapter.On("Forward", mock.Anything)

	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, api.JoinChannelMessage(&joinChanMsg{}),
		disabledMetrics, nil)
	defer gc.Stop()
	select {
	case <-gossipedEvents:
		assert.Fail(t, "Should not have gossiped because metadata has not been updated yet")
	case <-time.After(time.Millisecond * 500):
	}
	gc.UpdateLedgerHeight(0)
	select {
	case <-gossipedEvents:
	case <-time.After(time.Second):
		assert.Fail(t, "Didn't gossip within a timely manner")
	}
	select {
	case <-gossipedEvents:
	case <-time.After(time.Second):
		assert.Fail(t, "Should have gossiped a second time, because membership is empty")
	}

	lock.Lock()
	signal = true
	lock.Unlock()

	select {
	case <-gossipedEvents:
	case <-time.After(time.Second):
		assert.Fail(t, "Should have gossiped a third time")
	}
	select {
	case <-gossipedEvents:
		assert.Fail(t, "Should not have gossiped a fourth time, because dirty flag should have been turned off")
	case <-time.After(time.Millisecond * 500):
	}
	gc.UpdateLedgerHeight(1)
	select {
	case <-gossipedEvents:
	case <-time.After(time.Second):
		assert.Fail(t, "Should have gossiped a block now, because got a new StateInfo message")
	}
}

func TestChannelPullWithDigestsFilter(t *testing.T) {
	t.Parallel()
	cs := &cryptoService{}
	cs.On("VerifyBlock", mock.Anything).Return(nil)
	receivedBlocksChan := make(chan *protoext.SignedGossipMessage, 2)
	adapter := new(gossipAdapterMock)
	configureAdapter(adapter, discovery.NetworkMember{PKIid: pkiIDInOrg1})
	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("DeMultiplex", mock.Anything).Run(func(arg mock.Arguments) {
		msg := arg.Get(0).(*protoext.SignedGossipMessage)
		if !protoext.IsDataMsg(msg.GossipMessage) {
			return
		}
		
		assert.True(t, protoext.IsDataMsg(msg.GossipMessage))
		receivedBlocksChan <- msg
	})
	gc := NewGossipChannel(pkiIDInOrg1, orgInChannelA, cs, channelA, adapter, &joinChanMsg{}, disabledMetrics, nil)
	go gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(100, pkiIDInOrg1, channelA)})

	gc.UpdateLedgerHeight(11)

	var wg sync.WaitGroup
	wg.Add(1)

	pullPhase := simulatePullPhaseWithVariableDigest(gc, t, &wg, func(envelope *proto.Envelope) {}, [][]byte{[]byte("10"), []byte("11")}, []string{"11"}, 11)
	adapter.On("Send", mock.Anything, mock.Anything).Run(pullPhase)
	wg.Wait()

	select {
	case <-time.After(time.Second * 5):
		t.Fatal("Haven't received blocks on time")
	case msg := <-receivedBlocksChan:
		assert.Equal(t, uint64(11), msg.GetDataMsg().Payload.SeqNum)
	}

}

func TestFilterForeignOrgLeadershipMessages(t *testing.T) {
	t.Parallel()

	org1 := api.OrgIdentityType("org1")
	org2 := api.OrgIdentityType("org2")

	p1 := common.PKIidType("p1")
	p2 := common.PKIidType("p2")

	cs := &cryptoService{}
	adapter := &gossipAdapterMock{}

	relayedLeadershipMsgs := make(chan interface{}, 2)

	adapter.On("GetOrgOfPeer", p1).Return(org1)
	adapter.On("GetOrgOfPeer", p2).Return(org2)

	adapter.On("GetMembership").Return([]discovery.NetworkMember{})
	adapter.On("GetConf").Return(conf)
	adapter.On("DeMultiplex", mock.Anything).Run(func(args mock.Arguments) {
		relayedLeadershipMsgs <- args.Get(0)
	})

	joinMsg := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			string(org1): {},
			string(org2): {},
		},
	}

	loggedEntries := make(chan string, 1)
	logger := flogging.MustGetLogger("test").WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		loggedEntries <- entry.Message
		return nil
	}))
	assertLogged := func(s string) {
		assert.Len(t, loggedEntries, 1)
		loggedEntry := <-loggedEntries
		assert.Contains(t, loggedEntry, s)
	}

	gc := NewGossipChannel(pkiIDInOrg1, org1, cs, channelA, adapter, joinMsg, disabledMetrics, logger)

	leadershipMsg := func(sender common.PKIidType, creator common.PKIidType) protoext.ReceivedMessage {
		return &receivedMsg{PKIID: sender,
			msg: &protoext.SignedGossipMessage{
				GossipMessage: &proto.GossipMessage{
					Channel: common.ChannelID("A"),
					Tag:     proto.GossipMessage_CHAN_AND_ORG,
					Content: &proto.GossipMessage_LeadershipMsg{
						LeadershipMsg: &proto.LeadershipMessage{
							PkiId: creator,
						}},
				},
			},
		}
	}

	gc.HandleMessage(leadershipMsg(p1, p1))
	assert.Len(t, relayedLeadershipMsgs, 1, "should have relayed a message from p1 (same org)")
	assert.Len(t, loggedEntries, 0)

	gc.HandleMessage(leadershipMsg(p2, p1))
	assert.Len(t, relayedLeadershipMsgs, 1, "should not have relayed a message from p2 (foreign org)")
	assertLogged("Received leadership message from  that belongs to a foreign organization org2")

	gc.HandleMessage(leadershipMsg(p1, p2))
	assert.Len(t, relayedLeadershipMsgs, 1, "should not have relayed a message from p2 (foreign org)")
	assertLogged("Received leadership message created by a foreign organization org2")
}

func createDataUpdateMsg(nonce uint64, seqs ...uint64) *protoext.SignedGossipMessage {
	msg := &proto.GossipMessage{
		Nonce:   0,
		Channel: []byte(channelA),
		Tag:     proto.GossipMessage_CHAN_AND_ORG,
		Content: &proto.GossipMessage_DataUpdate{
			DataUpdate: &proto.DataUpdate{
				MsgType: proto.PullMsgType_BLOCK_MSG,
				Nonce:   nonce,
				Data:    []*proto.Envelope{},
			},
		},
	}
	for _, seq := range seqs {
		msg.GetDataUpdate().Data = append(msg.GetDataUpdate().Data, createDataMsg(seq, channelA).Envelope)
	}
	sMsg, _ := protoext.NoopSign(msg)
	return sMsg
}

func createHelloMsg(PKIID common.PKIidType) *receivedMsg {
	msg := &proto.GossipMessage{
		Channel: []byte(channelA),
		Tag:     proto.GossipMessage_CHAN_AND_ORG,
		Content: &proto.GossipMessage_Hello{
			Hello: &proto.GossipHello{
				Nonce:    500,
				Metadata: nil,
				MsgType:  proto.PullMsgType_BLOCK_MSG,
			},
		},
	}
	sMsg, _ := protoext.NoopSign(msg)
	return &receivedMsg{msg: sMsg, PKIID: PKIID}
}

func dataMsgOfChannel(seqnum uint64, channel common.ChannelID) *protoext.SignedGossipMessage {
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Channel: []byte(channel),
		Nonce:   0,
		Tag:     proto.GossipMessage_CHAN_AND_ORG,
		Content: &proto.GossipMessage_DataMsg{
			DataMsg: &proto.DataMessage{
				Payload: &proto.Payload{
					Data:   []byte{},
					SeqNum: seqnum,
				},
			},
		},
	})
	return sMsg
}

func createStateInfoMsg(ledgerHeight int, pkiID common.PKIidType, channel common.ChannelID) *protoext.SignedGossipMessage {
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Tag: proto.GossipMessage_CHAN_OR_ORG,
		Content: &proto.GossipMessage_StateInfo{
			StateInfo: &proto.StateInfo{
				Channel_MAC: GenerateMAC(pkiID, channel),
				Timestamp:   &proto.PeerTime{IncNum: uint64(time.Now().UnixNano()), SeqNum: 1},
				PkiId:       []byte(pkiID),
				Properties: &proto.Properties{
					LedgerHeight: uint64(ledgerHeight),
				},
			},
		},
	})
	return sMsg
}

func stateInfoSnapshotForChannel(channelID common.ChannelID, stateInfoMsgs ...*protoext.SignedGossipMessage) *protoext.SignedGossipMessage {
	envelopes := make([]*proto.Envelope, len(stateInfoMsgs))
	for i, sim := range stateInfoMsgs {
		envelopes[i] = sim.Envelope
	}
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Channel: channelID,
		Tag:     proto.GossipMessage_CHAN_OR_ORG,
		Nonce:   0,
		Content: &proto.GossipMessage_StateSnapshot{
			StateSnapshot: &proto.StateInfoSnapshot{
				Elements: envelopes,
			},
		},
	})
	return sMsg
}

func createDataMsg(seqnum uint64, channel common.ChannelID) *protoext.SignedGossipMessage {
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Nonce:   0,
		Tag:     proto.GossipMessage_CHAN_AND_ORG,
		Channel: []byte(channel),
		Content: &proto.GossipMessage_DataMsg{
			DataMsg: &proto.DataMessage{
				Payload: &proto.Payload{
					Data:   []byte{},
					SeqNum: seqnum,
				},
			},
		},
	})
	return sMsg
}

func simulatePullPhase(gc GossipChannel, t *testing.T, wg *sync.WaitGroup, mutator msgMutator, seqs ...uint64) func(args mock.Arguments) {
	return simulatePullPhaseWithVariableDigest(gc, t, wg, mutator, [][]byte{[]byte("10"), []byte("11")}, []string{"10", "11"}, seqs...)
}

func simulatePullPhaseWithVariableDigest(gc GossipChannel, t *testing.T, wg *sync.WaitGroup, mutator msgMutator, proposedDigestSeqs [][]byte, resultDigestSeqs []string, seqs ...uint64) func(args mock.Arguments) {
	var l sync.Mutex
	var sentHello bool
	var sentReq bool
	return func(args mock.Arguments) {
		msg := args.Get(0).(*protoext.SignedGossipMessage)
		l.Lock()
		defer l.Unlock()
		if protoext.IsHelloMsg(msg.GossipMessage) && !sentHello {
			sentHello = true
			
			sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
				Tag:     proto.GossipMessage_CHAN_AND_ORG,
				Channel: []byte(channelA),
				Content: &proto.GossipMessage_DataDig{
					DataDig: &proto.DataDigest{
						MsgType: proto.PullMsgType_BLOCK_MSG,
						Digests: proposedDigestSeqs,
						Nonce:   msg.GetHello().Nonce,
					},
				},
			})
			digestMsg := &receivedMsg{
				PKIID: pkiIDInOrg1,
				msg:   sMsg,
			}
			go gc.HandleMessage(digestMsg)
		}
		if protoext.IsDataReq(msg.GossipMessage) && !sentReq {
			sentReq = true
			dataReq := msg.GetDataReq()
			for _, expectedDigest := range util.StringsToBytes(resultDigestSeqs) {
				assert.Contains(t, dataReq.Digests, expectedDigest)
			}
			assert.Equal(t, len(resultDigestSeqs), len(dataReq.Digests))
			
			
			dataUpdateMsg := new(receivedMsg)
			dataUpdateMsg.PKIID = pkiIDInOrg1
			dataUpdateMsg.msg = createDataUpdateMsg(dataReq.Nonce, seqs...)
			mutator(dataUpdateMsg.msg.GetDataUpdate().Data[0])
			gc.HandleMessage(dataUpdateMsg)
			wg.Done()
		}
	}
}

func sequence(start uint64, end uint64) []uint64 {
	sequence := make([]uint64, end-start+1)
	i := 0
	for n := start; n <= end; n++ {
		sequence[i] = n
		i++
	}
	return sequence

}

func TestChangesInPeers(t *testing.T) {
	
	
	
	
	
	
	
	
	type testCase struct {
		name                     string
		oldMembers               map[string]struct{}
		newMembers               map[string]struct{}
		expected                 []string
		entryInChannel           func(chan string)
		expectedReportInvocation bool
		expectedTotal            float64
	}
	var cases = []testCase{
		{
			name:       "noChanges",
			oldMembers: map[string]struct{}{"pkiID11": {}, "pkiID22": {}, "pkiID33": {}},
			newMembers: map[string]struct{}{"pkiID11": {}, "pkiID22": {}, "pkiID33": {}},
			expected:   []string{""},
			entryInChannel: func(chStr chan string) {
				chStr <- ""
			},
			expectedReportInvocation: false,
			expectedTotal:            3,
		},
		{
			name:       "newPeerWasAdded",
			oldMembers: map[string]struct{}{"pkiID1": {}},
			newMembers: map[string]struct{}{"pkiID1": {}, "pkiID3": {}},
			expected: []string{"Membership view has changed. peers went online: [[pkiID3]], current view: [[pkiID1] [pkiID3]]",
				"Membership view has changed. peers went online: [[pkiID3]], current view: [[pkiID3] [pkiID1]]"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            2,
		},

		{
			name:       "newPeerAddedOldPeerDeleted",
			oldMembers: map[string]struct{}{"pkiID1": {}, "pkiID2": {}},
			newMembers: map[string]struct{}{"pkiID1": {}, "pkiID3": {}},
			expected: []string{"Membership view has changed. peers went offline: [[pkiID2]], peers went online: [[pkiID3]], current view: [[pkiID1] [pkiID3]]",
				"Membership view has changed. peers went offline: [[pkiID2]], peers went online: [[pkiID3]], current view: [[pkiID3] [pkiID1]]"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            2,
		},
		{
			name:       "newPeersAddedOldPeerStayed",
			oldMembers: map[string]struct{}{"pkiID1": {}},
			newMembers: map[string]struct{}{"pkiID2": {}},
			expected: []string{"Membership view has changed. peers went offline: [[pkiID1]], peers went online: [[pkiID2]], current view: [[pkiID2]]",
				"Membership view has changed. peers went offline: [[pkiID1]], peers went online: [[pkiID2]], current view: [[pkiID2]]"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            1,
		},
		{
			name:                     "newPeersAddedNoOldPeers",
			oldMembers:               map[string]struct{}{},
			newMembers:               map[string]struct{}{"pkiID1": {}},
			expected:                 []string{"Membership view has changed. peers went online: [[pkiID1]], current view: [[pkiID1]]"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            1,
		},
		{
			name:                     "PeerWasDeletedNoNewPeers",
			oldMembers:               map[string]struct{}{"pkiID1": {}},
			newMembers:               map[string]struct{}{},
			expected:                 []string{"Membership view has changed. peers went offline: [[pkiID1]], current view: []"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            0,
		},
		{
			name:       "onePeerWasDeletedRestStayed",
			oldMembers: map[string]struct{}{"pkiID01": {}, "pkiID02": {}, "pkiID03": {}},
			newMembers: map[string]struct{}{"pkiID01": {}, "pkiID02": {}},
			expected: []string{"Membership view has changed. peers went offline: [[pkiID03]], current view: [[pkiID01] [pkiID02]]",
				"Membership view has changed. peers went offline: [[pkiID03]], current view: [[pkiID02] [pkiID01]]"},
			entryInChannel:           func(chStr chan string) {},
			expectedReportInvocation: true,
			expectedTotal:            2,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			
			chForString := make(chan string, 1)
			
			funcLogger := func(a ...interface{}) {
				chForString <- fmt.Sprint(a...)
			}

			tickChan := make(chan time.Time)

			buildMembers := func(rangeMembers map[string]struct{}) []discovery.NetworkMember {
				var members []discovery.NetworkMember
				for peerID := range rangeMembers {
					peer := discovery.NetworkMember{
						Endpoint:         peerID,
						InternalEndpoint: peerID,
					}
					peer.PKIid = common.PKIidType(peerID)
					members = append(members, peer)
				}
				return members
			}

			stopChan := make(chan struct{})

			getPeersToTrackCallCount := 0
			getListOfPeers := func() []discovery.NetworkMember {
				var members []discovery.NetworkMember
				if getPeersToTrackCallCount == 0 {
					members = buildMembers(test.oldMembers)
					getPeersToTrackCallCount++
				} else if getPeersToTrackCallCount == 1 {
					members = buildMembers(test.newMembers)
					getPeersToTrackCallCount++
					close(stopChan) 
				} else {
					t.Fatal("getPeersToTrack called too many times")
				}
				return members
			}

			testMetricProvider := mocks.TestUtilConstructMetricProvider()
			metrics := metrics.NewGossipMetrics(testMetricProvider.FakeProvider).MembershipMetrics

			mt := &membershipTracker{
				getPeersToTrack: getListOfPeers,
				report:          funcLogger,
				stopChan:        stopChan,
				tickerChannel:   tickChan,
				metrics:         metrics,
				chainID:         common.ChannelID("test"),
			}

			wgMT := sync.WaitGroup{}
			wgMT.Add(1)
			go func() {
				mt.trackMembershipChanges()
				wgMT.Done()
			}()

			tickChan <- time.Time{}

			test.entryInChannel(chForString) 
			actual := <-chForString

			
			assert.Contains(t, test.expected, actual)

			
			wgMT.Wait()
			if testMetricProvider.FakeTotalGauge.WithCallCount() < 1 {
				t.Fatal("did not get With() call")
			}
			if testMetricProvider.FakeTotalGauge.SetCallCount() < 1 {
				t.Fatal("did not get Set() call")
			}
			assert.Equal(t, []string{"channel", "test"}, testMetricProvider.FakeTotalGauge.WithArgsForCall(0))
			assert.EqualValues(t, test.expectedTotal, testMetricProvider.FakeTotalGauge.SetArgsForCall(0))
		})
	}
}

func TestMembershiptrackerStopWhenGCStops(t *testing.T) {
	
	
	
	
	t.Parallel()
	membershipReported := make(chan struct{}, 1)
	cs := &cryptoService{}
	pkiID1 := common.PKIidType("1")
	adapter := new(gossipAdapterMock)

	jcm := &joinChanMsg{}

	peerA := discovery.NetworkMember{
		PKIid:            pkiIDInOrg1,
		Endpoint:         "a",
		InternalEndpoint: "a",
	}
	peerB := discovery.NetworkMember{
		PKIid:            pkiIDinOrg2,
		Endpoint:         "b",
		InternalEndpoint: "b",
	}

	conf := conf
	conf.RequestStateInfoInterval = time.Hour
	conf.PullInterval = time.Hour
	conf.TimeForMembershipTracker = time.Millisecond * 10

	adapter.On("Gossip", mock.Anything)
	adapter.On("Forward", mock.Anything)
	adapter.On("Send", mock.Anything, mock.Anything)
	adapter.On("DeMultiplex", mock.Anything)
	adapter.On("GetConf").Return(conf)
	adapter.On("GetOrgOfPeer", pkiIDInOrg1).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", pkiIDinOrg2).Return(orgInChannelA)
	adapter.On("GetOrgOfPeer", mock.Anything).Return(api.OrgIdentityType(nil))

	waitForHandleMsgChan := make(chan struct{})

	adapter.On("GetMembership").Return([]discovery.NetworkMember{peerA}).Run(func(args mock.Arguments) {
		waitForHandleMsgChan <- struct{}{}
	}).Once()

	var check uint32
	atomic.StoreUint32(&check, 0)
	logger := util.GetLogger(util.ChannelLogger, adapter.GetConf().ID)
	logger = logger.(*flogging.FabricLogger).WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if atomic.LoadUint32(&check) == 1 {
			if !strings.Contains(entry.Message, "Membership view has changed. peers went offline:  [[a]] , peers went online:  [[b]] , current view:  [[b]]") {
				return nil
			}
			close(membershipReported)
			return nil
		}
		return nil
	}))

	gc := NewGossipChannel(pkiID1, orgInChannelA, cs, channelA, adapter, jcm, disabledMetrics, logger)

	gc.HandleMessage(&receivedMsg{PKIID: pkiIDInOrg1, msg: createStateInfoMsg(1, pkiIDInOrg1, channelA)})
	gc.HandleMessage(&receivedMsg{PKIID: pkiIDinOrg2, msg: createStateInfoMsg(1, pkiIDinOrg2, channelA)})
	<-waitForHandleMsgChan

	wg := sync.WaitGroup{}
	wg.Add(1)
	adapter.On("GetMembership").Return([]discovery.NetworkMember{peerB}).Run(func(args mock.Arguments) {
		defer wg.Done()
		gc.(*gossipChannel).Stop()
	}).Once()

	flogging.ActivateSpec("info")
	atomic.StoreUint32(&check, 1)
	<-membershipReported

	wg.Wait()
	adapter.On("GetMembership").Return([]discovery.NetworkMember{peerB}).Run(func(args mock.Arguments) {
		t.Fatalf("Membership tracker should have been stopped already.")
	})

	time.Sleep(conf.TimeForMembershipTracker * 2)
}
