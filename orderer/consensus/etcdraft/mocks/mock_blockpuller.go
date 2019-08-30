
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
)

type FakeBlockPuller struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	HeightsByEndpointsStub        func() (map[string]uint64, error)
	heightsByEndpointsMutex       sync.RWMutex
	heightsByEndpointsArgsForCall []struct {
	}
	heightsByEndpointsReturns struct {
		result1 map[string]uint64
		result2 error
	}
	heightsByEndpointsReturnsOnCall map[int]struct {
		result1 map[string]uint64
		result2 error
	}
	PullBlockStub        func(uint64) *common.Block
	pullBlockMutex       sync.RWMutex
	pullBlockArgsForCall []struct {
		arg1 uint64
	}
	pullBlockReturns struct {
		result1 *common.Block
	}
	pullBlockReturnsOnCall map[int]struct {
		result1 *common.Block
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBlockPuller) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeBlockPuller) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeBlockPuller) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeBlockPuller) HeightsByEndpoints() (map[string]uint64, error) {
	fake.heightsByEndpointsMutex.Lock()
	ret, specificReturn := fake.heightsByEndpointsReturnsOnCall[len(fake.heightsByEndpointsArgsForCall)]
	fake.heightsByEndpointsArgsForCall = append(fake.heightsByEndpointsArgsForCall, struct {
	}{})
	fake.recordInvocation("HeightsByEndpoints", []interface{}{})
	fake.heightsByEndpointsMutex.Unlock()
	if fake.HeightsByEndpointsStub != nil {
		return fake.HeightsByEndpointsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.heightsByEndpointsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBlockPuller) HeightsByEndpointsCallCount() int {
	fake.heightsByEndpointsMutex.RLock()
	defer fake.heightsByEndpointsMutex.RUnlock()
	return len(fake.heightsByEndpointsArgsForCall)
}

func (fake *FakeBlockPuller) HeightsByEndpointsCalls(stub func() (map[string]uint64, error)) {
	fake.heightsByEndpointsMutex.Lock()
	defer fake.heightsByEndpointsMutex.Unlock()
	fake.HeightsByEndpointsStub = stub
}

func (fake *FakeBlockPuller) HeightsByEndpointsReturns(result1 map[string]uint64, result2 error) {
	fake.heightsByEndpointsMutex.Lock()
	defer fake.heightsByEndpointsMutex.Unlock()
	fake.HeightsByEndpointsStub = nil
	fake.heightsByEndpointsReturns = struct {
		result1 map[string]uint64
		result2 error
	}{result1, result2}
}

func (fake *FakeBlockPuller) HeightsByEndpointsReturnsOnCall(i int, result1 map[string]uint64, result2 error) {
	fake.heightsByEndpointsMutex.Lock()
	defer fake.heightsByEndpointsMutex.Unlock()
	fake.HeightsByEndpointsStub = nil
	if fake.heightsByEndpointsReturnsOnCall == nil {
		fake.heightsByEndpointsReturnsOnCall = make(map[int]struct {
			result1 map[string]uint64
			result2 error
		})
	}
	fake.heightsByEndpointsReturnsOnCall[i] = struct {
		result1 map[string]uint64
		result2 error
	}{result1, result2}
}

func (fake *FakeBlockPuller) PullBlock(arg1 uint64) *common.Block {
	fake.pullBlockMutex.Lock()
	ret, specificReturn := fake.pullBlockReturnsOnCall[len(fake.pullBlockArgsForCall)]
	fake.pullBlockArgsForCall = append(fake.pullBlockArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("PullBlock", []interface{}{arg1})
	fake.pullBlockMutex.Unlock()
	if fake.PullBlockStub != nil {
		return fake.PullBlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pullBlockReturns
	return fakeReturns.result1
}

func (fake *FakeBlockPuller) PullBlockCallCount() int {
	fake.pullBlockMutex.RLock()
	defer fake.pullBlockMutex.RUnlock()
	return len(fake.pullBlockArgsForCall)
}

func (fake *FakeBlockPuller) PullBlockCalls(stub func(uint64) *common.Block) {
	fake.pullBlockMutex.Lock()
	defer fake.pullBlockMutex.Unlock()
	fake.PullBlockStub = stub
}

func (fake *FakeBlockPuller) PullBlockArgsForCall(i int) uint64 {
	fake.pullBlockMutex.RLock()
	defer fake.pullBlockMutex.RUnlock()
	argsForCall := fake.pullBlockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBlockPuller) PullBlockReturns(result1 *common.Block) {
	fake.pullBlockMutex.Lock()
	defer fake.pullBlockMutex.Unlock()
	fake.PullBlockStub = nil
	fake.pullBlockReturns = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeBlockPuller) PullBlockReturnsOnCall(i int, result1 *common.Block) {
	fake.pullBlockMutex.Lock()
	defer fake.pullBlockMutex.Unlock()
	fake.PullBlockStub = nil
	if fake.pullBlockReturnsOnCall == nil {
		fake.pullBlockReturnsOnCall = make(map[int]struct {
			result1 *common.Block
		})
	}
	fake.pullBlockReturnsOnCall[i] = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeBlockPuller) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.heightsByEndpointsMutex.RLock()
	defer fake.heightsByEndpointsMutex.RUnlock()
	fake.pullBlockMutex.RLock()
	defer fake.pullBlockMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBlockPuller) recordInvocation(key string, args []interface{}) {
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

var _ etcdraft.BlockPuller = new(FakeBlockPuller)
