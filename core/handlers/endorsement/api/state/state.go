/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
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
