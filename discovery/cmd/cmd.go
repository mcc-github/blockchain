/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"os"
	"time"

	"github.com/mcc-github/blockchain/cmd/common"
	discovery "github.com/mcc-github/blockchain/discovery/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	PeersCommand     = "peers"
	ConfigCommand    = "config"
	EndorsersCommand = "endorsers"
)

var (
	
	responseParserWriter = os.Stdout
)

const (
	defaultTimeout = time.Second * 10
)




type Stub interface {
	
	Send(server string, conf common.Config, req *discovery.Request) (ServiceResponse, error)
}




type ResponseParser interface {
	
	ParseResponse(channel string, response ServiceResponse) error
}




type CommandRegistrar interface {
	
	Command(name, help string, onCommand common.CLICommand) *kingpin.CmdClause
}


func AddCommands(cli CommandRegistrar) {
	peerCmd := NewPeerCmd(&ClientStub{}, &PeerResponseParser{Writer: responseParserWriter})
	peers := cli.Command(PeersCommand, "Discover peers", peerCmd.Execute)
	server := peers.Flag("server", "Sets the endpoint of the server to connect").String()
	channel := peers.Flag("channel", "Sets the channel the query is intended to").String()
	peerCmd.SetServer(server)
	peerCmd.SetChannel(channel)

	configCmd := NewConfigCmd(&ClientStub{}, &ConfigResponseParser{Writer: responseParserWriter})
	config := cli.Command(ConfigCommand, "Discover channel config", configCmd.Execute)
	server = config.Flag("server", "Sets the endpoint of the server to connect").String()
	channel = config.Flag("channel", "Sets the channel the query is intended to").String()
	configCmd.SetServer(server)
	configCmd.SetChannel(channel)

	endorserCmd := NewEndorsersCmd(&RawStub{}, &EndorserResponseParser{Writer: responseParserWriter})
	endorsers := cli.Command(EndorsersCommand, "Discover chaincode endorsers", endorserCmd.Execute)
	chaincodes := endorsers.Flag("chaincode", "Specifies the chaincode name(s)").Strings()
	collections := endorsers.Flag("collection", "Specifies the collection name(s) as a mapping from chaincode to a comma separated list of collections").PlaceHolder("CC:C1,C2").StringMap()
	server = endorsers.Flag("server", "Sets the endpoint of the server to connect").String()
	channel = endorsers.Flag("channel", "Sets the channel the query is intended to").String()
	endorserCmd.SetChannel(channel)
	endorserCmd.SetServer(server)
	endorserCmd.SetChaincodes(chaincodes)
	endorserCmd.SetCollections(collections)
}
