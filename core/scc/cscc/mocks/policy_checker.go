
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
)

type PolicyChecker struct {
	CheckPolicyStub        func(string, string, *peer.SignedProposal) error
	checkPolicyMutex       sync.RWMutex
	checkPolicyArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 *peer.SignedProposal
	}
	checkPolicyReturns struct {
		result1 error
	}
	checkPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	CheckPolicyBySignedDataStub        func(string, string, []*protoutil.SignedData) error
	checkPolicyBySignedDataMutex       sync.RWMutex
	checkPolicyBySignedDataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []*protoutil.SignedData
	}
	checkPolicyBySignedDataReturns struct {
		result1 error
	}
	checkPolicyBySignedDataReturnsOnCall map[int]struct {
		result1 error
	}
	CheckPolicyNoChannelStub        func(string, *peer.SignedProposal) error
	checkPolicyNoChannelMutex       sync.RWMutex
	checkPolicyNoChannelArgsForCall []struct {
		arg1 string
		arg2 *peer.SignedProposal
	}
	checkPolicyNoChannelReturns struct {
		result1 error
	}
	checkPolicyNoChannelReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PolicyChecker) CheckPolicy(arg1 string, arg2 string, arg3 *peer.SignedProposal) error {
	fake.checkPolicyMutex.Lock()
	ret, specificReturn := fake.checkPolicyReturnsOnCall[len(fake.checkPolicyArgsForCall)]
	fake.checkPolicyArgsForCall = append(fake.checkPolicyArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 *peer.SignedProposal
	}{arg1, arg2, arg3})
	fake.recordInvocation("CheckPolicy", []interface{}{arg1, arg2, arg3})
	fake.checkPolicyMutex.Unlock()
	if fake.CheckPolicyStub != nil {
		return fake.CheckPolicyStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkPolicyReturns
	return fakeReturns.result1
}

func (fake *PolicyChecker) CheckPolicyCallCount() int {
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
	return len(fake.checkPolicyArgsForCall)
}

func (fake *PolicyChecker) CheckPolicyCalls(stub func(string, string, *peer.SignedProposal) error) {
	fake.checkPolicyMutex.Lock()
	defer fake.checkPolicyMutex.Unlock()
	fake.CheckPolicyStub = stub
}

