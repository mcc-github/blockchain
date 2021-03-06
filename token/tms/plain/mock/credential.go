
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/tms/plain"
)

type Credential struct {
	PublicStub        func() []byte
	publicMutex       sync.RWMutex
	publicArgsForCall []struct{}
	publicReturns     struct {
		result1 []byte
	}
	publicReturnsOnCall map[int]struct {
		result1 []byte
	}
	PrivateStub        func() []byte
	privateMutex       sync.RWMutex
	privateArgsForCall []struct{}
	privateReturns     struct {
		result1 []byte
	}
	privateReturnsOnCall map[int]struct {
		result1 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Credential) Public() []byte {
	fake.publicMutex.Lock()
	ret, specificReturn := fake.publicReturnsOnCall[len(fake.publicArgsForCall)]
	fake.publicArgsForCall = append(fake.publicArgsForCall, struct{}{})
	fake.recordInvocation("Public", []interface{}{})
	fake.publicMutex.Unlock()
	if fake.PublicStub != nil {
		return fake.PublicStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.publicReturns.result1
}

func (fake *Credential) PublicCallCount() int {
	fake.publicMutex.RLock()
	defer fake.publicMutex.RUnlock()
	return len(fake.publicArgsForCall)
}

func (fake *Credential) PublicReturns(result1 []byte) {
	fake.PublicStub = nil
	fake.publicReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *Credential) PublicReturnsOnCall(i int, result1 []byte) {
	fake.PublicStub = nil
	if fake.publicReturnsOnCall == nil {
		fake.publicReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.publicReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *Credential) Private() []byte {
	fake.privateMutex.Lock()
	ret, specificReturn := fake.privateReturnsOnCall[len(fake.privateArgsForCall)]
	fake.privateArgsForCall = append(fake.privateArgsForCall, struct{}{})
	fake.recordInvocation("Private", []interface{}{})
	fake.privateMutex.Unlock()
	if fake.PrivateStub != nil {
		return fake.PrivateStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.privateReturns.result1
}

func (fake *Credential) PrivateCallCount() int {
	fake.privateMutex.RLock()
	defer fake.privateMutex.RUnlock()
	return len(fake.privateArgsForCall)
}

func (fake *Credential) PrivateReturns(result1 []byte) {
	fake.PrivateStub = nil
	fake.privateReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *Credential) PrivateReturnsOnCall(i int, result1 []byte) {
	fake.PrivateStub = nil
	if fake.privateReturnsOnCall == nil {
		fake.privateReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.privateReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *Credential) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.publicMutex.RLock()
	defer fake.publicMutex.RUnlock()
	fake.privateMutex.RLock()
	defer fake.privateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Credential) recordInvocation(key string, args []interface{}) {
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

var _ plain.Credential = new(Credential)
