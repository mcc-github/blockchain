
package mock

import (
	"os"
	"sync"
)

type IOReadWriter struct {
	ExistsStub        func(string) (bool, error)
	existsMutex       sync.RWMutex
	existsArgsForCall []struct {
		arg1 string
	}
	existsReturns struct {
		result1 bool
		result2 error
	}
	existsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	MakeDirStub        func(string, os.FileMode) error
	makeDirMutex       sync.RWMutex
	makeDirArgsForCall []struct {
		arg1 string
		arg2 os.FileMode
	}
	makeDirReturns struct {
		result1 error
	}
	makeDirReturnsOnCall map[int]struct {
		result1 error
	}
	ReadDirStub        func(string) ([]os.FileInfo, error)
	readDirMutex       sync.RWMutex
	readDirArgsForCall []struct {
		arg1 string
	}
	readDirReturns struct {
		result1 []os.FileInfo
		result2 error
	}
	readDirReturnsOnCall map[int]struct {
		result1 []os.FileInfo
		result2 error
	}
	ReadFileStub        func(string) ([]byte, error)
	readFileMutex       sync.RWMutex
	readFileArgsForCall []struct {
		arg1 string
	}
	readFileReturns struct {
		result1 []byte
		result2 error
	}
	readFileReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	RemoveStub        func(string) error
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
		arg1 string
	}
	removeReturns struct {
		result1 error
	}
	removeReturnsOnCall map[int]struct {
		result1 error
	}
	WriteFileStub        func(string, string, []byte) error
	writeFileMutex       sync.RWMutex
	writeFileArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []byte
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

