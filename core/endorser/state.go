/*
Copyright IBM Corp. 2018 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/pkg/errors"
)



type QueryCreator interface {
	NewQueryExecutor() (ledger.QueryExecutor, error)
}


type ChannelState struct {
	transientstore.Store
	QueryCreator
}


func (cs *ChannelState) FetchState() (endorsement.State, error) {
	qe, err := cs.NewQueryExecutor()
	if err != nil {
		return nil, err
	}

	return &StateContext{
		QueryExecutor: qe,
		Store:         cs.Store,
	}, nil
}


type StateContext struct {
	transientstore.Store
	ledger.QueryExecutor
}


func (sc *StateContext) GetTransientByTXID(txID string) ([]*rwset.TxPvtReadWriteSet, error) {
	scanner, err := sc.Store.GetTxPvtRWSetByTxid(txID, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer scanner.Close()
	var data []*rwset.TxPvtReadWriteSet
	for {
		res, err := scanner.Next()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if res == nil {
			break
		}
		if res.PvtSimulationResultsWithConfig == nil {
			continue
		}
		data = append(data, res.PvtSimulationResultsWithConfig.PvtRwset)
	}
	return data, nil
}
