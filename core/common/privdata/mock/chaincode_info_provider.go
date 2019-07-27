
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)

type ChaincodeInfoProvider struct {
	ChaincodeInfoStub        func(string, string, ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error)
	chaincodeInfoMutex       sync.RWMutex
	chaincodeInfoArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}
	chaincodeInfoReturns struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	chaincodeInfoReturnsOnCall map[int]struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	CollectionInfoStub        func(string, string, string, ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error)
	collectionInfoMutex       sync.RWMutex
	collectionInfoArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 ledger.SimpleQueryExecutor
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

func (fake *ChaincodeInfoProvider) ChaincodeInfo(arg1 string, arg2 string, arg3 ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error) {
	fake.chaincodeInfoMutex.Lock()
	ret, specificReturn := fake.chaincodeInfoReturnsOnCall[len(fake.chaincodeInfoArgsForCall)]
	fake.chaincodeInfoArgsForCall = append(fake.chaincodeInfoArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ChaincodeInfo", []interface{}{arg1, arg2, arg3})
	fake.chaincodeInfoMutex.Unlock()
	if fake.ChaincodeInfoStub != nil {
		return fake.ChaincodeInfoStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeInfoProvider) ChaincodeInfoCallCount() int {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	return len(fake.chaincodeInfoArgsForCall)
}

func (fake *ChaincodeInfoProvider) ChaincodeInfoCalls(stub func(string, string, ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error)) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = stub
}

func (fake *ChaincodeInfoProvider) ChaincodeInfoArgsForCall(i int) (string, string, ledger.SimpleQueryExecutor) {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	argsForCall := fake.chaincodeInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ChaincodeInfoProvider) ChaincodeInfoReturns(result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	fake.chaincodeInfoReturns = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoProvider) ChaincodeInfoReturnsOnCall(i int, result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	if fake.chaincodeInfoReturnsOnCall == nil {
		fake.chaincodeInfoReturnsOnCall = make(map[int]struct {
			result1 *ledger.DeployedChaincodeInfo
			result2 error
		})
	}
	fake.chaincodeInfoReturnsOnCall[i] = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoProvider) CollectionInfo(arg1 string, arg2 string, arg3 string, arg4 ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error) {
	fake.collectionInfoMutex.Lock()
	ret, specificReturn := fake.collectionInfoReturnsOnCall[len(fake.collectionInfoArgsForCall)]
	fake.collectionInfoArgsForCall = append(fake.collectionInfoArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("CollectionInfo", []interface{}{arg1, arg2, arg3, arg4})
	fake.collectionInfoMutex.Unlock()
	if fake.CollectionInfoStub != nil {
		return fake.CollectionInfoStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.collectionInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeInfoProvider) CollectionInfoCallCount() int {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	return len(fake.collectionInfoArgsForCall)
}

func (fake *ChaincodeInfoProvider) CollectionInfoCalls(stub func(string, string, string, ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error)) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = stub
}

func (fake *ChaincodeInfoProvider) CollectionInfoArgsForCall(i int) (string, string, string, ledger.SimpleQueryExecutor) {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	argsForCall := fake.collectionInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *ChaincodeInfoProvider) CollectionInfoReturns(result1 *common.StaticCollectionConfig, result2 error) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = nil
	fake.collectionInfoReturns = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoProvider) CollectionInfoReturnsOnCall(i int, result1 *common.StaticCollectionConfig, result2 error) {
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

func (fake *ChaincodeInfoProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeInfoProvider) recordInvocation(key string, args []interface{}) {
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