func (fake *PolicyChecker) CheckPolicyArgsForCall(i int) (string, string, *peer.SignedProposal) {
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
	argsForCall := fake.checkPolicyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *PolicyChecker) CheckPolicyReturns(result1 error) {
	fake.checkPolicyMutex.Lock()
	defer fake.checkPolicyMutex.Unlock()
	fake.CheckPolicyStub = nil
	fake.checkPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyReturnsOnCall(i int, result1 error) {
	fake.checkPolicyMutex.Lock()
	defer fake.checkPolicyMutex.Unlock()
	fake.CheckPolicyStub = nil
	if fake.checkPolicyReturnsOnCall == nil {
		fake.checkPolicyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkPolicyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyBySignedData(arg1 string, arg2 string, arg3 []*protoutil.SignedData) error {
	var arg3Copy []*protoutil.SignedData
	if arg3 != nil {
		arg3Copy = make([]*protoutil.SignedData, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.checkPolicyBySignedDataMutex.Lock()
	ret, specificReturn := fake.checkPolicyBySignedDataReturnsOnCall[len(fake.checkPolicyBySignedDataArgsForCall)]
	fake.checkPolicyBySignedDataArgsForCall = append(fake.checkPolicyBySignedDataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []*protoutil.SignedData
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("CheckPolicyBySignedData", []interface{}{arg1, arg2, arg3Copy})
	fake.checkPolicyBySignedDataMutex.Unlock()
	if fake.CheckPolicyBySignedDataStub != nil {
		return fake.CheckPolicyBySignedDataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkPolicyBySignedDataReturns
	return fakeReturns.result1
}

func (fake *PolicyChecker) CheckPolicyBySignedDataCallCount() int {
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	return len(fake.checkPolicyBySignedDataArgsForCall)
}

func (fake *PolicyChecker) CheckPolicyBySignedDataCalls(stub func(string, string, []*protoutil.SignedData) error) {
	fake.checkPolicyBySignedDataMutex.Lock()
	defer fake.checkPolicyBySignedDataMutex.Unlock()
	fake.CheckPolicyBySignedDataStub = stub
}

func (fake *PolicyChecker) CheckPolicyBySignedDataArgsForCall(i int) (string, string, []*protoutil.SignedData) {
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	argsForCall := fake.checkPolicyBySignedDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *PolicyChecker) CheckPolicyBySignedDataReturns(result1 error) {
	fake.checkPolicyBySignedDataMutex.Lock()
	defer fake.checkPolicyBySignedDataMutex.Unlock()
	fake.CheckPolicyBySignedDataStub = nil
	fake.checkPolicyBySignedDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyBySignedDataReturnsOnCall(i int, result1 error) {
	fake.checkPolicyBySignedDataMutex.Lock()
	defer fake.checkPolicyBySignedDataMutex.Unlock()
	fake.CheckPolicyBySignedDataStub = nil
	if fake.checkPolicyBySignedDataReturnsOnCall == nil {
		fake.checkPolicyBySignedDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkPolicyBySignedDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyNoChannel(arg1 string, arg2 *peer.SignedProposal) error {
	fake.checkPolicyNoChannelMutex.Lock()
	ret, specificReturn := fake.checkPolicyNoChannelReturnsOnCall[len(fake.checkPolicyNoChannelArgsForCall)]
	fake.checkPolicyNoChannelArgsForCall = append(fake.checkPolicyNoChannelArgsForCall, struct {
		arg1 string
		arg2 *peer.SignedProposal
	}{arg1, arg2})
	fake.recordInvocation("CheckPolicyNoChannel", []interface{}{arg1, arg2})
	fake.checkPolicyNoChannelMutex.Unlock()
	if fake.CheckPolicyNoChannelStub != nil {
		return fake.CheckPolicyNoChannelStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkPolicyNoChannelReturns
	return fakeReturns.result1
}

func (fake *PolicyChecker) CheckPolicyNoChannelCallCount() int {
	fake.checkPolicyNoChannelMutex.RLock()
	defer fake.checkPolicyNoChannelMutex.RUnlock()
	return len(fake.checkPolicyNoChannelArgsForCall)
}

func (fake *PolicyChecker) CheckPolicyNoChannelCalls(stub func(string, *peer.SignedProposal) error) {
	fake.checkPolicyNoChannelMutex.Lock()
	defer fake.checkPolicyNoChannelMutex.Unlock()
	fake.CheckPolicyNoChannelStub = stub
}

func (fake *PolicyChecker) CheckPolicyNoChannelArgsForCall(i int) (string, *peer.SignedProposal) {
	fake.checkPolicyNoChannelMutex.RLock()
	defer fake.checkPolicyNoChannelMutex.RUnlock()
	argsForCall := fake.checkPolicyNoChannelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PolicyChecker) CheckPolicyNoChannelReturns(result1 error) {
	fake.checkPolicyNoChannelMutex.Lock()
	defer fake.checkPolicyNoChannelMutex.Unlock()
	fake.CheckPolicyNoChannelStub = nil
	fake.checkPolicyNoChannelReturns = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyNoChannelReturnsOnCall(i int, result1 error) {
	fake.checkPolicyNoChannelMutex.Lock()
	defer fake.checkPolicyNoChannelMutex.Unlock()
	fake.CheckPolicyNoChannelStub = nil
	if fake.checkPolicyNoChannelReturnsOnCall == nil {
		fake.checkPolicyNoChannelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkPolicyNoChannelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	fake.checkPolicyNoChannelMutex.RLock()
	defer fake.checkPolicyNoChannelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PolicyChecker) recordInvocation(key string, args []interface{}) {
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
