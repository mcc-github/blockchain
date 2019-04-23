
package mock

import (
	"context"
	"sync"

	"github.com/mcc-github/blockchain/protos/peer"
	"google.golang.org/grpc"
)

type PeerDeliverClient struct {
	DeliverStub        func(context.Context, ...grpc.CallOption) (peer.Deliver_DeliverClient, error)
	deliverMutex       sync.RWMutex
	deliverArgsForCall []struct {
		arg1 context.Context
		arg2 []grpc.CallOption
	}
	deliverReturns struct {
		result1 peer.Deliver_DeliverClient
		result2 error
	}
	deliverReturnsOnCall map[int]struct {
		result1 peer.Deliver_DeliverClient
		result2 error
	}
	DeliverFilteredStub        func(context.Context, ...grpc.CallOption) (peer.Deliver_DeliverFilteredClient, error)
	deliverFilteredMutex       sync.RWMutex
	deliverFilteredArgsForCall []struct {
		arg1 context.Context
		arg2 []grpc.CallOption
	}
	deliverFilteredReturns struct {
		result1 peer.Deliver_DeliverFilteredClient
		result2 error
	}
	deliverFilteredReturnsOnCall map[int]struct {
		result1 peer.Deliver_DeliverFilteredClient
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PeerDeliverClient) Deliver(arg1 context.Context, arg2 ...grpc.CallOption) (peer.Deliver_DeliverClient, error) {
	fake.deliverMutex.Lock()
	ret, specificReturn := fake.deliverReturnsOnCall[len(fake.deliverArgsForCall)]
	fake.deliverArgsForCall = append(fake.deliverArgsForCall, struct {
		arg1 context.Context
		arg2 []grpc.CallOption
	}{arg1, arg2})
	fake.recordInvocation("Deliver", []interface{}{arg1, arg2})
	fake.deliverMutex.Unlock()
	if fake.DeliverStub != nil {
		return fake.DeliverStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deliverReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerDeliverClient) DeliverCallCount() int {
	fake.deliverMutex.RLock()
	defer fake.deliverMutex.RUnlock()
	return len(fake.deliverArgsForCall)
}

func (fake *PeerDeliverClient) DeliverCalls(stub func(context.Context, ...grpc.CallOption) (peer.Deliver_DeliverClient, error)) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = stub
}

func (fake *PeerDeliverClient) DeliverArgsForCall(i int) (context.Context, []grpc.CallOption) {
	fake.deliverMutex.RLock()
	defer fake.deliverMutex.RUnlock()
	argsForCall := fake.deliverArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerDeliverClient) DeliverReturns(result1 peer.Deliver_DeliverClient, result2 error) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = nil
	fake.deliverReturns = struct {
		result1 peer.Deliver_DeliverClient
		result2 error
	}{result1, result2}
}

func (fake *PeerDeliverClient) DeliverReturnsOnCall(i int, result1 peer.Deliver_DeliverClient, result2 error) {
	fake.deliverMutex.Lock()
	defer fake.deliverMutex.Unlock()
	fake.DeliverStub = nil
	if fake.deliverReturnsOnCall == nil {
		fake.deliverReturnsOnCall = make(map[int]struct {
			result1 peer.Deliver_DeliverClient
			result2 error
		})
	}
	fake.deliverReturnsOnCall[i] = struct {
		result1 peer.Deliver_DeliverClient
		result2 error
	}{result1, result2}
}

func (fake *PeerDeliverClient) DeliverFiltered(arg1 context.Context, arg2 ...grpc.CallOption) (peer.Deliver_DeliverFilteredClient, error) {
	fake.deliverFilteredMutex.Lock()
	ret, specificReturn := fake.deliverFilteredReturnsOnCall[len(fake.deliverFilteredArgsForCall)]
	fake.deliverFilteredArgsForCall = append(fake.deliverFilteredArgsForCall, struct {
		arg1 context.Context
		arg2 []grpc.CallOption
	}{arg1, arg2})
	fake.recordInvocation("DeliverFiltered", []interface{}{arg1, arg2})
	fake.deliverFilteredMutex.Unlock()
	if fake.DeliverFilteredStub != nil {
		return fake.DeliverFilteredStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deliverFilteredReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PeerDeliverClient) DeliverFilteredCallCount() int {
	fake.deliverFilteredMutex.RLock()
	defer fake.deliverFilteredMutex.RUnlock()
	return len(fake.deliverFilteredArgsForCall)
}

func (fake *PeerDeliverClient) DeliverFilteredCalls(stub func(context.Context, ...grpc.CallOption) (peer.Deliver_DeliverFilteredClient, error)) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = stub
}

func (fake *PeerDeliverClient) DeliverFilteredArgsForCall(i int) (context.Context, []grpc.CallOption) {
	fake.deliverFilteredMutex.RLock()
	defer fake.deliverFilteredMutex.RUnlock()
	argsForCall := fake.deliverFilteredArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PeerDeliverClient) DeliverFilteredReturns(result1 peer.Deliver_DeliverFilteredClient, result2 error) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = nil
	fake.deliverFilteredReturns = struct {
		result1 peer.Deliver_DeliverFilteredClient
		result2 error
	}{result1, result2}
}

func (fake *PeerDeliverClient) DeliverFilteredReturnsOnCall(i int, result1 peer.Deliver_DeliverFilteredClient, result2 error) {
	fake.deliverFilteredMutex.Lock()
	defer fake.deliverFilteredMutex.Unlock()
	fake.DeliverFilteredStub = nil
	if fake.deliverFilteredReturnsOnCall == nil {
		fake.deliverFilteredReturnsOnCall = make(map[int]struct {
			result1 peer.Deliver_DeliverFilteredClient
			result2 error
		})
	}
	fake.deliverFilteredReturnsOnCall[i] = struct {
		result1 peer.Deliver_DeliverFilteredClient
		result2 error
	}{result1, result2}
}

func (fake *PeerDeliverClient) Invocations() map[string][][]interface{} {
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

func (fake *PeerDeliverClient) recordInvocation(key string, args []interface{}) {
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
