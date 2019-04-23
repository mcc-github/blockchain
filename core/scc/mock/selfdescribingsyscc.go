
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
)

type SelfDescribingSysCC struct {
	ChaincodeStub        func() shim.Chaincode
	chaincodeMutex       sync.RWMutex
	chaincodeArgsForCall []struct {
	}
	chaincodeReturns struct {
		result1 shim.Chaincode
	}
	chaincodeReturnsOnCall map[int]struct {
		result1 shim.Chaincode
	}
	EnabledStub        func() bool
	enabledMutex       sync.RWMutex
	enabledArgsForCall []struct {
	}
	enabledReturns struct {
		result1 bool
	}
	enabledReturnsOnCall map[int]struct {
		result1 bool
	}
	InitArgsStub        func() [][]byte
	initArgsMutex       sync.RWMutex
	initArgsArgsForCall []struct {
	}
	initArgsReturns struct {
		result1 [][]byte
	}
	initArgsReturnsOnCall map[int]struct {
		result1 [][]byte
	}
	InvokableCC2CCStub        func() bool
	invokableCC2CCMutex       sync.RWMutex
	invokableCC2CCArgsForCall []struct {
	}
	invokableCC2CCReturns struct {
		result1 bool
	}
	invokableCC2CCReturnsOnCall map[int]struct {
		result1 bool
	}
	InvokableExternalStub        func() bool
	invokableExternalMutex       sync.RWMutex
	invokableExternalArgsForCall []struct {
	}
	invokableExternalReturns struct {
		result1 bool
	}
	invokableExternalReturnsOnCall map[int]struct {
		result1 bool
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
	PathStub        func() string
	pathMutex       sync.RWMutex
	pathArgsForCall []struct {
	}
	pathReturns struct {
		result1 string
	}
	pathReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SelfDescribingSysCC) Chaincode() shim.Chaincode {
	fake.chaincodeMutex.Lock()
	ret, specificReturn := fake.chaincodeReturnsOnCall[len(fake.chaincodeArgsForCall)]
	fake.chaincodeArgsForCall = append(fake.chaincodeArgsForCall, struct {
	}{})
	fake.recordInvocation("Chaincode", []interface{}{})
	fake.chaincodeMutex.Unlock()
	if fake.ChaincodeStub != nil {
		return fake.ChaincodeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.chaincodeReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) ChaincodeCallCount() int {
	fake.chaincodeMutex.RLock()
	defer fake.chaincodeMutex.RUnlock()
	return len(fake.chaincodeArgsForCall)
}

func (fake *SelfDescribingSysCC) ChaincodeCalls(stub func() shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = stub
}

func (fake *SelfDescribingSysCC) ChaincodeReturns(result1 shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = nil
	fake.chaincodeReturns = struct {
		result1 shim.Chaincode
	}{result1}
}

func (fake *SelfDescribingSysCC) ChaincodeReturnsOnCall(i int, result1 shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = nil
	if fake.chaincodeReturnsOnCall == nil {
		fake.chaincodeReturnsOnCall = make(map[int]struct {
			result1 shim.Chaincode
		})
	}
	fake.chaincodeReturnsOnCall[i] = struct {
		result1 shim.Chaincode
	}{result1}
}

func (fake *SelfDescribingSysCC) Enabled() bool {
	fake.enabledMutex.Lock()
	ret, specificReturn := fake.enabledReturnsOnCall[len(fake.enabledArgsForCall)]
	fake.enabledArgsForCall = append(fake.enabledArgsForCall, struct {
	}{})
	fake.recordInvocation("Enabled", []interface{}{})
	fake.enabledMutex.Unlock()
	if fake.EnabledStub != nil {
		return fake.EnabledStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.enabledReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) EnabledCallCount() int {
	fake.enabledMutex.RLock()
	defer fake.enabledMutex.RUnlock()
	return len(fake.enabledArgsForCall)
}

func (fake *SelfDescribingSysCC) EnabledCalls(stub func() bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = stub
}

func (fake *SelfDescribingSysCC) EnabledReturns(result1 bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = nil
	fake.enabledReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) EnabledReturnsOnCall(i int, result1 bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = nil
	if fake.enabledReturnsOnCall == nil {
		fake.enabledReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.enabledReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) InitArgs() [][]byte {
	fake.initArgsMutex.Lock()
	ret, specificReturn := fake.initArgsReturnsOnCall[len(fake.initArgsArgsForCall)]
	fake.initArgsArgsForCall = append(fake.initArgsArgsForCall, struct {
	}{})
	fake.recordInvocation("InitArgs", []interface{}{})
	fake.initArgsMutex.Unlock()
	if fake.InitArgsStub != nil {
		return fake.InitArgsStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.initArgsReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) InitArgsCallCount() int {
	fake.initArgsMutex.RLock()
	defer fake.initArgsMutex.RUnlock()
	return len(fake.initArgsArgsForCall)
}

func (fake *SelfDescribingSysCC) InitArgsCalls(stub func() [][]byte) {
	fake.initArgsMutex.Lock()
	defer fake.initArgsMutex.Unlock()
	fake.InitArgsStub = stub
}

func (fake *SelfDescribingSysCC) InitArgsReturns(result1 [][]byte) {
	fake.initArgsMutex.Lock()
	defer fake.initArgsMutex.Unlock()
	fake.InitArgsStub = nil
	fake.initArgsReturns = struct {
		result1 [][]byte
	}{result1}
}

func (fake *SelfDescribingSysCC) InitArgsReturnsOnCall(i int, result1 [][]byte) {
	fake.initArgsMutex.Lock()
	defer fake.initArgsMutex.Unlock()
	fake.InitArgsStub = nil
	if fake.initArgsReturnsOnCall == nil {
		fake.initArgsReturnsOnCall = make(map[int]struct {
			result1 [][]byte
		})
	}
	fake.initArgsReturnsOnCall[i] = struct {
		result1 [][]byte
	}{result1}
}

func (fake *SelfDescribingSysCC) InvokableCC2CC() bool {
	fake.invokableCC2CCMutex.Lock()
	ret, specificReturn := fake.invokableCC2CCReturnsOnCall[len(fake.invokableCC2CCArgsForCall)]
	fake.invokableCC2CCArgsForCall = append(fake.invokableCC2CCArgsForCall, struct {
	}{})
	fake.recordInvocation("InvokableCC2CC", []interface{}{})
	fake.invokableCC2CCMutex.Unlock()
	if fake.InvokableCC2CCStub != nil {
		return fake.InvokableCC2CCStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.invokableCC2CCReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) InvokableCC2CCCallCount() int {
	fake.invokableCC2CCMutex.RLock()
	defer fake.invokableCC2CCMutex.RUnlock()
	return len(fake.invokableCC2CCArgsForCall)
}

func (fake *SelfDescribingSysCC) InvokableCC2CCCalls(stub func() bool) {
	fake.invokableCC2CCMutex.Lock()
	defer fake.invokableCC2CCMutex.Unlock()
	fake.InvokableCC2CCStub = stub
}

func (fake *SelfDescribingSysCC) InvokableCC2CCReturns(result1 bool) {
	fake.invokableCC2CCMutex.Lock()
	defer fake.invokableCC2CCMutex.Unlock()
	fake.InvokableCC2CCStub = nil
	fake.invokableCC2CCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) InvokableCC2CCReturnsOnCall(i int, result1 bool) {
	fake.invokableCC2CCMutex.Lock()
	defer fake.invokableCC2CCMutex.Unlock()
	fake.InvokableCC2CCStub = nil
	if fake.invokableCC2CCReturnsOnCall == nil {
		fake.invokableCC2CCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.invokableCC2CCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) InvokableExternal() bool {
	fake.invokableExternalMutex.Lock()
	ret, specificReturn := fake.invokableExternalReturnsOnCall[len(fake.invokableExternalArgsForCall)]
	fake.invokableExternalArgsForCall = append(fake.invokableExternalArgsForCall, struct {
	}{})
	fake.recordInvocation("InvokableExternal", []interface{}{})
	fake.invokableExternalMutex.Unlock()
	if fake.InvokableExternalStub != nil {
		return fake.InvokableExternalStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.invokableExternalReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) InvokableExternalCallCount() int {
	fake.invokableExternalMutex.RLock()
	defer fake.invokableExternalMutex.RUnlock()
	return len(fake.invokableExternalArgsForCall)
}

func (fake *SelfDescribingSysCC) InvokableExternalCalls(stub func() bool) {
	fake.invokableExternalMutex.Lock()
	defer fake.invokableExternalMutex.Unlock()
	fake.InvokableExternalStub = stub
}

func (fake *SelfDescribingSysCC) InvokableExternalReturns(result1 bool) {
	fake.invokableExternalMutex.Lock()
	defer fake.invokableExternalMutex.Unlock()
	fake.InvokableExternalStub = nil
	fake.invokableExternalReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) InvokableExternalReturnsOnCall(i int, result1 bool) {
	fake.invokableExternalMutex.Lock()
	defer fake.invokableExternalMutex.Unlock()
	fake.InvokableExternalStub = nil
	if fake.invokableExternalReturnsOnCall == nil {
		fake.invokableExternalReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.invokableExternalReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SelfDescribingSysCC) Name() string {
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

func (fake *SelfDescribingSysCC) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *SelfDescribingSysCC) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *SelfDescribingSysCC) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *SelfDescribingSysCC) NameReturnsOnCall(i int, result1 string) {
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

func (fake *SelfDescribingSysCC) Path() string {
	fake.pathMutex.Lock()
	ret, specificReturn := fake.pathReturnsOnCall[len(fake.pathArgsForCall)]
	fake.pathArgsForCall = append(fake.pathArgsForCall, struct {
	}{})
	fake.recordInvocation("Path", []interface{}{})
	fake.pathMutex.Unlock()
	if fake.PathStub != nil {
		return fake.PathStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pathReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) PathCallCount() int {
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	return len(fake.pathArgsForCall)
}

func (fake *SelfDescribingSysCC) PathCalls(stub func() string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = stub
}

func (fake *SelfDescribingSysCC) PathReturns(result1 string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = nil
	fake.pathReturns = struct {
		result1 string
	}{result1}
}

func (fake *SelfDescribingSysCC) PathReturnsOnCall(i int, result1 string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = nil
	if fake.pathReturnsOnCall == nil {
		fake.pathReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.pathReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *SelfDescribingSysCC) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeMutex.RLock()
	defer fake.chaincodeMutex.RUnlock()
	fake.enabledMutex.RLock()
	defer fake.enabledMutex.RUnlock()
	fake.initArgsMutex.RLock()
	defer fake.initArgsMutex.RUnlock()
	fake.invokableCC2CCMutex.RLock()
	defer fake.invokableCC2CCMutex.RUnlock()
	fake.invokableExternalMutex.RLock()
	defer fake.invokableExternalMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SelfDescribingSysCC) recordInvocation(key string, args []interface{}) {
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
