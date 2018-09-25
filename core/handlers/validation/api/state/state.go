/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"github.com/mcc-github/blockchain/core/handlers/validation/api"
)


type State interface {
	
	GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error)

	
	
	
	
	
	GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ResultsIterator, error)

	
	GetStateMetadata(namespace, key string) (map[string][]byte, error)

	
	GetPrivateDataMetadataByHash(namespace, collection string, keyhash []byte) (map[string][]byte, error)

	
	Done()
}


type StateFetcher interface {
	validation.Dependency

	
	FetchState() (State, error)
}


type ResultsIterator interface {
	
	
	Next() (QueryResult, error)
	
	Close()
}


type QueryResult interface{}
