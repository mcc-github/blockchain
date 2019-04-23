/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package token

import (
	"os"
	"time"

	cmdcommon "github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/protos/token"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	IssueCommand     = "issue"
	TransferCommand  = "transfer"
	ListTokensCommad = "list"
	RedeemCommand    = "redeem"
)

var (
	
	responseParserWriter = os.Stdout
)




type CommandRegistrar interface {
	
	Command(name, help string, onCommand cmdcommon.CLICommand) *kingpin.CmdClause
}




type ResponseParser interface {
	
	ParseResponse(response StubResponse) error
}

type StubResponse interface {
}




type Stub interface {
	
	Setup(configFilePath, channel, mspPath, mspID string) error

	
	Issue(tokensToIssue []*token.Token, waitTimeout time.Duration) (StubResponse, error)

	
	Transfer(tokenIDs []*token.TokenId, shares []*token.RecipientShare, waitTimeout time.Duration) (StubResponse, error)

	
	Redeem(tokenIDs []*token.TokenId, quantity string, waitTimeout time.Duration) (StubResponse, error)

	
	ListTokens() (StubResponse, error)
}




type Loader interface {
	
	TokenOwner(s string) (*token.TokenOwner, error)

	
	TokenIDs(s string) ([]*token.TokenId, error)

	
	Shares(s string) ([]*token.RecipientShare, error)
}


type BaseCmd struct {
	channel *string
	mspPath *string
	mspID   *string
}


func (cmd *BaseCmd) SetChannel(channel *string) {
	cmd.channel = channel
}


func (cmd *BaseCmd) SetMSPPath(mspPath *string) {
	cmd.mspPath = mspPath
}


func (cmd *BaseCmd) SetMSPId(mspID *string) {
	cmd.mspID = mspID
}


func (cmd *BaseCmd) GetArgs() (string, string, string) {
	var channel string
	if cmd.channel == nil {
		channel = ""
	} else {
		channel = *cmd.channel
	}
	var mspPath string
	if cmd.mspPath == nil {
		mspPath = ""
	} else {
		mspPath = *cmd.mspPath
	}
	var mspID string
	if cmd.mspID == nil {
		mspID = ""
	} else {
		mspID = *cmd.mspID
	}
	return channel, mspPath, mspID
}

func addBaseFlags(cli *kingpin.CmdClause, baseCmd *BaseCmd) {
	channel := cli.Flag("channel", "Overrides channel configuration").String()
	mspPath := cli.Flag("mspPath", "Overrides msp path configuration").String()
	mspID := cli.Flag("mspId", "Overrides msp id configuration").String()

	baseCmd.SetChannel(channel)
	baseCmd.SetMSPPath(mspPath)
	baseCmd.SetMSPId(mspID)
}


func AddCommands(cli CommandRegistrar) {
	
	issueCmd := NewIssueCmd(&TokenClientStub{}, &JsonLoader{}, &OperationResponseParser{responseParserWriter})
	importCli := cli.Command(IssueCommand, "Import token command", issueCmd.Execute)
	addBaseFlags(importCli, issueCmd.BaseCmd)
	configPath := importCli.Flag("config", "Sets the client configuration path").String()
	ttype := importCli.Flag("type", "Sets the token type to issue").String()
	quantity := importCli.Flag("quantity", "Sets the quantity of tokens to issue").String()
	recipient := importCli.Flag("recipient", "Sets the recipient of tokens to issue").String()

	issueCmd.SetClientConfigPath(configPath)
	issueCmd.SetType(ttype)
	issueCmd.SetQuantity(quantity)
	issueCmd.SetRecipient(recipient)

	
	listTokensCmd := NewListTokensCmd(&TokenClientStub{}, &UnspentTokenResponseParser{responseParserWriter})
	listTokensCli := cli.Command(ListTokensCommad, "List tokens command", listTokensCmd.Execute)
	addBaseFlags(listTokensCli, listTokensCmd.BaseCmd)
	configPath = listTokensCli.Flag("config", "Sets the client configuration path").String()
	listTokensCmd.SetClientConfigPath(configPath)

	
	transferCmd := NewTransferCmd(&TokenClientStub{}, &JsonLoader{}, &OperationResponseParser{responseParserWriter})
	transferCli := cli.Command(TransferCommand, "Transfer tokens command", transferCmd.Execute)
	addBaseFlags(transferCli, transferCmd.BaseCmd)
	configPath = transferCli.Flag("config", "Sets the client configuration path").String()
	tokenIDs := transferCli.Flag("tokenIDs", "Sets the token IDs to transfer").String()
	shares := transferCli.Flag("shares", "Sets the shares of the recipients").String()
	transferCmd.SetClientConfigPath(configPath)
	transferCmd.SetTokenIDs(tokenIDs)
	transferCmd.SetShares(shares)

	
	redeemCmd := NewRedeemCmd(&TokenClientStub{}, &JsonLoader{}, &OperationResponseParser{responseParserWriter})
	redeemCli := cli.Command(RedeemCommand, "Redeem tokens command", redeemCmd.Execute)
	addBaseFlags(redeemCli, redeemCmd.BaseCmd)
	configPath = redeemCli.Flag("config", "Sets the client configuration path").String()
	tokenIDs = redeemCli.Flag("tokenIDs", "Sets the token IDs to redeem").String()
	quantity = redeemCli.Flag("quantity", "Sets the quantity of tokens to redeem").String()
	redeemCmd.SetClientConfigPath(configPath)
	redeemCmd.SetTokenIDs(tokenIDs)
	redeemCmd.SetQuantity(quantity)
}
