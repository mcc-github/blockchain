
package mock

import (
	"sync"
)

type WriteSyncer struct {
	WriteStub        func(p []byte) (n int, err error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		p []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	SyncStub        func() error
	syncMutex       sync.RWMutex
	syncArgsForCall []struct{}
	syncReturns     struct {
		result1 error
	}
	syncReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *WriteSyncer) Write(p []byte) (n int, err error) {
	var pCopy []byte
	if p != nil {
		pCopy = make([]byte, len(p))
		copy(pCopy, p)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		p []byte
	}{pCopy})
	fake.recordInvocation("Write", []interface{}{pCopy})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(p)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.writeReturns.result1, fake.writeReturns.result2
}

func (fake *WriteSyncer) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *WriteSyncer) WriteArgsForCall(i int) []byte {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return fake.writeArgsForCall[i].p
}

func (fake *WriteSyncer) WriteReturns(result1 int, result2 error) {
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *WriteSyncer) WriteReturnsOnCall(i int, result1 int, result2 error) {
	fake.WriteStub = nil
	if fake.writeReturnsOnCall == nil {
		fake.writeReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.writeReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *WriteSyncer) Sync() error {
	fake.syncMutex.Lock()
	ret, specificReturn := fake.syncReturnsOnCall[len(fake.syncArgsForCall)]
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct{}{})
	fake.recordInvocation("Sync", []interface{}{})
	fake.syncMutex.Unlock()
	if fake.SyncStub != nil {
		return fake.SyncStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.syncReturns.result1
}

func (fake *WriteSyncer) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *WriteSyncer) SyncReturns(result1 error) {
	fake.SyncStub = nil
	fake.syncReturns = struct {
		result1 error
	}{result1}
}

func (fake *WriteSyncer) SyncReturnsOnCall(i int, result1 error) {
	fake.SyncStub = nil
	if fake.syncReturnsOnCall == nil {
		fake.syncReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.syncReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *WriteSyncer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *WriteSyncer) recordInvocation(key string, args []interface{}) {
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
