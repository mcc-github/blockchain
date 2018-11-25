
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type User struct {
	NewKeyStub        func() (handlers.Big, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct{}
	newKeyReturns     struct {
		result1 handlers.Big
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 handlers.Big
		result2 error
	}
	MakeNymStub        func(sk handlers.Big, key handlers.IssuerPublicKey) (handlers.Ecp, handlers.Big, error)
	makeNymMutex       sync.RWMutex
	makeNymArgsForCall []struct {
		sk  handlers.Big
		key handlers.IssuerPublicKey
	}
	makeNymReturns struct {
		result1 handlers.Ecp
		result2 handlers.Big
		result3 error
	}
	makeNymReturnsOnCall map[int]struct {
		result1 handlers.Ecp
		result2 handlers.Big
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *User) NewKey() (handlers.Big, error) {
	fake.newKeyMutex.Lock()
	ret, specificReturn := fake.newKeyReturnsOnCall[len(fake.newKeyArgsForCall)]
	fake.newKeyArgsForCall = append(fake.newKeyArgsForCall, struct{}{})
	fake.recordInvocation("NewKey", []interface{}{})
	fake.newKeyMutex.Unlock()
	if fake.NewKeyStub != nil {
		return fake.NewKeyStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newKeyReturns.result1, fake.newKeyReturns.result2
}

func (fake *User) NewKeyCallCount() int {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	return len(fake.newKeyArgsForCall)
}

func (fake *User) NewKeyReturns(result1 handlers.Big, result2 error) {
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 handlers.Big
		result2 error
	}{result1, result2}
}

func (fake *User) NewKeyReturnsOnCall(i int, result1 handlers.Big, result2 error) {
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 handlers.Big
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 handlers.Big
		result2 error
	}{result1, result2}
}

func (fake *User) MakeNym(sk handlers.Big, key handlers.IssuerPublicKey) (handlers.Ecp, handlers.Big, error) {
	fake.makeNymMutex.Lock()
	ret, specificReturn := fake.makeNymReturnsOnCall[len(fake.makeNymArgsForCall)]
	fake.makeNymArgsForCall = append(fake.makeNymArgsForCall, struct {
		sk  handlers.Big
		key handlers.IssuerPublicKey
	}{sk, key})
	fake.recordInvocation("MakeNym", []interface{}{sk, key})
	fake.makeNymMutex.Unlock()
	if fake.MakeNymStub != nil {
		return fake.MakeNymStub(sk, key)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.makeNymReturns.result1, fake.makeNymReturns.result2, fake.makeNymReturns.result3
}

func (fake *User) MakeNymCallCount() int {
	fake.makeNymMutex.RLock()
	defer fake.makeNymMutex.RUnlock()
	return len(fake.makeNymArgsForCall)
}

func (fake *User) MakeNymArgsForCall(i int) (handlers.Big, handlers.IssuerPublicKey) {
	fake.makeNymMutex.RLock()
	defer fake.makeNymMutex.RUnlock()
	return fake.makeNymArgsForCall[i].sk, fake.makeNymArgsForCall[i].key
}

func (fake *User) MakeNymReturns(result1 handlers.Ecp, result2 handlers.Big, result3 error) {
	fake.MakeNymStub = nil
	fake.makeNymReturns = struct {
		result1 handlers.Ecp
		result2 handlers.Big
		result3 error
	}{result1, result2, result3}
}

func (fake *User) MakeNymReturnsOnCall(i int, result1 handlers.Ecp, result2 handlers.Big, result3 error) {
	fake.MakeNymStub = nil
	if fake.makeNymReturnsOnCall == nil {
		fake.makeNymReturnsOnCall = make(map[int]struct {
			result1 handlers.Ecp
			result2 handlers.Big
			result3 error
		})
	}
	fake.makeNymReturnsOnCall[i] = struct {
		result1 handlers.Ecp
		result2 handlers.Big
		result3 error
	}{result1, result2, result3}
}

func (fake *User) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	fake.makeNymMutex.RLock()
	defer fake.makeNymMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *User) recordInvocation(key string, args []interface{}) {
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

var _ handlers.User = new(User)
