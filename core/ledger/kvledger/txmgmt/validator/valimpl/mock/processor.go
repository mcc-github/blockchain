
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
)

type Processor struct {
	GenerateSimulationResultsStub        func(txEnvelop *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error
	generateSimulationResultsMutex       sync.RWMutex
	generateSimulationResultsArgsForCall []struct {
		txEnvelop          *common.Envelope
		simulator          ledger.TxSimulator
		initializingLedger bool
	}
	generateSimulationResultsReturns struct {
		result1 error
	}
	generateSimulationResultsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Processor) GenerateSimulationResults(txEnvelop *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error {
	fake.generateSimulationResultsMutex.Lock()
	ret, specificReturn := fake.generateSimulationResultsReturnsOnCall[len(fake.generateSimulationResultsArgsForCall)]
	fake.generateSimulationResultsArgsForCall = append(fake.generateSimulationResultsArgsForCall, struct {
		txEnvelop          *common.Envelope
		simulator          ledger.TxSimulator
		initializingLedger bool
	}{txEnvelop, simulator, initializingLedger})
	fake.recordInvocation("GenerateSimulationResults", []interface{}{txEnvelop, simulator, initializingLedger})
	fake.generateSimulationResultsMutex.Unlock()
	if fake.GenerateSimulationResultsStub != nil {
		return fake.GenerateSimulationResultsStub(txEnvelop, simulator, initializingLedger)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.generateSimulationResultsReturns.result1
}

func (fake *Processor) GenerateSimulationResultsCallCount() int {
	fake.generateSimulationResultsMutex.RLock()
	defer fake.generateSimulationResultsMutex.RUnlock()
	return len(fake.generateSimulationResultsArgsForCall)
}

func (fake *Processor) GenerateSimulationResultsArgsForCall(i int) (*common.Envelope, ledger.TxSimulator, bool) {
	fake.generateSimulationResultsMutex.RLock()
	defer fake.generateSimulationResultsMutex.RUnlock()
	return fake.generateSimulationResultsArgsForCall[i].txEnvelop, fake.generateSimulationResultsArgsForCall[i].simulator, fake.generateSimulationResultsArgsForCall[i].initializingLedger
}

func (fake *Processor) GenerateSimulationResultsReturns(result1 error) {
	fake.GenerateSimulationResultsStub = nil
	fake.generateSimulationResultsReturns = struct {
		result1 error
	}{result1}
}

func (fake *Processor) GenerateSimulationResultsReturnsOnCall(i int, result1 error) {
	fake.GenerateSimulationResultsStub = nil
	if fake.generateSimulationResultsReturnsOnCall == nil {
		fake.generateSimulationResultsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.generateSimulationResultsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Processor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateSimulationResultsMutex.RLock()
	defer fake.generateSimulationResultsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Processor) recordInvocation(key string, args []interface{}) {
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
