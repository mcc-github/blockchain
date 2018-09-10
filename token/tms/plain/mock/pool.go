
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/tms"
	"github.com/mcc-github/blockchain/token/tms/plain"
)

type Pool struct {
	CommitUpdateStub        func(transactionData []tms.TransactionData) error
	commitUpdateMutex       sync.RWMutex
	commitUpdateArgsForCall []struct {
		transactionData []tms.TransactionData
	}
	commitUpdateReturns struct {
		result1 error
	}
	commitUpdateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Pool) CommitUpdate(transactionData []tms.TransactionData) error {
	var transactionDataCopy []tms.TransactionData
	if transactionData != nil {
		transactionDataCopy = make([]tms.TransactionData, len(transactionData))
		copy(transactionDataCopy, transactionData)
	}
	fake.commitUpdateMutex.Lock()
	ret, specificReturn := fake.commitUpdateReturnsOnCall[len(fake.commitUpdateArgsForCall)]
	fake.commitUpdateArgsForCall = append(fake.commitUpdateArgsForCall, struct {
		transactionData []tms.TransactionData
	}{transactionDataCopy})
	fake.recordInvocation("CommitUpdate", []interface{}{transactionDataCopy})
	fake.commitUpdateMutex.Unlock()
	if fake.CommitUpdateStub != nil {
		return fake.CommitUpdateStub(transactionData)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.commitUpdateReturns.result1
}

func (fake *Pool) CommitUpdateCallCount() int {
	fake.commitUpdateMutex.RLock()
	defer fake.commitUpdateMutex.RUnlock()
	return len(fake.commitUpdateArgsForCall)
}

func (fake *Pool) CommitUpdateArgsForCall(i int) []tms.TransactionData {
	fake.commitUpdateMutex.RLock()
	defer fake.commitUpdateMutex.RUnlock()
	return fake.commitUpdateArgsForCall[i].transactionData
}

func (fake *Pool) CommitUpdateReturns(result1 error) {
	fake.CommitUpdateStub = nil
	fake.commitUpdateReturns = struct {
		result1 error
	}{result1}
}

func (fake *Pool) CommitUpdateReturnsOnCall(i int, result1 error) {
	fake.CommitUpdateStub = nil
	if fake.commitUpdateReturnsOnCall == nil {
		fake.commitUpdateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commitUpdateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Pool) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.commitUpdateMutex.RLock()
	defer fake.commitUpdateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Pool) recordInvocation(key string, args []interface{}) {
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

var _ plain.Pool = new(Pool)