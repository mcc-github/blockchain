
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type Runtime struct {
	StartStub        func(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error
	startMutex       sync.RWMutex
	startArgsForCall []struct {
		ccci        *ccprovider.ChaincodeContainerInfo
		codePackage []byte
	}
	startReturns struct {
		result1 error
	}
	startReturnsOnCall map[int]struct {
		result1 error
	}
	StopStub        func(ccci *ccprovider.ChaincodeContainerInfo) error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
		ccci *ccprovider.ChaincodeContainerInfo
	}
	stopReturns struct {
		result1 error
	}
	stopReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Runtime) Start(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error {
	var codePackageCopy []byte
	if codePackage != nil {
		codePackageCopy = make([]byte, len(codePackage))
		copy(codePackageCopy, codePackage)
	}
	fake.startMutex.Lock()
	ret, specificReturn := fake.startReturnsOnCall[len(fake.startArgsForCall)]
	fake.startArgsForCall = append(fake.startArgsForCall, struct {
		ccci        *ccprovider.ChaincodeContainerInfo
		codePackage []byte
	}{ccci, codePackageCopy})
	fake.recordInvocation("Start", []interface{}{ccci, codePackageCopy})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub(ccci, codePackage)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.startReturns.result1
}

func (fake *Runtime) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *Runtime) StartArgsForCall(i int) (*ccprovider.ChaincodeContainerInfo, []byte) {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return fake.startArgsForCall[i].ccci, fake.startArgsForCall[i].codePackage
}

func (fake *Runtime) StartReturns(result1 error) {
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *Runtime) StartReturnsOnCall(i int, result1 error) {
	fake.StartStub = nil
	if fake.startReturnsOnCall == nil {
		fake.startReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.startReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Runtime) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	fake.stopMutex.Lock()
	ret, specificReturn := fake.stopReturnsOnCall[len(fake.stopArgsForCall)]
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
		ccci *ccprovider.ChaincodeContainerInfo
	}{ccci})
	fake.recordInvocation("Stop", []interface{}{ccci})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub(ccci)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.stopReturns.result1
}

func (fake *Runtime) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *Runtime) StopArgsForCall(i int) *ccprovider.ChaincodeContainerInfo {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return fake.stopArgsForCall[i].ccci
}

func (fake *Runtime) StopReturns(result1 error) {
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

func (fake *Runtime) StopReturnsOnCall(i int, result1 error) {
	fake.StopStub = nil
	if fake.stopReturnsOnCall == nil {
		fake.stopReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.stopReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Runtime) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Runtime) recordInvocation(key string, args []interface{}) {
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
