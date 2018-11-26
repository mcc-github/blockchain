/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import (
	"github.com/mcc-github/blockchain/protos/token"
)




type Issuer interface {
	
	RequestImport(tokensToIssue []*token.TokenToIssue) (*token.TokenTransaction, error)
}




type Transactor interface {
	
	
	
	RequestTransfer(request *token.TransferRequest) (*token.TokenTransaction, error)

	
	
	
	
	RequestRedeem(request *token.RedeemRequest) (*token.TokenTransaction, error)

	
	ListTokens() (*token.UnspentTokens, error)

	
	
	RequestApprove(request *token.ApproveRequest) (*token.TokenTransaction, error)

	
	
	
	RequestTransferFrom(request *token.TransferRequest) (*token.TokenTransaction, error)
}



type TMSManager interface {
	
	
	GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error)

	
	
	GetTransactor(channel string, privateCredential, publicCredential []byte) (Transactor, error)
}
