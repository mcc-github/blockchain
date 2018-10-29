
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/transaction"
)

type TMSManager struct {
	GetTxProcessorStub        func(channel string) (transaction.TMSTxProcessor, error)
	getTxProcessorMutex       sync.RWMutex
	getTxProcessorArgsForCall []struct {
		channel string
	}
	getTxProcessorReturns struct {
		result1 transaction.TMSTxProcessor
		result2 error
	}
	getTxProcessorReturnsOnCall map[int]struct {
		result1 transaction.TMSTxProcessor
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TMSManager) GetTxProcessor(channel string) (transaction.TMSTxProcessor, error) {
	fake.getTxProcessorMutex.Lock()
	ret, specificReturn := fake.getTxProcessorReturnsOnCall[len(fake.getTxProcessorArgsForCall)]
	fake.getTxProcessorArgsForCall = append(fake.getTxProcessorArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("GetTxProcessor", []interface{}{channel})
	fake.getTxProcessorMutex.Unlock()
	if fake.GetTxProcessorStub != nil {
		return fake.GetTxProcessorStub(channel)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTxProcessorReturns.result1, fake.getTxProcessorReturns.result2
}

func (fake *TMSManager) GetTxProcessorCallCount() int {
	fake.getTxProcessorMutex.RLock()
	defer fake.getTxProcessorMutex.RUnlock()
	return len(fake.getTxProcessorArgsForCall)
}

func (fake *TMSManager) GetTxProcessorArgsForCall(i int) string {
	fake.getTxProcessorMutex.RLock()
	defer fake.getTxProcessorMutex.RUnlock()
	return fake.getTxProcessorArgsForCall[i].channel
}

func (fake *TMSManager) GetTxProcessorReturns(result1 transaction.TMSTxProcessor, result2 error) {
	fake.GetTxProcessorStub = nil
	fake.getTxProcessorReturns = struct {
		result1 transaction.TMSTxProcessor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetTxProcessorReturnsOnCall(i int, result1 transaction.TMSTxProcessor, result2 error) {
	fake.GetTxProcessorStub = nil
	if fake.getTxProcessorReturnsOnCall == nil {
		fake.getTxProcessorReturnsOnCall = make(map[int]struct {
			result1 transaction.TMSTxProcessor
			result2 error
		})
	}
	fake.getTxProcessorReturnsOnCall[i] = struct {
		result1 transaction.TMSTxProcessor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getTxProcessorMutex.RLock()
	defer fake.getTxProcessorMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TMSManager) recordInvocation(key string, args []interface{}) {
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

var _ transaction.TMSManager = new(TMSManager)
