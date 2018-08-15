
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type LegacyPackageProvider struct {
	GetChaincodeCodePackageStub        func(name, version string) (codePackage []byte, err error)
	getChaincodeCodePackageMutex       sync.RWMutex
	getChaincodeCodePackageArgsForCall []struct {
		name    string
		version string
	}
	getChaincodeCodePackageReturns struct {
		result1 []byte
		result2 error
	}
	getChaincodeCodePackageReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	ListInstalledChaincodesStub        func(dir string, de ccprovider.DirEnumerator, ce ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error)
	listInstalledChaincodesMutex       sync.RWMutex
	listInstalledChaincodesArgsForCall []struct {
		dir string
		de  ccprovider.DirEnumerator
		ce  ccprovider.ChaincodeExtractor
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

func (fake *LegacyPackageProvider) GetChaincodeCodePackage(name string, version string) (codePackage []byte, err error) {
	fake.getChaincodeCodePackageMutex.Lock()
	ret, specificReturn := fake.getChaincodeCodePackageReturnsOnCall[len(fake.getChaincodeCodePackageArgsForCall)]
	fake.getChaincodeCodePackageArgsForCall = append(fake.getChaincodeCodePackageArgsForCall, struct {
		name    string
		version string
	}{name, version})
	fake.recordInvocation("GetChaincodeCodePackage", []interface{}{name, version})
	fake.getChaincodeCodePackageMutex.Unlock()
	if fake.GetChaincodeCodePackageStub != nil {
		return fake.GetChaincodeCodePackageStub(name, version)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeCodePackageReturns.result1, fake.getChaincodeCodePackageReturns.result2
}

func (fake *LegacyPackageProvider) GetChaincodeCodePackageCallCount() int {
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
	return len(fake.getChaincodeCodePackageArgsForCall)
}

func (fake *LegacyPackageProvider) GetChaincodeCodePackageArgsForCall(i int) (string, string) {
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
	return fake.getChaincodeCodePackageArgsForCall[i].name, fake.getChaincodeCodePackageArgsForCall[i].version
}

func (fake *LegacyPackageProvider) GetChaincodeCodePackageReturns(result1 []byte, result2 error) {
	fake.GetChaincodeCodePackageStub = nil
	fake.getChaincodeCodePackageReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *LegacyPackageProvider) GetChaincodeCodePackageReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetChaincodeCodePackageStub = nil
	if fake.getChaincodeCodePackageReturnsOnCall == nil {
		fake.getChaincodeCodePackageReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getChaincodeCodePackageReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *LegacyPackageProvider) ListInstalledChaincodes(dir string, de ccprovider.DirEnumerator, ce ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error) {
	fake.listInstalledChaincodesMutex.Lock()
	ret, specificReturn := fake.listInstalledChaincodesReturnsOnCall[len(fake.listInstalledChaincodesArgsForCall)]
	fake.listInstalledChaincodesArgsForCall = append(fake.listInstalledChaincodesArgsForCall, struct {
		dir string
		de  ccprovider.DirEnumerator
		ce  ccprovider.ChaincodeExtractor
	}{dir, de, ce})
	fake.recordInvocation("ListInstalledChaincodes", []interface{}{dir, de, ce})
	fake.listInstalledChaincodesMutex.Unlock()
	if fake.ListInstalledChaincodesStub != nil {
		return fake.ListInstalledChaincodesStub(dir, de, ce)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.listInstalledChaincodesReturns.result1, fake.listInstalledChaincodesReturns.result2
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesCallCount() int {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return len(fake.listInstalledChaincodesArgsForCall)
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesArgsForCall(i int) (string, ccprovider.DirEnumerator, ccprovider.ChaincodeExtractor) {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return fake.listInstalledChaincodesArgsForCall[i].dir, fake.listInstalledChaincodesArgsForCall[i].de, fake.listInstalledChaincodesArgsForCall[i].ce
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesReturns(result1 []chaincode.InstalledChaincode, result2 error) {
	fake.ListInstalledChaincodesStub = nil
	fake.listInstalledChaincodesReturns = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *LegacyPackageProvider) ListInstalledChaincodesReturnsOnCall(i int, result1 []chaincode.InstalledChaincode, result2 error) {
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
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
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
