
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix"
)

type User struct {
	NewKeyStub        func() (idemix.Big, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct{}
	newKeyReturns     struct {
		result1 idemix.Big
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 idemix.Big
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *User) NewKey() (idemix.Big, error) {
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

func (fake *User) NewKeyReturns(result1 idemix.Big, result2 error) {
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 idemix.Big
		result2 error
	}{result1, result2}
}

func (fake *User) NewKeyReturnsOnCall(i int, result1 idemix.Big, result2 error) {
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 idemix.Big
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 idemix.Big
		result2 error
	}{result1, result2}
}

func (fake *User) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
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

var _ idemix.User = new(User)
