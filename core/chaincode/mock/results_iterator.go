
package mock

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
)

type QueryResultsIterator struct {
	NextStub        func() (commonledger.QueryResult, error)
	nextMutex       sync.RWMutex
	nextArgsForCall []struct{}
	nextReturns     struct {
		result1 commonledger.QueryResult
		result2 error
	}
	nextReturnsOnCall map[int]struct {
		result1 commonledger.QueryResult
		result2 error
	}
	CloseStub                      func()
	closeMutex                     sync.RWMutex
	closeArgsForCall               []struct{}
	GetBookmarkAndCloseStub        func() string
	getBookmarkAndCloseMutex       sync.RWMutex
	getBookmarkAndCloseArgsForCall []struct{}
	getBookmarkAndCloseReturns     struct {
		result1 string
	}
	getBookmarkAndCloseReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *QueryResultsIterator) Next() (commonledger.QueryResult, error) {
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

func (fake *QueryResultsIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *QueryResultsIterator) NextReturns(result1 commonledger.QueryResult, result2 error) {
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 commonledger.QueryResult
		result2 error
	}{result1, result2}
}

func (fake *QueryResultsIterator) NextReturnsOnCall(i int, result1 commonledger.QueryResult, result2 error) {
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 commonledger.QueryResult
			result2 error
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 commonledger.QueryResult
		result2 error
	}{result1, result2}
}

func (fake *QueryResultsIterator) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *QueryResultsIterator) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *QueryResultsIterator) GetBookmarkAndClose() string {
	fake.getBookmarkAndCloseMutex.Lock()
	ret, specificReturn := fake.getBookmarkAndCloseReturnsOnCall[len(fake.getBookmarkAndCloseArgsForCall)]
	fake.getBookmarkAndCloseArgsForCall = append(fake.getBookmarkAndCloseArgsForCall, struct{}{})
	fake.recordInvocation("GetBookmarkAndClose", []interface{}{})
	fake.getBookmarkAndCloseMutex.Unlock()
	if fake.GetBookmarkAndCloseStub != nil {
		return fake.GetBookmarkAndCloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getBookmarkAndCloseReturns.result1
}

func (fake *QueryResultsIterator) GetBookmarkAndCloseCallCount() int {
	fake.getBookmarkAndCloseMutex.RLock()
	defer fake.getBookmarkAndCloseMutex.RUnlock()
	return len(fake.getBookmarkAndCloseArgsForCall)
}

func (fake *QueryResultsIterator) GetBookmarkAndCloseReturns(result1 string) {
	fake.GetBookmarkAndCloseStub = nil
	fake.getBookmarkAndCloseReturns = struct {
		result1 string
	}{result1}
}

func (fake *QueryResultsIterator) GetBookmarkAndCloseReturnsOnCall(i int, result1 string) {
	fake.GetBookmarkAndCloseStub = nil
	if fake.getBookmarkAndCloseReturnsOnCall == nil {
		fake.getBookmarkAndCloseReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getBookmarkAndCloseReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *QueryResultsIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.getBookmarkAndCloseMutex.RLock()
	defer fake.getBookmarkAndCloseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *QueryResultsIterator) recordInvocation(key string, args []interface{}) {
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
