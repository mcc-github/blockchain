/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"github.com/mcc-github/blockchain-protos-go/discovery"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	
	ErrNotFound = errors.New("not found")
)



type Signer func(msg []byte) ([]byte, error)


type Dialer func() (*grpc.ClientConn, error)


type Response interface {
	
	ForChannel(string) ChannelResponse

	
	ForLocal() LocalResponse
}


type ChannelResponse interface {
	
	Config() (*discovery.ConfigResult, error)

	
	Peers(invocationChain ...*discovery.ChaincodeCall) ([]*Peer, error)

	
	
	
	
	
	
	
	
	Endorsers(invocationChain InvocationChain, f Filter) (Endorsers, error)
}


type LocalResponse interface {
	
	Peers() ([]*Peer, error)
}



type Endorsers []*Peer



type Peer struct {
	MSPID            string
	AliveMessage     *protoext.SignedGossipMessage
	StateInfoMessage *protoext.SignedGossipMessage
	Identity         []byte
}
