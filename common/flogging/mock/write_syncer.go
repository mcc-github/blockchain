
package mock

import (
	"sync"
)

type WriteSyncer struct {
	SyncStub        func() error
	syncMutex       sync.RWMutex
	syncArgsForCall []struct {
	}
	syncReturns struct {
		result1 error
	}
	syncReturnsOnCall map[int]struct {
		result1 error
	}
	WriteStub        func([]byte) (int, error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		arg1 []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *WriteSyncer) Sync() error {
	fake.syncMutex.Lock()
	ret, specificReturn := fake.syncReturnsOnCall[len(fake.syncArgsForCall)]
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct {
	}{})
	fake.recordInvocation("Sync", []interface{}{})
	fake.syncMutex.Unlock()
	if fake.SyncStub != nil {
		return fake.SyncStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.syncReturns
	return fakeReturns.result1
}

func (fake *WriteSyncer) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *WriteSyncer) SyncCalls(stub func() error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = stub
}

func (fake *WriteSyncer) SyncReturns(result1 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
	fake.SyncStub = nil
	fake.syncReturns = struct {
		result1 error
	}{result1}
}

func (fake *WriteSyncer) SyncReturnsOnCall(i int, result1 error) {
	fake.syncMutex.Lock()
	defer fake.syncMutex.Unlock()
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

func (fake *WriteSyncer) Write(arg1 []byte) (int, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("Write", []interface{}{arg1Copy})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.writeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *WriteSyncer) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *WriteSyncer) WriteCalls(stub func([]byte) (int, error)) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = stub
}

func (fake *WriteSyncer) WriteArgsForCall(i int) []byte {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	argsForCall := fake.writeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *WriteSyncer) WriteReturns(result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *WriteSyncer) WriteReturnsOnCall(i int, result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
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

func (fake *WriteSyncer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
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
