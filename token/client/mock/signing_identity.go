
package mock

import (
	sync "sync"

	token "github.com/mcc-github/blockchain/token"
)

type SigningIdentity struct {
	SerializeStub        func() ([]byte, error)
	serializeMutex       sync.RWMutex
	serializeArgsForCall []struct {
	}
	serializeReturns struct {
		result1 []byte
		result2 error
	}
	serializeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SignStub        func([]byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		arg1 []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SigningIdentity) Serialize() ([]byte, error) {
	fake.serializeMutex.Lock()
	ret, specificReturn := fake.serializeReturnsOnCall[len(fake.serializeArgsForCall)]
	fake.serializeArgsForCall = append(fake.serializeArgsForCall, struct {
	}{})
	fake.recordInvocation("Serialize", []interface{}{})
	fake.serializeMutex.Unlock()
	if fake.SerializeStub != nil {
		return fake.SerializeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.serializeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SigningIdentity) SerializeCallCount() int {
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	return len(fake.serializeArgsForCall)
}

func (fake *SigningIdentity) SerializeCalls(stub func() ([]byte, error)) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
	fake.SerializeStub = stub
}

func (fake *SigningIdentity) SerializeReturns(result1 []byte, result2 error) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
	fake.SerializeStub = nil
	fake.serializeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SigningIdentity) SerializeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
	fake.SerializeStub = nil
	if fake.serializeReturnsOnCall == nil {
		fake.serializeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.serializeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SigningIdentity) Sign(arg1 []byte) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("Sign", []interface{}{arg1Copy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.signReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SigningIdentity) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *SigningIdentity) SignCalls(stub func([]byte) ([]byte, error)) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = stub
}

func (fake *SigningIdentity) SignArgsForCall(i int) []byte {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	argsForCall := fake.signArgsForCall[i]
	return argsForCall.arg1
}

func (fake *SigningIdentity) SignReturns(result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SigningIdentity) SignReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = nil
	if fake.signReturnsOnCall == nil {
		fake.signReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.signReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SigningIdentity) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SigningIdentity) recordInvocation(key string, args []interface{}) {
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

var _ token.SigningIdentity = new(SigningIdentity)
