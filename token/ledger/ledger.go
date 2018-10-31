/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import "github.com/mcc-github/blockchain/common/ledger"





type LedgerManager interface {
	
	GetLedgerReader(channel string) (LedgerReader, error)
}


type LedgerReader interface {
	
	GetState(namespace string, key string) ([]byte, error)

	
	
	
	
	
	GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ledger.ResultsIterator, error)
}




type LedgerWriter interface {
	LedgerReader
	
	SetState(namespace string, key string, value []byte) error
}




type ResultsIterator interface {
	
	
	Next() (ledger.QueryResult, error)
	
	Close()
}
