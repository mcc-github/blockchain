/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)



type Processor struct {
	TMSManager TMSManager
}

func (p *Processor) GenerateSimulationResults(txEnv *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error {
	
	ch, ttx, ci, err := UnmarshalTokenTransaction(txEnv.Payload)
	if err != nil {
		return errors.WithMessage(err, "failed unmarshalling token transaction")
	}

	
	txProcessor, err := p.TMSManager.GetTxProcessor(ch.ChannelId)
	if err != nil {
		return errors.WithMessage(err, "failed getting committer")
	}

	
	err = txProcessor.ProcessTx(ch.TxId, ci, ttx, simulator)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed committing transaction for channel %s", ch.ChannelId))
	}

	return err
}
