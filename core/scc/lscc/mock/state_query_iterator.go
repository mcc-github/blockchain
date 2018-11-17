
package mock

import (
	sync "sync"

	queryresult "github.com/mcc-github/blockchain/protos/ledger/queryresult"
)

type StateQueryIterator struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	HasNextStub        func() bool
	hasNextMutex       sync.RWMutex
	hasNextArgsForCall []struct {
	}
	hasNextReturns struct {
		result1 bool
	}
	hasNextReturnsOnCall map[int]struct {
		result1 bool
	}
	NextStub        func() (*queryresult.KV, error)
	nextMutex       sync.RWMutex
	nextArgsForCall []struct {
	}
	nextReturns struct {
		result1 *queryresult.KV
		result2 error
	}
	nextReturnsOnCall map[int]struct {
		result1 *queryresult.KV
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StateQueryIterator) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.closeReturns
	return fakeReturns.result1
}

func (fake *StateQueryIterator) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *StateQueryIterator) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *StateQueryIterator) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateQueryIterator) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateQueryIterator) HasNext() bool {
	fake.hasNextMutex.Lock()
	ret, specificReturn := fake.hasNextReturnsOnCall[len(fake.hasNextArgsForCall)]
	fake.hasNextArgsForCall = append(fake.hasNextArgsForCall, struct {
	}{})
	fake.recordInvocation("HasNext", []interface{}{})
	fake.hasNextMutex.Unlock()
	if fake.HasNextStub != nil {
		return fake.HasNextStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.hasNextReturns
	return fakeReturns.result1
}

func (fake *StateQueryIterator) HasNextCallCount() int {
	fake.hasNextMutex.RLock()
	defer fake.hasNextMutex.RUnlock()
	return len(fake.hasNextArgsForCall)
}

func (fake *StateQueryIterator) HasNextCalls(stub func() bool) {
	fake.hasNextMutex.Lock()
	defer fake.hasNextMutex.Unlock()
	fake.HasNextStub = stub
}

func (fake *StateQueryIterator) HasNextReturns(result1 bool) {
	fake.hasNextMutex.Lock()
	defer fake.hasNextMutex.Unlock()
	fake.HasNextStub = nil
	fake.hasNextReturns = struct {
		result1 bool
	}{result1}
}

func (fake *StateQueryIterator) HasNextReturnsOnCall(i int, result1 bool) {
	fake.hasNextMutex.Lock()
	defer fake.hasNextMutex.Unlock()
	fake.HasNextStub = nil
	if fake.hasNextReturnsOnCall == nil {
		fake.hasNextReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.hasNextReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *StateQueryIterator) Next() (*queryresult.KV, error) {
	fake.nextMutex.Lock()
	ret, specificReturn := fake.nextReturnsOnCall[len(fake.nextArgsForCall)]
	fake.nextArgsForCall = append(fake.nextArgsForCall, struct {
	}{})
	fake.recordInvocation("Next", []interface{}{})
	fake.nextMutex.Unlock()
	if fake.NextStub != nil {
		return fake.NextStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.nextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *StateQueryIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *StateQueryIterator) NextCalls(stub func() (*queryresult.KV, error)) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = stub
}

func (fake *StateQueryIterator) NextReturns(result1 *queryresult.KV, result2 error) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 *queryresult.KV
		result2 error
	}{result1, result2}
}

func (fake *StateQueryIterator) NextReturnsOnCall(i int, result1 *queryresult.KV, result2 error) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 *queryresult.KV
			result2 error
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 *queryresult.KV
		result2 error
	}{result1, result2}
}

func (fake *StateQueryIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.hasNextMutex.RLock()
	defer fake.hasNextMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StateQueryIterator) recordInvocation(key string, args []interface{}) {
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
