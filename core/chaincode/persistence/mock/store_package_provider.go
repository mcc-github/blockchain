
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
)

type StorePackageProvider struct {
	GetChaincodeInstallPathStub        func() string
	getChaincodeInstallPathMutex       sync.RWMutex
	getChaincodeInstallPathArgsForCall []struct{}
	getChaincodeInstallPathReturns     struct {
		result1 string
	}
	getChaincodeInstallPathReturnsOnCall map[int]struct {
		result1 string
	}
	ListInstalledChaincodesStub        func() ([]chaincode.InstalledChaincode, error)
	listInstalledChaincodesMutex       sync.RWMutex
	listInstalledChaincodesArgsForCall []struct{}
	listInstalledChaincodesReturns     struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	listInstalledChaincodesReturnsOnCall map[int]struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}
	LoadStub        func(packageID persistence.PackageID) ([]byte, error)
	loadMutex       sync.RWMutex
	loadArgsForCall []struct {
		packageID persistence.PackageID
	}
	loadReturns struct {
		result1 []byte
		result2 error
	}
	loadReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StorePackageProvider) GetChaincodeInstallPath() string {
	fake.getChaincodeInstallPathMutex.Lock()
	ret, specificReturn := fake.getChaincodeInstallPathReturnsOnCall[len(fake.getChaincodeInstallPathArgsForCall)]
	fake.getChaincodeInstallPathArgsForCall = append(fake.getChaincodeInstallPathArgsForCall, struct{}{})
	fake.recordInvocation("GetChaincodeInstallPath", []interface{}{})
	fake.getChaincodeInstallPathMutex.Unlock()
	if fake.GetChaincodeInstallPathStub != nil {
		return fake.GetChaincodeInstallPathStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getChaincodeInstallPathReturns.result1
}

func (fake *StorePackageProvider) GetChaincodeInstallPathCallCount() int {
	fake.getChaincodeInstallPathMutex.RLock()
	defer fake.getChaincodeInstallPathMutex.RUnlock()
	return len(fake.getChaincodeInstallPathArgsForCall)
}

func (fake *StorePackageProvider) GetChaincodeInstallPathReturns(result1 string) {
	fake.GetChaincodeInstallPathStub = nil
	fake.getChaincodeInstallPathReturns = struct {
		result1 string
	}{result1}
}

func (fake *StorePackageProvider) GetChaincodeInstallPathReturnsOnCall(i int, result1 string) {
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

func (fake *StorePackageProvider) ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error) {
	fake.listInstalledChaincodesMutex.Lock()
	ret, specificReturn := fake.listInstalledChaincodesReturnsOnCall[len(fake.listInstalledChaincodesArgsForCall)]
	fake.listInstalledChaincodesArgsForCall = append(fake.listInstalledChaincodesArgsForCall, struct{}{})
	fake.recordInvocation("ListInstalledChaincodes", []interface{}{})
	fake.listInstalledChaincodesMutex.Unlock()
	if fake.ListInstalledChaincodesStub != nil {
		return fake.ListInstalledChaincodesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.listInstalledChaincodesReturns.result1, fake.listInstalledChaincodesReturns.result2
}

func (fake *StorePackageProvider) ListInstalledChaincodesCallCount() int {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return len(fake.listInstalledChaincodesArgsForCall)
}

func (fake *StorePackageProvider) ListInstalledChaincodesReturns(result1 []chaincode.InstalledChaincode, result2 error) {
	fake.ListInstalledChaincodesStub = nil
	fake.listInstalledChaincodesReturns = struct {
		result1 []chaincode.InstalledChaincode
		result2 error
	}{result1, result2}
}

func (fake *StorePackageProvider) ListInstalledChaincodesReturnsOnCall(i int, result1 []chaincode.InstalledChaincode, result2 error) {
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

func (fake *StorePackageProvider) Load(packageID persistence.PackageID) ([]byte, error) {
	fake.loadMutex.Lock()
	ret, specificReturn := fake.loadReturnsOnCall[len(fake.loadArgsForCall)]
	fake.loadArgsForCall = append(fake.loadArgsForCall, struct {
		packageID persistence.PackageID
	}{packageID})
	fake.recordInvocation("Load", []interface{}{packageID})
	fake.loadMutex.Unlock()
	if fake.LoadStub != nil {
		return fake.LoadStub(packageID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.loadReturns.result1, fake.loadReturns.result2
}

func (fake *StorePackageProvider) LoadCallCount() int {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return len(fake.loadArgsForCall)
}

func (fake *StorePackageProvider) LoadArgsForCall(i int) persistence.PackageID {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return fake.loadArgsForCall[i].packageID
}

func (fake *StorePackageProvider) LoadReturns(result1 []byte, result2 error) {
	fake.LoadStub = nil
	fake.loadReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *StorePackageProvider) LoadReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *StorePackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodeInstallPathMutex.RLock()
	defer fake.getChaincodeInstallPathMutex.RUnlock()
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StorePackageProvider) recordInvocation(key string, args []interface{}) {
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
