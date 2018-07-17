/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"fmt"
	"time"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/filter"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)


type Gossip interface {

	
	SelfMembershipInfo() discovery.NetworkMember

	
	SelfChannelInfo(common.ChainID) *proto.SignedGossipMessage

	
	Send(msg *proto.GossipMessage, peers ...*comm.RemotePeer)

	
	SendByCriteria(*proto.SignedGossipMessage, SendCriteria) error

	
	Peers() []discovery.NetworkMember

	
	
	PeersOfChannel(common.ChainID) []discovery.NetworkMember

	
	
	UpdateMetadata(metadata []byte)

	
	
	UpdateLedgerHeight(height uint64, chainID common.ChainID)

	
	
	UpdateChaincodes(chaincode []*proto.Chaincode, chainID common.ChainID)

	
	Gossip(msg *proto.GossipMessage)

	
	
	PeerFilter(channel common.ChainID, messagePredicate api.SubChannelSelectionCriteria) (filter.RoutingFilter, error)

	
	
	
	
	Accept(acceptor common.MessageAcceptor, passThrough bool) (<-chan *proto.GossipMessage, <-chan proto.ReceivedMessage)

	
	JoinChan(joinMsg api.JoinChannelMessage, chainID common.ChainID)

	
	
	
	
	LeaveChan(chainID common.ChainID)

	
	
	SuspectPeers(s api.PeerSuspector)

	
	IdentityInfo() api.PeerIdentitySet

	
	Stop()
}



type emittedGossipMessage struct {
	*proto.SignedGossipMessage
	filter func(id common.PKIidType) bool
}


type SendCriteria struct {
	Timeout    time.Duration        
	MinAck     int                  
	MaxPeers   int                  
	IsEligible filter.RoutingFilter 
	Channel    common.ChainID       
	
}


func (sc SendCriteria) String() string {
	return fmt.Sprintf("channel: %s, tout: %v, minAck: %d, maxPeers: %d", sc.Channel, sc.Timeout, sc.MinAck, sc.MaxPeers)
}


type Config struct {
	BindPort            int      
	ID                  string   
	BootstrapPeers      []string 
	PropagateIterations int      
	PropagatePeerNum    int      

	MaxBlockCountToStore int 

	MaxPropagationBurstSize    int           
	MaxPropagationBurstLatency time.Duration 

	PullInterval time.Duration 
	PullPeerNum  int           

	SkipBlockVerification bool 

	PublishCertPeriod        time.Duration 
	PublishStateInfoInterval time.Duration 
	RequestStateInfoInterval time.Duration 

	TLSCerts *common.TLSCertificates 

	InternalEndpoint string 
	ExternalEndpoint string 
}
