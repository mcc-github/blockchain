/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/asn1"
	"encoding/hex"
	"sync"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)

const (
	defaultMaxCacheSize   = 1000
	defaultRetentionRatio = 0.75
)

var (
	
	asBytes = asn1.Marshal
)

type acSupport interface {
	
	
	EligibleForService(channel string, data common.SignedData) error

	
	ConfigSequence(channel string) uint64
}

type authCacheConfig struct {
	enabled bool
	
	
	maxCacheSize int
	
	
	purgeRetentionRatio float64
}



type authCache struct {
	credentialCache map[string]*accessCache
	acSupport
	sync.RWMutex
	conf authCacheConfig
}

func newAuthCache(s acSupport, conf authCacheConfig) *authCache {
	return &authCache{
		acSupport:       s,
		credentialCache: make(map[string]*accessCache),
		conf:            conf,
	}
}



func (ac *authCache) EligibleForService(channel string, data common.SignedData) error {
	if !ac.conf.enabled {
		return ac.acSupport.EligibleForService(channel, data)
	}
	
	ac.RLock()
	cache := ac.credentialCache[channel]
	ac.RUnlock()
	if cache == nil {
		
		ac.Lock()
		cache = ac.newAccessCache(channel)
		
		ac.credentialCache[channel] = cache
		ac.Unlock()
	}
	return cache.EligibleForService(data)
}

type accessCache struct {
	sync.RWMutex
	channel      string
	ac           *authCache
	lastSequence uint64
	entries      map[string]error
}

func (ac *authCache) newAccessCache(channel string) *accessCache {
	return &accessCache{
		channel: channel,
		ac:      ac,
		entries: make(map[string]error),
	}
}

func (cache *accessCache) EligibleForService(data common.SignedData) error {
	key, err := signedDataToKey(data)
	if err != nil {
		logger.Warningf("Failed computing key of signed data: +%v", err)
		return errors.Wrap(err, "failed computing key of signed data")
	}
	currSeq := cache.ac.acSupport.ConfigSequence(cache.channel)
	if cache.isValid(currSeq) {
		foundInCache, isEligibleErr := cache.lookup(key)
		if foundInCache {
			return isEligibleErr
		}
	} else {
		cache.configChange(currSeq)
	}

	
	
	
	cache.purgeEntriesIfNeeded()

	
	err = cache.ac.acSupport.EligibleForService(cache.channel, data)
	cache.Lock()
	defer cache.Unlock()
	
	if currSeq != cache.ac.acSupport.ConfigSequence(cache.channel) {
		
		
		
		
		return err
	}
	
	
	cache.entries[key] = err
	return err
}

func (cache *accessCache) isPurgeNeeded() bool {
	cache.RLock()
	defer cache.RUnlock()
	return len(cache.entries)+1 > cache.ac.conf.maxCacheSize
}

func (cache *accessCache) purgeEntriesIfNeeded() {
	if !cache.isPurgeNeeded() {
		return
	}

	cache.Lock()
	defer cache.Unlock()

	maxCacheSize := cache.ac.conf.maxCacheSize
	purgeRatio := cache.ac.conf.purgeRetentionRatio
	entries2evict := maxCacheSize - int(purgeRatio*float64(maxCacheSize))

	for key := range cache.entries {
		if entries2evict == 0 {
			return
		}
		entries2evict--
		delete(cache.entries, key)
	}
}

func (cache *accessCache) isValid(currSeq uint64) bool {
	cache.RLock()
	defer cache.RUnlock()
	return currSeq == cache.lastSequence
}

func (cache *accessCache) configChange(currSeq uint64) {
	cache.Lock()
	defer cache.Unlock()
	cache.lastSequence = currSeq
	
	cache.entries = make(map[string]error)
}

func (cache *accessCache) lookup(key string) (cacheHit bool, lookupResult error) {
	cache.RLock()
	defer cache.RUnlock()

	lookupResult, cacheHit = cache.entries[key]
	return
}

func signedDataToKey(data common.SignedData) (string, error) {
	b, err := asBytes(data)
	if err != nil {
		return "", errors.Wrap(err, "failed marshaling signed data")
	}
	return hex.EncodeToString(util.ComputeSHA256(b)), nil
}
