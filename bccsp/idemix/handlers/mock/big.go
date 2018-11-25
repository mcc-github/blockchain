
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type Big struct {
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
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Big) Bytes() ([]byte, error) {
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

func (fake *Big) BytesCallCount() int {
	fake.bytesMutex.RLock()
	defer fake.bytesMutex.RUnlock()
	return len(fake.bytesArgsForCall)
}

func (fake *Big) BytesReturns(result1 []byte, result2 error) {
	fake.BytesStub = nil
	fake.bytesReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Big) BytesReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Big) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bytesMutex.RLock()
	defer fake.bytesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Big) recordInvocation(key string, args []interface{}) {
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

var _ handlers.Big = new(Big)
