
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/core/transientstore"
)

type StoreProvider struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	OpenStoreStub        func(string) (transientstore.Store, error)
	openStoreMutex       sync.RWMutex
	openStoreArgsForCall []struct {
		arg1 string
	}
	openStoreReturns struct {
		result1 transientstore.Store
		result2 error
	}
	openStoreReturnsOnCall map[int]struct {
		result1 transientstore.Store
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StoreProvider) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *StoreProvider) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *StoreProvider) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *StoreProvider) OpenStore(arg1 string) (transientstore.Store, error) {
	fake.openStoreMutex.Lock()
	ret, specificReturn := fake.openStoreReturnsOnCall[len(fake.openStoreArgsForCall)]
	fake.openStoreArgsForCall = append(fake.openStoreArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("OpenStore", []interface{}{arg1})
	fake.openStoreMutex.Unlock()
	if fake.OpenStoreStub != nil {
		return fake.OpenStoreStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.openStoreReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *StoreProvider) OpenStoreCallCount() int {
	fake.openStoreMutex.RLock()
	defer fake.openStoreMutex.RUnlock()
	return len(fake.openStoreArgsForCall)
}

func (fake *StoreProvider) OpenStoreCalls(stub func(string) (transientstore.Store, error)) {
	fake.openStoreMutex.Lock()
	defer fake.openStoreMutex.Unlock()
	fake.OpenStoreStub = stub
}

func (fake *StoreProvider) OpenStoreArgsForCall(i int) string {
	fake.openStoreMutex.RLock()
	defer fake.openStoreMutex.RUnlock()
	argsForCall := fake.openStoreArgsForCall[i]
	return argsForCall.arg1
}

func (fake *StoreProvider) OpenStoreReturns(result1 transientstore.Store, result2 error) {
	fake.openStoreMutex.Lock()
	defer fake.openStoreMutex.Unlock()
	fake.OpenStoreStub = nil
	fake.openStoreReturns = struct {
		result1 transientstore.Store
		result2 error
	}{result1, result2}
}

func (fake *StoreProvider) OpenStoreReturnsOnCall(i int, result1 transientstore.Store, result2 error) {
	fake.openStoreMutex.Lock()
	defer fake.openStoreMutex.Unlock()
	fake.OpenStoreStub = nil
	if fake.openStoreReturnsOnCall == nil {
		fake.openStoreReturnsOnCall = make(map[int]struct {
			result1 transientstore.Store
			result2 error
		})
	}
	fake.openStoreReturnsOnCall[i] = struct {
		result1 transientstore.Store
		result2 error
	}{result1, result2}
}

func (fake *StoreProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.openStoreMutex.RLock()
	defer fake.openStoreMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StoreProvider) recordInvocation(key string, args []interface{}) {
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
