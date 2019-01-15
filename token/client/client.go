/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"time"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	tk "github.com/mcc-github/blockchain/token"
)



type Prover interface {

	
	
	
	
	RequestImport(tokensToIssue []*token.TokenToIssue, signingIdentity tk.SigningIdentity) ([]byte, error)

	
	
	
	
	
	RequestTransfer(tokenIDs [][]byte, shares []*token.RecipientTransferShare, signingIdentity tk.SigningIdentity) ([]byte, error)
}



type FabricTxSubmitter interface {

	
	
	
	
	
	
	
	Submit(txEnvelope *common.Envelope, waitTimeout time.Duration) (*common.Status, bool, error)

	
	
	CreateTxEnvelope(tokenTx []byte) (*common.Envelope, string, error)
}


type Client struct {
	Config          *ClientConfig
	SigningIdentity tk.SigningIdentity
	Prover          Prover
	TxSubmitter     FabricTxSubmitter
}



func NewClient(config ClientConfig, signingIdentity tk.SigningIdentity) (*Client, error) {
	err := ValidateClientConfig(config)
	if err != nil {
		return nil, err
	}

	prover, err := NewProverPeer(&config)
	if err != nil {
		return nil, err
	}

	txSubmitter, err := NewTxSubmitter(&config, signingIdentity)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:          &config,
		SigningIdentity: signingIdentity,
		Prover:          prover,
		TxSubmitter:     txSubmitter,
	}, nil
}











func (c *Client) Issue(tokensToIssue []*token.TokenToIssue, waitTimeout time.Duration) (*common.Envelope, string, *common.Status, bool, error) {
	serializedTokenTx, err := c.Prover.RequestImport(tokensToIssue, c.SigningIdentity)
	if err != nil {
		return nil, "", nil, false, err
	}

	txEnvelope, txid, err := c.TxSubmitter.CreateTxEnvelope(serializedTokenTx)
	if err != nil {
		return nil, "", nil, false, err
	}

	ordererStatus, committed, err := c.TxSubmitter.Submit(txEnvelope, waitTimeout)
	return txEnvelope, txid, ordererStatus, committed, err
}












func (c *Client) Transfer(tokenIDs [][]byte, shares []*token.RecipientTransferShare, waitTimeout time.Duration) (*common.Envelope, string, *common.Status, bool, error) {
	serializedTokenTx, err := c.Prover.RequestTransfer(tokenIDs, shares, c.SigningIdentity)
	if err != nil {
		return nil, "", nil, false, err
	}

	txEnvelope, txid, err := c.TxSubmitter.CreateTxEnvelope(serializedTokenTx)
	if err != nil {
		return nil, "", nil, false, err
	}

	ordererStatus, committed, err := c.TxSubmitter.Submit(txEnvelope, waitTimeout)
	return txEnvelope, txid, ordererStatus, committed, err
}
