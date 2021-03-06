
package mock

import (
	"io"
	"sync"

	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/container"
)

type PackageProvider struct {
	GetChaincodePackageStub        func(string) (*persistence.ChaincodePackageMetadata, io.ReadCloser, error)
	getChaincodePackageMutex       sync.RWMutex
	getChaincodePackageArgsForCall []struct {
		arg1 string
	}
	getChaincodePackageReturns struct {
		result1 *persistence.ChaincodePackageMetadata
		result2 io.ReadCloser
		result3 error
	}
	getChaincodePackageReturnsOnCall map[int]struct {
		result1 *persistence.ChaincodePackageMetadata
		result2 io.ReadCloser
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PackageProvider) GetChaincodePackage(arg1 string) (*persistence.ChaincodePackageMetadata, io.ReadCloser, error) {
	fake.getChaincodePackageMutex.Lock()
	ret, specificReturn := fake.getChaincodePackageReturnsOnCall[len(fake.getChaincodePackageArgsForCall)]
	fake.getChaincodePackageArgsForCall = append(fake.getChaincodePackageArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetChaincodePackage", []interface{}{arg1})
	fake.getChaincodePackageMutex.Unlock()
	if fake.GetChaincodePackageStub != nil {
		return fake.GetChaincodePackageStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getChaincodePackageReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *PackageProvider) GetChaincodePackageCallCount() int {
	fake.getChaincodePackageMutex.RLock()
	defer fake.getChaincodePackageMutex.RUnlock()
	return len(fake.getChaincodePackageArgsForCall)
}

func (fake *PackageProvider) GetChaincodePackageCalls(stub func(string) (*persistence.ChaincodePackageMetadata, io.ReadCloser, error)) {
	fake.getChaincodePackageMutex.Lock()
	defer fake.getChaincodePackageMutex.Unlock()
	fake.GetChaincodePackageStub = stub
}

func (fake *PackageProvider) GetChaincodePackageArgsForCall(i int) string {
	fake.getChaincodePackageMutex.RLock()
	defer fake.getChaincodePackageMutex.RUnlock()
	argsForCall := fake.getChaincodePackageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PackageProvider) GetChaincodePackageReturns(result1 *persistence.ChaincodePackageMetadata, result2 io.ReadCloser, result3 error) {
	fake.getChaincodePackageMutex.Lock()
	defer fake.getChaincodePackageMutex.Unlock()
	fake.GetChaincodePackageStub = nil
	fake.getChaincodePackageReturns = struct {
		result1 *persistence.ChaincodePackageMetadata
		result2 io.ReadCloser
		result3 error
	}{result1, result2, result3}
}

func (fake *PackageProvider) GetChaincodePackageReturnsOnCall(i int, result1 *persistence.ChaincodePackageMetadata, result2 io.ReadCloser, result3 error) {
	fake.getChaincodePackageMutex.Lock()
	defer fake.getChaincodePackageMutex.Unlock()
	fake.GetChaincodePackageStub = nil
	if fake.getChaincodePackageReturnsOnCall == nil {
		fake.getChaincodePackageReturnsOnCall = make(map[int]struct {
			result1 *persistence.ChaincodePackageMetadata
			result2 io.ReadCloser
			result3 error
		})
	}
	fake.getChaincodePackageReturnsOnCall[i] = struct {
		result1 *persistence.ChaincodePackageMetadata
		result2 io.ReadCloser
		result3 error
	}{result1, result2, result3}
}

func (fake *PackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodePackageMutex.RLock()
	defer fake.getChaincodePackageMutex.RUnlock()
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

var _ container.PackageProvider = new(PackageProvider)
