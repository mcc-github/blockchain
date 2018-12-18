
package mock

import (
	sync "sync"

	peer "github.com/mcc-github/blockchain/protos/peer"
)

type DeliverServer struct {
	DeliverStub        func(peer.Deliver_DeliverServer) error
	deliverMutex       sync.RWMutex
	deliverArgsForCall []struct {
		arg1 peer.Deliver_DeliverServer
	}
	deliverReturns struct {
		result1 error
	}
	deliverReturnsOnCall map[int]struct {
		result1 error
	}
	DeliverFilteredStub        func(peer.Deliver_DeliverFilteredServer) error
	deliverFilteredMutex       sync.RWMutex
	deliverFilteredArgsForCall []struct {
		arg1 peer.Deliver_DeliverFilteredServer
	}
	deliverFilteredReturns struct {
		result1 error
	}
	deliverFilteredReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DeliverServer) Deliver(arg1 peer.Deliver_DeliverServer) error {
	fake.deliverMutex.Lock()
	ret, specificReturn := fake.deliverReturnsOnCall[len(fake.deliverArgsForCall)]
	fake.deliverArgsForCall = append(fake.deliverArgsForCall, struct {
		arg1 peer.Deliver_DeliverServer
	}{arg1})
	fake.recordInvocation("Deliver", []interface{}{arg1})
	fake.deliverMutex.Unlock()
	if fake.DeliverStub != nil {
		return fake.DeliverStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deliverReturns
	return fakeReturns.result1
}

func (fake *DeliverServer) DeliverCallCount() int {
	fake.deliverMutex.RLock()
	defer fake.deliverMutex.RUnlock()
	return len(fake.deliverArgsForCall)
}

func (fake *DeliverServer) DeliverCalls(stub func(peer.Deliver_DeliverServer) error) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = stub
}

func (fake *DeliverServer) DeliverArgsForCall(i int) peer.Deliver_DeliverServer {
	fake.deliverMutex.RLock()
	defer fake.deliverMutex.RUnlock()
	argsForCall := fake.deliverArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverServer) DeliverReturns(result1 error) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = nil
	fake.deliverReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverServer) DeliverReturnsOnCall(i int, result1 error) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = nil
	if fake.deliverReturnsOnCall == nil {
		fake.deliverReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deliverReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverServer) DeliverFiltered(arg1 peer.Deliver_DeliverFilteredServer) error {
	fake.deliverFilteredMutex.Lock()
	ret, specificReturn := fake.deliverFilteredReturnsOnCall[len(fake.deliverFilteredArgsForCall)]
	fake.deliverFilteredArgsForCall = append(fake.deliverFilteredArgsForCall, struct {
		arg1 peer.Deliver_DeliverFilteredServer
	}{arg1})
	fake.recordInvocation("DeliverFiltered", []interface{}{arg1})
	fake.deliverFilteredMutex.Unlock()
	if fake.DeliverFilteredStub != nil {
		return fake.DeliverFilteredStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deliverFilteredReturns
	return fakeReturns.result1
}

func (fake *DeliverServer) DeliverFilteredCallCount() int {
	fake.deliverFilteredMutex.RLock()
	defer fake.deliverFilteredMutex.RUnlock()
	return len(fake.deliverFilteredArgsForCall)
}

func (fake *DeliverServer) DeliverFilteredCalls(stub func(peer.Deliver_DeliverFilteredServer) error) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = stub
}

func (fake *DeliverServer) DeliverFilteredArgsForCall(i int) peer.Deliver_DeliverFilteredServer {
	fake.deliverFilteredMutex.RLock()
	defer fake.deliverFilteredMutex.RUnlock()
	argsForCall := fake.deliverFilteredArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverServer) DeliverFilteredReturns(result1 error) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = nil
	fake.deliverFilteredReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverServer) DeliverFilteredReturnsOnCall(i int, result1 error) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = nil
	if fake.deliverFilteredReturnsOnCall == nil {
		fake.deliverFilteredReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deliverFilteredReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverServer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deliverMutex.RLock()
	defer fake.deliverMutex.RUnlock()
	fake.deliverFilteredMutex.RLock()
	defer fake.deliverFilteredMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DeliverServer) recordInvocation(key string, args []interface{}) {
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
