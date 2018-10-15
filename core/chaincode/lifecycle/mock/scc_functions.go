
package mock

import (
	"sync"
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

func (fake *SCCFunctions) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installChaincodeMutex.RLock()
	defer fake.installChaincodeMutex.RUnlock()
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
