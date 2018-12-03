
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type IssuerPublicKey struct {
	BytesStub        func() ([]byte, error)
	bytesMutex       sync.RWMutex
	bytesArgsForCall []struct{}
	bytesReturns     struct {
		result1 []byte
		result2 error
	}
	bytesReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	HashStub        func() []byte
	hashMutex       sync.RWMutex
	hashArgsForCall []struct{}
	hashReturns     struct {
		result1 []byte
	}
	hashReturnsOnCall map[int]struct {
		result1 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *IssuerPublicKey) Bytes() ([]byte, error) {
	fake.bytesMutex.Lock()
	ret, specificReturn := fake.bytesReturnsOnCall[len(fake.bytesArgsForCall)]
	fake.bytesArgsForCall = append(fake.bytesArgsForCall, struct{}{})
	fake.recordInvocation("Bytes", []interface{}{})
	fake.bytesMutex.Unlock()
	if fake.BytesStub != nil {
		return fake.BytesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.bytesReturns.result1, fake.bytesReturns.result2
}

func (fake *IssuerPublicKey) BytesCallCount() int {
	fake.bytesMutex.RLock()
	defer fake.bytesMutex.RUnlock()
	return len(fake.bytesArgsForCall)
}

func (fake *IssuerPublicKey) BytesReturns(result1 []byte, result2 error) {
	fake.BytesStub = nil
	fake.bytesReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *IssuerPublicKey) BytesReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.BytesStub = nil
	if fake.bytesReturnsOnCall == nil {
		fake.bytesReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.bytesReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *IssuerPublicKey) Hash() []byte {
	fake.hashMutex.Lock()
	ret, specificReturn := fake.hashReturnsOnCall[len(fake.hashArgsForCall)]
	fake.hashArgsForCall = append(fake.hashArgsForCall, struct{}{})
	fake.recordInvocation("Hash", []interface{}{})
	fake.hashMutex.Unlock()
	if fake.HashStub != nil {
		return fake.HashStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.hashReturns.result1
}

func (fake *IssuerPublicKey) HashCallCount() int {
	fake.hashMutex.RLock()
	defer fake.hashMutex.RUnlock()
	return len(fake.hashArgsForCall)
}

func (fake *IssuerPublicKey) HashReturns(result1 []byte) {
	fake.HashStub = nil
	fake.hashReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *IssuerPublicKey) HashReturnsOnCall(i int, result1 []byte) {
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

func (fake *IssuerPublicKey) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bytesMutex.RLock()
	defer fake.bytesMutex.RUnlock()
	fake.hashMutex.RLock()
	defer fake.hashMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *IssuerPublicKey) recordInvocation(key string, args []interface{}) {
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

var _ handlers.IssuerPublicKey = new(IssuerPublicKey)