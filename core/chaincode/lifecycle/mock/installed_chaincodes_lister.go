
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
)

type InstalledChaincodesLister struct {
	ListInstalledChaincodesStub        func() []*chaincode.InstalledChaincode
	listInstalledChaincodesMutex       sync.RWMutex
	listInstalledChaincodesArgsForCall []struct {
	}
	listInstalledChaincodesReturns struct {
		result1 []*chaincode.InstalledChaincode
	}
	listInstalledChaincodesReturnsOnCall map[int]struct {
		result1 []*chaincode.InstalledChaincode
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *InstalledChaincodesLister) ListInstalledChaincodes() []*chaincode.InstalledChaincode {
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
		return ret.result1
	}
	fakeReturns := fake.listInstalledChaincodesReturns
	return fakeReturns.result1
}

func (fake *InstalledChaincodesLister) ListInstalledChaincodesCallCount() int {
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	return len(fake.listInstalledChaincodesArgsForCall)
}

func (fake *InstalledChaincodesLister) ListInstalledChaincodesCalls(stub func() []*chaincode.InstalledChaincode) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = stub
}

func (fake *InstalledChaincodesLister) ListInstalledChaincodesReturns(result1 []*chaincode.InstalledChaincode) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = nil
	fake.listInstalledChaincodesReturns = struct {
		result1 []*chaincode.InstalledChaincode
	}{result1}
}

func (fake *InstalledChaincodesLister) ListInstalledChaincodesReturnsOnCall(i int, result1 []*chaincode.InstalledChaincode) {
	fake.listInstalledChaincodesMutex.Lock()
	defer fake.listInstalledChaincodesMutex.Unlock()
	fake.ListInstalledChaincodesStub = nil
	if fake.listInstalledChaincodesReturnsOnCall == nil {
		fake.listInstalledChaincodesReturnsOnCall = make(map[int]struct {
			result1 []*chaincode.InstalledChaincode
		})
	}
	fake.listInstalledChaincodesReturnsOnCall[i] = struct {
		result1 []*chaincode.InstalledChaincode
	}{result1}
}

func (fake *InstalledChaincodesLister) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listInstalledChaincodesMutex.RLock()
	defer fake.listInstalledChaincodesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *InstalledChaincodesLister) recordInvocation(key string, args []interface{}) {
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

var _ lifecycle.InstalledChaincodesLister = new(InstalledChaincodesLister)
