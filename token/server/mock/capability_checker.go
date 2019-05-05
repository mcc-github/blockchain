
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/server"
)

type CapabilityChecker struct {
	FabTokenStub        func(string) (bool, error)
	fabTokenMutex       sync.RWMutex
	fabTokenArgsForCall []struct {
		arg1 string
	}
	fabTokenReturns struct {
		result1 bool
		result2 error
	}
	fabTokenReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CapabilityChecker) FabToken(arg1 string) (bool, error) {
	fake.fabTokenMutex.Lock()
	ret, specificReturn := fake.fabTokenReturnsOnCall[len(fake.fabTokenArgsForCall)]
	fake.fabTokenArgsForCall = append(fake.fabTokenArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FabToken", []interface{}{arg1})
	fake.fabTokenMutex.Unlock()
	if fake.FabTokenStub != nil {
		return fake.FabTokenStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.fabTokenReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CapabilityChecker) FabTokenCallCount() int {
	fake.fabTokenMutex.RLock()
	defer fake.fabTokenMutex.RUnlock()
	return len(fake.fabTokenArgsForCall)
}

func (fake *CapabilityChecker) FabTokenCalls(stub func(string) (bool, error)) {
	fake.fabTokenMutex.Lock()
	defer fake.fabTokenMutex.Unlock()
	fake.FabTokenStub = stub
}

func (fake *CapabilityChecker) FabTokenArgsForCall(i int) string {
	fake.fabTokenMutex.RLock()
	defer fake.fabTokenMutex.RUnlock()
	argsForCall := fake.fabTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *CapabilityChecker) FabTokenReturns(result1 bool, result2 error) {
	fake.fabTokenMutex.Lock()
	defer fake.fabTokenMutex.Unlock()
	fake.FabTokenStub = nil
	fake.fabTokenReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *CapabilityChecker) FabTokenReturnsOnCall(i int, result1 bool, result2 error) {
	fake.fabTokenMutex.Lock()
	defer fake.fabTokenMutex.Unlock()
	fake.FabTokenStub = nil
	if fake.fabTokenReturnsOnCall == nil {
		fake.fabTokenReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.fabTokenReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *CapabilityChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.fabTokenMutex.RLock()
	defer fake.fabTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CapabilityChecker) recordInvocation(key string, args []interface{}) {
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

var _ server.CapabilityChecker = new(CapabilityChecker)
