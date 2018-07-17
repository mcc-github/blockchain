
package mock

import (
	"sync"
)

type TransactionRegistry struct {
	AddStub        func(channelID, txID string) bool
	addMutex       sync.RWMutex
	addArgsForCall []struct {
		channelID string
		txID      string
	}
	addReturns struct {
		result1 bool
	}
	addReturnsOnCall map[int]struct {
		result1 bool
	}
	RemoveStub        func(channelID, txID string)
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
		channelID string
		txID      string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TransactionRegistry) Add(channelID string, txID string) bool {
	fake.addMutex.Lock()
	ret, specificReturn := fake.addReturnsOnCall[len(fake.addArgsForCall)]
	fake.addArgsForCall = append(fake.addArgsForCall, struct {
		channelID string
		txID      string
	}{channelID, txID})
	fake.recordInvocation("Add", []interface{}{channelID, txID})
	fake.addMutex.Unlock()
	if fake.AddStub != nil {
		return fake.AddStub(channelID, txID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.addReturns.result1
}

func (fake *TransactionRegistry) AddCallCount() int {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return len(fake.addArgsForCall)
}

func (fake *TransactionRegistry) AddArgsForCall(i int) (string, string) {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return fake.addArgsForCall[i].channelID, fake.addArgsForCall[i].txID
}

func (fake *TransactionRegistry) AddReturns(result1 bool) {
	fake.AddStub = nil
	fake.addReturns = struct {
		result1 bool
	}{result1}
}

func (fake *TransactionRegistry) AddReturnsOnCall(i int, result1 bool) {
	fake.AddStub = nil
	if fake.addReturnsOnCall == nil {
		fake.addReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.addReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *TransactionRegistry) Remove(channelID string, txID string) {
	fake.removeMutex.Lock()
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
		channelID string
		txID      string
	}{channelID, txID})
	fake.recordInvocation("Remove", []interface{}{channelID, txID})
	fake.removeMutex.Unlock()
	if fake.RemoveStub != nil {
		fake.RemoveStub(channelID, txID)
	}
}

func (fake *TransactionRegistry) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *TransactionRegistry) RemoveArgsForCall(i int) (string, string) {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return fake.removeArgsForCall[i].channelID, fake.removeArgsForCall[i].txID
}

func (fake *TransactionRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TransactionRegistry) recordInvocation(key string, args []interface{}) {
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
