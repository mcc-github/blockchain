
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type InstantiationPolicyChecker struct {
	CheckInstantiationPolicyStub        func(string, string, *ccprovider.ChaincodeData) error
	checkInstantiationPolicyMutex       sync.RWMutex
	checkInstantiationPolicyArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 *ccprovider.ChaincodeData
	}
	checkInstantiationPolicyReturns struct {
		result1 error
	}
	checkInstantiationPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicy(arg1 string, arg2 string, arg3 *ccprovider.ChaincodeData) error {
	fake.checkInstantiationPolicyMutex.Lock()
	ret, specificReturn := fake.checkInstantiationPolicyReturnsOnCall[len(fake.checkInstantiationPolicyArgsForCall)]
	fake.checkInstantiationPolicyArgsForCall = append(fake.checkInstantiationPolicyArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 *ccprovider.ChaincodeData
	}{arg1, arg2, arg3})
	fake.recordInvocation("CheckInstantiationPolicy", []interface{}{arg1, arg2, arg3})
	fake.checkInstantiationPolicyMutex.Unlock()
	if fake.CheckInstantiationPolicyStub != nil {
		return fake.CheckInstantiationPolicyStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkInstantiationPolicyReturns
	return fakeReturns.result1
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicyCallCount() int {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	return len(fake.checkInstantiationPolicyArgsForCall)
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicyCalls(stub func(string, string, *ccprovider.ChaincodeData) error) {
	fake.checkInstantiationPolicyMutex.Lock()
	defer fake.checkInstantiationPolicyMutex.Unlock()
	fake.CheckInstantiationPolicyStub = stub
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicyArgsForCall(i int) (string, string, *ccprovider.ChaincodeData) {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	argsForCall := fake.checkInstantiationPolicyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicyReturns(result1 error) {
	fake.checkInstantiationPolicyMutex.Lock()
	defer fake.checkInstantiationPolicyMutex.Unlock()
	fake.CheckInstantiationPolicyStub = nil
	fake.checkInstantiationPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *InstantiationPolicyChecker) CheckInstantiationPolicyReturnsOnCall(i int, result1 error) {
	fake.checkInstantiationPolicyMutex.Lock()
	defer fake.checkInstantiationPolicyMutex.Unlock()
	fake.CheckInstantiationPolicyStub = nil
	if fake.checkInstantiationPolicyReturnsOnCall == nil {
		fake.checkInstantiationPolicyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkInstantiationPolicyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *InstantiationPolicyChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *InstantiationPolicyChecker) recordInvocation(key string, args []interface{}) {
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
