
package mock

import (
	"sync"

	ledgera "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/token/ledger"
)

type LedgerWriter struct {
	DeleteStateStub        func(string, string) error
	deleteStateMutex       sync.RWMutex
	deleteStateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteStateReturns struct {
		result1 error
	}
	deleteStateReturnsOnCall map[int]struct {
		result1 error
	}
	DoneStub        func()
	doneMutex       sync.RWMutex
	doneArgsForCall []struct {
	}
	GetStateStub        func(string, string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateRangeScanIteratorStub        func(string, string, string) (ledgera.ResultsIterator, error)
	getStateRangeScanIteratorMutex       sync.RWMutex
	getStateRangeScanIteratorArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getStateRangeScanIteratorReturns struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	getStateRangeScanIteratorReturnsOnCall map[int]struct {
		result1 ledgera.ResultsIterator
		result2 error
	}
	SetStateStub        func(string, string, []byte) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []byte
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

func (fake *LedgerWriter) DeleteState(arg1 string, arg2 string) error {
	fake.deleteStateMutex.Lock()
	ret, specificReturn := fake.deleteStateReturnsOnCall[len(fake.deleteStateArgsForCall)]
	fake.deleteStateArgsForCall = append(fake.deleteStateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DeleteState", []interface{}{arg1, arg2})
	fake.deleteStateMutex.Unlock()
	if fake.DeleteStateStub != nil {
		return fake.DeleteStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteStateReturns
	return fakeReturns.result1
}

func (fake *LedgerWriter) DeleteStateCallCount() int {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	return len(fake.deleteStateArgsForCall)
}

func (fake *LedgerWriter) DeleteStateCalls(stub func(string, string) error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = stub
}

func (fake *LedgerWriter) DeleteStateArgsForCall(i int) (string, string) {
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	argsForCall := fake.deleteStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *LedgerWriter) DeleteStateReturns(result1 error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = nil
	fake.deleteStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *LedgerWriter) DeleteStateReturnsOnCall(i int, result1 error) {
	fake.deleteStateMutex.Lock()
	defer fake.deleteStateMutex.Unlock()
	fake.DeleteStateStub = nil
	if fake.deleteStateReturnsOnCall == nil {
		fake.deleteStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *LedgerWriter) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct {
	}{})
	fake.recordInvocation("Done", []interface{}{})
	fake.doneMutex.Unlock()
	if fake.DoneStub != nil {
		fake.DoneStub()
	}
}

func (fake *LedgerWriter) DoneCallCount() int {
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	return len(fake.doneArgsForCall)
}

func (fake *LedgerWriter) DoneCalls(stub func()) {
	fake.doneMutex.Lock()
	defer fake.doneMutex.Unlock()
	fake.DoneStub = stub
}

func (fake *LedgerWriter) GetState(arg1 string, arg2 string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetState", []interface{}{arg1, arg2})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *LedgerWriter) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *LedgerWriter) GetStateCalls(stub func(string, string) ([]byte, error)) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = stub
}

func (fake *LedgerWriter) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	argsForCall := fake.getStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *LedgerWriter) GetStateReturns(result1 []byte, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *LedgerWriter) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *LedgerWriter) GetStateRangeScanIterator(arg1 string, arg2 string, arg3 string) (ledgera.ResultsIterator, error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorReturnsOnCall[len(fake.getStateRangeScanIteratorArgsForCall)]
	fake.getStateRangeScanIteratorArgsForCall = append(fake.getStateRangeScanIteratorArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetStateRangeScanIterator", []interface{}{arg1, arg2, arg3})
	fake.getStateRangeScanIteratorMutex.Unlock()
	if fake.GetStateRangeScanIteratorStub != nil {
		return fake.GetStateRangeScanIteratorStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeScanIteratorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *LedgerWriter) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *LedgerWriter) GetStateRangeScanIteratorCalls(stub func(string, string, string) (ledgera.ResultsIterator, error)) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = stub
}

func (fake *LedgerWriter) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *LedgerWriter) GetStateRangeScanIteratorReturns(result1 ledgera.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *LedgerWriter) GetStateRangeScanIteratorReturnsOnCall(i int, result1 ledgera.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	if fake.getStateRangeScanIteratorReturnsOnCall == nil {
		fake.getStateRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 ledgera.ResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorReturnsOnCall[i] = struct {
		result1 ledgera.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *LedgerWriter) SetState(arg1 string, arg2 string, arg3 []byte) error {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("SetState", []interface{}{arg1, arg2, arg3Copy})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setStateReturns
	return fakeReturns.result1
}

func (fake *LedgerWriter) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *LedgerWriter) SetStateCalls(stub func(string, string, []byte) error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = stub
}

func (fake *LedgerWriter) SetStateArgsForCall(i int) (string, string, []byte) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	argsForCall := fake.setStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *LedgerWriter) SetStateReturns(result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *LedgerWriter) SetStateReturnsOnCall(i int, result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
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
	fake.deleteStateMutex.RLock()
	defer fake.deleteStateMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
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

var _ ledger.LedgerWriter = new(LedgerWriter)
