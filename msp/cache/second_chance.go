/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cache

import (
	"sync"
	"sync/atomic"
)







type secondChanceCache struct {
	
	table map[string]*cacheItem

	
	items []*cacheItem

	
	position int

	
	rwlock sync.RWMutex
}

type cacheItem struct {
	key   string
	value interface{}
	
	referenced int32
}

func newSecondChanceCache(cacheSize int) *secondChanceCache {
	var cache secondChanceCache
	cache.position = 0
	cache.items = make([]*cacheItem, cacheSize)
	cache.table = make(map[string]*cacheItem)

	return &cache
}

func (cache *secondChanceCache) len() int {
	cache.rwlock.RLock()
	defer cache.rwlock.RUnlock()

	return len(cache.table)
}

func (cache *secondChanceCache) get(key string) (interface{}, bool) {
	cache.rwlock.RLock()
	defer cache.rwlock.RUnlock()

	item, ok := cache.table[key]
	if !ok {
		return nil, false
	}

	
	atomic.StoreInt32(&item.referenced, 1)

	return item.value, true
}

func (cache *secondChanceCache) add(key string, value interface{}) {
	cache.rwlock.Lock()
	defer cache.rwlock.Unlock()

	if old, ok := cache.table[key]; ok {
		old.value = value
		atomic.StoreInt32(&old.referenced, 1)
		return
	}

	var item cacheItem
	item.key = key
	item.value = value
	atomic.StoreInt32(&item.referenced, 1)

	size := len(cache.items)
	num := len(cache.table)
	if num < size {
		
		cache.table[key] = &item
		cache.items[num] = &item
		return
	}

	
	for {
		
		victim := cache.items[cache.position]
		if atomic.LoadInt32(&victim.referenced) == 0 {
			
			delete(cache.table, victim.key)
			cache.table[key] = &item
			cache.items[cache.position] = &item
			cache.position = (cache.position + 1) % size
			return
		}

		
		
		atomic.StoreInt32(&victim.referenced, 0)
		cache.position = (cache.position + 1) % size
	}
}
