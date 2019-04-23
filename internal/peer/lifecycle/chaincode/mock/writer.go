
package mock

import (
	"os"
	"sync"
)

type Writer struct {
	WriteFileStub        func(string, []byte, os.FileMode) error
	writeFileMutex       sync.RWMutex
	writeFileArgsForCall []struct {
		arg1 string
		arg2 []byte
		arg3 os.FileMode
	}
	writeFileReturns struct {
		result1 error
	}
	writeFileReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Writer) WriteFile(arg1 string, arg2 []byte, arg3 os.FileMode) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.writeFileMutex.Lock()
	ret, specificReturn := fake.writeFileReturnsOnCall[len(fake.writeFileArgsForCall)]
	fake.writeFileArgsForCall = append(fake.writeFileArgsForCall, struct {
		arg1 string
		arg2 []byte
		arg3 os.FileMode
	}{arg1, arg2Copy, arg3})
	fake.recordInvocation("WriteFile", []interface{}{arg1, arg2Copy, arg3})
	fake.writeFileMutex.Unlock()
	if fake.WriteFileStub != nil {
		return fake.WriteFileStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.writeFileReturns
	return fakeReturns.result1
}

func (fake *Writer) WriteFileCallCount() int {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	return len(fake.writeFileArgsForCall)
}

func (fake *Writer) WriteFileCalls(stub func(string, []byte, os.FileMode) error) {
	fake.writeFileMutex.Lock()
	defer fake.writeFileMutex.Unlock()
	fake.WriteFileStub = stub
}

func (fake *Writer) WriteFileArgsForCall(i int) (string, []byte, os.FileMode) {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	argsForCall := fake.writeFileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Writer) WriteFileReturns(result1 error) {
	fake.writeFileMutex.Lock()
	defer fake.writeFileMutex.Unlock()
	fake.WriteFileStub = nil
	fake.writeFileReturns = struct {
		result1 error
	}{result1}
}

func (fake *Writer) WriteFileReturnsOnCall(i int, result1 error) {
	fake.writeFileMutex.Lock()
	defer fake.writeFileMutex.Unlock()
	fake.WriteFileStub = nil
	if fake.writeFileReturnsOnCall == nil {
		fake.writeFileReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.writeFileReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Writer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Writer) recordInvocation(key string, args []interface{}) {
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
