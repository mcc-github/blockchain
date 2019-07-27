
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
)

type ConvertiblePolicy struct {
	ConvertStub        func() (*common.SignaturePolicyEnvelope, error)
	convertMutex       sync.RWMutex
	convertArgsForCall []struct {
	}
	convertReturns struct {
		result1 *common.SignaturePolicyEnvelope
		result2 error
	}
	convertReturnsOnCall map[int]struct {
		result1 *common.SignaturePolicyEnvelope
		result2 error
	}
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

func (fake *ConvertiblePolicy) Convert() (*common.SignaturePolicyEnvelope, error) {
	fake.convertMutex.Lock()
	ret, specificReturn := fake.convertReturnsOnCall[len(fake.convertArgsForCall)]
	fake.convertArgsForCall = append(fake.convertArgsForCall, struct {
	}{})
	fake.recordInvocation("Convert", []interface{}{})
	fake.convertMutex.Unlock()
	if fake.ConvertStub != nil {
		return fake.ConvertStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.convertReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ConvertiblePolicy) ConvertCallCount() int {
	fake.convertMutex.RLock()
	defer fake.convertMutex.RUnlock()
	return len(fake.convertArgsForCall)
}

func (fake *ConvertiblePolicy) ConvertCalls(stub func() (*common.SignaturePolicyEnvelope, error)) {
	fake.convertMutex.Lock()
	defer fake.convertMutex.Unlock()
	fake.ConvertStub = stub
}

func (fake *ConvertiblePolicy) ConvertReturns(result1 *common.SignaturePolicyEnvelope, result2 error) {
	fake.convertMutex.Lock()
	defer fake.convertMutex.Unlock()
	fake.ConvertStub = nil
	fake.convertReturns = struct {
		result1 *common.SignaturePolicyEnvelope
		result2 error
	}{result1, result2}
}

func (fake *ConvertiblePolicy) ConvertReturnsOnCall(i int, result1 *common.SignaturePolicyEnvelope, result2 error) {
	fake.convertMutex.Lock()
	defer fake.convertMutex.Unlock()
	fake.ConvertStub = nil
	if fake.convertReturnsOnCall == nil {
		fake.convertReturnsOnCall = make(map[int]struct {
			result1 *common.SignaturePolicyEnvelope
			result2 error
		})
	}
	fake.convertReturnsOnCall[i] = struct {
		result1 *common.SignaturePolicyEnvelope
		result2 error
	}{result1, result2}
}

func (fake *ConvertiblePolicy) Evaluate(arg1 []*protoutil.SignedData) error {
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

func (fake *ConvertiblePolicy) EvaluateCallCount() int {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	return len(fake.evaluateArgsForCall)
}

func (fake *ConvertiblePolicy) EvaluateCalls(stub func([]*protoutil.SignedData) error) {
	fake.evaluateMutex.Lock()
	defer fake.evaluateMutex.Unlock()
	fake.EvaluateStub = stub
}

func (fake *ConvertiblePolicy) EvaluateArgsForCall(i int) []*protoutil.SignedData {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	argsForCall := fake.evaluateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ConvertiblePolicy) EvaluateReturns(result1 error) {
	fake.evaluateMutex.Lock()
	defer fake.evaluateMutex.Unlock()
	fake.EvaluateStub = nil
	fake.evaluateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ConvertiblePolicy) EvaluateReturnsOnCall(i int, result1 error) {
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

func (fake *ConvertiblePolicy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.convertMutex.RLock()
	defer fake.convertMutex.RUnlock()
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ConvertiblePolicy) recordInvocation(key string, args []interface{}) {
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
