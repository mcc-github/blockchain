
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/orderer/common/broadcast"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type ChannelSupportRegistrar struct {
	BroadcastChannelSupportStub        func(msg *cb.Envelope) (*cb.ChannelHeader, bool, broadcast.ChannelSupport, error)
	broadcastChannelSupportMutex       sync.RWMutex
	broadcastChannelSupportArgsForCall []struct {
		msg *cb.Envelope
	}
	broadcastChannelSupportReturns struct {
		result1 *cb.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}
	broadcastChannelSupportReturnsOnCall map[int]struct {
		result1 *cb.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, broadcast.ChannelSupport, error) {
	fake.broadcastChannelSupportMutex.Lock()
	ret, specificReturn := fake.broadcastChannelSupportReturnsOnCall[len(fake.broadcastChannelSupportArgsForCall)]
	fake.broadcastChannelSupportArgsForCall = append(fake.broadcastChannelSupportArgsForCall, struct {
		msg *cb.Envelope
	}{msg})
	fake.recordInvocation("BroadcastChannelSupport", []interface{}{msg})
	fake.broadcastChannelSupportMutex.Unlock()
	if fake.BroadcastChannelSupportStub != nil {
		return fake.BroadcastChannelSupportStub(msg)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fake.broadcastChannelSupportReturns.result1, fake.broadcastChannelSupportReturns.result2, fake.broadcastChannelSupportReturns.result3, fake.broadcastChannelSupportReturns.result4
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportCallCount() int {
	fake.broadcastChannelSupportMutex.RLock()
	defer fake.broadcastChannelSupportMutex.RUnlock()
	return len(fake.broadcastChannelSupportArgsForCall)
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportArgsForCall(i int) *cb.Envelope {
	fake.broadcastChannelSupportMutex.RLock()
	defer fake.broadcastChannelSupportMutex.RUnlock()
	return fake.broadcastChannelSupportArgsForCall[i].msg
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportReturns(result1 *cb.ChannelHeader, result2 bool, result3 broadcast.ChannelSupport, result4 error) {
	fake.BroadcastChannelSupportStub = nil
	fake.broadcastChannelSupportReturns = struct {
		result1 *cb.ChannelHeader
		result2 bool
		result3 broadcast.ChannelSupport
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *ChannelSupportRegistrar) BroadcastChannelSupportReturnsOnCall(i int, result1 *cb.ChannelHeader, result2 bool, result3 broadcast.ChannelSupport, result4 error) {
	fake.BroadcastChannelSupportStub = nil
	if fake.broadcastChannelSupportReturnsOnCall == nil {
		fake.broadcastChannelSupportReturnsOnCall = make(map[int]struct {
			result1 *cb.ChannelHeader
			result2 bool
			result3 broadcast.ChannelSupport
			result4 error
		})
	}
	fake.broadcastChannelSupportReturnsOnCall[i] = struct {
		result1 *cb.ChannelHeader
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