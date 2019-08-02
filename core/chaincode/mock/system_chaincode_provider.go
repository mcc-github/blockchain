
package mock

import (
	"sync"
)

type SystemCCProvider struct {
	IsSysCCStub        func(string) bool
	isSysCCMutex       sync.RWMutex
	isSysCCArgsForCall []struct {
		arg1 string
	}
	isSysCCReturns struct {
		result1 bool
	}
	isSysCCReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SystemCCProvider) IsSysCC(arg1 string) bool {
	fake.isSysCCMutex.Lock()
	ret, specificReturn := fake.isSysCCReturnsOnCall[len(fake.isSysCCArgsForCall)]
	fake.isSysCCArgsForCall = append(fake.isSysCCArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("IsSysCC", []interface{}{arg1})
	fake.isSysCCMutex.Unlock()
	if fake.IsSysCCStub != nil {
		return fake.IsSysCCStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isSysCCReturns
	return fakeReturns.result1
}

func (fake *SystemCCProvider) IsSysCCCallCount() int {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return len(fake.isSysCCArgsForCall)
}

func (fake *SystemCCProvider) IsSysCCCalls(stub func(string) bool) {
	fake.isSysCCMutex.Lock()
	defer fake.isSysCCMutex.Unlock()
	fake.IsSysCCStub = stub
}

func (fake *SystemCCProvider) IsSysCCArgsForCall(i int) string {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	argsForCall := fake.isSysCCArgsForCall[i]
	return argsForCall.arg1
}

func (fake *SystemCCProvider) IsSysCCReturns(result1 bool) {
	fake.isSysCCMutex.Lock()
	defer fake.isSysCCMutex.Unlock()
	fake.IsSysCCStub = nil
	fake.isSysCCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemCCProvider) IsSysCCReturnsOnCall(i int, result1 bool) {
	fake.isSysCCMutex.Lock()
	defer fake.isSysCCMutex.Unlock()
	fake.IsSysCCStub = nil
	if fake.isSysCCReturnsOnCall == nil {
		fake.isSysCCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SystemCCProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SystemCCProvider) recordInvocation(key string, args []interface{}) {
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
