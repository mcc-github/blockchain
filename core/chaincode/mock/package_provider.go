
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type PackageProvider struct {
	GetChaincodeStub        func(ccname string, ccversion string) (ccprovider.CCPackage, error)
	getChaincodeMutex       sync.RWMutex
	getChaincodeArgsForCall []struct {
		ccname    string
		ccversion string
	}
	getChaincodeReturns struct {
		result1 ccprovider.CCPackage
		result2 error
	}
	getChaincodeReturnsOnCall map[int]struct {
		result1 ccprovider.CCPackage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PackageProvider) GetChaincode(ccname string, ccversion string) (ccprovider.CCPackage, error) {
	fake.getChaincodeMutex.Lock()
	ret, specificReturn := fake.getChaincodeReturnsOnCall[len(fake.getChaincodeArgsForCall)]
	fake.getChaincodeArgsForCall = append(fake.getChaincodeArgsForCall, struct {
		ccname    string
		ccversion string
	}{ccname, ccversion})
	fake.recordInvocation("GetChaincode", []interface{}{ccname, ccversion})
	fake.getChaincodeMutex.Unlock()
	if fake.GetChaincodeStub != nil {
		return fake.GetChaincodeStub(ccname, ccversion)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeReturns.result1, fake.getChaincodeReturns.result2
}

func (fake *PackageProvider) GetChaincodeCallCount() int {
	fake.getChaincodeMutex.RLock()
	defer fake.getChaincodeMutex.RUnlock()
	return len(fake.getChaincodeArgsForCall)
}

func (fake *PackageProvider) GetChaincodeArgsForCall(i int) (string, string) {
	fake.getChaincodeMutex.RLock()
	defer fake.getChaincodeMutex.RUnlock()
	return fake.getChaincodeArgsForCall[i].ccname, fake.getChaincodeArgsForCall[i].ccversion
}

func (fake *PackageProvider) GetChaincodeReturns(result1 ccprovider.CCPackage, result2 error) {
	fake.GetChaincodeStub = nil
	fake.getChaincodeReturns = struct {
		result1 ccprovider.CCPackage
		result2 error
	}{result1, result2}
}

func (fake *PackageProvider) GetChaincodeReturnsOnCall(i int, result1 ccprovider.CCPackage, result2 error) {
	fake.GetChaincodeStub = nil
	if fake.getChaincodeReturnsOnCall == nil {
		fake.getChaincodeReturnsOnCall = make(map[int]struct {
			result1 ccprovider.CCPackage
			result2 error
		})
	}
	fake.getChaincodeReturnsOnCall[i] = struct {
		result1 ccprovider.CCPackage
		result2 error
	}{result1, result2}
}

func (fake *PackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodeMutex.RLock()
	defer fake.getChaincodeMutex.RUnlock()
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
