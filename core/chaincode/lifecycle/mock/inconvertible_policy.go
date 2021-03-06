
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protoutil"
)

type InconvertiblePolicy struct {
	EvaluateStub        func([]*protoutil.SignedData) error
	evaluateMutex       sync.RWMutex
	evaluateArgsForCall []struct {
		arg1 []*protoutil.SignedData
	}
	evaluateReturns struct {
		result1 error
	}
	evaluateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *InconvertiblePolicy) Evaluate(arg1 []*protoutil.SignedData) error {
	var arg1Copy []*protoutil.SignedData
	if arg1 != nil {
		arg1Copy = make([]*protoutil.SignedData, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.evaluateMutex.Lock()
	ret, specificReturn := fake.evaluateReturnsOnCall[len(fake.evaluateArgsForCall)]
	fake.evaluateArgsForCall = append(fake.evaluateArgsForCall, struct {
		arg1 []*protoutil.SignedData
	}{arg1Copy})
	fake.recordInvocation("Evaluate", []interface{}{arg1Copy})
	fake.evaluateMutex.Unlock()
	if fake.EvaluateStub != nil {
		return fake.EvaluateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.evaluateReturns
	return fakeReturns.result1
}

func (fake *InconvertiblePolicy) EvaluateCallCount() int {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	return len(fake.evaluateArgsForCall)
}

func (fake *InconvertiblePolicy) EvaluateCalls(stub func([]*protoutil.SignedData) error) {
	fake.evaluateMutex.Lock()
	defer fake.evaluateMutex.Unlock()
	fake.EvaluateStub = stub
}

func (fake *InconvertiblePolicy) EvaluateArgsForCall(i int) []*protoutil.SignedData {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	argsForCall := fake.evaluateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *InconvertiblePolicy) EvaluateReturns(result1 error) {
	fake.evaluateMutex.Lock()
	defer fake.evaluateMutex.Unlock()
	fake.EvaluateStub = nil
	fake.evaluateReturns = struct {
		result1 error
	}{result1}
}

func (fake *InconvertiblePolicy) EvaluateReturnsOnCall(i int, result1 error) {
	fake.evaluateMutex.Lock()
	defer fake.evaluateMutex.Unlock()
	fake.EvaluateStub = nil
	if fake.evaluateReturnsOnCall == nil {
		fake.evaluateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.evaluateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *InconvertiblePolicy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *InconvertiblePolicy) recordInvocation(key string, args []interface{}) {
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
