/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger


type LedgerReader interface {
	
	GetState(namespace string, key string) ([]byte, error)
}




type LedgerWriter interface {
	LedgerReader
	
	SetState(namespace string, key string, value []byte) error
}
