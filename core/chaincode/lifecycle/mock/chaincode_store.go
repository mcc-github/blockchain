
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
)

type ChaincodeStore struct {
	ListInstalledChaincodesStub        func() ([]chaincode.InstalledChaincode, error)
	listInstalledChaincodesMutex       sync.RWMutex
	listInstalledChaincodesArgsForCall []struct {
	}
	listInstalledChaincodesReturns struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	listInstalledChaincodesReturnsOnCall map[int]struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	LoadStub        func(persistence.PackageID) ([]byte, error)
	loadMutex       sync.RWMutex
	loadArgsForCall []struct {
		arg1 persistence.PackageID
	}
	loadReturns struct {
		result1 []byte
		result2 error
	}
	loadReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SaveStub        func(string, []byte) (persistence.PackageID, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	saveReturns struct {
		result1 persistence.PackageID
		result2 error
	}
	saveReturnsOnCall map[int]struct {
		result1 persistence.PackageID
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeStore) ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error) {
	fake.listInstalledChaincodesMutex.Lock()
	ret, specificReturn := fake.listInstalledChaincodesReturnsOnCall[len(fake.listInstalledChaincodesArgsForCall)]
	fake.listInstalledChaincodesArgsForCall = append(fake.listInstalledChaincodesArgsForCall, struct {
	}{})
	fake.recordInvocation("ListInstalledChaincodes", []interface{}{})
	fake.listInstalledChaincodesMutex.Unlock()
	if fake.ListInstalledChaincodesStub != nil {
		return fake.ListInstalledChaincodesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listInstalledChaincodesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeStore) ListInstalledChaincodesCallCount() int {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return len(fake.listInstalledChaincodesArgsForCall)
}

func (fake *ChaincodeStore) ListInstalledChaincodesCalls(stub func() ([]chaincode.InstalledChaincode, error)) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = stub
}

func (fake *ChaincodeStore) ListInstalledChaincodesReturns(result1 []chaincode.InstalledChaincode, result2 error) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = nil
	fake.listInstalledChaincodesReturns = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) ListInstalledChaincodesReturnsOnCall(i int, result1 []chaincode.InstalledChaincode, result2 error) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = nil
	if fake.listInstalledChaincodesReturnsOnCall == nil {
		fake.listInstalledChaincodesReturnsOnCall = make(map[int]struct {
			result1 []chaincode.InstalledChaincode
			result2 error
		})
	}
	fake.listInstalledChaincodesReturnsOnCall[i] = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) Load(arg1 persistence.PackageID) ([]byte, error) {
	fake.loadMutex.Lock()
	ret, specificReturn := fake.loadReturnsOnCall[len(fake.loadArgsForCall)]
	fake.loadArgsForCall = append(fake.loadArgsForCall, struct {
		arg1 persistence.PackageID
	}{arg1})
	fake.recordInvocation("Load", []interface{}{arg1})
	fake.loadMutex.Unlock()
	if fake.LoadStub != nil {
		return fake.LoadStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.loadReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeStore) LoadCallCount() int {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return len(fake.loadArgsForCall)
}

func (fake *ChaincodeStore) LoadCalls(stub func(persistence.PackageID) ([]byte, error)) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = stub
}

func (fake *ChaincodeStore) LoadArgsForCall(i int) persistence.PackageID {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	argsForCall := fake.loadArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChaincodeStore) LoadReturns(result1 []byte, result2 error) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = nil
	fake.loadReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) LoadReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = nil
	if fake.loadReturnsOnCall == nil {
		fake.loadReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.loadReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) Save(arg1 string, arg2 []byte) (persistence.PackageID, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.saveMutex.Lock()
	ret, specificReturn := fake.saveReturnsOnCall[len(fake.saveArgsForCall)]
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("Save", []interface{}{arg1, arg2Copy})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.saveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeStore) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *ChaincodeStore) SaveCalls(stub func(string, []byte) (persistence.PackageID, error)) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = stub
}

func (fake *ChaincodeStore) SaveArgsForCall(i int) (string, []byte) {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	argsForCall := fake.saveArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeStore) SaveReturns(result1 persistence.PackageID, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 persistence.PackageID
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) SaveReturnsOnCall(i int, result1 persistence.PackageID, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	if fake.saveReturnsOnCall == nil {
		fake.saveReturnsOnCall = make(map[int]struct {
			result1 persistence.PackageID
			result2 error
		})
	}
	fake.saveReturnsOnCall[i] = struct {
		result1 persistence.PackageID
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeStore) recordInvocation(key string, args []interface{}) {
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
