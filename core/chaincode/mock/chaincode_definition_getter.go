
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
)

type ChaincodeDefinitionGetter struct {
	ChaincodeDefinitionStub        func(string, string, ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)
	chaincodeDefinitionMutex       sync.RWMutex
	chaincodeDefinitionArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
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

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinition(arg1 string, arg2 string, arg3 ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	fake.chaincodeDefinitionMutex.Lock()
	ret, specificReturn := fake.chaincodeDefinitionReturnsOnCall[len(fake.chaincodeDefinitionArgsForCall)]
	fake.chaincodeDefinitionArgsForCall = append(fake.chaincodeDefinitionArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ChaincodeDefinition", []interface{}{arg1, arg2, arg3})
	fake.chaincodeDefinitionMutex.Unlock()
	if fake.ChaincodeDefinitionStub != nil {
		return fake.ChaincodeDefinitionStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeDefinitionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinitionCallCount() int {
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	return len(fake.chaincodeDefinitionArgsForCall)
}

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinitionCalls(stub func(string, string, ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)) {
	fake.chaincodeDefinitionMutex.Lock()
	defer fake.chaincodeDefinitionMutex.Unlock()
	fake.ChaincodeDefinitionStub = stub
}

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinitionArgsForCall(i int) (string, string, ledger.SimpleQueryExecutor) {
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	argsForCall := fake.chaincodeDefinitionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinitionReturns(result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.chaincodeDefinitionMutex.Lock()
	defer fake.chaincodeDefinitionMutex.Unlock()
	fake.ChaincodeDefinitionStub = nil
	fake.chaincodeDefinitionReturns = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeDefinitionGetter) ChaincodeDefinitionReturnsOnCall(i int, result1 ccprovider.ChaincodeDefinition, result2 error) {
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

func (fake *ChaincodeDefinitionGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeDefinitionMutex.RLock()
	defer fake.chaincodeDefinitionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeDefinitionGetter) recordInvocation(key string, args []interface{}) {
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
