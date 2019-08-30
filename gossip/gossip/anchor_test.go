/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/protoext"
	utilgossip "github.com/mcc-github/blockchain/gossip/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type peerMock struct {
	pkiID                common.PKIidType
	selfCertHash         []byte
	gRGCserv             *grpc.Server
	finishedSignal       sync.WaitGroup
	expectedMsgs2Receive uint32
	msgReceivedCount     uint32
	msgAssertions        []msgInspection
	t                    *testing.T
}

func (p *peerMock) GossipStream(stream proto.Gossip_GossipStreamServer) error {
	sessionCounter := 0
	for {
		envelope, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		gMsg, err := protoext.EnvelopeToGossipMessage(envelope)
		if err != nil {
			panic(err)
		}
		if sessionCounter == 0 {
			connEstablishMsg := p.connEstablishMsg(p.pkiID, p.selfCertHash, api.PeerIdentityType(p.pkiID))
			stream.Send(connEstablishMsg.Envelope)
		}
		for _, assertion := range p.msgAssertions {
			assertion(p.t, sessionCounter, &receivedMsg{stream: stream, SignedGossipMessage: gMsg})
		}
		p.t.Log("sessionCounter:", sessionCounter, string(p.pkiID), "got msg:", gMsg)
		sessionCounter++
		atomic.AddUint32(&p.msgReceivedCount, uint32(1))
		if atomic.LoadUint32(&p.msgReceivedCount) == p.expectedMsgs2Receive {
			p.finishedSignal.Done()
		}
	}
}

