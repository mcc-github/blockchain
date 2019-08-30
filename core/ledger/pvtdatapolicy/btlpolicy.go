/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatapolicy

import (
	"math"
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/common/privdata"
)

var defaultBTL uint64 = math.MaxUint64


type BTLPolicy interface {
	
	GetBTL(ns string, coll string) (uint64, error)
	
	GetExpiringBlock(namesapce string, collection string, committingBlock uint64) (uint64, error)
}




type LSCCBasedBTLPolicy struct {
	collInfoProvider collectionInfoProvider
	cache            map[btlkey]uint64
	lock             sync.Mutex
}

type btlkey struct {
	ns   string
	coll string
}


func ConstructBTLPolicy(collInfoProvider collectionInfoProvider) BTLPolicy {
	return &LSCCBasedBTLPolicy{
		collInfoProvider: collInfoProvider,
		cache:            make(map[btlkey]uint64),
	}
}


func (p *LSCCBasedBTLPolicy) GetBTL(namesapce string, collection string) (uint64, error) {
	var btl uint64
	var ok bool
	key := btlkey{namesapce, collection}
	p.lock.Lock()
	defer p.lock.Unlock()
	btl, ok = p.cache[key]
	if !ok {
		collConfig, err := p.collInfoProvider.CollectionInfo(namesapce, collection)
		if err != nil {
			return 0, err
		}
		if collConfig == nil {
			return 0, privdata.NoSuchCollectionError{Namespace: namesapce, Collection: collection}
		}
		btlConfigured := collConfig.BlockToLive
		if btlConfigured > 0 {
			btl = uint64(btlConfigured)
		} else {
			btl = defaultBTL
		}
		p.cache[key] = btl
	}
	return btl, nil
}


func (p *LSCCBasedBTLPolicy) GetExpiringBlock(namesapce string, collection string, committingBlock uint64) (uint64, error) {
	btl, err := p.GetBTL(namesapce, collection)
	if err != nil {
		return 0, err
	}
	expiryBlk := committingBlock + btl + uint64(1)
	if expiryBlk <= committingBlock { 
		expiryBlk = math.MaxUint64
	}
	return expiryBlk, nil
}

type collectionInfoProvider interface {
	CollectionInfo(chaincodeName, collectionName string) (*common.StaticCollectionConfig, error)
}


