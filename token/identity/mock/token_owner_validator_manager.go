
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/identity"
)

type TokenOwnerValidatorManager struct {
	GetStub        func(channel string) (identity.TokenOwnerValidator, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		channel string
	}
	getReturns struct {
		result1 identity.TokenOwnerValidator
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 identity.TokenOwnerValidator
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TokenOwnerValidatorManager) Get(channel string) (identity.TokenOwnerValidator, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("Get", []interface{}{channel})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(channel)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getReturns.result1, fake.getReturns.result2
}

func (fake *TokenOwnerValidatorManager) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *TokenOwnerValidatorManager) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].channel
}

func (fake *TokenOwnerValidatorManager) GetReturns(result1 identity.TokenOwnerValidator, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 identity.TokenOwnerValidator
		result2 error
	}{result1, result2}
}

func (fake *TokenOwnerValidatorManager) GetReturnsOnCall(i int, result1 identity.TokenOwnerValidator, result2 error) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 identity.TokenOwnerValidator
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 identity.TokenOwnerValidator
		result2 error
	}{result1, result2}
}

func (fake *TokenOwnerValidatorManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TokenOwnerValidatorManager) recordInvocation(key string, args []interface{}) {
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

var _ identity.TokenOwnerValidatorManager = new(TokenOwnerValidatorManager)
