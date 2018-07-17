
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/deliver"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type PolicyChecker struct {
	CheckPolicyStub        func(envelope *cb.Envelope, channelID string) error
	checkPolicyMutex       sync.RWMutex
	checkPolicyArgsForCall []struct {
		envelope  *cb.Envelope
		channelID string
	}
	checkPolicyReturns struct {
		result1 error
	}
	checkPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PolicyChecker) CheckPolicy(envelope *cb.Envelope, channelID string) error {
	fake.checkPolicyMutex.Lock()
	ret, specificReturn := fake.checkPolicyReturnsOnCall[len(fake.checkPolicyArgsForCall)]
	fake.checkPolicyArgsForCall = append(fake.checkPolicyArgsForCall, struct {
		envelope  *cb.Envelope
		channelID string
	}{envelope, channelID})
	fake.recordInvocation("CheckPolicy", []interface{}{envelope, channelID})
	fake.checkPolicyMutex.Unlock()
	if fake.CheckPolicyStub != nil {
		return fake.CheckPolicyStub(envelope, channelID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.checkPolicyReturns.result1
}

func (fake *PolicyChecker) CheckPolicyCallCount() int {
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
	return len(fake.checkPolicyArgsForCall)
}

func (fake *PolicyChecker) CheckPolicyArgsForCall(i int) (*cb.Envelope, string) {
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
	return fake.checkPolicyArgsForCall[i].envelope, fake.checkPolicyArgsForCall[i].channelID
}

func (fake *PolicyChecker) CheckPolicyReturns(result1 error) {
	fake.CheckPolicyStub = nil
	fake.checkPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *PolicyChecker) CheckPolicyReturnsOnCall(i int, result1 error) {
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

func (fake *PolicyChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkPolicyMutex.RLock()
	defer fake.checkPolicyMutex.RUnlock()
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

var _ deliver.PolicyChecker = new(PolicyChecker)
