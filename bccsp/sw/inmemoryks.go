/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sw

import (
	"encoding/hex"
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/pkg/errors"
)


func NewInMemoryKeyStore() bccsp.KeyStore {
	eks := &inmemoryKeyStore{}
	eks.keys = make(map[string]bccsp.Key)
	return eks
}

type inmemoryKeyStore struct {
	
	keys map[string]bccsp.Key
	m    sync.RWMutex
}


func (ks *inmemoryKeyStore) ReadOnly() bool {
	return false
}


func (ks *inmemoryKeyStore) GetKey(ski []byte) (bccsp.Key, error) {
	if len(ski) == 0 {
		return nil, errors.New("ski is nil or empty")
	}

	skiStr := hex.EncodeToString(ski)

	ks.m.RLock()
	defer ks.m.RUnlock()
	if key, found := ks.keys[skiStr]; found {
		return key, nil
	}
	return nil, errors.Errorf("no key found for ski %x", ski)
}


func (ks *inmemoryKeyStore) StoreKey(k bccsp.Key) error {
	if k == nil {
		return errors.New("key is nil")
	}

	ski := hex.EncodeToString(k.SKI())

	ks.m.Lock()
	defer ks.m.Unlock()

	if _, found := ks.keys[ski]; found {
		return errors.Errorf("ski %x already exists in the keystore", k.SKI())
	}
	ks.keys[ski] = k

	return nil
}
