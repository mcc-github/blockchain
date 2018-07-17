
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/discovery/support/acl"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type Verifier struct {
	VerifyByChannelStub        func(channel string, sd *cb.SignedData) error
	verifyByChannelMutex       sync.RWMutex
	verifyByChannelArgsForCall []struct {
		channel string
		sd      *cb.SignedData
	}
	verifyByChannelReturns struct {
		result1 error
	}
	verifyByChannelReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Verifier) VerifyByChannel(channel string, sd *cb.SignedData) error {
	fake.verifyByChannelMutex.Lock()
	ret, specificReturn := fake.verifyByChannelReturnsOnCall[len(fake.verifyByChannelArgsForCall)]
	fake.verifyByChannelArgsForCall = append(fake.verifyByChannelArgsForCall, struct {
		channel string
		sd      *cb.SignedData
	}{channel, sd})
	fake.recordInvocation("VerifyByChannel", []interface{}{channel, sd})
	fake.verifyByChannelMutex.Unlock()
	if fake.VerifyByChannelStub != nil {
		return fake.VerifyByChannelStub(channel, sd)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyByChannelReturns.result1
}

func (fake *Verifier) VerifyByChannelCallCount() int {
	fake.verifyByChannelMutex.RLock()
	defer fake.verifyByChannelMutex.RUnlock()
	return len(fake.verifyByChannelArgsForCall)
}

func (fake *Verifier) VerifyByChannelArgsForCall(i int) (string, *cb.SignedData) {
	fake.verifyByChannelMutex.RLock()
	defer fake.verifyByChannelMutex.RUnlock()
	return fake.verifyByChannelArgsForCall[i].channel, fake.verifyByChannelArgsForCall[i].sd
}

func (fake *Verifier) VerifyByChannelReturns(result1 error) {
	fake.VerifyByChannelStub = nil
	fake.verifyByChannelReturns = struct {
		result1 error
	}{result1}
}

func (fake *Verifier) VerifyByChannelReturnsOnCall(i int, result1 error) {
	fake.VerifyByChannelStub = nil
	if fake.verifyByChannelReturnsOnCall == nil {
		fake.verifyByChannelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyByChannelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Verifier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.verifyByChannelMutex.RLock()
	defer fake.verifyByChannelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Verifier) recordInvocation(key string, args []interface{}) {
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

var _ acl.Verifier = new(Verifier)
