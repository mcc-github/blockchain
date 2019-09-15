
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
)

type CollectionInfoProvider struct {
	CollectionInfoStub        func(string, string) (*common.StaticCollectionConfig, error)
	collectionInfoMutex       sync.RWMutex
	collectionInfoArgsForCall []struct {
		arg1 string
		arg2 string
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

func (fake *CollectionInfoProvider) CollectionInfo(arg1 string, arg2 string) (*common.StaticCollectionConfig, error) {
	fake.collectionInfoMutex.Lock()
	ret, specificReturn := fake.collectionInfoReturnsOnCall[len(fake.collectionInfoArgsForCall)]
	fake.collectionInfoArgsForCall = append(fake.collectionInfoArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("CollectionInfo", []interface{}{arg1, arg2})
	fake.collectionInfoMutex.Unlock()
	if fake.CollectionInfoStub != nil {
		return fake.CollectionInfoStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.collectionInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CollectionInfoProvider) CollectionInfoCallCount() int {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	return len(fake.collectionInfoArgsForCall)
}

func (fake *CollectionInfoProvider) CollectionInfoCalls(stub func(string, string) (*common.StaticCollectionConfig, error)) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = stub
}

func (fake *CollectionInfoProvider) CollectionInfoArgsForCall(i int) (string, string) {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	argsForCall := fake.collectionInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *CollectionInfoProvider) CollectionInfoReturns(result1 *common.StaticCollectionConfig, result2 error) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = nil
	fake.collectionInfoReturns = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *CollectionInfoProvider) CollectionInfoReturnsOnCall(i int, result1 *common.StaticCollectionConfig, result2 error) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
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
