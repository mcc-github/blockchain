
package mock

import (
	"sync"

	container_test "github.com/mcc-github/blockchain/core/container"
)

type VMProvider struct {
	NewVMStub        func() container_test.VM
	newVMMutex       sync.RWMutex
	newVMArgsForCall []struct{}
	newVMReturns     struct {
		result1 container_test.VM
	}
	newVMReturnsOnCall map[int]struct {
		result1 container_test.VM
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *VMProvider) NewVM() container_test.VM {
	fake.newVMMutex.Lock()
	ret, specificReturn := fake.newVMReturnsOnCall[len(fake.newVMArgsForCall)]
	fake.newVMArgsForCall = append(fake.newVMArgsForCall, struct{}{})
	fake.recordInvocation("NewVM", []interface{}{})
	fake.newVMMutex.Unlock()
	if fake.NewVMStub != nil {
		return fake.NewVMStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.newVMReturns.result1
}

func (fake *VMProvider) NewVMCallCount() int {
	fake.newVMMutex.RLock()
	defer fake.newVMMutex.RUnlock()
	return len(fake.newVMArgsForCall)
}

func (fake *VMProvider) NewVMReturns(result1 container_test.VM) {
	fake.NewVMStub = nil
	fake.newVMReturns = struct {
		result1 container_test.VM
	}{result1}
}

func (fake *VMProvider) NewVMReturnsOnCall(i int, result1 container_test.VM) {
	fake.NewVMStub = nil
	if fake.newVMReturnsOnCall == nil {
		fake.newVMReturnsOnCall = make(map[int]struct {
			result1 container_test.VM
		})
	}
	fake.newVMReturnsOnCall[i] = struct {
		result1 container_test.VM
	}{result1}
}

func (fake *VMProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newVMMutex.RLock()
	defer fake.newVMMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *VMProvider) recordInvocation(key string, args []interface{}) {
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
