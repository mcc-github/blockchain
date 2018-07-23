
package mock

import (
	"sync"
)

type PackageProvider struct {
	GetChaincodeCodePackageStub        func(ccname string, ccversion string) ([]byte, error)
	getChaincodeCodePackageMutex       sync.RWMutex
	getChaincodeCodePackageArgsForCall []struct {
		ccname    string
		ccversion string
	}
	getChaincodeCodePackageReturns struct {
		result1 []byte
		result2 error
	}
	getChaincodeCodePackageReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PackageProvider) GetChaincodeCodePackage(ccname string, ccversion string) ([]byte, error) {
	fake.getChaincodeCodePackageMutex.Lock()
	ret, specificReturn := fake.getChaincodeCodePackageReturnsOnCall[len(fake.getChaincodeCodePackageArgsForCall)]
	fake.getChaincodeCodePackageArgsForCall = append(fake.getChaincodeCodePackageArgsForCall, struct {
		ccname    string
		ccversion string
	}{ccname, ccversion})
	fake.recordInvocation("GetChaincodeCodePackage", []interface{}{ccname, ccversion})
	fake.getChaincodeCodePackageMutex.Unlock()
	if fake.GetChaincodeCodePackageStub != nil {
		return fake.GetChaincodeCodePackageStub(ccname, ccversion)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeCodePackageReturns.result1, fake.getChaincodeCodePackageReturns.result2
}

func (fake *PackageProvider) GetChaincodeCodePackageCallCount() int {
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
	return len(fake.getChaincodeCodePackageArgsForCall)
}

func (fake *PackageProvider) GetChaincodeCodePackageArgsForCall(i int) (string, string) {
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
	return fake.getChaincodeCodePackageArgsForCall[i].ccname, fake.getChaincodeCodePackageArgsForCall[i].ccversion
}

func (fake *PackageProvider) GetChaincodeCodePackageReturns(result1 []byte, result2 error) {
	fake.GetChaincodeCodePackageStub = nil
	fake.getChaincodeCodePackageReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *PackageProvider) GetChaincodeCodePackageReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *PackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodeCodePackageMutex.RLock()
	defer fake.getChaincodeCodePackageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PackageProvider) recordInvocation(key string, args []interface{}) {
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