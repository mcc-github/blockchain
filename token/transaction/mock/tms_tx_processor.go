
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/transaction"
)

type TMSTxProcessor struct {
	ProcessTxStub        func(txID string, creator transaction.CreatorInfo, ttx *token.TokenTransaction, simulator transaction.LedgerWriter) error
	processTxMutex       sync.RWMutex
	processTxArgsForCall []struct {
		txID      string
		creator   transaction.CreatorInfo
		ttx       *token.TokenTransaction
		simulator transaction.LedgerWriter
	}
	processTxReturns struct {
		result1 error
	}
	processTxReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TMSTxProcessor) ProcessTx(txID string, creator transaction.CreatorInfo, ttx *token.TokenTransaction, simulator transaction.LedgerWriter) error {
	fake.processTxMutex.Lock()
	ret, specificReturn := fake.processTxReturnsOnCall[len(fake.processTxArgsForCall)]
	fake.processTxArgsForCall = append(fake.processTxArgsForCall, struct {
		txID      string
		creator   transaction.CreatorInfo
		ttx       *token.TokenTransaction
		simulator transaction.LedgerWriter
	}{txID, creator, ttx, simulator})
	fake.recordInvocation("ProcessTx", []interface{}{txID, creator, ttx, simulator})
	fake.processTxMutex.Unlock()
	if fake.ProcessTxStub != nil {
		return fake.ProcessTxStub(txID, creator, ttx, simulator)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.processTxReturns.result1
}

func (fake *TMSTxProcessor) ProcessTxCallCount() int {
	fake.processTxMutex.RLock()
	defer fake.processTxMutex.RUnlock()
	return len(fake.processTxArgsForCall)
}

func (fake *TMSTxProcessor) ProcessTxArgsForCall(i int) (string, transaction.CreatorInfo, *token.TokenTransaction, transaction.LedgerWriter) {
	fake.processTxMutex.RLock()
	defer fake.processTxMutex.RUnlock()
	return fake.processTxArgsForCall[i].txID, fake.processTxArgsForCall[i].creator, fake.processTxArgsForCall[i].ttx, fake.processTxArgsForCall[i].simulator
}

func (fake *TMSTxProcessor) ProcessTxReturns(result1 error) {
	fake.ProcessTxStub = nil
	fake.processTxReturns = struct {
		result1 error
	}{result1}
}

func (fake *TMSTxProcessor) ProcessTxReturnsOnCall(i int, result1 error) {
	fake.ProcessTxStub = nil
	if fake.processTxReturnsOnCall == nil {
		fake.processTxReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.processTxReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *TMSTxProcessor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.processTxMutex.RLock()
	defer fake.processTxMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TMSTxProcessor) recordInvocation(key string, args []interface{}) {
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

var _ transaction.TMSTxProcessor = new(TMSTxProcessor)
