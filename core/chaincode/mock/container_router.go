
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
)

type ContainerRouter struct {
	BuildStub        func(*ccprovider.ChaincodeContainerInfo, []byte) error
	buildMutex       sync.RWMutex
	buildArgsForCall []struct {
		arg1 *ccprovider.ChaincodeContainerInfo
		arg2 []byte
	}
	buildReturns struct {
		result1 error
	}
	buildReturnsOnCall map[int]struct {
		result1 error
	}
	StartStub        func(ccintf.CCID, *ccintf.PeerConnection) error
	startMutex       sync.RWMutex
	startArgsForCall []struct {
		arg1 ccintf.CCID
		arg2 *ccintf.PeerConnection
	}
	startReturns struct {
		result1 error
	}
	startReturnsOnCall map[int]struct {
		result1 error
	}
	StopStub        func(ccintf.CCID) error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
		arg1 ccintf.CCID
	}
	stopReturns struct {
		result1 error
	}
	stopReturnsOnCall map[int]struct {
		result1 error
	}
	WaitStub        func(ccintf.CCID) (int, error)
	waitMutex       sync.RWMutex
	waitArgsForCall []struct {
		arg1 ccintf.CCID
	}
	waitReturns struct {
		result1 int
		result2 error
	}
	waitReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ContainerRouter) Build(arg1 *ccprovider.ChaincodeContainerInfo, arg2 []byte) error {
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
		return ret.result1
	}
	fakeReturns := fake.buildReturns
	return fakeReturns.result1
}

func (fake *ContainerRouter) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *ContainerRouter) BuildCalls(stub func(*ccprovider.ChaincodeContainerInfo, []byte) error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = stub
}

func (fake *ContainerRouter) BuildArgsForCall(i int) (*ccprovider.ChaincodeContainerInfo, []byte) {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	argsForCall := fake.buildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ContainerRouter) BuildReturns(result1 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 error
	}{result1}
}

func (fake *ContainerRouter) BuildReturnsOnCall(i int, result1 error) {
	fake.buildMutex.Lock()
	defer fake.buildMutex.Unlock()
	fake.BuildStub = nil
	if fake.buildReturnsOnCall == nil {
		fake.buildReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.buildReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ContainerRouter) Start(arg1 ccintf.CCID, arg2 *ccintf.PeerConnection) error {
	fake.startMutex.Lock()
	ret, specificReturn := fake.startReturnsOnCall[len(fake.startArgsForCall)]
	fake.startArgsForCall = append(fake.startArgsForCall, struct {
		arg1 ccintf.CCID
		arg2 *ccintf.PeerConnection
	}{arg1, arg2})
	fake.recordInvocation("Start", []interface{}{arg1, arg2})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.startReturns
	return fakeReturns.result1
}

func (fake *ContainerRouter) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *ContainerRouter) StartCalls(stub func(ccintf.CCID, *ccintf.PeerConnection) error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = stub
}

func (fake *ContainerRouter) StartArgsForCall(i int) (ccintf.CCID, *ccintf.PeerConnection) {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	argsForCall := fake.startArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ContainerRouter) StartReturns(result1 error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *ContainerRouter) StartReturnsOnCall(i int, result1 error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
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

func (fake *ContainerRouter) Stop(arg1 ccintf.CCID) error {
	fake.stopMutex.Lock()
	ret, specificReturn := fake.stopReturnsOnCall[len(fake.stopArgsForCall)]
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
		arg1 ccintf.CCID
	}{arg1})
	fake.recordInvocation("Stop", []interface{}{arg1})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.stopReturns
	return fakeReturns.result1
}

func (fake *ContainerRouter) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *ContainerRouter) StopCalls(stub func(ccintf.CCID) error) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = stub
}

func (fake *ContainerRouter) StopArgsForCall(i int) ccintf.CCID {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	argsForCall := fake.stopArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ContainerRouter) StopReturns(result1 error) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

func (fake *ContainerRouter) StopReturnsOnCall(i int, result1 error) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
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

func (fake *ContainerRouter) Wait(arg1 ccintf.CCID) (int, error) {
	fake.waitMutex.Lock()
	ret, specificReturn := fake.waitReturnsOnCall[len(fake.waitArgsForCall)]
	fake.waitArgsForCall = append(fake.waitArgsForCall, struct {
		arg1 ccintf.CCID
	}{arg1})
	fake.recordInvocation("Wait", []interface{}{arg1})
	fake.waitMutex.Unlock()
	if fake.WaitStub != nil {
		return fake.WaitStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.waitReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ContainerRouter) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *ContainerRouter) WaitCalls(stub func(ccintf.CCID) (int, error)) {
	fake.waitMutex.Lock()
	defer fake.waitMutex.Unlock()
	fake.WaitStub = stub
}

func (fake *ContainerRouter) WaitArgsForCall(i int) ccintf.CCID {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	argsForCall := fake.waitArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ContainerRouter) WaitReturns(result1 int, result2 error) {
	fake.waitMutex.Lock()
	defer fake.waitMutex.Unlock()
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *ContainerRouter) WaitReturnsOnCall(i int, result1 int, result2 error) {
	fake.waitMutex.Lock()
	defer fake.waitMutex.Unlock()
	fake.WaitStub = nil
	if fake.waitReturnsOnCall == nil {
		fake.waitReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.waitReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *ContainerRouter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ContainerRouter) recordInvocation(key string, args []interface{}) {
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
