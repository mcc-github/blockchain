
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/common"
)

type CollectionInfoProvider struct {
	CollectionInfoStub        func(chaincodeName, collectionName string) (*common.StaticCollectionConfig, error)
	collectionInfoMutex       sync.RWMutex
	collectionInfoArgsForCall []struct {
		chaincodeName  string
		collectionName string
	}
	collectionInfoReturns struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}
	collectionInfoReturnsOnCall map[int]struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CollectionInfoProvider) CollectionInfo(chaincodeName string, collectionName string) (*common.StaticCollectionConfig, error) {
	fake.collectionInfoMutex.Lock()
	ret, specificReturn := fake.collectionInfoReturnsOnCall[len(fake.collectionInfoArgsForCall)]
	fake.collectionInfoArgsForCall = append(fake.collectionInfoArgsForCall, struct {
		chaincodeName  string
		collectionName string
	}{chaincodeName, collectionName})
	fake.recordInvocation("CollectionInfo", []interface{}{chaincodeName, collectionName})
	fake.collectionInfoMutex.Unlock()
	if fake.CollectionInfoStub != nil {
		return fake.CollectionInfoStub(chaincodeName, collectionName)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.collectionInfoReturns.result1, fake.collectionInfoReturns.result2
}

func (fake *CollectionInfoProvider) CollectionInfoCallCount() int {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	return len(fake.collectionInfoArgsForCall)
}

func (fake *CollectionInfoProvider) CollectionInfoArgsForCall(i int) (string, string) {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	return fake.collectionInfoArgsForCall[i].chaincodeName, fake.collectionInfoArgsForCall[i].collectionName
}

func (fake *CollectionInfoProvider) CollectionInfoReturns(result1 *common.StaticCollectionConfig, result2 error) {
	fake.CollectionInfoStub = nil
	fake.collectionInfoReturns = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *CollectionInfoProvider) CollectionInfoReturnsOnCall(i int, result1 *common.StaticCollectionConfig, result2 error) {
	fake.CollectionInfoStub = nil
	if fake.collectionInfoReturnsOnCall == nil {
		fake.collectionInfoReturnsOnCall = make(map[int]struct {
			result1 *common.StaticCollectionConfig
			result2 error
		})
	}
	fake.collectionInfoReturnsOnCall[i] = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *CollectionInfoProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CollectionInfoProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
