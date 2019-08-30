
package fake

import (
	"sync"

	mspa "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/msp"
)

type IdentityDeserializer struct {
	DeserializeIdentityStub        func([]byte) (msp.Identity, error)
	deserializeIdentityMutex       sync.RWMutex
	deserializeIdentityArgsForCall []struct {
		arg1 []byte
	}
	deserializeIdentityReturns struct {
		result1 msp.Identity
		result2 error
	}
	deserializeIdentityReturnsOnCall map[int]struct {
		result1 msp.Identity
		result2 error
	}
	IsWellFormedStub        func(*mspa.SerializedIdentity) error
	isWellFormedMutex       sync.RWMutex
	isWellFormedArgsForCall []struct {
		arg1 *mspa.SerializedIdentity
	}
	isWellFormedReturns struct {
		result1 error
	}
	isWellFormedReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *IdentityDeserializer) DeserializeIdentity(arg1 []byte) (msp.Identity, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.deserializeIdentityMutex.Lock()
	ret, specificReturn := fake.deserializeIdentityReturnsOnCall[len(fake.deserializeIdentityArgsForCall)]
	fake.deserializeIdentityArgsForCall = append(fake.deserializeIdentityArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("DeserializeIdentity", []interface{}{arg1Copy})
	fake.deserializeIdentityMutex.Unlock()
	if fake.DeserializeIdentityStub != nil {
		return fake.DeserializeIdentityStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deserializeIdentityReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *IdentityDeserializer) DeserializeIdentityCallCount() int {
	fake.deserializeIdentityMutex.RLock()
	defer fake.deserializeIdentityMutex.RUnlock()
	return len(fake.deserializeIdentityArgsForCall)
}

func (fake *IdentityDeserializer) DeserializeIdentityCalls(stub func([]byte) (msp.Identity, error)) {
	fake.deserializeIdentityMutex.Lock()
	defer fake.deserializeIdentityMutex.Unlock()
	fake.DeserializeIdentityStub = stub
}

func (fake *IdentityDeserializer) DeserializeIdentityArgsForCall(i int) []byte {
	fake.deserializeIdentityMutex.RLock()
	defer fake.deserializeIdentityMutex.RUnlock()
	argsForCall := fake.deserializeIdentityArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IdentityDeserializer) DeserializeIdentityReturns(result1 msp.Identity, result2 error) {
	fake.deserializeIdentityMutex.Lock()
	defer fake.deserializeIdentityMutex.Unlock()
	fake.DeserializeIdentityStub = nil
	fake.deserializeIdentityReturns = struct {
		result1 msp.Identity
		result2 error
	}{result1, result2}
}

func (fake *IdentityDeserializer) DeserializeIdentityReturnsOnCall(i int, result1 msp.Identity, result2 error) {
	fake.deserializeIdentityMutex.Lock()
	defer fake.deserializeIdentityMutex.Unlock()
	fake.DeserializeIdentityStub = nil
	if fake.deserializeIdentityReturnsOnCall == nil {
		fake.deserializeIdentityReturnsOnCall = make(map[int]struct {
			result1 msp.Identity
			result2 error
		})
	}
	fake.deserializeIdentityReturnsOnCall[i] = struct {
		result1 msp.Identity
		result2 error
	}{result1, result2}
}

func (fake *IdentityDeserializer) IsWellFormed(arg1 *mspa.SerializedIdentity) error {
	fake.isWellFormedMutex.Lock()
	ret, specificReturn := fake.isWellFormedReturnsOnCall[len(fake.isWellFormedArgsForCall)]
	fake.isWellFormedArgsForCall = append(fake.isWellFormedArgsForCall, struct {
		arg1 *mspa.SerializedIdentity
	}{arg1})
	fake.recordInvocation("IsWellFormed", []interface{}{arg1})
	fake.isWellFormedMutex.Unlock()
	if fake.IsWellFormedStub != nil {
		return fake.IsWellFormedStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isWellFormedReturns
	return fakeReturns.result1
}

func (fake *IdentityDeserializer) IsWellFormedCallCount() int {
	fake.isWellFormedMutex.RLock()
	defer fake.isWellFormedMutex.RUnlock()
	return len(fake.isWellFormedArgsForCall)
}

func (fake *IdentityDeserializer) IsWellFormedCalls(stub func(*mspa.SerializedIdentity) error) {
	fake.isWellFormedMutex.Lock()
	defer fake.isWellFormedMutex.Unlock()
	fake.IsWellFormedStub = stub
}

func (fake *IdentityDeserializer) IsWellFormedArgsForCall(i int) *mspa.SerializedIdentity {
	fake.isWellFormedMutex.RLock()
	defer fake.isWellFormedMutex.RUnlock()
	argsForCall := fake.isWellFormedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IdentityDeserializer) IsWellFormedReturns(result1 error) {
	fake.isWellFormedMutex.Lock()
	defer fake.isWellFormedMutex.Unlock()
	fake.IsWellFormedStub = nil
	fake.isWellFormedReturns = struct {
		result1 error
	}{result1}
}

func (fake *IdentityDeserializer) IsWellFormedReturnsOnCall(i int, result1 error) {
	fake.isWellFormedMutex.Lock()
	defer fake.isWellFormedMutex.Unlock()
	fake.IsWellFormedStub = nil
	if fake.isWellFormedReturnsOnCall == nil {
		fake.isWellFormedReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.isWellFormedReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *IdentityDeserializer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deserializeIdentityMutex.RLock()
	defer fake.deserializeIdentityMutex.RUnlock()
	fake.isWellFormedMutex.RLock()
	defer fake.isWellFormedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *IdentityDeserializer) recordInvocation(key string, args []interface{}) {
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
