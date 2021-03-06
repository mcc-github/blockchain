
package mock

import (
	"sync"
)

type ChaincodeDefinition struct {
	CCNameStub        func() string
	cCNameMutex       sync.RWMutex
	cCNameArgsForCall []struct {
	}
	cCNameReturns struct {
		result1 string
	}
	cCNameReturnsOnCall map[int]struct {
		result1 string
	}
	CCVersionStub        func() string
	cCVersionMutex       sync.RWMutex
	cCVersionArgsForCall []struct {
	}
	cCVersionReturns struct {
		result1 string
	}
	cCVersionReturnsOnCall map[int]struct {
		result1 string
	}
	EndorsementStub        func() string
	endorsementMutex       sync.RWMutex
	endorsementArgsForCall []struct {
	}
	endorsementReturns struct {
		result1 string
	}
	endorsementReturnsOnCall map[int]struct {
		result1 string
	}
	HashStub        func() []byte
	hashMutex       sync.RWMutex
	hashArgsForCall []struct {
	}
	hashReturns struct {
		result1 []byte
	}
	hashReturnsOnCall map[int]struct {
		result1 []byte
	}
	RequiresInitStub        func() bool
	requiresInitMutex       sync.RWMutex
	requiresInitArgsForCall []struct {
	}
	requiresInitReturns struct {
		result1 bool
	}
	requiresInitReturnsOnCall map[int]struct {
		result1 bool
	}
	ValidationStub        func() (string, []byte)
	validationMutex       sync.RWMutex
	validationArgsForCall []struct {
	}
	validationReturns struct {
		result1 string
		result2 []byte
	}
	validationReturnsOnCall map[int]struct {
		result1 string
		result2 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeDefinition) CCName() string {
	fake.cCNameMutex.Lock()
	ret, specificReturn := fake.cCNameReturnsOnCall[len(fake.cCNameArgsForCall)]
	fake.cCNameArgsForCall = append(fake.cCNameArgsForCall, struct {
	}{})
	fake.recordInvocation("CCName", []interface{}{})
	fake.cCNameMutex.Unlock()
	if fake.CCNameStub != nil {
		return fake.CCNameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.cCNameReturns
	return fakeReturns.result1
}

func (fake *ChaincodeDefinition) CCNameCallCount() int {
	fake.cCNameMutex.RLock()
	defer fake.cCNameMutex.RUnlock()
	return len(fake.cCNameArgsForCall)
}

func (fake *ChaincodeDefinition) CCNameCalls(stub func() string) {
	fake.cCNameMutex.Lock()
	defer fake.cCNameMutex.Unlock()
	fake.CCNameStub = stub
}

func (fake *ChaincodeDefinition) CCNameReturns(result1 string) {
	fake.cCNameMutex.Lock()
	defer fake.cCNameMutex.Unlock()
	fake.CCNameStub = nil
	fake.cCNameReturns = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) CCNameReturnsOnCall(i int, result1 string) {
	fake.cCNameMutex.Lock()
	defer fake.cCNameMutex.Unlock()
	fake.CCNameStub = nil
	if fake.cCNameReturnsOnCall == nil {
		fake.cCNameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.cCNameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) CCVersion() string {
	fake.cCVersionMutex.Lock()
	ret, specificReturn := fake.cCVersionReturnsOnCall[len(fake.cCVersionArgsForCall)]
	fake.cCVersionArgsForCall = append(fake.cCVersionArgsForCall, struct {
	}{})
	fake.recordInvocation("CCVersion", []interface{}{})
	fake.cCVersionMutex.Unlock()
	if fake.CCVersionStub != nil {
		return fake.CCVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.cCVersionReturns
	return fakeReturns.result1
}

func (fake *ChaincodeDefinition) CCVersionCallCount() int {
	fake.cCVersionMutex.RLock()
	defer fake.cCVersionMutex.RUnlock()
	return len(fake.cCVersionArgsForCall)
}

func (fake *ChaincodeDefinition) CCVersionCalls(stub func() string) {
	fake.cCVersionMutex.Lock()
	defer fake.cCVersionMutex.Unlock()
	fake.CCVersionStub = stub
}

func (fake *ChaincodeDefinition) CCVersionReturns(result1 string) {
	fake.cCVersionMutex.Lock()
	defer fake.cCVersionMutex.Unlock()
	fake.CCVersionStub = nil
	fake.cCVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) CCVersionReturnsOnCall(i int, result1 string) {
	fake.cCVersionMutex.Lock()
	defer fake.cCVersionMutex.Unlock()
	fake.CCVersionStub = nil
	if fake.cCVersionReturnsOnCall == nil {
		fake.cCVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.cCVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) Endorsement() string {
	fake.endorsementMutex.Lock()
	ret, specificReturn := fake.endorsementReturnsOnCall[len(fake.endorsementArgsForCall)]
	fake.endorsementArgsForCall = append(fake.endorsementArgsForCall, struct {
	}{})
	fake.recordInvocation("Endorsement", []interface{}{})
	fake.endorsementMutex.Unlock()
	if fake.EndorsementStub != nil {
		return fake.EndorsementStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.endorsementReturns
	return fakeReturns.result1
}

func (fake *ChaincodeDefinition) EndorsementCallCount() int {
	fake.endorsementMutex.RLock()
	defer fake.endorsementMutex.RUnlock()
	return len(fake.endorsementArgsForCall)
}

func (fake *ChaincodeDefinition) EndorsementCalls(stub func() string) {
	fake.endorsementMutex.Lock()
	defer fake.endorsementMutex.Unlock()
	fake.EndorsementStub = stub
}

func (fake *ChaincodeDefinition) EndorsementReturns(result1 string) {
	fake.endorsementMutex.Lock()
	defer fake.endorsementMutex.Unlock()
	fake.EndorsementStub = nil
	fake.endorsementReturns = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) EndorsementReturnsOnCall(i int, result1 string) {
	fake.endorsementMutex.Lock()
	defer fake.endorsementMutex.Unlock()
	fake.EndorsementStub = nil
	if fake.endorsementReturnsOnCall == nil {
		fake.endorsementReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.endorsementReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeDefinition) Hash() []byte {
	fake.hashMutex.Lock()
	ret, specificReturn := fake.hashReturnsOnCall[len(fake.hashArgsForCall)]
	fake.hashArgsForCall = append(fake.hashArgsForCall, struct {
	}{})
	fake.recordInvocation("Hash", []interface{}{})
	fake.hashMutex.Unlock()
	if fake.HashStub != nil {
		return fake.HashStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.hashReturns
	return fakeReturns.result1
}

func (fake *ChaincodeDefinition) HashCallCount() int {
	fake.hashMutex.RLock()
	defer fake.hashMutex.RUnlock()
	return len(fake.hashArgsForCall)
}

func (fake *ChaincodeDefinition) HashCalls(stub func() []byte) {
	fake.hashMutex.Lock()
	defer fake.hashMutex.Unlock()
	fake.HashStub = stub
}

func (fake *ChaincodeDefinition) HashReturns(result1 []byte) {
	fake.hashMutex.Lock()
	defer fake.hashMutex.Unlock()
	fake.HashStub = nil
	fake.hashReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *ChaincodeDefinition) HashReturnsOnCall(i int, result1 []byte) {
	fake.hashMutex.Lock()
	defer fake.hashMutex.Unlock()
	fake.HashStub = nil
	if fake.hashReturnsOnCall == nil {
		fake.hashReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.hashReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *ChaincodeDefinition) RequiresInit() bool {
	fake.requiresInitMutex.Lock()
	ret, specificReturn := fake.requiresInitReturnsOnCall[len(fake.requiresInitArgsForCall)]
	fake.requiresInitArgsForCall = append(fake.requiresInitArgsForCall, struct {
	}{})
	fake.recordInvocation("RequiresInit", []interface{}{})
	fake.requiresInitMutex.Unlock()
	if fake.RequiresInitStub != nil {
		return fake.RequiresInitStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.requiresInitReturns
	return fakeReturns.result1
}

func (fake *ChaincodeDefinition) RequiresInitCallCount() int {
	fake.requiresInitMutex.RLock()
	defer fake.requiresInitMutex.RUnlock()
	return len(fake.requiresInitArgsForCall)
}

func (fake *ChaincodeDefinition) RequiresInitCalls(stub func() bool) {
	fake.requiresInitMutex.Lock()
	defer fake.requiresInitMutex.Unlock()
	fake.RequiresInitStub = stub
}

func (fake *ChaincodeDefinition) RequiresInitReturns(result1 bool) {
	fake.requiresInitMutex.Lock()
	defer fake.requiresInitMutex.Unlock()
	fake.RequiresInitStub = nil
	fake.requiresInitReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ChaincodeDefinition) RequiresInitReturnsOnCall(i int, result1 bool) {
	fake.requiresInitMutex.Lock()
	defer fake.requiresInitMutex.Unlock()
	fake.RequiresInitStub = nil
	if fake.requiresInitReturnsOnCall == nil {
		fake.requiresInitReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.requiresInitReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ChaincodeDefinition) Validation() (string, []byte) {
	fake.validationMutex.Lock()
	ret, specificReturn := fake.validationReturnsOnCall[len(fake.validationArgsForCall)]
	fake.validationArgsForCall = append(fake.validationArgsForCall, struct {
	}{})
	fake.recordInvocation("Validation", []interface{}{})
	fake.validationMutex.Unlock()
	if fake.ValidationStub != nil {
		return fake.ValidationStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.validationReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeDefinition) ValidationCallCount() int {
	fake.validationMutex.RLock()
	defer fake.validationMutex.RUnlock()
	return len(fake.validationArgsForCall)
}

func (fake *ChaincodeDefinition) ValidationCalls(stub func() (string, []byte)) {
	fake.validationMutex.Lock()
	defer fake.validationMutex.Unlock()
	fake.ValidationStub = stub
}

func (fake *ChaincodeDefinition) ValidationReturns(result1 string, result2 []byte) {
	fake.validationMutex.Lock()
	defer fake.validationMutex.Unlock()
	fake.ValidationStub = nil
	fake.validationReturns = struct {
		result1 string
		result2 []byte
	}{result1, result2}
}

func (fake *ChaincodeDefinition) ValidationReturnsOnCall(i int, result1 string, result2 []byte) {
	fake.validationMutex.Lock()
	defer fake.validationMutex.Unlock()
	fake.ValidationStub = nil
	if fake.validationReturnsOnCall == nil {
		fake.validationReturnsOnCall = make(map[int]struct {
			result1 string
			result2 []byte
		})
	}
	fake.validationReturnsOnCall[i] = struct {
		result1 string
		result2 []byte
	}{result1, result2}
}

func (fake *ChaincodeDefinition) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cCNameMutex.RLock()
	defer fake.cCNameMutex.RUnlock()
	fake.cCVersionMutex.RLock()
	defer fake.cCVersionMutex.RUnlock()
	fake.endorsementMutex.RLock()
	defer fake.endorsementMutex.RUnlock()
	fake.hashMutex.RLock()
	defer fake.hashMutex.RUnlock()
	fake.requiresInitMutex.RLock()
	defer fake.requiresInitMutex.RUnlock()
	fake.validationMutex.RLock()
	defer fake.validationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeDefinition) recordInvocation(key string, args []interface{}) {
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
