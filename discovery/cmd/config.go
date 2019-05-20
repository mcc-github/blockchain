/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/mcc-github/blockchain/cmd/common"
	discovery "github.com/mcc-github/blockchain/discovery/client"
	"github.com/pkg/errors"
)


func NewConfigCmd(stub Stub, parser ResponseParser) *ConfigCmd {
	return &ConfigCmd{
		stub:   stub,
		parser: parser,
	}
}


type ConfigCmd struct {
	stub    Stub
	server  *string
	channel *string
	parser  ResponseParser
}


func (pc *ConfigCmd) SetServer(server *string) {
	pc.server = server
}


func (pc *ConfigCmd) SetChannel(channel *string) {
	pc.channel = channel
}


func (pc *ConfigCmd) Execute(conf common.Config) error {
	if pc.server == nil || *pc.server == "" {
		return errors.New("no server specified")
	}
	if pc.channel == nil || *pc.channel == "" {
		return errors.New("no channel specified")
	}

	server := *pc.server
	channel := *pc.channel

	req := discovery.NewRequest().OfChannel(channel).AddConfigQuery()
	res, err := pc.stub.Send(server, conf, req)
	if err != nil {
		return err
	}
	return pc.parser.ParseResponse(channel, res)
}


type ConfigResponseParser struct {
	io.Writer
}


func (parser *ConfigResponseParser) ParseResponse(channel string, res ServiceResponse) error {
	chanConf, err := res.ForChannel(channel).Config()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.MarshalIndent(chanConf, "", "\t")
	fmt.Fprintln(parser.Writer, string(jsonBytes))
	return nil
}
