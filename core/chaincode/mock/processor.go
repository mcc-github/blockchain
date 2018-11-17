
package mock

import (
	sync "sync"

	container "github.com/mcc-github/blockchain/core/container"
)

type Processor struct {
	ProcessStub        func(string, container.VMCReq) error
	processMutex       sync.RWMutex
	processArgsForCall []struct {
		arg1 string
		arg2 container.VMCReq
	}
	processReturns struct {
		result1 error
	}
	processReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Processor) Process(arg1 string, arg2 container.VMCReq) error {
	fake.processMutex.Lock()
	ret, specificReturn := fake.processReturnsOnCall[len(fake.processArgsForCall)]
	fake.processArgsForCall = append(fake.processArgsForCall, struct {
		arg1 string
		arg2 container.VMCReq
	}{arg1, arg2})
	fake.recordInvocation("Process", []interface{}{arg1, arg2})
	fake.processMutex.Unlock()
	if fake.ProcessStub != nil {
		return fake.ProcessStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.processReturns
	return fakeReturns.result1
}

func (fake *Processor) ProcessCallCount() int {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return len(fake.processArgsForCall)
}

func (fake *Processor) ProcessCalls(stub func(string, container.VMCReq) error) {
	fake.processMutex.Lock()
	defer fake.processMutex.Unlock()
	fake.ProcessStub = stub
}

func (fake *Processor) ProcessArgsForCall(i int) (string, container.VMCReq) {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	argsForCall := fake.processArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Processor) ProcessReturns(result1 error) {
	fake.processMutex.Lock()
	defer fake.processMutex.Unlock()
	fake.ProcessStub = nil
	fake.processReturns = struct {
		result1 error
	}{result1}
}

func (fake *Processor) ProcessReturnsOnCall(i int, result1 error) {
	fake.processMutex.Lock()
	defer fake.processMutex.Unlock()
	fake.ProcessStub = nil
	if fake.processReturnsOnCall == nil {
		fake.processReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.processReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Processor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
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
