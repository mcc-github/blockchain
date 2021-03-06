/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import (
	"github.com/mcc-github/blockchain/protos/token"
)




type Issuer interface {
	
	RequestIssue(tokensToIssue []*token.Token) (*token.TokenTransaction, error)

	
	RequestTokenOperation(op *token.TokenOperation) (*token.TokenTransaction, error)
}




type Transactor interface {
	
	
	
	RequestTransfer(request *token.TransferRequest) (*token.TokenTransaction, error)

	
	
	
	
	RequestRedeem(request *token.RedeemRequest) (*token.TokenTransaction, error)

	
	ListTokens() (*token.UnspentTokens, error)

	
	RequestTokenOperation(tokenIDs []*token.TokenId, op *token.TokenOperation) (*token.TokenTransaction, int, error)

	
	Done()
}



type TMSManager interface {
	
	
	GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error)

	
	
	GetTransactor(channel string, privateCredential, publicCredential []byte) (Transactor, error)
}
