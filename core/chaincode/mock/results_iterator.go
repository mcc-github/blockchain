
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/ledger"
)

type ResultsIterator struct {
	NextStub        func() (ledger.QueryResult, error)
	nextMutex       sync.RWMutex
	nextArgsForCall []struct{}
	nextReturns     struct {
		result1 ledger.QueryResult
		result2 error
	}
	nextReturnsOnCall map[int]struct {
		result1 ledger.QueryResult
		result2 error
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ResultsIterator) Next() (ledger.QueryResult, error) {
	fake.nextMutex.Lock()
	ret, specificReturn := fake.nextReturnsOnCall[len(fake.nextArgsForCall)]
	fake.nextArgsForCall = append(fake.nextArgsForCall, struct{}{})
	fake.recordInvocation("Next", []interface{}{})
	fake.nextMutex.Unlock()
	if fake.NextStub != nil {
		return fake.NextStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.nextReturns.result1, fake.nextReturns.result2
}

func (fake *ResultsIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *ResultsIterator) NextReturns(result1 ledger.QueryResult, result2 error) {
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 ledger.QueryResult
		result2 error
	}{result1, result2}
}

func (fake *ResultsIterator) NextReturnsOnCall(i int, result1 ledger.QueryResult, result2 error) {
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 ledger.QueryResult
			result2 error
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 ledger.QueryResult
		result2 error
	}{result1, result2}
}

func (fake *ResultsIterator) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *ResultsIterator) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *ResultsIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ResultsIterator) recordInvocation(key string, args []interface{}) {
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
