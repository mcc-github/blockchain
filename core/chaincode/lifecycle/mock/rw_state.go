
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
)

type ReadWritableState struct {
	GetStateStub        func(key string) (value []byte, err error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		key string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	PutStateStub        func(key string, value []byte) error
	putStateMutex       sync.RWMutex
	putStateArgsForCall []struct {
		key   string
		value []byte
	}
	putStateReturns struct {
		result1 error
	}
	putStateReturnsOnCall map[int]struct {
		result1 error
	}
	DelStateStub        func(key string) error
	delStateMutex       sync.RWMutex
	delStateArgsForCall []struct {
		key string
	}
	delStateReturns struct {
		result1 error
	}
	delStateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReadWritableState) GetState(key string) (value []byte, err error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("GetState", []interface{}{key})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateReturns.result1, fake.getStateReturns.result2
}

func (fake *ReadWritableState) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *ReadWritableState) GetStateArgsForCall(i int) string {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return fake.getStateArgsForCall[i].key
}

func (fake *ReadWritableState) GetStateReturns(result1 []byte, result2 error) {
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *ReadWritableState) PutState(key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.putStateMutex.Lock()
	ret, specificReturn := fake.putStateReturnsOnCall[len(fake.putStateArgsForCall)]
	fake.putStateArgsForCall = append(fake.putStateArgsForCall, struct {
		key   string
		value []byte
	}{key, valueCopy})
	fake.recordInvocation("PutState", []interface{}{key, valueCopy})
	fake.putStateMutex.Unlock()
	if fake.PutStateStub != nil {
		return fake.PutStateStub(key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putStateReturns.result1
}

func (fake *ReadWritableState) PutStateCallCount() int {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	return len(fake.putStateArgsForCall)
}

func (fake *ReadWritableState) PutStateArgsForCall(i int) (string, []byte) {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	return fake.putStateArgsForCall[i].key, fake.putStateArgsForCall[i].value
}

func (fake *ReadWritableState) PutStateReturns(result1 error) {
	fake.PutStateStub = nil
	fake.putStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) PutStateReturnsOnCall(i int, result1 error) {
	fake.PutStateStub = nil
	if fake.putStateReturnsOnCall == nil {
		fake.putStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) DelState(key string) error {
	fake.delStateMutex.Lock()
	ret, specificReturn := fake.delStateReturnsOnCall[len(fake.delStateArgsForCall)]
	fake.delStateArgsForCall = append(fake.delStateArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("DelState", []interface{}{key})
	fake.delStateMutex.Unlock()
	if fake.DelStateStub != nil {
		return fake.DelStateStub(key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.delStateReturns.result1
}

func (fake *ReadWritableState) DelStateCallCount() int {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	return len(fake.delStateArgsForCall)
}

func (fake *ReadWritableState) DelStateArgsForCall(i int) string {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	return fake.delStateArgsForCall[i].key
}

func (fake *ReadWritableState) DelStateReturns(result1 error) {
	fake.DelStateStub = nil
	fake.delStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) DelStateReturnsOnCall(i int, result1 error) {
	fake.DelStateStub = nil
	if fake.delStateReturnsOnCall == nil {
		fake.delStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ReadWritableState) recordInvocation(key string, args []interface{}) {
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

var _ lifecycle.ReadWritableState = new(ReadWritableState)
