/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
)


type State interface {
	
	GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([][]byte, error)

	
	GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error)

	
	GetTransientByTXID(txID string) ([]*rwset.TxPvtReadWriteSet, error)

	
	Done()
}


type StateFetcher interface {
	endorsement.Dependency

	
	FetchState() (State, error)
}
