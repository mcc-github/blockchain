
package fake

import (
	"sync"

	"github.com/mcc-github/blockchain/internal/pkg/peer/blocksprovider"
	"github.com/mcc-github/blockchain/internal/pkg/peer/orderers"
)

type OrdererConnectionSource struct {
	RandomEndpointStub        func() (*orderers.Endpoint, error)
	randomEndpointMutex       sync.RWMutex
	randomEndpointArgsForCall []struct {
	}
	randomEndpointReturns struct {
		result1 *orderers.Endpoint
		result2 error
	}
	randomEndpointReturnsOnCall map[int]struct {
		result1 *orderers.Endpoint
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OrdererConnectionSource) RandomEndpoint() (*orderers.Endpoint, error) {
	fake.randomEndpointMutex.Lock()
	ret, specificReturn := fake.randomEndpointReturnsOnCall[len(fake.randomEndpointArgsForCall)]
	fake.randomEndpointArgsForCall = append(fake.randomEndpointArgsForCall, struct {
	}{})
	fake.recordInvocation("RandomEndpoint", []interface{}{})
	fake.randomEndpointMutex.Unlock()
	if fake.RandomEndpointStub != nil {
		return fake.RandomEndpointStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.randomEndpointReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *OrdererConnectionSource) RandomEndpointCallCount() int {
	fake.randomEndpointMutex.RLock()
	defer fake.randomEndpointMutex.RUnlock()
	return len(fake.randomEndpointArgsForCall)
}

func (fake *OrdererConnectionSource) RandomEndpointCalls(stub func() (*orderers.Endpoint, error)) {
	fake.randomEndpointMutex.Lock()
	defer fake.randomEndpointMutex.Unlock()
	fake.RandomEndpointStub = stub
}

func (fake *OrdererConnectionSource) RandomEndpointReturns(result1 *orderers.Endpoint, result2 error) {
	fake.randomEndpointMutex.Lock()
	defer fake.randomEndpointMutex.Unlock()
	fake.RandomEndpointStub = nil
	fake.randomEndpointReturns = struct {
		result1 *orderers.Endpoint
		result2 error
	}{result1, result2}
}

func (fake *OrdererConnectionSource) RandomEndpointReturnsOnCall(i int, result1 *orderers.Endpoint, result2 error) {
	fake.randomEndpointMutex.Lock()
	defer fake.randomEndpointMutex.Unlock()
	fake.RandomEndpointStub = nil
	if fake.randomEndpointReturnsOnCall == nil {
		fake.randomEndpointReturnsOnCall = make(map[int]struct {
			result1 *orderers.Endpoint
			result2 error
		})
	}
	fake.randomEndpointReturnsOnCall[i] = struct {
		result1 *orderers.Endpoint
		result2 error
	}{result1, result2}
}

func (fake *OrdererConnectionSource) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.randomEndpointMutex.RLock()
	defer fake.randomEndpointMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OrdererConnectionSource) recordInvocation(key string, args []interface{}) {
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

var _ blocksprovider.OrdererConnectionSource = new(OrdererConnectionSource)
