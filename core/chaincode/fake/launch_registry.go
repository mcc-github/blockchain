
package fake

import (
	"sync"

	chaincode_test "github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/container/ccintf"
)

type LaunchRegistry struct {
	LaunchingStub        func(packageID ccintf.CCID) (launchState *chaincode_test.LaunchState, started bool)
	launchingMutex       sync.RWMutex
	launchingArgsForCall []struct {
		packageID ccintf.CCID
	}
	launchingReturns struct {
		result1 *chaincode_test.LaunchState
		result2 bool
	}
	launchingReturnsOnCall map[int]struct {
		result1 *chaincode_test.LaunchState
		result2 bool
	}
	DeregisterStub        func(packageID ccintf.CCID) error
	deregisterMutex       sync.RWMutex
	deregisterArgsForCall []struct {
		packageID ccintf.CCID
	}
	deregisterReturns struct {
		result1 error
	}
	deregisterReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LaunchRegistry) Launching(packageID ccintf.CCID) (launchState *chaincode_test.LaunchState, started bool) {
	fake.launchingMutex.Lock()
	ret, specificReturn := fake.launchingReturnsOnCall[len(fake.launchingArgsForCall)]
	fake.launchingArgsForCall = append(fake.launchingArgsForCall, struct {
		packageID ccintf.CCID
	}{packageID})
	fake.recordInvocation("Launching", []interface{}{packageID})
	fake.launchingMutex.Unlock()
	if fake.LaunchingStub != nil {
		return fake.LaunchingStub(packageID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.launchingReturns.result1, fake.launchingReturns.result2
}

func (fake *LaunchRegistry) LaunchingCallCount() int {
	fake.launchingMutex.RLock()
	defer fake.launchingMutex.RUnlock()
	return len(fake.launchingArgsForCall)
}

func (fake *LaunchRegistry) LaunchingArgsForCall(i int) ccintf.CCID {
	fake.launchingMutex.RLock()
	defer fake.launchingMutex.RUnlock()
	return fake.launchingArgsForCall[i].packageID
}

func (fake *LaunchRegistry) LaunchingReturns(result1 *chaincode_test.LaunchState, result2 bool) {
	fake.LaunchingStub = nil
	fake.launchingReturns = struct {
		result1 *chaincode_test.LaunchState
		result2 bool
	}{result1, result2}
}

func (fake *LaunchRegistry) LaunchingReturnsOnCall(i int, result1 *chaincode_test.LaunchState, result2 bool) {
	fake.LaunchingStub = nil
	if fake.launchingReturnsOnCall == nil {
		fake.launchingReturnsOnCall = make(map[int]struct {
			result1 *chaincode_test.LaunchState
			result2 bool
		})
	}
	fake.launchingReturnsOnCall[i] = struct {
		result1 *chaincode_test.LaunchState
		result2 bool
	}{result1, result2}
}

func (fake *LaunchRegistry) Deregister(packageID ccintf.CCID) error {
	fake.deregisterMutex.Lock()
	ret, specificReturn := fake.deregisterReturnsOnCall[len(fake.deregisterArgsForCall)]
	fake.deregisterArgsForCall = append(fake.deregisterArgsForCall, struct {
		packageID ccintf.CCID
	}{packageID})
	fake.recordInvocation("Deregister", []interface{}{packageID})
	fake.deregisterMutex.Unlock()
	if fake.DeregisterStub != nil {
		return fake.DeregisterStub(packageID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deregisterReturns.result1
}

func (fake *LaunchRegistry) DeregisterCallCount() int {
	fake.deregisterMutex.RLock()
	defer fake.deregisterMutex.RUnlock()
	return len(fake.deregisterArgsForCall)
}

func (fake *LaunchRegistry) DeregisterArgsForCall(i int) ccintf.CCID {
	fake.deregisterMutex.RLock()
	defer fake.deregisterMutex.RUnlock()
	return fake.deregisterArgsForCall[i].packageID
}

func (fake *LaunchRegistry) DeregisterReturns(result1 error) {
	fake.DeregisterStub = nil
	fake.deregisterReturns = struct {
		result1 error
	}{result1}
}

func (fake *LaunchRegistry) DeregisterReturnsOnCall(i int, result1 error) {
	fake.DeregisterStub = nil
	if fake.deregisterReturnsOnCall == nil {
		fake.deregisterReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deregisterReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *LaunchRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.launchingMutex.RLock()
	defer fake.launchingMutex.RUnlock()
	fake.deregisterMutex.RLock()
	defer fake.deregisterMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LaunchRegistry) recordInvocation(key string, args []interface{}) {
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