func (fake *IOReadWriter) Exists(arg1 string) (bool, error) {
	fake.existsMutex.Lock()
	ret, specificReturn := fake.existsReturnsOnCall[len(fake.existsArgsForCall)]
	fake.existsArgsForCall = append(fake.existsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Exists", []interface{}{arg1})
	fake.existsMutex.Unlock()
	if fake.ExistsStub != nil {
		return fake.ExistsStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.existsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *IOReadWriter) ExistsCallCount() int {
	fake.existsMutex.RLock()
	defer fake.existsMutex.RUnlock()
	return len(fake.existsArgsForCall)
}

func (fake *IOReadWriter) ExistsCalls(stub func(string) (bool, error)) {
	fake.existsMutex.Lock()
	defer fake.existsMutex.Unlock()
	fake.ExistsStub = stub
}

func (fake *IOReadWriter) ExistsArgsForCall(i int) string {
	fake.existsMutex.RLock()
	defer fake.existsMutex.RUnlock()
	argsForCall := fake.existsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IOReadWriter) ExistsReturns(result1 bool, result2 error) {
	fake.existsMutex.Lock()
	defer fake.existsMutex.Unlock()
	fake.ExistsStub = nil
	fake.existsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) ExistsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.existsMutex.Lock()
	defer fake.existsMutex.Unlock()
	fake.ExistsStub = nil
	if fake.existsReturnsOnCall == nil {
		fake.existsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.existsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) MakeDir(arg1 string, arg2 os.FileMode) error {
	fake.makeDirMutex.Lock()
	ret, specificReturn := fake.makeDirReturnsOnCall[len(fake.makeDirArgsForCall)]
	fake.makeDirArgsForCall = append(fake.makeDirArgsForCall, struct {
		arg1 string
		arg2 os.FileMode
	}{arg1, arg2})
	fake.recordInvocation("MakeDir", []interface{}{arg1, arg2})
	fake.makeDirMutex.Unlock()
	if fake.MakeDirStub != nil {
		return fake.MakeDirStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.makeDirReturns
	return fakeReturns.result1
}

func (fake *IOReadWriter) MakeDirCallCount() int {
	fake.makeDirMutex.RLock()
	defer fake.makeDirMutex.RUnlock()
	return len(fake.makeDirArgsForCall)
}

func (fake *IOReadWriter) MakeDirCalls(stub func(string, os.FileMode) error) {
	fake.makeDirMutex.Lock()
	defer fake.makeDirMutex.Unlock()
	fake.MakeDirStub = stub
}

func (fake *IOReadWriter) MakeDirArgsForCall(i int) (string, os.FileMode) {
	fake.makeDirMutex.RLock()
	defer fake.makeDirMutex.RUnlock()
	argsForCall := fake.makeDirArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *IOReadWriter) MakeDirReturns(result1 error) {
	fake.makeDirMutex.Lock()
	defer fake.makeDirMutex.Unlock()
	fake.MakeDirStub = nil
	fake.makeDirReturns = struct {
		result1 error
	}{result1}
}

func (fake *IOReadWriter) MakeDirReturnsOnCall(i int, result1 error) {
	fake.makeDirMutex.Lock()
	defer fake.makeDirMutex.Unlock()
	fake.MakeDirStub = nil
	if fake.makeDirReturnsOnCall == nil {
		fake.makeDirReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.makeDirReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *IOReadWriter) ReadDir(arg1 string) ([]os.FileInfo, error) {
	fake.readDirMutex.Lock()
	ret, specificReturn := fake.readDirReturnsOnCall[len(fake.readDirArgsForCall)]
	fake.readDirArgsForCall = append(fake.readDirArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ReadDir", []interface{}{arg1})
	fake.readDirMutex.Unlock()
	if fake.ReadDirStub != nil {
		return fake.ReadDirStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.readDirReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *IOReadWriter) ReadDirCallCount() int {
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	return len(fake.readDirArgsForCall)
}

func (fake *IOReadWriter) ReadDirCalls(stub func(string) ([]os.FileInfo, error)) {
	fake.readDirMutex.Lock()
	defer fake.readDirMutex.Unlock()
	fake.ReadDirStub = stub
}

func (fake *IOReadWriter) ReadDirArgsForCall(i int) string {
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	argsForCall := fake.readDirArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IOReadWriter) ReadDirReturns(result1 []os.FileInfo, result2 error) {
	fake.readDirMutex.Lock()
	defer fake.readDirMutex.Unlock()
	fake.ReadDirStub = nil
	fake.readDirReturns = struct {
		result1 []os.FileInfo
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) ReadDirReturnsOnCall(i int, result1 []os.FileInfo, result2 error) {
	fake.readDirMutex.Lock()
	defer fake.readDirMutex.Unlock()
	fake.ReadDirStub = nil
	if fake.readDirReturnsOnCall == nil {
		fake.readDirReturnsOnCall = make(map[int]struct {
			result1 []os.FileInfo
			result2 error
		})
	}
	fake.readDirReturnsOnCall[i] = struct {
		result1 []os.FileInfo
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) ReadFile(arg1 string) ([]byte, error) {
	fake.readFileMutex.Lock()
	ret, specificReturn := fake.readFileReturnsOnCall[len(fake.readFileArgsForCall)]
	fake.readFileArgsForCall = append(fake.readFileArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ReadFile", []interface{}{arg1})
	fake.readFileMutex.Unlock()
	if fake.ReadFileStub != nil {
		return fake.ReadFileStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.readFileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *IOReadWriter) ReadFileCallCount() int {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return len(fake.readFileArgsForCall)
}

func (fake *IOReadWriter) ReadFileCalls(stub func(string) ([]byte, error)) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = stub
}

func (fake *IOReadWriter) ReadFileArgsForCall(i int) string {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	argsForCall := fake.readFileArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IOReadWriter) ReadFileReturns(result1 []byte, result2 error) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = nil
	fake.readFileReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) ReadFileReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.readFileMutex.Lock()
	defer fake.readFileMutex.Unlock()
	fake.ReadFileStub = nil
	if fake.readFileReturnsOnCall == nil {
		fake.readFileReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.readFileReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *IOReadWriter) Remove(arg1 string) error {
	fake.removeMutex.Lock()
	ret, specificReturn := fake.removeReturnsOnCall[len(fake.removeArgsForCall)]
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Remove", []interface{}{arg1})
	fake.removeMutex.Unlock()
	if fake.RemoveStub != nil {
		return fake.RemoveStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeReturns
	return fakeReturns.result1
}

func (fake *IOReadWriter) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *IOReadWriter) RemoveCalls(stub func(string) error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = stub
}

func (fake *IOReadWriter) RemoveArgsForCall(i int) string {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	argsForCall := fake.removeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *IOReadWriter) RemoveReturns(result1 error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = nil
	fake.removeReturns = struct {
		result1 error
	}{result1}
}

func (fake *IOReadWriter) RemoveReturnsOnCall(i int, result1 error) {
	fake.removeMutex.Lock()
	defer fake.removeMutex.Unlock()
	fake.RemoveStub = nil
	if fake.removeReturnsOnCall == nil {
		fake.removeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *IOReadWriter) WriteFile(arg1 string, arg2 string, arg3 []byte) error {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.writeFileMutex.Lock()
	ret, specificReturn := fake.writeFileReturnsOnCall[len(fake.writeFileArgsForCall)]
	fake.writeFileArgsForCall = append(fake.writeFileArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("WriteFile", []interface{}{arg1, arg2, arg3Copy})
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

func (fake *IOReadWriter) WriteFileCallCount() int {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	return len(fake.writeFileArgsForCall)
}

func (fake *IOReadWriter) WriteFileCalls(stub func(string, string, []byte) error) {
	fake.writeFileMutex.Lock()
	defer fake.writeFileMutex.Unlock()
	fake.WriteFileStub = stub
}

func (fake *IOReadWriter) WriteFileArgsForCall(i int) (string, string, []byte) {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	argsForCall := fake.writeFileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *IOReadWriter) WriteFileReturns(result1 error) {
	fake.writeFileMutex.Lock()
	defer fake.writeFileMutex.Unlock()
	fake.WriteFileStub = nil
	fake.writeFileReturns = struct {
		result1 error
	}{result1}
}

func (fake *IOReadWriter) WriteFileReturnsOnCall(i int, result1 error) {
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

func (fake *IOReadWriter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.existsMutex.RLock()
	defer fake.existsMutex.RUnlock()
	fake.makeDirMutex.RLock()
	defer fake.makeDirMutex.RUnlock()
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *IOReadWriter) recordInvocation(key string, args []interface{}) {
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
