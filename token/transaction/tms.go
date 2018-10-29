/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/ledger"
)








type TMSTxProcessor interface {
	
	ProcessTx(txID string, creator identity.PublicInfo, ttx *token.TokenTransaction, simulator ledger.LedgerWriter) error
}

type TMSManager interface {
	
	GetTxProcessor(channel string) (TMSTxProcessor, error)
}

type TxCreatorInfo struct {
	public []byte
}

func (t *TxCreatorInfo) Public() []byte {
	return t.public
}
