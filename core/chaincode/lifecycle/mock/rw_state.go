
package mock

import (
	"sync"
)

type ReadWritableState struct {
	DelStateStub        func(string) error
	delStateMutex       sync.RWMutex
	delStateArgsForCall []struct {
		arg1 string
	}
	delStateReturns struct {
		result1 error
	}
	delStateReturnsOnCall map[int]struct {
		result1 error
	}
	GetStateStub        func(string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		arg1 string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateHashStub        func(string) ([]byte, error)
	getStateHashMutex       sync.RWMutex
	getStateHashArgsForCall []struct {
		arg1 string
	}
	getStateHashReturns struct {
		result1 []byte
		result2 error
	}
	getStateHashReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateRangeStub        func(string) (map[string][]byte, error)
	getStateRangeMutex       sync.RWMutex
	getStateRangeArgsForCall []struct {
		arg1 string
	}
	getStateRangeReturns struct {
		result1 map[string][]byte
		result2 error
	}
	getStateRangeReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	PutStateStub        func(string, []byte) error
	putStateMutex       sync.RWMutex
	putStateArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	putStateReturns struct {
		result1 error
	}
	putStateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReadWritableState) DelState(arg1 string) error {
	fake.delStateMutex.Lock()
	ret, specificReturn := fake.delStateReturnsOnCall[len(fake.delStateArgsForCall)]
	fake.delStateArgsForCall = append(fake.delStateArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DelState", []interface{}{arg1})
	fake.delStateMutex.Unlock()
	if fake.DelStateStub != nil {
		return fake.DelStateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.delStateReturns
	return fakeReturns.result1
}

func (fake *ReadWritableState) DelStateCallCount() int {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	return len(fake.delStateArgsForCall)
}

func (fake *ReadWritableState) DelStateCalls(stub func(string) error) {
	fake.delStateMutex.Lock()
	defer fake.delStateMutex.Unlock()
	fake.DelStateStub = stub
}

func (fake *ReadWritableState) DelStateArgsForCall(i int) string {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	argsForCall := fake.delStateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ReadWritableState) DelStateReturns(result1 error) {
	fake.delStateMutex.Lock()
	defer fake.delStateMutex.Unlock()
	fake.DelStateStub = nil
	fake.delStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) DelStateReturnsOnCall(i int, result1 error) {
	fake.delStateMutex.Lock()
	defer fake.delStateMutex.Unlock()
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

func (fake *ReadWritableState) GetState(arg1 string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetState", []interface{}{arg1})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReadWritableState) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *ReadWritableState) GetStateCalls(stub func(string) ([]byte, error)) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = stub
}

func (fake *ReadWritableState) GetStateArgsForCall(i int) string {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	argsForCall := fake.getStateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ReadWritableState) GetStateReturns(result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
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

func (fake *ReadWritableState) GetStateHash(arg1 string) ([]byte, error) {
	fake.getStateHashMutex.Lock()
	ret, specificReturn := fake.getStateHashReturnsOnCall[len(fake.getStateHashArgsForCall)]
	fake.getStateHashArgsForCall = append(fake.getStateHashArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetStateHash", []interface{}{arg1})
	fake.getStateHashMutex.Unlock()
	if fake.GetStateHashStub != nil {
		return fake.GetStateHashStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateHashReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReadWritableState) GetStateHashCallCount() int {
	fake.getStateHashMutex.RLock()
	defer fake.getStateHashMutex.RUnlock()
	return len(fake.getStateHashArgsForCall)
}

func (fake *ReadWritableState) GetStateHashCalls(stub func(string) ([]byte, error)) {
	fake.getStateHashMutex.Lock()
	defer fake.getStateHashMutex.Unlock()
	fake.GetStateHashStub = stub
}

func (fake *ReadWritableState) GetStateHashArgsForCall(i int) string {
	fake.getStateHashMutex.RLock()
	defer fake.getStateHashMutex.RUnlock()
	argsForCall := fake.getStateHashArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ReadWritableState) GetStateHashReturns(result1 []byte, result2 error) {
	fake.getStateHashMutex.Lock()
	defer fake.getStateHashMutex.Unlock()
	fake.GetStateHashStub = nil
	fake.getStateHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) GetStateHashReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getStateHashMutex.Lock()
	defer fake.getStateHashMutex.Unlock()
	fake.GetStateHashStub = nil
	if fake.getStateHashReturnsOnCall == nil {
		fake.getStateHashReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getStateHashReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) GetStateRange(arg1 string) (map[string][]byte, error) {
	fake.getStateRangeMutex.Lock()
	ret, specificReturn := fake.getStateRangeReturnsOnCall[len(fake.getStateRangeArgsForCall)]
	fake.getStateRangeArgsForCall = append(fake.getStateRangeArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetStateRange", []interface{}{arg1})
	fake.getStateRangeMutex.Unlock()
	if fake.GetStateRangeStub != nil {
		return fake.GetStateRangeStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReadWritableState) GetStateRangeCallCount() int {
	fake.getStateRangeMutex.RLock()
	defer fake.getStateRangeMutex.RUnlock()
	return len(fake.getStateRangeArgsForCall)
}

func (fake *ReadWritableState) GetStateRangeCalls(stub func(string) (map[string][]byte, error)) {
	fake.getStateRangeMutex.Lock()
	defer fake.getStateRangeMutex.Unlock()
	fake.GetStateRangeStub = stub
}

func (fake *ReadWritableState) GetStateRangeArgsForCall(i int) string {
	fake.getStateRangeMutex.RLock()
	defer fake.getStateRangeMutex.RUnlock()
	argsForCall := fake.getStateRangeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ReadWritableState) GetStateRangeReturns(result1 map[string][]byte, result2 error) {
	fake.getStateRangeMutex.Lock()
	defer fake.getStateRangeMutex.Unlock()
	fake.GetStateRangeStub = nil
	fake.getStateRangeReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) GetStateRangeReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
	fake.getStateRangeMutex.Lock()
	defer fake.getStateRangeMutex.Unlock()
	fake.GetStateRangeStub = nil
	if fake.getStateRangeReturnsOnCall == nil {
		fake.getStateRangeReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
			result2 error
		})
	}
	fake.getStateRangeReturnsOnCall[i] = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *ReadWritableState) PutState(arg1 string, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.putStateMutex.Lock()
	ret, specificReturn := fake.putStateReturnsOnCall[len(fake.putStateArgsForCall)]
	fake.putStateArgsForCall = append(fake.putStateArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("PutState", []interface{}{arg1, arg2Copy})
	fake.putStateMutex.Unlock()
	if fake.PutStateStub != nil {
		return fake.PutStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.putStateReturns
	return fakeReturns.result1
}

func (fake *ReadWritableState) PutStateCallCount() int {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	return len(fake.putStateArgsForCall)
}

func (fake *ReadWritableState) PutStateCalls(stub func(string, []byte) error) {
	fake.putStateMutex.Lock()
	defer fake.putStateMutex.Unlock()
	fake.PutStateStub = stub
}

func (fake *ReadWritableState) PutStateArgsForCall(i int) (string, []byte) {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	argsForCall := fake.putStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ReadWritableState) PutStateReturns(result1 error) {
	fake.putStateMutex.Lock()
	defer fake.putStateMutex.Unlock()
	fake.PutStateStub = nil
	fake.putStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ReadWritableState) PutStateReturnsOnCall(i int, result1 error) {
	fake.putStateMutex.Lock()
	defer fake.putStateMutex.Unlock()
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

func (fake *ReadWritableState) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.getStateHashMutex.RLock()
	defer fake.getStateHashMutex.RUnlock()
	fake.getStateRangeMutex.RLock()
	defer fake.getStateRangeMutex.RUnlock()
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
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
