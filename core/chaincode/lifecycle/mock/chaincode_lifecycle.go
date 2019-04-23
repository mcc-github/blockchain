
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
)

type ChaincodeLifecycle struct {
	ChaincodeContainerInfoStub        func(string, ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error)
	chaincodeContainerInfoMutex       sync.RWMutex
	chaincodeContainerInfoArgsForCall []struct {
		arg1 string
		arg2 ledger.SimpleQueryExecutor
	}
	chaincodeContainerInfoReturns struct {
		result1 *ccprovider.ChaincodeContainerInfo
		result2 error
	}
	chaincodeContainerInfoReturnsOnCall map[int]struct {
		result1 *ccprovider.ChaincodeContainerInfo
		result2 error
	}
	ChaincodeDefinitionStub        func(string, ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)
	chaincodeDefinitionMutex       sync.RWMutex
	chaincodeDefinitionArgsForCall []struct {
		arg1 string
		arg2 ledger.SimpleQueryExecutor
	}
	chaincodeDefinitionReturns struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	chaincodeDefinitionReturnsOnCall map[int]struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfo(arg1 string, arg2 ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error) {
	fake.chaincodeContainerInfoMutex.Lock()
	ret, specificReturn := fake.chaincodeContainerInfoReturnsOnCall[len(fake.chaincodeContainerInfoArgsForCall)]
	fake.chaincodeContainerInfoArgsForCall = append(fake.chaincodeContainerInfoArgsForCall, struct {
		arg1 string
		arg2 ledger.SimpleQueryExecutor
	}{arg1, arg2})
	fake.recordInvocation("ChaincodeContainerInfo", []interface{}{arg1, arg2})
	fake.chaincodeContainerInfoMutex.Unlock()
	if fake.ChaincodeContainerInfoStub != nil {
		return fake.ChaincodeContainerInfoStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeContainerInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfoCallCount() int {
	fake.chaincodeContainerInfoMutex.RLock()
	defer fake.chaincodeContainerInfoMutex.RUnlock()
	return len(fake.chaincodeContainerInfoArgsForCall)
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfoCalls(stub func(string, ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error)) {
	fake.chaincodeContainerInfoMutex.Lock()
	defer fake.chaincodeContainerInfoMutex.Unlock()
	fake.ChaincodeContainerInfoStub = stub
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfoArgsForCall(i int) (string, ledger.SimpleQueryExecutor) {
	fake.chaincodeContainerInfoMutex.RLock()
	defer fake.chaincodeContainerInfoMutex.RUnlock()
	argsForCall := fake.chaincodeContainerInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfoReturns(result1 *ccprovider.ChaincodeContainerInfo, result2 error) {
	fake.chaincodeContainerInfoMutex.Lock()
	defer fake.chaincodeContainerInfoMutex.Unlock()
	fake.ChaincodeContainerInfoStub = nil
	fake.chaincodeContainerInfoReturns = struct {
		result1 *ccprovider.ChaincodeContainerInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeLifecycle) ChaincodeContainerInfoReturnsOnCall(i int, result1 *ccprovider.ChaincodeContainerInfo, result2 error) {
	fake.chaincodeContainerInfoMutex.Lock()
	defer fake.chaincodeContainerInfoMutex.Unlock()
	fake.ChaincodeContainerInfoStub = nil
	if fake.chaincodeContainerInfoReturnsOnCall == nil {
		fake.chaincodeContainerInfoReturnsOnCall = make(map[int]struct {
			result1 *ccprovider.ChaincodeContainerInfo
			result2 error
		})
	}
	fake.chaincodeContainerInfoReturnsOnCall[i] = struct {
		result1 *ccprovider.ChaincodeContainerInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeLifecycle) ChaincodeDefinition(arg1 string, arg2 ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	fake.chaincodeDefinitionMutex.Lock()
	ret, specificReturn := fake.chaincodeDefinitionReturnsOnCall[len(fake.chaincodeDefinitionArgsForCall)]
	fake.chaincodeDefinitionArgsForCall = append(fake.chaincodeDefinitionArgsForCall, struct {
		arg1 string
		arg2 ledger.SimpleQueryExecutor
	}{arg1, arg2})
	fake.recordInvocation("ChaincodeDefinition", []interface{}{arg1, arg2})
	fake.chaincodeDefinitionMutex.Unlock()
	if fake.ChaincodeDefinitionStub != nil {
		return fake.ChaincodeDefinitionStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeDefinitionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeLifecycle) ChaincodeDefinitionCallCount() int {
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	return len(fake.chaincodeDefinitionArgsForCall)
}

func (fake *ChaincodeLifecycle) ChaincodeDefinitionCalls(stub func(string, ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)) {
	fake.chaincodeDefinitionMutex.Lock()
	defer fake.chaincodeDefinitionMutex.Unlock()
	fake.ChaincodeDefinitionStub = stub
}

func (fake *ChaincodeLifecycle) ChaincodeDefinitionArgsForCall(i int) (string, ledger.SimpleQueryExecutor) {
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	argsForCall := fake.chaincodeDefinitionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeLifecycle) ChaincodeDefinitionReturns(result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.chaincodeDefinitionMutex.Lock()
	defer fake.chaincodeDefinitionMutex.Unlock()
	fake.ChaincodeDefinitionStub = nil
	fake.chaincodeDefinitionReturns = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeLifecycle) ChaincodeDefinitionReturnsOnCall(i int, result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.chaincodeDefinitionMutex.Lock()
	defer fake.chaincodeDefinitionMutex.Unlock()
	fake.ChaincodeDefinitionStub = nil
	if fake.chaincodeDefinitionReturnsOnCall == nil {
		fake.chaincodeDefinitionReturnsOnCall = make(map[int]struct {
			result1 ccprovider.ChaincodeDefinition
			result2 error
		})
	}
	fake.chaincodeDefinitionReturnsOnCall[i] = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeLifecycle) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeContainerInfoMutex.RLock()
	defer fake.chaincodeContainerInfoMutex.RUnlock()
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeLifecycle) recordInvocation(key string, args []interface{}) {
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
