
package fakes

import (
	sync "sync"

	operations "github.com/mcc-github/blockchain/core/operations"
)

type Logger struct {
	WarnStub        func(...interface{})
	warnMutex       sync.RWMutex
	warnArgsForCall []struct {
		arg1 []interface{}
	}
	WarnfStub        func(string, ...interface{})
	warnfMutex       sync.RWMutex
	warnfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Logger) Warn(arg1 ...interface{}) {
	fake.warnMutex.Lock()
	fake.warnArgsForCall = append(fake.warnArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Warn", []interface{}{arg1})
	fake.warnMutex.Unlock()
	if fake.WarnStub != nil {
		fake.WarnStub(arg1...)
	}
}

func (fake *Logger) WarnCallCount() int {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	return len(fake.warnArgsForCall)
}

func (fake *Logger) WarnCalls(stub func(...interface{})) {
	fake.warnMutex.Lock()
	defer fake.warnMutex.Unlock()
	fake.WarnStub = stub
}

func (fake *Logger) WarnArgsForCall(i int) []interface{} {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	argsForCall := fake.warnArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Logger) Warnf(arg1 string, arg2 ...interface{}) {
	fake.warnfMutex.Lock()
	fake.warnfArgsForCall = append(fake.warnfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Warnf", []interface{}{arg1, arg2})
	fake.warnfMutex.Unlock()
	if fake.WarnfStub != nil {
		fake.WarnfStub(arg1, arg2...)
	}
}

func (fake *Logger) WarnfCallCount() int {
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	return len(fake.warnfArgsForCall)
}

func (fake *Logger) WarnfCalls(stub func(string, ...interface{})) {
	fake.warnfMutex.Lock()
	defer fake.warnfMutex.Unlock()
	fake.WarnfStub = stub
}

func (fake *Logger) WarnfArgsForCall(i int) (string, []interface{}) {
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	argsForCall := fake.warnfArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Logger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Logger) recordInvocation(key string, args []interface{}) {
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

var _ operations.Logger = new(Logger)
