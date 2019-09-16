
package fake

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/internal/pkg/peer/blocksprovider"
)

type GossipServiceAdapter struct {
	AddPayloadStub        func(string, *gossip.Payload) error
	addPayloadMutex       sync.RWMutex
	addPayloadArgsForCall []struct {
		arg1 string
		arg2 *gossip.Payload
	}
	addPayloadReturns struct {
		result1 error
	}
	addPayloadReturnsOnCall map[int]struct {
		result1 error
	}
	GossipStub        func(*gossip.GossipMessage)
	gossipMutex       sync.RWMutex
	gossipArgsForCall []struct {
		arg1 *gossip.GossipMessage
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *GossipServiceAdapter) AddPayload(arg1 string, arg2 *gossip.Payload) error {
	fake.addPayloadMutex.Lock()
	ret, specificReturn := fake.addPayloadReturnsOnCall[len(fake.addPayloadArgsForCall)]
	fake.addPayloadArgsForCall = append(fake.addPayloadArgsForCall, struct {
		arg1 string
		arg2 *gossip.Payload
	}{arg1, arg2})
	fake.recordInvocation("AddPayload", []interface{}{arg1, arg2})
	fake.addPayloadMutex.Unlock()
	if fake.AddPayloadStub != nil {
		return fake.AddPayloadStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addPayloadReturns
	return fakeReturns.result1
}

func (fake *GossipServiceAdapter) AddPayloadCallCount() int {
	fake.addPayloadMutex.RLock()
	defer fake.addPayloadMutex.RUnlock()
	return len(fake.addPayloadArgsForCall)
}

func (fake *GossipServiceAdapter) AddPayloadCalls(stub func(string, *gossip.Payload) error) {
	fake.addPayloadMutex.Lock()
	defer fake.addPayloadMutex.Unlock()
	fake.AddPayloadStub = stub
}

func (fake *GossipServiceAdapter) AddPayloadArgsForCall(i int) (string, *gossip.Payload) {
	fake.addPayloadMutex.RLock()
	defer fake.addPayloadMutex.RUnlock()
	argsForCall := fake.addPayloadArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *GossipServiceAdapter) AddPayloadReturns(result1 error) {
	fake.addPayloadMutex.Lock()
	defer fake.addPayloadMutex.Unlock()
	fake.AddPayloadStub = nil
	fake.addPayloadReturns = struct {
		result1 error
	}{result1}
}

func (fake *GossipServiceAdapter) AddPayloadReturnsOnCall(i int, result1 error) {
	fake.addPayloadMutex.Lock()
	defer fake.addPayloadMutex.Unlock()
	fake.AddPayloadStub = nil
	if fake.addPayloadReturnsOnCall == nil {
		fake.addPayloadReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addPayloadReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *GossipServiceAdapter) Gossip(arg1 *gossip.GossipMessage) {
	fake.gossipMutex.Lock()
	fake.gossipArgsForCall = append(fake.gossipArgsForCall, struct {
		arg1 *gossip.GossipMessage
	}{arg1})
	fake.recordInvocation("Gossip", []interface{}{arg1})
	fake.gossipMutex.Unlock()
	if fake.GossipStub != nil {
		fake.GossipStub(arg1)
	}
}

func (fake *GossipServiceAdapter) GossipCallCount() int {
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	return len(fake.gossipArgsForCall)
}

func (fake *GossipServiceAdapter) GossipCalls(stub func(*gossip.GossipMessage)) {
	fake.gossipMutex.Lock()
	defer fake.gossipMutex.Unlock()
	fake.GossipStub = stub
}

func (fake *GossipServiceAdapter) GossipArgsForCall(i int) *gossip.GossipMessage {
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	argsForCall := fake.gossipArgsForCall[i]
	return argsForCall.arg1
}

func (fake *GossipServiceAdapter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addPayloadMutex.RLock()
	defer fake.addPayloadMutex.RUnlock()
	fake.gossipMutex.RLock()
	defer fake.gossipMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *GossipServiceAdapter) recordInvocation(key string, args []interface{}) {
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

var _ blocksprovider.GossipServiceAdapter = new(GossipServiceAdapter)
