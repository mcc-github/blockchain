
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	ab "github.com/mcc-github/blockchain/protos/orderer"
)

type BlockReader struct {
	IteratorStub        func(startType *ab.SeekPosition) (blockledger.Iterator, uint64)
	iteratorMutex       sync.RWMutex
	iteratorArgsForCall []struct {
		startType *ab.SeekPosition
	}
	iteratorReturns struct {
		result1 blockledger.Iterator
		result2 uint64
	}
	iteratorReturnsOnCall map[int]struct {
		result1 blockledger.Iterator
		result2 uint64
	}
	HeightStub        func() uint64
	heightMutex       sync.RWMutex
	heightArgsForCall []struct{}
	heightReturns     struct {
		result1 uint64
	}
	heightReturnsOnCall map[int]struct {
		result1 uint64
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BlockReader) Iterator(startType *ab.SeekPosition) (blockledger.Iterator, uint64) {
	fake.iteratorMutex.Lock()
	ret, specificReturn := fake.iteratorReturnsOnCall[len(fake.iteratorArgsForCall)]
	fake.iteratorArgsForCall = append(fake.iteratorArgsForCall, struct {
		startType *ab.SeekPosition
	}{startType})
	fake.recordInvocation("Iterator", []interface{}{startType})
	fake.iteratorMutex.Unlock()
	if fake.IteratorStub != nil {
		return fake.IteratorStub(startType)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.iteratorReturns.result1, fake.iteratorReturns.result2
}

func (fake *BlockReader) IteratorCallCount() int {
	fake.iteratorMutex.RLock()
	defer fake.iteratorMutex.RUnlock()
	return len(fake.iteratorArgsForCall)
}

func (fake *BlockReader) IteratorArgsForCall(i int) *ab.SeekPosition {
	fake.iteratorMutex.RLock()
	defer fake.iteratorMutex.RUnlock()
	return fake.iteratorArgsForCall[i].startType
}

func (fake *BlockReader) IteratorReturns(result1 blockledger.Iterator, result2 uint64) {
	fake.IteratorStub = nil
	fake.iteratorReturns = struct {
		result1 blockledger.Iterator
		result2 uint64
	}{result1, result2}
}

func (fake *BlockReader) IteratorReturnsOnCall(i int, result1 blockledger.Iterator, result2 uint64) {
	fake.IteratorStub = nil
	if fake.iteratorReturnsOnCall == nil {
		fake.iteratorReturnsOnCall = make(map[int]struct {
			result1 blockledger.Iterator
			result2 uint64
		})
	}
	fake.iteratorReturnsOnCall[i] = struct {
		result1 blockledger.Iterator
		result2 uint64
	}{result1, result2}
}

func (fake *BlockReader) Height() uint64 {
	fake.heightMutex.Lock()
	ret, specificReturn := fake.heightReturnsOnCall[len(fake.heightArgsForCall)]
	fake.heightArgsForCall = append(fake.heightArgsForCall, struct{}{})
	fake.recordInvocation("Height", []interface{}{})
	fake.heightMutex.Unlock()
	if fake.HeightStub != nil {
		return fake.HeightStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.heightReturns.result1
}

func (fake *BlockReader) HeightCallCount() int {
	fake.heightMutex.RLock()
	defer fake.heightMutex.RUnlock()
	return len(fake.heightArgsForCall)
}

func (fake *BlockReader) HeightReturns(result1 uint64) {
	fake.HeightStub = nil
	fake.heightReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *BlockReader) HeightReturnsOnCall(i int, result1 uint64) {
	fake.HeightStub = nil
	if fake.heightReturnsOnCall == nil {
		fake.heightReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.heightReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *BlockReader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.iteratorMutex.RLock()
	defer fake.iteratorMutex.RUnlock()
	fake.heightMutex.RLock()
	defer fake.heightMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BlockReader) recordInvocation(key string, args []interface{}) {
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
