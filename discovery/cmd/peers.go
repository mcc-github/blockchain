/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/discovery/client"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)


func NewPeerCmd(stub Stub, parser ResponseParser) *PeerCmd {
	return &PeerCmd{
		stub:   stub,
		parser: parser,
	}
}


type PeerCmd struct {
	stub    Stub
	server  *string
	channel *string
	parser  ResponseParser
}


func (pc *PeerCmd) SetServer(server *string) {
	pc.server = server
}


func (pc *PeerCmd) SetChannel(channel *string) {
	pc.channel = channel
}


func (pc *PeerCmd) Execute(conf common.Config) error {
	channel := ""

	if pc.channel != nil {
		channel = *pc.channel
	}

	if pc.server == nil || *pc.server == "" {
		return errors.New("no server specified")
	}

	server := *pc.server

	req := discovery.NewRequest()
	if channel != "" {
		req = req.OfChannel(channel)
		req = req.AddPeersQuery()
	} else {
		req = req.AddLocalPeersQuery()
	}
	res, err := pc.stub.Send(server, conf, req)
	if err != nil {
		return err
	}
	return pc.parser.ParseResponse(channel, res)
}


type PeerResponseParser struct {
	io.Writer
}


func (parser *PeerResponseParser) ParseResponse(channel string, res ServiceResponse) error {
	var listPeers peerLister
	if channel == "" {
		listPeers = res.ForLocal()
	} else {
		listPeers = res.ForChannel(channel)
	}
	peers, err := listPeers.Peers()
	if err != nil {
		return err
	}

	var peerSlices []peer
	for _, p := range peers {
		peerSlices = append(peerSlices, rawPeerToPeer(p))
	}

	b, _ := json.MarshalIndent(peerSlices, "", "\t")
	fmt.Fprintln(parser.Writer, string(b))
	return nil
}

type peer struct {
	MSPID        string
	LedgerHeight uint64
	Endpoint     string
	Identity     string
	Chaincodes   []string
}

type peerLister interface {
	Peers() ([]*discovery.Peer, error)
}

func rawPeerToPeer(p *discovery.Peer) peer {
	var ledgerHeight uint64
	var ccs []string
	if p.StateInfoMessage != nil && p.StateInfoMessage.GetStateInfo() != nil && p.StateInfoMessage.GetStateInfo().Properties != nil {
		properties := p.StateInfoMessage.GetStateInfo().Properties
		ledgerHeight = properties.LedgerHeight
		for _, cc := range properties.Chaincodes {
			if cc == nil {
				continue
			}
			ccs = append(ccs, cc.Name)
		}
	}
	var endpoint string
	if p.AliveMessage != nil && p.AliveMessage.GetAliveMsg() != nil && p.AliveMessage.GetAliveMsg().Membership != nil {
		endpoint = p.AliveMessage.GetAliveMsg().Membership.Endpoint
	}
	sID := &msp.SerializedIdentity{}
	proto.Unmarshal(p.Identity, sID)
	return peer{
		MSPID:        p.MSPID,
		Endpoint:     endpoint,
		LedgerHeight: ledgerHeight,
		Identity:     string(sID.IdBytes),
		Chaincodes:   ccs,
	}
}
