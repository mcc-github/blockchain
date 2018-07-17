
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/deliver"
)

type ChainManager struct {
	GetChainStub        func(chainID string) (deliver.Chain, bool)
	getChainMutex       sync.RWMutex
	getChainArgsForCall []struct {
		chainID string
	}
	getChainReturns struct {
		result1 deliver.Chain
		result2 bool
	}
	getChainReturnsOnCall map[int]struct {
		result1 deliver.Chain
		result2 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChainManager) GetChain(chainID string) (deliver.Chain, bool) {
	fake.getChainMutex.Lock()
	ret, specificReturn := fake.getChainReturnsOnCall[len(fake.getChainArgsForCall)]
	fake.getChainArgsForCall = append(fake.getChainArgsForCall, struct {
		chainID string
	}{chainID})
	fake.recordInvocation("GetChain", []interface{}{chainID})
	fake.getChainMutex.Unlock()
	if fake.GetChainStub != nil {
		return fake.GetChainStub(chainID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChainReturns.result1, fake.getChainReturns.result2
}

func (fake *ChainManager) GetChainCallCount() int {
	fake.getChainMutex.RLock()
	defer fake.getChainMutex.RUnlock()
	return len(fake.getChainArgsForCall)
}

func (fake *ChainManager) GetChainArgsForCall(i int) string {
	fake.getChainMutex.RLock()
	defer fake.getChainMutex.RUnlock()
	return fake.getChainArgsForCall[i].chainID
}

func (fake *ChainManager) GetChainReturns(result1 deliver.Chain, result2 bool) {
	fake.GetChainStub = nil
	fake.getChainReturns = struct {
		result1 deliver.Chain
		result2 bool
	}{result1, result2}
}

func (fake *ChainManager) GetChainReturnsOnCall(i int, result1 deliver.Chain, result2 bool) {
	fake.GetChainStub = nil
	if fake.getChainReturnsOnCall == nil {
		fake.getChainReturnsOnCall = make(map[int]struct {
			result1 deliver.Chain
			result2 bool
		})
	}
	fake.getChainReturnsOnCall[i] = struct {
		result1 deliver.Chain
		result2 bool
	}{result1, result2}
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
