
package mock

import (
	sync "sync"

	deliver "github.com/mcc-github/blockchain/common/deliver"
)

type ChainManager struct {
	GetChainStub        func(string) deliver.Chain
	getChainMutex       sync.RWMutex
	getChainArgsForCall []struct {
		arg1 string
	}
	getChainReturns struct {
		result1 deliver.Chain
	}
	getChainReturnsOnCall map[int]struct {
		result1 deliver.Chain
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChainManager) GetChain(arg1 string) deliver.Chain {
	fake.getChainMutex.Lock()
	ret, specificReturn := fake.getChainReturnsOnCall[len(fake.getChainArgsForCall)]
	fake.getChainArgsForCall = append(fake.getChainArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetChain", []interface{}{arg1})
	fake.getChainMutex.Unlock()
	if fake.GetChainStub != nil {
		return fake.GetChainStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getChainReturns
	return fakeReturns.result1
}

func (fake *ChainManager) GetChainCallCount() int {
	fake.getChainMutex.RLock()
	defer fake.getChainMutex.RUnlock()
	return len(fake.getChainArgsForCall)
}

func (fake *ChainManager) GetChainCalls(stub func(string) deliver.Chain) {
	fake.getChainMutex.Lock()
	defer fake.getChainMutex.Unlock()
	fake.GetChainStub = stub
}

func (fake *ChainManager) GetChainArgsForCall(i int) string {
	fake.getChainMutex.RLock()
	defer fake.getChainMutex.RUnlock()
	argsForCall := fake.getChainArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChainManager) GetChainReturns(result1 deliver.Chain) {
	fake.getChainMutex.Lock()
	defer fake.getChainMutex.Unlock()
	fake.GetChainStub = nil
	fake.getChainReturns = struct {
		result1 deliver.Chain
	}{result1}
}

func (fake *ChainManager) GetChainReturnsOnCall(i int, result1 deliver.Chain) {
	fake.getChainMutex.Lock()
	defer fake.getChainMutex.Unlock()
	fake.GetChainStub = nil
	if fake.getChainReturnsOnCall == nil {
		fake.getChainReturnsOnCall = make(map[int]struct {
			result1 deliver.Chain
		})
	}
	fake.getChainReturnsOnCall[i] = struct {
		result1 deliver.Chain
	}{result1}
}

func (fake *ChainManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChainMutex.RLock()
	defer fake.getChainMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChainManager) recordInvocation(key string, args []interface{}) {
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

var _ deliver.ChainManager = new(ChainManager)
