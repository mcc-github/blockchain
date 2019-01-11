
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
)

type SCCFunctions struct {
	InstallChaincodeStub        func(name, version string, chaincodePackage []byte) (hash []byte, err error)
	installChaincodeMutex       sync.RWMutex
	installChaincodeArgsForCall []struct {
		name             string
		version          string
		chaincodePackage []byte
	}
	installChaincodeReturns struct {
		result1 []byte
		result2 error
	}
	installChaincodeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	QueryInstalledChaincodeStub        func(name, version string) (hash []byte, err error)
	queryInstalledChaincodeMutex       sync.RWMutex
	queryInstalledChaincodeArgsForCall []struct {
		name    string
		version string
	}
	queryInstalledChaincodeReturns struct {
		result1 []byte
		result2 error
	}
	queryInstalledChaincodeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	QueryInstalledChaincodesStub        func() (chaincodes []chaincode.InstalledChaincode, err error)
	queryInstalledChaincodesMutex       sync.RWMutex
	queryInstalledChaincodesArgsForCall []struct{}
	queryInstalledChaincodesReturns     struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	queryInstalledChaincodesReturnsOnCall map[int]struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SCCFunctions) InstallChaincode(name string, version string, chaincodePackage []byte) (hash []byte, err error) {
	var chaincodePackageCopy []byte
	if chaincodePackage != nil {
		chaincodePackageCopy = make([]byte, len(chaincodePackage))
		copy(chaincodePackageCopy, chaincodePackage)
	}
	fake.installChaincodeMutex.Lock()
	ret, specificReturn := fake.installChaincodeReturnsOnCall[len(fake.installChaincodeArgsForCall)]
	fake.installChaincodeArgsForCall = append(fake.installChaincodeArgsForCall, struct {
		name             string
		version          string
		chaincodePackage []byte
	}{name, version, chaincodePackageCopy})
	fake.recordInvocation("InstallChaincode", []interface{}{name, version, chaincodePackageCopy})
	fake.installChaincodeMutex.Unlock()
	if fake.InstallChaincodeStub != nil {
		return fake.InstallChaincodeStub(name, version, chaincodePackage)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.installChaincodeReturns.result1, fake.installChaincodeReturns.result2
}

func (fake *SCCFunctions) InstallChaincodeCallCount() int {
	fake.installChaincodeMutex.RLock()
	defer fake.installChaincodeMutex.RUnlock()
	return len(fake.installChaincodeArgsForCall)
}

func (fake *SCCFunctions) InstallChaincodeArgsForCall(i int) (string, string, []byte) {
	fake.installChaincodeMutex.RLock()
	defer fake.installChaincodeMutex.RUnlock()
	return fake.installChaincodeArgsForCall[i].name, fake.installChaincodeArgsForCall[i].version, fake.installChaincodeArgsForCall[i].chaincodePackage
}

func (fake *SCCFunctions) InstallChaincodeReturns(result1 []byte, result2 error) {
	fake.InstallChaincodeStub = nil
	fake.installChaincodeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) InstallChaincodeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.InstallChaincodeStub = nil
	if fake.installChaincodeReturnsOnCall == nil {
		fake.installChaincodeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.installChaincodeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) QueryInstalledChaincode(name string, version string) (hash []byte, err error) {
	fake.queryInstalledChaincodeMutex.Lock()
	ret, specificReturn := fake.queryInstalledChaincodeReturnsOnCall[len(fake.queryInstalledChaincodeArgsForCall)]
	fake.queryInstalledChaincodeArgsForCall = append(fake.queryInstalledChaincodeArgsForCall, struct {
		name    string
		version string
	}{name, version})
	fake.recordInvocation("QueryInstalledChaincode", []interface{}{name, version})
	fake.queryInstalledChaincodeMutex.Unlock()
	if fake.QueryInstalledChaincodeStub != nil {
		return fake.QueryInstalledChaincodeStub(name, version)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.queryInstalledChaincodeReturns.result1, fake.queryInstalledChaincodeReturns.result2
}

func (fake *SCCFunctions) QueryInstalledChaincodeCallCount() int {
	fake.queryInstalledChaincodeMutex.RLock()
	defer fake.queryInstalledChaincodeMutex.RUnlock()
	return len(fake.queryInstalledChaincodeArgsForCall)
}

func (fake *SCCFunctions) QueryInstalledChaincodeArgsForCall(i int) (string, string) {
	fake.queryInstalledChaincodeMutex.RLock()
	defer fake.queryInstalledChaincodeMutex.RUnlock()
	return fake.queryInstalledChaincodeArgsForCall[i].name, fake.queryInstalledChaincodeArgsForCall[i].version
}

func (fake *SCCFunctions) QueryInstalledChaincodeReturns(result1 []byte, result2 error) {
	fake.QueryInstalledChaincodeStub = nil
	fake.queryInstalledChaincodeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) QueryInstalledChaincodeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.QueryInstalledChaincodeStub = nil
	if fake.queryInstalledChaincodeReturnsOnCall == nil {
		fake.queryInstalledChaincodeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.queryInstalledChaincodeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) QueryInstalledChaincodes() (chaincodes []chaincode.InstalledChaincode, err error) {
	fake.queryInstalledChaincodesMutex.Lock()
	ret, specificReturn := fake.queryInstalledChaincodesReturnsOnCall[len(fake.queryInstalledChaincodesArgsForCall)]
	fake.queryInstalledChaincodesArgsForCall = append(fake.queryInstalledChaincodesArgsForCall, struct{}{})
	fake.recordInvocation("QueryInstalledChaincodes", []interface{}{})
	fake.queryInstalledChaincodesMutex.Unlock()
	if fake.QueryInstalledChaincodesStub != nil {
		return fake.QueryInstalledChaincodesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.queryInstalledChaincodesReturns.result1, fake.queryInstalledChaincodesReturns.result2
}

func (fake *SCCFunctions) QueryInstalledChaincodesCallCount() int {
	fake.queryInstalledChaincodesMutex.RLock()
	defer fake.queryInstalledChaincodesMutex.RUnlock()
	return len(fake.queryInstalledChaincodesArgsForCall)
}

func (fake *SCCFunctions) QueryInstalledChaincodesReturns(result1 []chaincode.InstalledChaincode, result2 error) {
	fake.QueryInstalledChaincodesStub = nil
	fake.queryInstalledChaincodesReturns = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) QueryInstalledChaincodesReturnsOnCall(i int, result1 []chaincode.InstalledChaincode, result2 error) {
	fake.QueryInstalledChaincodesStub = nil
	if fake.queryInstalledChaincodesReturnsOnCall == nil {
		fake.queryInstalledChaincodesReturnsOnCall = make(map[int]struct {
			result1 []chaincode.InstalledChaincode
			result2 error
		})
	}
	fake.queryInstalledChaincodesReturnsOnCall[i] = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *SCCFunctions) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installChaincodeMutex.RLock()
	defer fake.installChaincodeMutex.RUnlock()
	fake.queryInstalledChaincodeMutex.RLock()
	defer fake.queryInstalledChaincodeMutex.RUnlock()
	fake.queryInstalledChaincodesMutex.RLock()
	defer fake.queryInstalledChaincodesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SCCFunctions) recordInvocation(key string, args []interface{}) {
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
