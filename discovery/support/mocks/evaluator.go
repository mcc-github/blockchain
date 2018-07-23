/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/discovery/support/acl"
	common2 "github.com/mcc-github/blockchain/protos/common"
)

type Evaluator struct {
	EvaluateStub        func(signatureSet []*common2.SignedData) error
	evaluateMutex       sync.RWMutex
	evaluateArgsForCall []struct {
		signatureSet []*common2.SignedData
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

func (fake *Evaluator) Evaluate(signatureSet []*common2.SignedData) error {
	var signatureSetCopy []*common2.SignedData
	if signatureSet != nil {
		signatureSetCopy = make([]*common2.SignedData, len(signatureSet))
		copy(signatureSetCopy, signatureSet)
	}
	fake.evaluateMutex.Lock()
	ret, specificReturn := fake.evaluateReturnsOnCall[len(fake.evaluateArgsForCall)]
	fake.evaluateArgsForCall = append(fake.evaluateArgsForCall, struct {
		signatureSet []*common2.SignedData
	}{signatureSetCopy})
	fake.recordInvocation("Evaluate", []interface{}{signatureSetCopy})
	fake.evaluateMutex.Unlock()
	if fake.EvaluateStub != nil {
		return fake.EvaluateStub(signatureSet)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.evaluateReturns.result1
}

func (fake *Evaluator) EvaluateCallCount() int {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	return len(fake.evaluateArgsForCall)
}

func (fake *Evaluator) EvaluateArgsForCall(i int) []*common2.SignedData {
	fake.evaluateMutex.RLock()
	defer fake.evaluateMutex.RUnlock()
	return fake.evaluateArgsForCall[i].signatureSet
}

func (fake *Evaluator) EvaluateReturns(result1 error) {
	fake.EvaluateStub = nil
	fake.evaluateReturns = struct {
		result1 error
	}{result1}
}

func (fake *Evaluator) EvaluateReturnsOnCall(i int, result1 error) {
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

func (fake *Evaluator) Invocations() map[string][][]interface{} {
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

func (fake *Evaluator) recordInvocation(key string, args []interface{}) {
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

var _ acl.Evaluator = new(Evaluator)