func (p *peerMock) Ping(context.Context, *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func newPeerMockWithGRPC(port int, gRPCServer *comm.GRPCServer, certs *common.TLSCertificates,
	expectedMsgs2Receive int, t *testing.T, msgAssertions ...msgInspection) *peerMock {
	p := &peerMock{
		gRGCserv:             gRPCServer.Server(),
		msgAssertions:        msgAssertions,
		t:                    t,
		pkiID:                common.PKIidType(fmt.Sprintf("127.0.0.1:%d", port)),
		selfCertHash:         util.ComputeSHA256(certs.TLSServerCert.Load().(*tls.Certificate).Certificate[0]),
		expectedMsgs2Receive: uint32(expectedMsgs2Receive),
	}
	p.finishedSignal.Add(1)
	proto.RegisterGossipServer(gRPCServer.Server(), p)
	go func() {
		gRPCServer.Start()
	}()
	return p
}

func (p *peerMock) connEstablishMsg(pkiID common.PKIidType, hash []byte, cert api.PeerIdentityType) *protoext.SignedGossipMessage {
	m := &proto.GossipMessage{
		Tag:   proto.GossipMessage_EMPTY,
		Nonce: 0,
		Content: &proto.GossipMessage_Conn{
			Conn: &proto.ConnEstablish{
				TlsCertHash: hash,
				Identity:    cert,
				PkiId:       pkiID,
			},
		},
	}
	gMsg := &protoext.SignedGossipMessage{
		GossipMessage: m,
	}
	gMsg.Sign((&configurableCryptoService{}).Sign)
	return gMsg
}

func (p *peerMock) stop() {
	p.gRGCserv.Stop()
}

type receivedMsg struct {
	*protoext.SignedGossipMessage
	stream proto.Gossip_GossipStreamServer
}

func (msg *receivedMsg) respond(message *protoext.SignedGossipMessage) {
	msg.stream.Send(message.Envelope)
}

func memResp(nonce uint64, endpoint string) *protoext.SignedGossipMessage {
	fakePeerAliveMsg := &protoext.SignedGossipMessage{
		GossipMessage: &proto.GossipMessage{
			Tag: proto.GossipMessage_EMPTY,
			Content: &proto.GossipMessage_AliveMsg{
				AliveMsg: &proto.AliveMessage{
					Membership: &proto.Member{
						Endpoint: endpoint,
						PkiId:    []byte(endpoint),
					},
					Identity: []byte(endpoint),
					Timestamp: &proto.PeerTime{
						IncNum: uint64(time.Now().UnixNano()),
						SeqNum: 0,
					},
				},
			},
		},
	}

	m, _ := fakePeerAliveMsg.Sign((&configurableCryptoService{}).Sign)
	sMsg, _ := protoext.NoopSign(&proto.GossipMessage{
		Tag:   proto.GossipMessage_EMPTY,
		Nonce: nonce,
		Content: &proto.GossipMessage_MemRes{
			MemRes: &proto.MembershipResponse{
				Alive: []*proto.Envelope{m},
				Dead:  []*proto.Envelope{},
			},
		},
	})
	return sMsg
}

type msgInspection func(t *testing.T, index int, m *receivedMsg)

func TestAnchorPeer(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	cs := &configurableCryptoService{m: make(map[string]api.OrgIdentityType)}
	orgA := "orgA"
	orgB := "orgB"

	port, grpc, cert, secDialOpt, _ := utilgossip.CreateGRPCLayer()
	port1, grpc1, cert1, _, _ := utilgossip.CreateGRPCLayer()
	port2, grpc2, cert2, _, _ := utilgossip.CreateGRPCLayer()
	port3, grpc3, cert3, _, _ := utilgossip.CreateGRPCLayer()
	port4, grpc4, cert4, _, _ := utilgossip.CreateGRPCLayer()

	cs.putInOrg(port, orgA)  
	cs.putInOrg(port1, orgA) 
	cs.putInOrg(port2, orgB) 
	cs.putInOrg(port3, orgA) 
	cs.putInOrg(port4, orgB) 

	
	handshake := func(t *testing.T, index int, m *receivedMsg) {
		if index != 0 {
			return
		}
		assert.NotNil(t, m.GetConn())
	}

	memReqWithInternalEndpoint := func(t *testing.T, index int, m *receivedMsg) {
		if m.GetMemReq() == nil {
			return
		}
		assert.True(t, index > 0)
		req := m.GetMemReq()
		am, err := protoext.EnvelopeToGossipMessage(req.SelfInformation)
		assert.NoError(t, err)
		assert.NotEmpty(t, protoext.InternalEndpoint(am.GetSecretEnvelope()))
		m.respond(memResp(m.Nonce, fmt.Sprintf("127.0.0.1:%d", port3)))
	}

	memReqWithoutInternalEndpoint := func(t *testing.T, index int, m *receivedMsg) {
		if m.GetMemReq() == nil {
			return
		}
		assert.True(t, index > 0)
		req := m.GetMemReq()
		am, err := protoext.EnvelopeToGossipMessage(req.SelfInformation)
		assert.NoError(t, err)
		assert.Nil(t, am.GetSecretEnvelope())
		m.respond(memResp(m.Nonce, fmt.Sprintf("127.0.0.1:%d", port4)))
	}

	
	pm1 := newPeerMockWithGRPC(port3, grpc3, cert3, 1, t, handshake)
	defer pm1.stop()
	pm2 := newPeerMockWithGRPC(port4, grpc4, cert4, 1, t, handshake)
	defer pm2.stop()
	jcm := &joinChanMsg{
		members2AnchorPeers: map[string][]api.AnchorPeer{
			orgA: {
				{Host: "127.0.0.1", Port: port1},
			},
			orgB: {
				{Host: "127.0.0.1", Port: port2},
			},
		},
	}
	channel := common.ChannelID("TEST")
	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	
	p := newGossipInstanceWithGRPCWithExternalEndpoint(0, port, grpc, cert, secDialOpt, cs, endpoint)
	defer p.Stop()
	p.JoinChan(jcm, channel)
	p.UpdateLedgerHeight(1, channel)

	time.Sleep(time.Second * 5)

	
	ap1 := newPeerMockWithGRPC(port1, grpc1, cert1, 3, t, handshake, memReqWithInternalEndpoint)
	defer ap1.stop()
	ap2 := newPeerMockWithGRPC(port2, grpc2, cert2, 3, t, handshake, memReqWithoutInternalEndpoint)
	defer ap2.stop()

	
	ap1.finishedSignal.Wait()
	ap2.finishedSignal.Wait()
	pm1.finishedSignal.Wait()
	pm2.finishedSignal.Wait()
}

func TestBootstrapPeerMisConfiguration(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	
	

	cs := &configurableCryptoService{m: make(map[string]api.OrgIdentityType)}
	orgA := "orgA"
	orgB := "orgB"

	port, grpc, cert, secDialOpt, _ := utilgossip.CreateGRPCLayer()
	fmt.Printf("port %d\n", port)
	port1, grpc1, cert1, secDialOpt, _ := utilgossip.CreateGRPCLayer()
	fmt.Printf("port1 %d\n", port1)
	port2, grpc2, cert2, secDialOpt, _ := utilgossip.CreateGRPCLayer()
	fmt.Printf("port2 %d\n", port2)

	cs.putInOrg(port, orgA)
	cs.putInOrg(port1, orgB)
	cs.putInOrg(port2, orgA)

	onlyHandshakes := func(t *testing.T, index int, m *receivedMsg) {
		
		
		assert.NotNil(t, m.GetConn())
		
		
		
		assert.Nil(t, m.GetMemReq())
	}
	
	bs1 := newPeerMockWithGRPC(port1, grpc1, cert1, 3, t, onlyHandshakes)
	defer bs1.stop()

	membershipRequestsSent := make(chan struct{}, 100)
	detectMembershipRequest := func(t *testing.T, index int, m *receivedMsg) {
		if m.GetMemReq() != nil {
			membershipRequestsSent <- struct{}{}
		}
	}

	bs2 := newPeerMockWithGRPC(port2, grpc2, cert2, 0, t, detectMembershipRequest)
	defer bs2.stop()

	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	p := newGossipInstanceWithGRPCWithExternalEndpoint(0, port, grpc, cert, secDialOpt, cs, endpoint, port1, port2)
	defer p.Stop()

	
	
	got3Handshakes := make(chan struct{})
	go func() {
		bs1.finishedSignal.Wait()
		got3Handshakes <- struct{}{}
	}()

	select {
	case <-got3Handshakes:
	case <-time.After(time.Second * 15):
		assert.Fail(t, "Didn't detect 3 handshake attempts to the bootstrap peer from orgB")
	}

	select {
	case <-membershipRequestsSent:
	case <-time.After(time.Second * 15):
		assert.Fail(t, "Bootstrap peer didn't receive a membership request from the peer within a timely manner")
	}
}
