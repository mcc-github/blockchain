/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package token

import (
	"time"

	"github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/pkg/errors"
)

type IssueCmd struct {
	*BaseCmd
	clientConfigPath *string
	recipient        *string
	quantity         *string
	ttype            *string

	stub   Stub
	loader Loader
	parser ResponseParser
}

func NewIssueCmd(stub Stub, loader Loader, parser ResponseParser) *IssueCmd {
	return &IssueCmd{BaseCmd: &BaseCmd{}, stub: stub, loader: loader, parser: parser}
}


func (cmd *IssueCmd) SetClientConfigPath(clientConfigPath *string) {
	cmd.clientConfigPath = clientConfigPath
}


func (cmd *IssueCmd) SetRecipient(recipient *string) {
	cmd.recipient = recipient
}


func (cmd *IssueCmd) SetQuantity(quantity *string) {
	cmd.quantity = quantity
}


func (cmd *IssueCmd) SetType(ttype *string) {
	cmd.ttype = ttype
}

func (cmd *IssueCmd) Execute(conf common.Config) error {
	if cmd.clientConfigPath == nil || len(*cmd.clientConfigPath) == 0 {
		return errors.New("no client config path specified")
	}
	if cmd.recipient == nil || len(*cmd.recipient) == 0 {
		return errors.New("no recipient specified")
	}
	if cmd.quantity == nil || len(*cmd.quantity) == 0 {
		return errors.New("no quantity specified")
	}
	if cmd.ttype == nil || len(*cmd.ttype) == 0 {
		return errors.New("no type specified")
	}

	clientConfigPath := *cmd.clientConfigPath
	recipient := *cmd.recipient
	quantity := *cmd.quantity
	ttype := *cmd.ttype
	channel, mspPath, mspID := cmd.BaseCmd.GetArgs()

	
	recipientBytes, err := cmd.loader.TokenOwner(recipient)
	if err != nil {
		return err
	}

	tti := &token.Token{
		Owner:    recipientBytes,
		Quantity: quantity,
		Type:     ttype,
	}

	
	err = cmd.stub.Setup(clientConfigPath, channel, mspPath, mspID)
	if err != nil {
		return err
	}
	response, err := cmd.stub.Issue([]*token.Token{tti}, 30*time.Second)
	if err != nil {
		return err
	}

	return cmd.parser.ParseResponse(response)
}
