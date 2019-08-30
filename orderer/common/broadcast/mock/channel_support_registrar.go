
package mock

import (
	sync "sync"

	common "github.com/mcc-github/blockchain-protos-go/common"
	broadcast "github.com/mcc-github/blockchain/orderer/common/broadcast"
)

type ChannelSupportRegistrar struct {
	BroadcastChannelSupportStub        func(*common.Envelope) (*common.ChannelHeader, bool, broadcast.ChannelSupport, error)
	broadcastChannelSupportMutex       sync.RWMutex
	broadcastChannelSupportArgsForCall []struct {
		arg1 *common.Envelope
	}
	broadcastChannelSupportReturns struct {
		result1 *common.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}
	broadcastChannelSupportReturnsOnCall map[int]struct {
		result1 *common.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupport(arg1 *common.Envelope) (*common.ChannelHeader, bool, broadcast.ChannelSupport, error) {
	fake.broadcastChannelSupportMutex.Lock()
	ret, specificReturn := fake.broadcastChannelSupportReturnsOnCall[len(fake.broadcastChannelSupportArgsForCall)]
	fake.broadcastChannelSupportArgsForCall = append(fake.broadcastChannelSupportArgsForCall, struct {
		arg1 *common.Envelope
	}{arg1})
	fake.recordInvocation("BroadcastChannelSupport", []interface{}{arg1})
	fake.broadcastChannelSupportMutex.Unlock()
	if fake.BroadcastChannelSupportStub != nil {
		return fake.BroadcastChannelSupportStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	fakeReturns := fake.broadcastChannelSupportReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportCallCount() int {
	fake.broadcastChannelSupportMutex.RLock()
	defer fake.broadcastChannelSupportMutex.RUnlock()
	return len(fake.broadcastChannelSupportArgsForCall)
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportCalls(stub func(*common.Envelope) (*common.ChannelHeader, bool, broadcast.ChannelSupport, error)) {
	fake.broadcastChannelSupportMutex.Lock()
	defer fake.broadcastChannelSupportMutex.Unlock()
	fake.BroadcastChannelSupportStub = stub
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportArgsForCall(i int) *common.Envelope {
	fake.broadcastChannelSupportMutex.RLock()
	defer fake.broadcastChannelSupportMutex.RUnlock()
	argsForCall := fake.broadcastChannelSupportArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportReturns(result1 *common.ChannelHeader, result2 bool, result3 broadcast.ChannelSupport, result4 error) {
	fake.broadcastChannelSupportMutex.Lock()
	defer fake.broadcastChannelSupportMutex.Unlock()
	fake.BroadcastChannelSupportStub = nil
	fake.broadcastChannelSupportReturns = struct {
		result1 *common.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportReturnsOnCall(i int, result1 *common.ChannelHeader, result2 bool, result3 broadcast.ChannelSupport, result4 error) {
	fake.broadcastChannelSupportMutex.Lock()
	defer fake.broadcastChannelSupportMutex.Unlock()
	fake.BroadcastChannelSupportStub = nil
	if fake.broadcastChannelSupportReturnsOnCall == nil {
		fake.broadcastChannelSupportReturnsOnCall = make(map[int]struct {
			result1 *common.ChannelHeader
			result2 bool
			result3 broadcast.ChannelSupport
			result4 error
		})
	}
	fake.broadcastChannelSupportReturnsOnCall[i] = struct {
		result1 *common.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *ChannelSupportRegistrar) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.broadcastChannelSupportMutex.RLock()
	defer fake.broadcastChannelSupportMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChannelSupportRegistrar) recordInvocation(key string, args []interface{}) {
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

var _ broadcast.ChannelSupportRegistrar = new(ChannelSupportRegistrar)
