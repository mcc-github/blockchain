
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/common/policies"
)

type ChannelPolicyManagerGetter struct {
	ManagerStub        func(string) policies.Manager
	managerMutex       sync.RWMutex
	managerArgsForCall []struct {
		arg1 string
	}
	managerReturns struct {
		result1 policies.Manager
	}
	managerReturnsOnCall map[int]struct {
		result1 policies.Manager
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelPolicyManagerGetter) Manager(arg1 string) policies.Manager {
	fake.managerMutex.Lock()
	ret, specificReturn := fake.managerReturnsOnCall[len(fake.managerArgsForCall)]
	fake.managerArgsForCall = append(fake.managerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Manager", []interface{}{arg1})
	fake.managerMutex.Unlock()
	if fake.ManagerStub != nil {
		return fake.ManagerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.managerReturns
	return fakeReturns.result1
}

func (fake *ChannelPolicyManagerGetter) ManagerCallCount() int {
	fake.managerMutex.RLock()
	defer fake.managerMutex.RUnlock()
	return len(fake.managerArgsForCall)
}

func (fake *ChannelPolicyManagerGetter) ManagerCalls(stub func(string) policies.Manager) {
	fake.managerMutex.Lock()
	defer fake.managerMutex.Unlock()
	fake.ManagerStub = stub
}

func (fake *ChannelPolicyManagerGetter) ManagerArgsForCall(i int) string {
	fake.managerMutex.RLock()
	defer fake.managerMutex.RUnlock()
	argsForCall := fake.managerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChannelPolicyManagerGetter) ManagerReturns(result1 policies.Manager) {
	fake.managerMutex.Lock()
	defer fake.managerMutex.Unlock()
	fake.ManagerStub = nil
	fake.managerReturns = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *ChannelPolicyManagerGetter) ManagerReturnsOnCall(i int, result1 policies.Manager) {
	fake.managerMutex.Lock()
	defer fake.managerMutex.Unlock()
	fake.ManagerStub = nil
	if fake.managerReturnsOnCall == nil {
		fake.managerReturnsOnCall = make(map[int]struct {
			result1 policies.Manager
		})
	}
	fake.managerReturnsOnCall[i] = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *ChannelPolicyManagerGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.managerMutex.RLock()
	defer fake.managerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChannelPolicyManagerGetter) recordInvocation(key string, args []interface{}) {
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
