
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
)

type SelfDescribingSysCC struct {
	ChaincodeStub        func() shim.Chaincode
	chaincodeMutex       sync.RWMutex
	chaincodeArgsForCall []struct {
	}
	chaincodeReturns struct {
		result1 shim.Chaincode
	}
	chaincodeReturnsOnCall map[int]struct {
		result1 shim.Chaincode
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SelfDescribingSysCC) Chaincode() shim.Chaincode {
	fake.chaincodeMutex.Lock()
	ret, specificReturn := fake.chaincodeReturnsOnCall[len(fake.chaincodeArgsForCall)]
	fake.chaincodeArgsForCall = append(fake.chaincodeArgsForCall, struct {
	}{})
	fake.recordInvocation("Chaincode", []interface{}{})
	fake.chaincodeMutex.Unlock()
	if fake.ChaincodeStub != nil {
		return fake.ChaincodeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.chaincodeReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) ChaincodeCallCount() int {
	fake.chaincodeMutex.RLock()
	defer fake.chaincodeMutex.RUnlock()
	return len(fake.chaincodeArgsForCall)
}

func (fake *SelfDescribingSysCC) ChaincodeCalls(stub func() shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = stub
}

func (fake *SelfDescribingSysCC) ChaincodeReturns(result1 shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = nil
	fake.chaincodeReturns = struct {
		result1 shim.Chaincode
	}{result1}
}

func (fake *SelfDescribingSysCC) ChaincodeReturnsOnCall(i int, result1 shim.Chaincode) {
	fake.chaincodeMutex.Lock()
	defer fake.chaincodeMutex.Unlock()
	fake.ChaincodeStub = nil
	if fake.chaincodeReturnsOnCall == nil {
		fake.chaincodeReturnsOnCall = make(map[int]struct {
			result1 shim.Chaincode
		})
	}
	fake.chaincodeReturnsOnCall[i] = struct {
		result1 shim.Chaincode
	}{result1}
}

func (fake *SelfDescribingSysCC) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *SelfDescribingSysCC) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *SelfDescribingSysCC) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *SelfDescribingSysCC) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *SelfDescribingSysCC) NameReturnsOnCall(i int, result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *SelfDescribingSysCC) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeMutex.RLock()
	defer fake.chaincodeMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SelfDescribingSysCC) recordInvocation(key string, args []interface{}) {
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
