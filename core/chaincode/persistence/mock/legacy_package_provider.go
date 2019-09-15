
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type LegacyPackageProvider struct {
	GetChaincodeInstallPathStub        func() string
	getChaincodeInstallPathMutex       sync.RWMutex
	getChaincodeInstallPathArgsForCall []struct {
	}
	getChaincodeInstallPathReturns struct {
		result1 string
	}
	getChaincodeInstallPathReturnsOnCall map[int]struct {
		result1 string
	}
	ListInstalledChaincodesStub        func(string, ccprovider.DirEnumerator, ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error)
	listInstalledChaincodesMutex       sync.RWMutex
	listInstalledChaincodesArgsForCall []struct {
		arg1 string
		arg2 ccprovider.DirEnumerator
		arg3 ccprovider.ChaincodeExtractor
	}
	listInstalledChaincodesReturns struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	listInstalledChaincodesReturnsOnCall map[int]struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LegacyPackageProvider) GetChaincodeInstallPath() string {
	fake.getChaincodeInstallPathMutex.Lock()
	ret, specificReturn := fake.getChaincodeInstallPathReturnsOnCall[len(fake.getChaincodeInstallPathArgsForCall)]
	fake.getChaincodeInstallPathArgsForCall = append(fake.getChaincodeInstallPathArgsForCall, struct {
	}{})
	fake.recordInvocation("GetChaincodeInstallPath", []interface{}{})
	fake.getChaincodeInstallPathMutex.Unlock()
	if fake.GetChaincodeInstallPathStub != nil {
		return fake.GetChaincodeInstallPathStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getChaincodeInstallPathReturns
	return fakeReturns.result1
}

func (fake *LegacyPackageProvider) GetChaincodeInstallPathCallCount() int {
	fake.getChaincodeInstallPathMutex.RLock()
	defer fake.getChaincodeInstallPathMutex.RUnlock()
	return len(fake.getChaincodeInstallPathArgsForCall)
}

func (fake *LegacyPackageProvider) GetChaincodeInstallPathCalls(stub func() string) {
	fake.getChaincodeInstallPathMutex.Lock()
	defer fake.getChaincodeInstallPathMutex.Unlock()
	fake.GetChaincodeInstallPathStub = stub
}

func (fake *LegacyPackageProvider) GetChaincodeInstallPathReturns(result1 string) {
	fake.getChaincodeInstallPathMutex.Lock()
	defer fake.getChaincodeInstallPathMutex.Unlock()
	fake.GetChaincodeInstallPathStub = nil
	fake.getChaincodeInstallPathReturns = struct {
		result1 string
	}{result1}
}

func (fake *LegacyPackageProvider) GetChaincodeInstallPathReturnsOnCall(i int, result1 string) {
	fake.getChaincodeInstallPathMutex.Lock()
	defer fake.getChaincodeInstallPathMutex.Unlock()
	fake.GetChaincodeInstallPathStub = nil
	if fake.getChaincodeInstallPathReturnsOnCall == nil {
		fake.getChaincodeInstallPathReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getChaincodeInstallPathReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *LegacyPackageProvider) ListInstalledChaincodes(arg1 string, arg2 ccprovider.DirEnumerator, arg3 ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error) {
	fake.listInstalledChaincodesMutex.Lock()
	ret, specificReturn := fake.listInstalledChaincodesReturnsOnCall[len(fake.listInstalledChaincodesArgsForCall)]
	fake.listInstalledChaincodesArgsForCall = append(fake.listInstalledChaincodesArgsForCall, struct {
		arg1 string
		arg2 ccprovider.DirEnumerator
		arg3 ccprovider.ChaincodeExtractor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ListInstalledChaincodes", []interface{}{arg1, arg2, arg3})
	fake.listInstalledChaincodesMutex.Unlock()
	if fake.ListInstalledChaincodesStub != nil {
		return fake.ListInstalledChaincodesStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listInstalledChaincodesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesCallCount() int {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return len(fake.listInstalledChaincodesArgsForCall)
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesCalls(stub func(string, ccprovider.DirEnumerator, ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error)) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = stub
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesArgsForCall(i int) (string, ccprovider.DirEnumerator, ccprovider.ChaincodeExtractor) {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	argsForCall := fake.listInstalledChaincodesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesReturns(result1 []chaincode.InstalledChaincode, result2 error) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = nil
	fake.listInstalledChaincodesReturns = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesReturnsOnCall(i int, result1 []chaincode.InstalledChaincode, result2 error) {
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

func (fake *LegacyPackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodeInstallPathMutex.RLock()
	defer fake.getChaincodeInstallPathMutex.RUnlock()
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LegacyPackageProvider) recordInvocation(key string, args []interface{}) {
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
