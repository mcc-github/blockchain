/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	tk "github.com/mcc-github/blockchain/token"
)



type Prover interface {

	
	
	
	
	RequestImport(tokensToIssue []*token.TokenToIssue, signingIdentity tk.SigningIdentity) ([]byte, error)

	
	
	
	
	
	RequestTransfer(tokenIDs [][]byte, shares []*token.RecipientTransferShare, signingIdentity tk.SigningIdentity) ([]byte, error)
}



type FabricTxSubmitter interface {

	
	
	
	
	Submit(tx []byte) error
}


type Client struct {
	SigningIdentity tk.SigningIdentity
	Prover          Prover
	TxSubmitter     FabricTxSubmitter
}





func (c *Client) Issue(tokensToIssue []*token.TokenToIssue) ([]byte, error) {
	serializedTokenTx, err := c.Prover.RequestImport(tokensToIssue, c.SigningIdentity)
	if err != nil {
		return nil, err
	}

	tx, err := c.createTx(serializedTokenTx)
	if err != nil {
		return nil, err
	}

	return tx, c.TxSubmitter.Submit(tx)
}




func (c *Client) Transfer(tokenIDs [][]byte, shares []*token.RecipientTransferShare) ([]byte, error) {
	serializedTokenTx, err := c.Prover.RequestTransfer(tokenIDs, shares, c.SigningIdentity)
	if err != nil {
		return nil, err
	}
	tx, err := c.createTx(serializedTokenTx)
	if err != nil {
		return nil, err
	}

	return tx, c.TxSubmitter.Submit(tx)
}



func (c *Client) createTx(tokenTx []byte) ([]byte, error) {
	payload := &common.Payload{Data: tokenTx}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}
	signature, err := c.SigningIdentity.Sign(payloadBytes)
	if err != nil {
		return nil, err
	}
	envelope := &common.Envelope{Payload: payloadBytes, Signature: signature}
	return proto.Marshal(envelope)
}
