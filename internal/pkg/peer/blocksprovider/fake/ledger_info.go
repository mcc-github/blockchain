
package fake

import (
	"sync"

	"github.com/mcc-github/blockchain/internal/pkg/peer/blocksprovider"
)

type LedgerInfo struct {
	LedgerHeightStub        func() (uint64, error)
	ledgerHeightMutex       sync.RWMutex
	ledgerHeightArgsForCall []struct {
	}
	ledgerHeightReturns struct {
		result1 uint64
		result2 error
	}
	ledgerHeightReturnsOnCall map[int]struct {
		result1 uint64
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LedgerInfo) LedgerHeight() (uint64, error) {
	fake.ledgerHeightMutex.Lock()
	ret, specificReturn := fake.ledgerHeightReturnsOnCall[len(fake.ledgerHeightArgsForCall)]
	fake.ledgerHeightArgsForCall = append(fake.ledgerHeightArgsForCall, struct {
	}{})
	fake.recordInvocation("LedgerHeight", []interface{}{})
	fake.ledgerHeightMutex.Unlock()
	if fake.LedgerHeightStub != nil {
		return fake.LedgerHeightStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.ledgerHeightReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *LedgerInfo) LedgerHeightCallCount() int {
	fake.ledgerHeightMutex.RLock()
	defer fake.ledgerHeightMutex.RUnlock()
	return len(fake.ledgerHeightArgsForCall)
}

func (fake *LedgerInfo) LedgerHeightCalls(stub func() (uint64, error)) {
	fake.ledgerHeightMutex.Lock()
	defer fake.ledgerHeightMutex.Unlock()
	fake.LedgerHeightStub = stub
}

func (fake *LedgerInfo) LedgerHeightReturns(result1 uint64, result2 error) {
	fake.ledgerHeightMutex.Lock()
	defer fake.ledgerHeightMutex.Unlock()
	fake.LedgerHeightStub = nil
	fake.ledgerHeightReturns = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *LedgerInfo) LedgerHeightReturnsOnCall(i int, result1 uint64, result2 error) {
	fake.ledgerHeightMutex.Lock()
	defer fake.ledgerHeightMutex.Unlock()
	fake.LedgerHeightStub = nil
	if fake.ledgerHeightReturnsOnCall == nil {
		fake.ledgerHeightReturnsOnCall = make(map[int]struct {
			result1 uint64
			result2 error
		})
	}
	fake.ledgerHeightReturnsOnCall[i] = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *LedgerInfo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.ledgerHeightMutex.RLock()
	defer fake.ledgerHeightMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LedgerInfo) recordInvocation(key string, args []interface{}) {
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

var _ blocksprovider.LedgerInfo = new(LedgerInfo)
