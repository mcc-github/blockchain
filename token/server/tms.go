/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import "github.com/mcc-github/blockchain/protos/token"




type Issuer interface {
	
	RequestImport(tokensToIssue []*token.TokenToIssue) (*token.TokenTransaction, error)
}



type TMSManager interface {
	
	
	GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error)
}
