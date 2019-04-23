/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"github.com/mcc-github/blockchain/common/ledger"
)



type MemoryLedger struct {
	entries map[string][]byte
}


func NewMemoryLedger() *MemoryLedger {
	return &MemoryLedger{
		entries: make(map[string][]byte),
	}
}


func (p *MemoryLedger) GetState(namespace string, key string) ([]byte, error) {
	value := p.entries[key]

	return value, nil
}


func (p *MemoryLedger) SetState(namespace string, key string, value []byte) error {
	p.entries[key] = value
	return nil
}


func (p *MemoryLedger) DeleteState(namespace string, key string) error {
	delete(p.entries, key)
	return nil
}



func (p *MemoryLedger) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ledger.ResultsIterator, error) {
	return nil, nil
}


func (p *MemoryLedger) Done() {
	
}
