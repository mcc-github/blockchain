
package mock

import (
	os "os"
	sync "sync"
	time "time"
)

type OSFileInfo struct {
	IsDirStub        func() bool
	isDirMutex       sync.RWMutex
	isDirArgsForCall []struct {
	}
	isDirReturns struct {
		result1 bool
	}
	isDirReturnsOnCall map[int]struct {
		result1 bool
	}
	ModTimeStub        func() time.Time
	modTimeMutex       sync.RWMutex
	modTimeArgsForCall []struct {
	}
	modTimeReturns struct {
		result1 time.Time
	}
	modTimeReturnsOnCall map[int]struct {
		result1 time.Time
	}
	ModeStub        func() os.FileMode
	modeMutex       sync.RWMutex
	modeArgsForCall []struct {
	}
	modeReturns struct {
		result1 os.FileMode
	}
	modeReturnsOnCall map[int]struct {
		result1 os.FileMode
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	SizeStub        func() int64
	sizeMutex       sync.RWMutex
	sizeArgsForCall []struct {
	}
	sizeReturns struct {
		result1 int64
	}
	sizeReturnsOnCall map[int]struct {
		result1 int64
	}
	SysStub        func() interface{}
	sysMutex       sync.RWMutex
	sysArgsForCall []struct {
	}
	sysReturns struct {
		result1 interface{}
	}
	sysReturnsOnCall map[int]struct {
		result1 interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OSFileInfo) IsDir() bool {
	fake.isDirMutex.Lock()
	ret, specificReturn := fake.isDirReturnsOnCall[len(fake.isDirArgsForCall)]
	fake.isDirArgsForCall = append(fake.isDirArgsForCall, struct {
	}{})
	fake.recordInvocation("IsDir", []interface{}{})
	fake.isDirMutex.Unlock()
	if fake.IsDirStub != nil {
		return fake.IsDirStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isDirReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) IsDirCallCount() int {
	fake.isDirMutex.RLock()
	defer fake.isDirMutex.RUnlock()
	return len(fake.isDirArgsForCall)
}

func (fake *OSFileInfo) IsDirCalls(stub func() bool) {
	fake.isDirMutex.Lock()
	defer fake.isDirMutex.Unlock()
	fake.IsDirStub = stub
}

func (fake *OSFileInfo) IsDirReturns(result1 bool) {
	fake.isDirMutex.Lock()
	defer fake.isDirMutex.Unlock()
	fake.IsDirStub = nil
	fake.isDirReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OSFileInfo) IsDirReturnsOnCall(i int, result1 bool) {
	fake.isDirMutex.Lock()
	defer fake.isDirMutex.Unlock()
	fake.IsDirStub = nil
	if fake.isDirReturnsOnCall == nil {
		fake.isDirReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isDirReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *OSFileInfo) ModTime() time.Time {
	fake.modTimeMutex.Lock()
	ret, specificReturn := fake.modTimeReturnsOnCall[len(fake.modTimeArgsForCall)]
	fake.modTimeArgsForCall = append(fake.modTimeArgsForCall, struct {
	}{})
	fake.recordInvocation("ModTime", []interface{}{})
	fake.modTimeMutex.Unlock()
	if fake.ModTimeStub != nil {
		return fake.ModTimeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.modTimeReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) ModTimeCallCount() int {
	fake.modTimeMutex.RLock()
	defer fake.modTimeMutex.RUnlock()
	return len(fake.modTimeArgsForCall)
}

func (fake *OSFileInfo) ModTimeCalls(stub func() time.Time) {
	fake.modTimeMutex.Lock()
	defer fake.modTimeMutex.Unlock()
	fake.ModTimeStub = stub
}

func (fake *OSFileInfo) ModTimeReturns(result1 time.Time) {
	fake.modTimeMutex.Lock()
	defer fake.modTimeMutex.Unlock()
	fake.ModTimeStub = nil
	fake.modTimeReturns = struct {
		result1 time.Time
	}{result1}
}

func (fake *OSFileInfo) ModTimeReturnsOnCall(i int, result1 time.Time) {
	fake.modTimeMutex.Lock()
	defer fake.modTimeMutex.Unlock()
	fake.ModTimeStub = nil
	if fake.modTimeReturnsOnCall == nil {
		fake.modTimeReturnsOnCall = make(map[int]struct {
			result1 time.Time
		})
	}
	fake.modTimeReturnsOnCall[i] = struct {
		result1 time.Time
	}{result1}
}

func (fake *OSFileInfo) Mode() os.FileMode {
	fake.modeMutex.Lock()
	ret, specificReturn := fake.modeReturnsOnCall[len(fake.modeArgsForCall)]
	fake.modeArgsForCall = append(fake.modeArgsForCall, struct {
	}{})
	fake.recordInvocation("Mode", []interface{}{})
	fake.modeMutex.Unlock()
	if fake.ModeStub != nil {
		return fake.ModeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.modeReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) ModeCallCount() int {
	fake.modeMutex.RLock()
	defer fake.modeMutex.RUnlock()
	return len(fake.modeArgsForCall)
}

func (fake *OSFileInfo) ModeCalls(stub func() os.FileMode) {
	fake.modeMutex.Lock()
	defer fake.modeMutex.Unlock()
	fake.ModeStub = stub
}

func (fake *OSFileInfo) ModeReturns(result1 os.FileMode) {
	fake.modeMutex.Lock()
	defer fake.modeMutex.Unlock()
	fake.ModeStub = nil
	fake.modeReturns = struct {
		result1 os.FileMode
	}{result1}
}

func (fake *OSFileInfo) ModeReturnsOnCall(i int, result1 os.FileMode) {
	fake.modeMutex.Lock()
	defer fake.modeMutex.Unlock()
	fake.ModeStub = nil
	if fake.modeReturnsOnCall == nil {
		fake.modeReturnsOnCall = make(map[int]struct {
			result1 os.FileMode
		})
	}
	fake.modeReturnsOnCall[i] = struct {
		result1 os.FileMode
	}{result1}
}

func (fake *OSFileInfo) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *OSFileInfo) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *OSFileInfo) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *OSFileInfo) NameReturnsOnCall(i int, result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *OSFileInfo) Size() int64 {
	fake.sizeMutex.Lock()
	ret, specificReturn := fake.sizeReturnsOnCall[len(fake.sizeArgsForCall)]
	fake.sizeArgsForCall = append(fake.sizeArgsForCall, struct {
	}{})
	fake.recordInvocation("Size", []interface{}{})
	fake.sizeMutex.Unlock()
	if fake.SizeStub != nil {
		return fake.SizeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sizeReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) SizeCallCount() int {
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	return len(fake.sizeArgsForCall)
}

func (fake *OSFileInfo) SizeCalls(stub func() int64) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = stub
}

func (fake *OSFileInfo) SizeReturns(result1 int64) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = nil
	fake.sizeReturns = struct {
		result1 int64
	}{result1}
}

func (fake *OSFileInfo) SizeReturnsOnCall(i int, result1 int64) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = nil
	if fake.sizeReturnsOnCall == nil {
		fake.sizeReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.sizeReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *OSFileInfo) Sys() interface{} {
	fake.sysMutex.Lock()
	ret, specificReturn := fake.sysReturnsOnCall[len(fake.sysArgsForCall)]
	fake.sysArgsForCall = append(fake.sysArgsForCall, struct {
	}{})
	fake.recordInvocation("Sys", []interface{}{})
	fake.sysMutex.Unlock()
	if fake.SysStub != nil {
		return fake.SysStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sysReturns
	return fakeReturns.result1
}

func (fake *OSFileInfo) SysCallCount() int {
	fake.sysMutex.RLock()
	defer fake.sysMutex.RUnlock()
	return len(fake.sysArgsForCall)
}

func (fake *OSFileInfo) SysCalls(stub func() interface{}) {
	fake.sysMutex.Lock()
	defer fake.sysMutex.Unlock()
	fake.SysStub = stub
}

func (fake *OSFileInfo) SysReturns(result1 interface{}) {
	fake.sysMutex.Lock()
	defer fake.sysMutex.Unlock()
	fake.SysStub = nil
	fake.sysReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *OSFileInfo) SysReturnsOnCall(i int, result1 interface{}) {
	fake.sysMutex.Lock()
	defer fake.sysMutex.Unlock()
	fake.SysStub = nil
	if fake.sysReturnsOnCall == nil {
		fake.sysReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.sysReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *OSFileInfo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isDirMutex.RLock()
	defer fake.isDirMutex.RUnlock()
	fake.modTimeMutex.RLock()
	defer fake.modTimeMutex.RUnlock()
	fake.modeMutex.RLock()
	defer fake.modeMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	fake.sysMutex.RLock()
	defer fake.sysMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OSFileInfo) recordInvocation(key string, args []interface{}) {
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
