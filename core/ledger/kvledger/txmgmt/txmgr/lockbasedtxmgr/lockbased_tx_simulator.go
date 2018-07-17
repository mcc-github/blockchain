/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"errors"
	"fmt"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/txmgr"
)


type lockBasedTxSimulator struct {
	lockBasedQueryExecutor
	rwsetBuilder              *rwsetutil.RWSetBuilder
	writePerformed            bool
	pvtdataQueriesPerformed   bool
	simulationResultsComputed bool
}

func newLockBasedTxSimulator(txmgr *LockBasedTxMgr, txid string) (*lockBasedTxSimulator, error) {
	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	helper := newQueryHelper(txmgr, rwsetBuilder)
	logger.Debugf("constructing new tx simulator txid = [%s]", txid)
	return &lockBasedTxSimulator{lockBasedQueryExecutor{helper, txid}, rwsetBuilder, false, false, false}, nil
}


func (s *lockBasedTxSimulator) SetState(ns string, key string, value []byte) error {
	if err := s.helper.checkDone(); err != nil {
		return err
	}
	if err := s.checkBeforeWrite(); err != nil {
		return err
	}
	if err := s.helper.txmgr.db.ValidateKeyValue(key, value); err != nil {
		return err
	}
	s.rwsetBuilder.AddToWriteSet(ns, key, value)
	return nil
}


func (s *lockBasedTxSimulator) DeleteState(ns string, key string) error {
	return s.SetState(ns, key, nil)
}


func (s *lockBasedTxSimulator) SetStateMultipleKeys(namespace string, kvs map[string][]byte) error {
	for k, v := range kvs {
		if err := s.SetState(namespace, k, v); err != nil {
			return err
		}
	}
	return nil
}


func (s *lockBasedTxSimulator) SetStateMetadata(namespace, key string, metadata map[string][]byte) error {
	return errors.New("not implemented")
}


func (s *lockBasedTxSimulator) DeleteStateMetadata(namespace, key string) error {
	return errors.New("not implemented")
}


func (s *lockBasedTxSimulator) SetPrivateData(ns, coll, key string, value []byte) error {
	if err := s.helper.validateCollName(ns, coll); err != nil {
		return err
	}
	if err := s.helper.checkDone(); err != nil {
		return err
	}
	if err := s.checkBeforeWrite(); err != nil {
		return err
	}
	if err := s.helper.txmgr.db.ValidateKeyValue(key, value); err != nil {
		return err
	}
	s.writePerformed = true
	s.rwsetBuilder.AddToPvtAndHashedWriteSet(ns, coll, key, value)
	return nil
}


func (s *lockBasedTxSimulator) DeletePrivateData(ns, coll, key string) error {
	return s.SetPrivateData(ns, coll, key, nil)
}


func (s *lockBasedTxSimulator) SetPrivateDataMultipleKeys(ns, coll string, kvs map[string][]byte) error {
	for k, v := range kvs {
		if err := s.SetPrivateData(ns, coll, k, v); err != nil {
			return err
		}
	}
	return nil
}


func (s *lockBasedTxSimulator) GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (commonledger.ResultsIterator, error) {
	if err := s.checkBeforePvtdataQueries(); err != nil {
		return nil, err
	}
	return s.lockBasedQueryExecutor.GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey)
}


func (s *lockBasedTxSimulator) SetPrivateDataMetadata(namespace, collection, key string, metadata map[string][]byte) error {
	return errors.New("not implemented")
}


func (s *lockBasedTxSimulator) DeletePrivateDataMetadata(namespace, collection, key string) error {
	return errors.New("not implemented")
}


func (s *lockBasedTxSimulator) ExecuteQueryOnPrivateData(namespace, collection, query string) (commonledger.ResultsIterator, error) {
	if err := s.checkBeforePvtdataQueries(); err != nil {
		return nil, err
	}
	return s.lockBasedQueryExecutor.ExecuteQueryOnPrivateData(namespace, collection, query)
}


func (s *lockBasedTxSimulator) GetTxSimulationResults() (*ledger.TxSimulationResults, error) {
	if s.simulationResultsComputed {
		return nil, errors.New("the function GetTxSimulationResults() should only be called once on a transaction simulator instance")
	}
	defer func() { s.simulationResultsComputed = true }()
	logger.Debugf("Simulation completed, getting simulation results")
	if s.helper.err != nil {
		return nil, s.helper.err
	}
	s.helper.addRangeQueryInfo()
	return s.rwsetBuilder.GetTxSimulationResults()
}


func (s *lockBasedTxSimulator) ExecuteUpdate(query string) error {
	return errors.New("Not supported")
}

func (s *lockBasedTxSimulator) checkBeforeWrite() error {
	if s.pvtdataQueriesPerformed {
		return &txmgr.ErrUnsupportedTransaction{
			Msg: fmt.Sprintf("Tx [%s]: Transaction has already performed queries on pvt data. Writes are not allowed", s.txid),
		}
	}
	s.writePerformed = true
	return nil
}

func (s *lockBasedTxSimulator) checkBeforePvtdataQueries() error {
	if s.writePerformed {
		return &txmgr.ErrUnsupportedTransaction{
			Msg: fmt.Sprintf("Tx [%s]: Queries on pvt data is supported only in a read-only transaction", s.txid),
		}
	}
	s.pvtdataQueriesPerformed = true
	return nil
}
