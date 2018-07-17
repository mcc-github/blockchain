
package fake

import (
	"sync"

	chaincode_test "github.com/mcc-github/blockchain/core/chaincode"
)

type LaunchRegistry struct {
	LaunchingStub        func(cname string) (*chaincode_test.LaunchState, error)
	launchingMutex       sync.RWMutex
	launchingArgsForCall []struct {
		cname string
	}
	launchingReturns struct {
		result1 *chaincode_test.LaunchState
		result2 error
	}
	launchingReturnsOnCall map[int]struct {
		result1 *chaincode_test.LaunchState
		result2 error
	}
	DeregisterStub        func(cname string) error
	deregisterMutex       sync.RWMutex
	deregisterArgsForCall []struct {
		cname string
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

func (fake *LaunchRegistry) Launching(cname string) (*chaincode_test.LaunchState, error) {
	fake.launchingMutex.Lock()
	ret, specificReturn := fake.launchingReturnsOnCall[len(fake.launchingArgsForCall)]
	fake.launchingArgsForCall = append(fake.launchingArgsForCall, struct {
		cname string
	}{cname})
	fake.recordInvocation("Launching", []interface{}{cname})
	fake.launchingMutex.Unlock()
	if fake.LaunchingStub != nil {
		return fake.LaunchingStub(cname)
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

func (fake *LaunchRegistry) LaunchingArgsForCall(i int) string {
	fake.launchingMutex.RLock()
	defer fake.launchingMutex.RUnlock()
	return fake.launchingArgsForCall[i].cname
}

func (fake *LaunchRegistry) LaunchingReturns(result1 *chaincode_test.LaunchState, result2 error) {
	fake.LaunchingStub = nil
	fake.launchingReturns = struct {
		result1 *chaincode_test.LaunchState
		result2 error
	}{result1, result2}
}

func (fake *LaunchRegistry) LaunchingReturnsOnCall(i int, result1 *chaincode_test.LaunchState, result2 error) {
	fake.LaunchingStub = nil
	if fake.launchingReturnsOnCall == nil {
		fake.launchingReturnsOnCall = make(map[int]struct {
			result1 *chaincode_test.LaunchState
			result2 error
		})
	}
	fake.launchingReturnsOnCall[i] = struct {
		result1 *chaincode_test.LaunchState
		result2 error
	}{result1, result2}
}

func (fake *LaunchRegistry) Deregister(cname string) error {
	fake.deregisterMutex.Lock()
	ret, specificReturn := fake.deregisterReturnsOnCall[len(fake.deregisterArgsForCall)]
	fake.deregisterArgsForCall = append(fake.deregisterArgsForCall, struct {
		cname string
	}{cname})
	fake.recordInvocation("Deregister", []interface{}{cname})
	fake.deregisterMutex.Unlock()
	if fake.DeregisterStub != nil {
		return fake.DeregisterStub(cname)
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

func (fake *LaunchRegistry) DeregisterArgsForCall(i int) string {
	fake.deregisterMutex.RLock()
	defer fake.deregisterMutex.RUnlock()
	return fake.deregisterArgsForCall[i].cname
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
