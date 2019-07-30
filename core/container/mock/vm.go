
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container"
)

type VM struct {
	BuildStub        func(*ccprovider.ChaincodeContainerInfo, []byte) (container.Instance, error)
	buildMutex       sync.RWMutex
	buildArgsForCall []struct {
		arg1 *ccprovider.ChaincodeContainerInfo
		arg2 []byte
	}
	buildReturns struct {
		result1 container.Instance
		result2 error
	}
	buildReturnsOnCall map[int]struct {
		result1 container.Instance
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *VM) Build(arg1 *ccprovider.ChaincodeContainerInfo, arg2 []byte) (container.Instance, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.buildMutex.Lock()
	ret, specificReturn := fake.buildReturnsOnCall[len(fake.buildArgsForCall)]
	fake.buildArgsForCall = append(fake.buildArgsForCall, struct {
		arg1 *ccprovider.ChaincodeContainerInfo
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("Build", []interface{}{arg1, arg2Copy})
	fake.buildMutex.Unlock()
	if fake.BuildStub != nil {
		return fake.BuildStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.buildReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VM) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *VM) BuildCalls(stub func(*ccprovider.ChaincodeContainerInfo, []byte) (container.Instance, error)) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = stub
}

func (fake *VM) BuildArgsForCall(i int) (*ccprovider.ChaincodeContainerInfo, []byte) {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	argsForCall := fake.buildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VM) BuildReturns(result1 container.Instance, result2 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 container.Instance
		result2 error
	}{result1, result2}
}

func (fake *VM) BuildReturnsOnCall(i int, result1 container.Instance, result2 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	if fake.buildReturnsOnCall == nil {
		fake.buildReturnsOnCall = make(map[int]struct {
			result1 container.Instance
			result2 error
		})
	}
	fake.buildReturnsOnCall[i] = struct {
		result1 container.Instance
		result2 error
	}{result1, result2}
}

func (fake *VM) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *VM) recordInvocation(key string, args []interface{}) {
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

var _ container.VM = new(VM)
