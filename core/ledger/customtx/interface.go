/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package customtx

import (
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)



type InvalidTxError struct {
	Msg string
}

func (e *InvalidTxError) Error() string {
	return e.Msg
}










type Processor interface {
	GenerateSimulationResults(txEnvelop *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error
}


