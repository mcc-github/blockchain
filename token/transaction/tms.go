/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"github.com/mcc-github/blockchain/protos/token"
)








type TMSTxProcessor interface {
	
	ProcessTx(txID string, creator CreatorInfo, ttx *token.TokenTransaction, simulator LedgerWriter) error
}

type TMSManager interface {
	
	GetTxProcessor(channel string) (TMSTxProcessor, error)
	
	SetPolicyValidator(channel string, validator PolicyValidator)
}




type CreatorInfo interface {
	Public() []byte
}

type TxCreatorInfo struct {
	public []byte
}

func (t *TxCreatorInfo) Public() []byte {
	return t.public
}




type PolicyValidator interface {
	
	IsIssuer(creator CreatorInfo, tokenType string) error
}


type LedgerReader interface {
	
	GetState(namespace string, key string) ([]byte, error)
}




type LedgerWriter interface {
	LedgerReader
	
	SetState(namespace string, key string, value []byte) error
}
