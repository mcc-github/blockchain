
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/transaction"
)

type LedgerWriter struct {
	GetStateStub        func(namespace string, key string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		namespace string
		key       string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SetStateStub        func(namespace string, key string, value []byte) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		namespace string
		key       string
		value     []byte
	}
	setStateReturns struct {
		result1 error
	}
	setStateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LedgerWriter) GetState(namespace string, key string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		namespace string
		key       string
	}{namespace, key})
	fake.recordInvocation("GetState", []interface{}{namespace, key})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(namespace, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateReturns.result1, fake.getStateReturns.result2
}

func (fake *LedgerWriter) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *LedgerWriter) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return fake.getStateArgsForCall[i].namespace, fake.getStateArgsForCall[i].key
}

func (fake *LedgerWriter) GetStateReturns(result1 []byte, result2 error) {
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *LedgerWriter) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetStateStub = nil
	if fake.getStateReturnsOnCall == nil {
		fake.getStateReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getStateReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *LedgerWriter) SetState(namespace string, key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		namespace string
		key       string
		value     []byte
	}{namespace, key, valueCopy})
	fake.recordInvocation("SetState", []interface{}{namespace, key, valueCopy})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(namespace, key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateReturns.result1
}

func (fake *LedgerWriter) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *LedgerWriter) SetStateArgsForCall(i int) (string, string, []byte) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return fake.setStateArgsForCall[i].namespace, fake.setStateArgsForCall[i].key, fake.setStateArgsForCall[i].value
}

func (fake *LedgerWriter) SetStateReturns(result1 error) {
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *LedgerWriter) SetStateReturnsOnCall(i int, result1 error) {
	fake.SetStateStub = nil
	if fake.setStateReturnsOnCall == nil {
		fake.setStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *LedgerWriter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LedgerWriter) recordInvocation(key string, args []interface{}) {
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

var _ transaction.LedgerWriter = new(LedgerWriter)
