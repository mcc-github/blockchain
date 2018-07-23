
package mock

import (
	"sync"
)

type SystemCCProvider struct {
	IsSysCCStub        func(name string) bool
	isSysCCMutex       sync.RWMutex
	isSysCCArgsForCall []struct {
		name string
	}
	isSysCCReturns struct {
		result1 bool
	}
	isSysCCReturnsOnCall map[int]struct {
		result1 bool
	}
	IsSysCCAndNotInvokableCC2CCStub        func(name string) bool
	isSysCCAndNotInvokableCC2CCMutex       sync.RWMutex
	isSysCCAndNotInvokableCC2CCArgsForCall []struct {
		name string
	}
	isSysCCAndNotInvokableCC2CCReturns struct {
		result1 bool
	}
	isSysCCAndNotInvokableCC2CCReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SystemCCProvider) IsSysCC(name string) bool {
	fake.isSysCCMutex.Lock()
	ret, specificReturn := fake.isSysCCReturnsOnCall[len(fake.isSysCCArgsForCall)]
	fake.isSysCCArgsForCall = append(fake.isSysCCArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCC", []interface{}{name})
	fake.isSysCCMutex.Unlock()
	if fake.IsSysCCStub != nil {
		return fake.IsSysCCStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCReturns.result1
}

func (fake *SystemCCProvider) IsSysCCCallCount() int {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return len(fake.isSysCCArgsForCall)
}

func (fake *SystemCCProvider) IsSysCCArgsForCall(i int) string {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return fake.isSysCCArgsForCall[i].name
}

func (fake *SystemCCProvider) IsSysCCReturns(result1 bool) {
	fake.IsSysCCStub = nil
	fake.isSysCCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemCCProvider) IsSysCCReturnsOnCall(i int, result1 bool) {
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

func (fake *SystemCCProvider) IsSysCCAndNotInvokableCC2CC(name string) bool {
	fake.isSysCCAndNotInvokableCC2CCMutex.Lock()
	ret, specificReturn := fake.isSysCCAndNotInvokableCC2CCReturnsOnCall[len(fake.isSysCCAndNotInvokableCC2CCArgsForCall)]
	fake.isSysCCAndNotInvokableCC2CCArgsForCall = append(fake.isSysCCAndNotInvokableCC2CCArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCCAndNotInvokableCC2CC", []interface{}{name})
	fake.isSysCCAndNotInvokableCC2CCMutex.Unlock()
	if fake.IsSysCCAndNotInvokableCC2CCStub != nil {
		return fake.IsSysCCAndNotInvokableCC2CCStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCAndNotInvokableCC2CCReturns.result1
}

func (fake *SystemCCProvider) IsSysCCAndNotInvokableCC2CCCallCount() int {
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
	return len(fake.isSysCCAndNotInvokableCC2CCArgsForCall)
}

func (fake *SystemCCProvider) IsSysCCAndNotInvokableCC2CCArgsForCall(i int) string {
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
	return fake.isSysCCAndNotInvokableCC2CCArgsForCall[i].name
}

func (fake *SystemCCProvider) IsSysCCAndNotInvokableCC2CCReturns(result1 bool) {
	fake.IsSysCCAndNotInvokableCC2CCStub = nil
	fake.isSysCCAndNotInvokableCC2CCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemCCProvider) IsSysCCAndNotInvokableCC2CCReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCAndNotInvokableCC2CCStub = nil
	if fake.isSysCCAndNotInvokableCC2CCReturnsOnCall == nil {
		fake.isSysCCAndNotInvokableCC2CCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCAndNotInvokableCC2CCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SystemCCProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
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