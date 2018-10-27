
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/transaction"
)

type PolicyValidator struct {
	IsIssuerStub        func(creator transaction.CreatorInfo, tokenType string) error
	isIssuerMutex       sync.RWMutex
	isIssuerArgsForCall []struct {
		creator   transaction.CreatorInfo
		tokenType string
	}
	isIssuerReturns struct {
		result1 error
	}
	isIssuerReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PolicyValidator) IsIssuer(creator transaction.CreatorInfo, tokenType string) error {
	fake.isIssuerMutex.Lock()
	ret, specificReturn := fake.isIssuerReturnsOnCall[len(fake.isIssuerArgsForCall)]
	fake.isIssuerArgsForCall = append(fake.isIssuerArgsForCall, struct {
		creator   transaction.CreatorInfo
		tokenType string
	}{creator, tokenType})
	fake.recordInvocation("IsIssuer", []interface{}{creator, tokenType})
	fake.isIssuerMutex.Unlock()
	if fake.IsIssuerStub != nil {
		return fake.IsIssuerStub(creator, tokenType)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isIssuerReturns.result1
}

func (fake *PolicyValidator) IsIssuerCallCount() int {
	fake.isIssuerMutex.RLock()
	defer fake.isIssuerMutex.RUnlock()
	return len(fake.isIssuerArgsForCall)
}

func (fake *PolicyValidator) IsIssuerArgsForCall(i int) (transaction.CreatorInfo, string) {
	fake.isIssuerMutex.RLock()
	defer fake.isIssuerMutex.RUnlock()
	return fake.isIssuerArgsForCall[i].creator, fake.isIssuerArgsForCall[i].tokenType
}

func (fake *PolicyValidator) IsIssuerReturns(result1 error) {
	fake.IsIssuerStub = nil
	fake.isIssuerReturns = struct {
		result1 error
	}{result1}
}

func (fake *PolicyValidator) IsIssuerReturnsOnCall(i int, result1 error) {
	fake.IsIssuerStub = nil
	if fake.isIssuerReturnsOnCall == nil {
		fake.isIssuerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.isIssuerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PolicyValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isIssuerMutex.RLock()
	defer fake.isIssuerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PolicyValidator) recordInvocation(key string, args []interface{}) {
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

var _ transaction.PolicyValidator = new(PolicyValidator)
