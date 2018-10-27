/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain



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
