/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/hex"
	"sync"

	"github.com/mcc-github/blockchain/common/util"
)



type MemoizeSigner struct {
	maxMemorySize int
	sync.RWMutex
	memory map[string][]byte
	sign   Signer
}



func NewMemoizeSigner(signFunc Signer, maxMemorySize int) *MemoizeSigner {
	return &MemoizeSigner{
		maxMemorySize: maxMemorySize,
		memory:        make(map[string][]byte),
		sign:          signFunc,
	}
}



func (ms *MemoizeSigner) Sign(msg []byte) ([]byte, error) {
	sig, isInMemory := ms.lookup(msg)
	if isInMemory {
		return sig, nil
	}
	sig, err := ms.sign(msg)
	if err != nil {
		return nil, err
	}
	ms.memorize(msg, sig)
	return sig, nil
}



func (ms *MemoizeSigner) lookup(msg []byte) ([]byte, bool) {
	ms.RLock()
	defer ms.RUnlock()
	sig, exists := ms.memory[msgDigest(msg)]
	return sig, exists
}

func (ms *MemoizeSigner) memorize(msg, signature []byte) {
	ms.RLock()
	shouldShrink := len(ms.memory) >= ms.maxMemorySize
	ms.RUnlock()

	if shouldShrink {
		ms.shrinkMemory()
	}

	ms.Lock()
	defer ms.Unlock()
	ms.memory[msgDigest(msg)] = signature
}



func (ms *MemoizeSigner) shrinkMemory() {
	ms.Lock()
	defer ms.Unlock()
	for len(ms.memory) > ms.maxMemorySize {
		ms.evictFromMemory()
	}
}


func (ms *MemoizeSigner) evictFromMemory() {
	for dig := range ms.memory {
		delete(ms.memory, dig)
		return
	}
}


func msgDigest(msg []byte) string {
	return hex.EncodeToString(util.ComputeSHA256(msg))
}
