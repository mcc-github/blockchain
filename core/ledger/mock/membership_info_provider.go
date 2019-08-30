
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
)

type MembershipInfoProvider struct {
	AmMemberOfStub        func(string, *common.CollectionPolicyConfig) (bool, error)
	amMemberOfMutex       sync.RWMutex
	amMemberOfArgsForCall []struct {
		arg1 string
		arg2 *common.CollectionPolicyConfig
	}
	amMemberOfReturns struct {
		result1 bool
		result2 error
	}
	amMemberOfReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MembershipInfoProvider) AmMemberOf(arg1 string, arg2 *common.CollectionPolicyConfig) (bool, error) {
	fake.amMemberOfMutex.Lock()
	ret, specificReturn := fake.amMemberOfReturnsOnCall[len(fake.amMemberOfArgsForCall)]
	fake.amMemberOfArgsForCall = append(fake.amMemberOfArgsForCall, struct {
		arg1 string
		arg2 *common.CollectionPolicyConfig
	}{arg1, arg2})
	fake.recordInvocation("AmMemberOf", []interface{}{arg1, arg2})
	fake.amMemberOfMutex.Unlock()
	if fake.AmMemberOfStub != nil {
		return fake.AmMemberOfStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.amMemberOfReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MembershipInfoProvider) AmMemberOfCallCount() int {
	fake.amMemberOfMutex.RLock()
	defer fake.amMemberOfMutex.RUnlock()
	return len(fake.amMemberOfArgsForCall)
}

func (fake *MembershipInfoProvider) AmMemberOfCalls(stub func(string, *common.CollectionPolicyConfig) (bool, error)) {
	fake.amMemberOfMutex.Lock()
	defer fake.amMemberOfMutex.Unlock()
	fake.AmMemberOfStub = stub
}

func (fake *MembershipInfoProvider) AmMemberOfArgsForCall(i int) (string, *common.CollectionPolicyConfig) {
	fake.amMemberOfMutex.RLock()
	defer fake.amMemberOfMutex.RUnlock()
	argsForCall := fake.amMemberOfArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MembershipInfoProvider) AmMemberOfReturns(result1 bool, result2 error) {
	fake.amMemberOfMutex.Lock()
	defer fake.amMemberOfMutex.Unlock()
	fake.AmMemberOfStub = nil
	fake.amMemberOfReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *MembershipInfoProvider) AmMemberOfReturnsOnCall(i int, result1 bool, result2 error) {
	fake.amMemberOfMutex.Lock()
	defer fake.amMemberOfMutex.Unlock()
	fake.AmMemberOfStub = nil
	if fake.amMemberOfReturnsOnCall == nil {
		fake.amMemberOfReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.amMemberOfReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *MembershipInfoProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.amMemberOfMutex.RLock()
	defer fake.amMemberOfMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MembershipInfoProvider) recordInvocation(key string, args []interface{}) {
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

var _ ledger.MembershipInfoProvider = new(MembershipInfoProvider)
