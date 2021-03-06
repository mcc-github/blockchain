
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/token/server"
)

type SignedDataPolicyChecker struct {
	CheckPolicyBySignedDataStub        func(channelID, policyName string, sd []*common.SignedData) error
	checkPolicyBySignedDataMutex       sync.RWMutex
	checkPolicyBySignedDataArgsForCall []struct {
		channelID  string
		policyName string
		sd         []*common.SignedData
	}
	checkPolicyBySignedDataReturns struct {
		result1 error
	}
	checkPolicyBySignedDataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SignedDataPolicyChecker) CheckPolicyBySignedData(channelID string, policyName string, sd []*common.SignedData) error {
	var sdCopy []*common.SignedData
	if sd != nil {
		sdCopy = make([]*common.SignedData, len(sd))
		copy(sdCopy, sd)
	}
	fake.checkPolicyBySignedDataMutex.Lock()
	ret, specificReturn := fake.checkPolicyBySignedDataReturnsOnCall[len(fake.checkPolicyBySignedDataArgsForCall)]
	fake.checkPolicyBySignedDataArgsForCall = append(fake.checkPolicyBySignedDataArgsForCall, struct {
		channelID  string
		policyName string
		sd         []*common.SignedData
	}{channelID, policyName, sdCopy})
	fake.recordInvocation("CheckPolicyBySignedData", []interface{}{channelID, policyName, sdCopy})
	fake.checkPolicyBySignedDataMutex.Unlock()
	if fake.CheckPolicyBySignedDataStub != nil {
		return fake.CheckPolicyBySignedDataStub(channelID, policyName, sd)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.checkPolicyBySignedDataReturns.result1
}

func (fake *SignedDataPolicyChecker) CheckPolicyBySignedDataCallCount() int {
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	return len(fake.checkPolicyBySignedDataArgsForCall)
}

func (fake *SignedDataPolicyChecker) CheckPolicyBySignedDataArgsForCall(i int) (string, string, []*common.SignedData) {
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	return fake.checkPolicyBySignedDataArgsForCall[i].channelID, fake.checkPolicyBySignedDataArgsForCall[i].policyName, fake.checkPolicyBySignedDataArgsForCall[i].sd
}

func (fake *SignedDataPolicyChecker) CheckPolicyBySignedDataReturns(result1 error) {
	fake.CheckPolicyBySignedDataStub = nil
	fake.checkPolicyBySignedDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *SignedDataPolicyChecker) CheckPolicyBySignedDataReturnsOnCall(i int, result1 error) {
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

func (fake *SignedDataPolicyChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkPolicyBySignedDataMutex.RLock()
	defer fake.checkPolicyBySignedDataMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SignedDataPolicyChecker) recordInvocation(key string, args []interface{}) {
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

var _ server.SignedDataPolicyChecker = new(SignedDataPolicyChecker)
