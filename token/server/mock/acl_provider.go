
package mock

import (
	sync "sync"

	server "github.com/mcc-github/blockchain/token/server"
)

type ACLProvider struct {
	CheckACLStub        func(string, string, interface{}) error
	checkACLMutex       sync.RWMutex
	checkACLArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 interface{}
	}
	checkACLReturns struct {
		result1 error
	}
	checkACLReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ACLProvider) CheckACL(arg1 string, arg2 string, arg3 interface{}) error {
	fake.checkACLMutex.Lock()
	ret, specificReturn := fake.checkACLReturnsOnCall[len(fake.checkACLArgsForCall)]
	fake.checkACLArgsForCall = append(fake.checkACLArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 interface{}
	}{arg1, arg2, arg3})
	fake.recordInvocation("CheckACL", []interface{}{arg1, arg2, arg3})
	fake.checkACLMutex.Unlock()
	if fake.CheckACLStub != nil {
		return fake.CheckACLStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkACLReturns
	return fakeReturns.result1
}

func (fake *ACLProvider) CheckACLCallCount() int {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	return len(fake.checkACLArgsForCall)
}

func (fake *ACLProvider) CheckACLCalls(stub func(string, string, interface{}) error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = stub
}

func (fake *ACLProvider) CheckACLArgsForCall(i int) (string, string, interface{}) {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	argsForCall := fake.checkACLArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ACLProvider) CheckACLReturns(result1 error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = nil
	fake.checkACLReturns = struct {
		result1 error
	}{result1}
}

func (fake *ACLProvider) CheckACLReturnsOnCall(i int, result1 error) {
	fake.checkACLMutex.Lock()
	defer fake.checkACLMutex.Unlock()
	fake.CheckACLStub = nil
	if fake.checkACLReturnsOnCall == nil {
		fake.checkACLReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkACLReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ACLProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ACLProvider) recordInvocation(key string, args []interface{}) {
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

var _ server.ACLProvider = new(ACLProvider)
