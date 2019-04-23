
package mock

import (
	"context"
	"sync"

	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
)

type VM struct {
	HealthCheckStub        func(context.Context) error
	healthCheckMutex       sync.RWMutex
	healthCheckArgsForCall []struct {
		arg1 context.Context
	}
	healthCheckReturns struct {
		result1 error
	}
	healthCheckReturnsOnCall map[int]struct {
		result1 error
	}
	StartStub        func(ccintf.CCID, []string, []string, map[string][]byte, container.Builder) error
	startMutex       sync.RWMutex
	startArgsForCall []struct {
		arg1 ccintf.CCID
		arg2 []string
		arg3 []string
		arg4 map[string][]byte
		arg5 container.Builder
	}
	startReturns struct {
		result1 error
	}
	startReturnsOnCall map[int]struct {
		result1 error
	}
	StopStub        func(ccintf.CCID, uint, bool, bool) error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
		arg1 ccintf.CCID
		arg2 uint
		arg3 bool
		arg4 bool
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

func (fake *VM) HealthCheck(arg1 context.Context) error {
	fake.healthCheckMutex.Lock()
	ret, specificReturn := fake.healthCheckReturnsOnCall[len(fake.healthCheckArgsForCall)]
	fake.healthCheckArgsForCall = append(fake.healthCheckArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("HealthCheck", []interface{}{arg1})
	fake.healthCheckMutex.Unlock()
	if fake.HealthCheckStub != nil {
		return fake.HealthCheckStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.healthCheckReturns
	return fakeReturns.result1
}

func (fake *VM) HealthCheckCallCount() int {
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
	return len(fake.healthCheckArgsForCall)
}

func (fake *VM) HealthCheckCalls(stub func(context.Context) error) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = stub
}

func (fake *VM) HealthCheckArgsForCall(i int) context.Context {
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
	argsForCall := fake.healthCheckArgsForCall[i]
	return argsForCall.arg1
}

func (fake *VM) HealthCheckReturns(result1 error) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = nil
	fake.healthCheckReturns = struct {
		result1 error
	}{result1}
}

func (fake *VM) HealthCheckReturnsOnCall(i int, result1 error) {
	fake.healthCheckMutex.Lock()
	defer fake.healthCheckMutex.Unlock()
	fake.HealthCheckStub = nil
	if fake.healthCheckReturnsOnCall == nil {
		fake.healthCheckReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.healthCheckReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *VM) Start(arg1 ccintf.CCID, arg2 []string, arg3 []string, arg4 map[string][]byte, arg5 container.Builder) error {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	var arg3Copy []string
	if arg3 != nil {
		arg3Copy = make([]string, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.startMutex.Lock()
	ret, specificReturn := fake.startReturnsOnCall[len(fake.startArgsForCall)]
	fake.startArgsForCall = append(fake.startArgsForCall, struct {
		arg1 ccintf.CCID
		arg2 []string
		arg3 []string
		arg4 map[string][]byte
		arg5 container.Builder
	}{arg1, arg2Copy, arg3Copy, arg4, arg5})
	fake.recordInvocation("Start", []interface{}{arg1, arg2Copy, arg3Copy, arg4, arg5})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.startReturns
	return fakeReturns.result1
}

func (fake *VM) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *VM) StartCalls(stub func(ccintf.CCID, []string, []string, map[string][]byte, container.Builder) error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = stub
}

func (fake *VM) StartArgsForCall(i int) (ccintf.CCID, []string, []string, map[string][]byte, container.Builder) {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	argsForCall := fake.startArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *VM) StartReturns(result1 error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *VM) StartReturnsOnCall(i int, result1 error) {
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

func (fake *VM) Stop(arg1 ccintf.CCID, arg2 uint, arg3 bool, arg4 bool) error {
	fake.stopMutex.Lock()
	ret, specificReturn := fake.stopReturnsOnCall[len(fake.stopArgsForCall)]
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
		arg1 ccintf.CCID
		arg2 uint
		arg3 bool
		arg4 bool
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("Stop", []interface{}{arg1, arg2, arg3, arg4})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.stopReturns
	return fakeReturns.result1
}

func (fake *VM) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *VM) StopCalls(stub func(ccintf.CCID, uint, bool, bool) error) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = stub
}

func (fake *VM) StopArgsForCall(i int) (ccintf.CCID, uint, bool, bool) {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	argsForCall := fake.stopArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *VM) StopReturns(result1 error) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

func (fake *VM) StopReturnsOnCall(i int, result1 error) {
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

func (fake *VM) Wait(arg1 ccintf.CCID) (int, error) {
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

func (fake *VM) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *VM) WaitCalls(stub func(ccintf.CCID) (int, error)) {
	fake.waitMutex.Lock()
	defer fake.waitMutex.Unlock()
	fake.WaitStub = stub
}

func (fake *VM) WaitArgsForCall(i int) ccintf.CCID {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	argsForCall := fake.waitArgsForCall[i]
	return argsForCall.arg1
}

func (fake *VM) WaitReturns(result1 int, result2 error) {
	fake.waitMutex.Lock()
	defer fake.waitMutex.Unlock()
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *VM) WaitReturnsOnCall(i int, result1 int, result2 error) {
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

func (fake *VM) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.healthCheckMutex.RLock()
	defer fake.healthCheckMutex.RUnlock()
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
