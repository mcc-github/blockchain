
package mock

import (
	"sync"
)

type LegacyChaincodeDefinition struct {
	ExecuteLegacySecurityChecksStub        func() error
	executeLegacySecurityChecksMutex       sync.RWMutex
	executeLegacySecurityChecksArgsForCall []struct {
	}
	executeLegacySecurityChecksReturns struct {
		result1 error
	}
	executeLegacySecurityChecksReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LegacyChaincodeDefinition) ExecuteLegacySecurityChecks() error {
	fake.executeLegacySecurityChecksMutex.Lock()
	ret, specificReturn := fake.executeLegacySecurityChecksReturnsOnCall[len(fake.executeLegacySecurityChecksArgsForCall)]
	fake.executeLegacySecurityChecksArgsForCall = append(fake.executeLegacySecurityChecksArgsForCall, struct {
	}{})
	fake.recordInvocation("ExecuteLegacySecurityChecks", []interface{}{})
	fake.executeLegacySecurityChecksMutex.Unlock()
	if fake.ExecuteLegacySecurityChecksStub != nil {
		return fake.ExecuteLegacySecurityChecksStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.executeLegacySecurityChecksReturns
	return fakeReturns.result1
}

func (fake *LegacyChaincodeDefinition) ExecuteLegacySecurityChecksCallCount() int {
	fake.executeLegacySecurityChecksMutex.RLock()
	defer fake.executeLegacySecurityChecksMutex.RUnlock()
	return len(fake.executeLegacySecurityChecksArgsForCall)
}

func (fake *LegacyChaincodeDefinition) ExecuteLegacySecurityChecksCalls(stub func() error) {
	fake.executeLegacySecurityChecksMutex.Lock()
	defer fake.executeLegacySecurityChecksMutex.Unlock()
	fake.ExecuteLegacySecurityChecksStub = stub
}

func (fake *LegacyChaincodeDefinition) ExecuteLegacySecurityChecksReturns(result1 error) {
	fake.executeLegacySecurityChecksMutex.Lock()
	defer fake.executeLegacySecurityChecksMutex.Unlock()
	fake.ExecuteLegacySecurityChecksStub = nil
	fake.executeLegacySecurityChecksReturns = struct {
		result1 error
	}{result1}
}

func (fake *LegacyChaincodeDefinition) ExecuteLegacySecurityChecksReturnsOnCall(i int, result1 error) {
	fake.executeLegacySecurityChecksMutex.Lock()
	defer fake.executeLegacySecurityChecksMutex.Unlock()
	fake.ExecuteLegacySecurityChecksStub = nil
	if fake.executeLegacySecurityChecksReturnsOnCall == nil {
		fake.executeLegacySecurityChecksReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.executeLegacySecurityChecksReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *LegacyChaincodeDefinition) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeLegacySecurityChecksMutex.RLock()
	defer fake.executeLegacySecurityChecksMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LegacyChaincodeDefinition) recordInvocation(key string, args []interface{}) {
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
