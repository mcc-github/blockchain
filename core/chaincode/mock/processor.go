
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/container"
)

type Processor struct {
	ProcessStub        func(vmtype string, req container.VMCReq) error
	processMutex       sync.RWMutex
	processArgsForCall []struct {
		vmtype string
		req    container.VMCReq
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

func (fake *Processor) Process(vmtype string, req container.VMCReq) error {
	fake.processMutex.Lock()
	ret, specificReturn := fake.processReturnsOnCall[len(fake.processArgsForCall)]
	fake.processArgsForCall = append(fake.processArgsForCall, struct {
		vmtype string
		req    container.VMCReq
	}{vmtype, req})
	fake.recordInvocation("Process", []interface{}{vmtype, req})
	fake.processMutex.Unlock()
	if fake.ProcessStub != nil {
		return fake.ProcessStub(vmtype, req)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.processReturns.result1
}

func (fake *Processor) ProcessCallCount() int {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return len(fake.processArgsForCall)
}

func (fake *Processor) ProcessArgsForCall(i int) (string, container.VMCReq) {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return fake.processArgsForCall[i].vmtype, fake.processArgsForCall[i].req
}

func (fake *Processor) ProcessReturns(result1 error) {
	fake.ProcessStub = nil
	fake.processReturns = struct {
		result1 error
	}{result1}
}

func (fake *Processor) ProcessReturnsOnCall(i int, result1 error) {
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
