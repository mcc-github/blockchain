/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	. "github.com/mcc-github/blockchain-protos-go/discovery"
	"github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/cmd/common"
	discovery "github.com/mcc-github/blockchain/discovery/client"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/pkg/errors"
)


func NewEndorsersCmd(stub Stub, parser ResponseParser) *EndorsersCmd {
	return &EndorsersCmd{
		stub:   stub,
		parser: parser,
	}
}


type EndorsersCmd struct {
	stub        Stub
	server      *string
	channel     *string
	chaincodes  *[]string
	collections *map[string]string
	parser      ResponseParser
}


func (pc *EndorsersCmd) SetCollections(collections *map[string]string) {
	pc.collections = collections
}


func (pc *EndorsersCmd) SetChaincodes(chaincodes *[]string) {
	pc.chaincodes = chaincodes
}


func (pc *EndorsersCmd) SetServer(server *string) {
	pc.server = server
}


func (pc *EndorsersCmd) SetChannel(channel *string) {
	pc.channel = channel
}


func (pc *EndorsersCmd) Execute(conf common.Config) error {
	if pc.channel == nil || *pc.channel == "" {
		return errors.New("no channel specified")
	}

	if pc.server == nil || *pc.server == "" {
		return errors.New("no server specified")
	}

	server := *pc.server
	channel := *pc.channel

	ccAndCol := &chaincodesAndCollections{
		Chaincodes:  pc.chaincodes,
		Collections: pc.collections,
	}
	cc2collections, err := ccAndCol.parseInput()
	if err != nil {
		return err
	}

	var ccCalls []*ChaincodeCall

	for _, cc := range *ccAndCol.Chaincodes {
		ccCalls = append(ccCalls, &ChaincodeCall{
			Name:            cc,
			CollectionNames: cc2collections[cc],
		})
	}

	req, err := discovery.NewRequest().OfChannel(channel).AddEndorsersQuery(&ChaincodeInterest{Chaincodes: ccCalls})
	if err != nil {
		return errors.Wrap(err, "failed creating request")
	}

	res, err := pc.stub.Send(server, conf, req)
	if err != nil {
		return err
	}

	return pc.parser.ParseResponse(channel, res)
}


type EndorserResponseParser struct {
	io.Writer
}


func (parser *EndorserResponseParser) ParseResponse(channel string, res ServiceResponse) error {
	rawResponse := res.Raw()
	if len(rawResponse.Results) == 0 {
		return errors.New("empty results")
	}

	if e := rawResponse.Results[0].GetError(); e != nil {
		return errors.Errorf("server returned: %s", e.Content)
	}

	ccQueryRes := rawResponse.Results[0].GetCcQueryRes()
	if ccQueryRes == nil {
		return errors.Errorf("server returned response of unexpected type: %v", reflect.TypeOf(rawResponse.Results[0]))
	}

	jsonBytes, _ := json.MarshalIndent(parseEndorsementDescriptors(ccQueryRes.Content), "", "\t")
	fmt.Fprintln(parser.Writer, string(jsonBytes))
	return nil
}

type chaincodesAndCollections struct {
	Chaincodes  *[]string
	Collections *map[string]string
}

func (ec *chaincodesAndCollections) existsInChaincodes(chaincodeName string) bool {
	for _, cc := range *ec.Chaincodes {
		if chaincodeName == cc {
			return true
		}
	}
	return false
}

func (ec *chaincodesAndCollections) parseInput() (map[string][]string, error) {
	var emptyChaincodes []string
	if ec.Chaincodes == nil {
		ec.Chaincodes = &emptyChaincodes
	}
	var emptyCollections map[string]string
	if ec.Collections == nil {
		ec.Collections = &emptyCollections
	}

	res := make(map[string][]string)

	for _, cc := range *ec.Chaincodes {
		res[cc] = nil
	}

	for cc, collections := range *ec.Collections {
		if !ec.existsInChaincodes(cc) {
			return nil, errors.Errorf("a collection specified chaincode %s but it wasn't specified with a chaincode flag", cc)
		}
		res[cc] = strings.Split(collections, ",")
	}
	return res, nil
}

func parseEndorsementDescriptors(descriptors []*EndorsementDescriptor) []endorsermentDescriptor {
	var res []endorsermentDescriptor
	for _, desc := range descriptors {
		endorsersByGroups := make(map[string][]endorser)
		for grp, endorsers := range desc.EndorsersByGroups {
			for _, p := range endorsers.Peers {
				endorsersByGroups[grp] = append(endorsersByGroups[grp], endorserFromRaw(p))
			}
		}
		res = append(res, endorsermentDescriptor{
			Chaincode:         desc.Chaincode,
			Layouts:           desc.Layouts,
			EndorsersByGroups: endorsersByGroups,
		})
	}
	return res
}

type endorser struct {
	MSPID        string
	LedgerHeight uint64
	Endpoint     string
	Identity     string
}

type endorsermentDescriptor struct {
	Chaincode         string
	EndorsersByGroups map[string][]endorser
	Layouts           []*Layout
}

func endorserFromRaw(p *Peer) endorser {
	sId := &msp.SerializedIdentity{}
	proto.Unmarshal(p.Identity, sId)
	return endorser{
		MSPID:        sId.Mspid,
		Endpoint:     endpointFromEnvelope(p.MembershipInfo),
		LedgerHeight: ledgerHeightFromEnvelope(p.StateInfo),
		Identity:     string(sId.IdBytes),
	}
}

func endpointFromEnvelope(env *gossip.Envelope) string {
	if env == nil {
		return ""
	}
	aliveMsg, _ := protoext.EnvelopeToGossipMessage(env)
	if aliveMsg == nil {
		return ""
	}
	if !protoext.IsAliveMsg(aliveMsg.GossipMessage) {
		return ""
	}
	if aliveMsg.GetAliveMsg().Membership == nil {
		return ""
	}
	return aliveMsg.GetAliveMsg().Membership.Endpoint
}

func ledgerHeightFromEnvelope(env *gossip.Envelope) uint64 {
	if env == nil {
		return 0
	}
	stateInfoMsg, _ := protoext.EnvelopeToGossipMessage(env)
	if stateInfoMsg == nil {
		return 0
	}
	if !protoext.IsStateInfoMsg(stateInfoMsg.GossipMessage) {
		return 0
	}
	if stateInfoMsg.GetStateInfo().Properties == nil {
		return 0
	}
	return stateInfoMsg.GetStateInfo().Properties.LedgerHeight
}
