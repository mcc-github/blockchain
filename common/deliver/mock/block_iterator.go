
package mock

import (
	"sync"

	cb "github.com/mcc-github/blockchain/protos/common"
)

type BlockIterator struct {
	NextStub        func() (*cb.Block, cb.Status)
	nextMutex       sync.RWMutex
	nextArgsForCall []struct{}
	nextReturns     struct {
		result1 *cb.Block
		result2 cb.Status
	}
	nextReturnsOnCall map[int]struct {
		result1 *cb.Block
		result2 cb.Status
	}
	ReadyChanStub        func() <-chan struct{}
	readyChanMutex       sync.RWMutex
	readyChanArgsForCall []struct{}
	readyChanReturns     struct {
		result1 <-chan struct{}
	}
	readyChanReturnsOnCall map[int]struct {
		result1 <-chan struct{}
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BlockIterator) Next() (*cb.Block, cb.Status) {
	fake.nextMutex.Lock()
	ret, specificReturn := fake.nextReturnsOnCall[len(fake.nextArgsForCall)]
	fake.nextArgsForCall = append(fake.nextArgsForCall, struct{}{})
	fake.recordInvocation("Next", []interface{}{})
	fake.nextMutex.Unlock()
	if fake.NextStub != nil {
		return fake.NextStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.nextReturns.result1, fake.nextReturns.result2
}

func (fake *BlockIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *BlockIterator) NextReturns(result1 *cb.Block, result2 cb.Status) {
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 *cb.Block
		result2 cb.Status
	}{result1, result2}
}

func (fake *BlockIterator) NextReturnsOnCall(i int, result1 *cb.Block, result2 cb.Status) {
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 *cb.Block
			result2 cb.Status
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 *cb.Block
		result2 cb.Status
	}{result1, result2}
}

func (fake *BlockIterator) ReadyChan() <-chan struct{} {
	fake.readyChanMutex.Lock()
	ret, specificReturn := fake.readyChanReturnsOnCall[len(fake.readyChanArgsForCall)]
	fake.readyChanArgsForCall = append(fake.readyChanArgsForCall, struct{}{})
	fake.recordInvocation("ReadyChan", []interface{}{})
	fake.readyChanMutex.Unlock()
	if fake.ReadyChanStub != nil {
		return fake.ReadyChanStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.readyChanReturns.result1
}

func (fake *BlockIterator) ReadyChanCallCount() int {
	fake.readyChanMutex.RLock()
	defer fake.readyChanMutex.RUnlock()
	return len(fake.readyChanArgsForCall)
}

func (fake *BlockIterator) ReadyChanReturns(result1 <-chan struct{}) {
	fake.ReadyChanStub = nil
	fake.readyChanReturns = struct {
		result1 <-chan struct{}
	}{result1}
}

func (fake *BlockIterator) ReadyChanReturnsOnCall(i int, result1 <-chan struct{}) {
	fake.ReadyChanStub = nil
	if fake.readyChanReturnsOnCall == nil {
		fake.readyChanReturnsOnCall = make(map[int]struct {
			result1 <-chan struct{}
		})
	}
	fake.readyChanReturnsOnCall[i] = struct {
		result1 <-chan struct{}
	}{result1}
}

func (fake *BlockIterator) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *BlockIterator) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *BlockIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	fake.readyChanMutex.RLock()
	defer fake.readyChanMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BlockIterator) recordInvocation(key string, args []interface{}) {
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
