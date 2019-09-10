
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	"github.com/mcc-github/blockchain/core/ledger"
)

type LegacyLifecycle struct {
	ChaincodeEndorsementInfoStub        func(string, string, ledger.SimpleQueryExecutor) (*lifecycle.ChaincodeEndorsementInfo, error)
	chaincodeEndorsementInfoMutex       sync.RWMutex
	chaincodeEndorsementInfoArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}
	chaincodeEndorsementInfoReturns struct {
		result1 *lifecycle.ChaincodeEndorsementInfo
		result2 error
	}
	chaincodeEndorsementInfoReturnsOnCall map[int]struct {
		result1 *lifecycle.ChaincodeEndorsementInfo
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfo(arg1 string, arg2 string, arg3 ledger.SimpleQueryExecutor) (*lifecycle.ChaincodeEndorsementInfo, error) {
	fake.chaincodeEndorsementInfoMutex.Lock()
	ret, specificReturn := fake.chaincodeEndorsementInfoReturnsOnCall[len(fake.chaincodeEndorsementInfoArgsForCall)]
	fake.chaincodeEndorsementInfoArgsForCall = append(fake.chaincodeEndorsementInfoArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ChaincodeEndorsementInfo", []interface{}{arg1, arg2, arg3})
	fake.chaincodeEndorsementInfoMutex.Unlock()
	if fake.ChaincodeEndorsementInfoStub != nil {
		return fake.ChaincodeEndorsementInfoStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeEndorsementInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfoCallCount() int {
	fake.chaincodeEndorsementInfoMutex.RLock()
	defer fake.chaincodeEndorsementInfoMutex.RUnlock()
	return len(fake.chaincodeEndorsementInfoArgsForCall)
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfoCalls(stub func(string, string, ledger.SimpleQueryExecutor) (*lifecycle.ChaincodeEndorsementInfo, error)) {
	fake.chaincodeEndorsementInfoMutex.Lock()
	defer fake.chaincodeEndorsementInfoMutex.Unlock()
	fake.ChaincodeEndorsementInfoStub = stub
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfoArgsForCall(i int) (string, string, ledger.SimpleQueryExecutor) {
	fake.chaincodeEndorsementInfoMutex.RLock()
	defer fake.chaincodeEndorsementInfoMutex.RUnlock()
	argsForCall := fake.chaincodeEndorsementInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfoReturns(result1 *lifecycle.ChaincodeEndorsementInfo, result2 error) {
	fake.chaincodeEndorsementInfoMutex.Lock()
	defer fake.chaincodeEndorsementInfoMutex.Unlock()
	fake.ChaincodeEndorsementInfoStub = nil
	fake.chaincodeEndorsementInfoReturns = struct {
		result1 *lifecycle.ChaincodeEndorsementInfo
		result2 error
	}{result1, result2}
}

func (fake *LegacyLifecycle) ChaincodeEndorsementInfoReturnsOnCall(i int, result1 *lifecycle.ChaincodeEndorsementInfo, result2 error) {
	fake.chaincodeEndorsementInfoMutex.Lock()
	defer fake.chaincodeEndorsementInfoMutex.Unlock()
	fake.ChaincodeEndorsementInfoStub = nil
	if fake.chaincodeEndorsementInfoReturnsOnCall == nil {
		fake.chaincodeEndorsementInfoReturnsOnCall = make(map[int]struct {
			result1 *lifecycle.ChaincodeEndorsementInfo
			result2 error
		})
	}
	fake.chaincodeEndorsementInfoReturnsOnCall[i] = struct {
		result1 *lifecycle.ChaincodeEndorsementInfo
		result2 error
	}{result1, result2}
}

func (fake *LegacyLifecycle) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeEndorsementInfoMutex.RLock()
	defer fake.chaincodeEndorsementInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LegacyLifecycle) recordInvocation(key string, args []interface{}) {
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

var _ lifecycle.Lifecycle = new(LegacyLifecycle)
