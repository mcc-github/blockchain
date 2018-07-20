
package mock

import (
	"sync"

	container_test "github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
)

type VM struct {
	StartStub        func(ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte, builder container_test.Builder) error
	startMutex       sync.RWMutex
	startArgsForCall []struct {
		ccid          ccintf.CCID
		args          []string
		env           []string
		filesToUpload map[string][]byte
		builder       container_test.Builder
	}
	startReturns struct {
		result1 error
	}
	startReturnsOnCall map[int]struct {
		result1 error
	}
	StopStub        func(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
		ccid       ccintf.CCID
		timeout    uint
		dontkill   bool
		dontremove bool
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

func (fake *VM) Start(ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte, builder container_test.Builder) error {
	var argsCopy []string
	if args != nil {
		argsCopy = make([]string, len(args))
		copy(argsCopy, args)
	}
	var envCopy []string
	if env != nil {
		envCopy = make([]string, len(env))
		copy(envCopy, env)
	}
	fake.startMutex.Lock()
	ret, specificReturn := fake.startReturnsOnCall[len(fake.startArgsForCall)]
	fake.startArgsForCall = append(fake.startArgsForCall, struct {
		ccid          ccintf.CCID
		args          []string
		env           []string
		filesToUpload map[string][]byte
		builder       container_test.Builder
	}{ccid, argsCopy, envCopy, filesToUpload, builder})
	fake.recordInvocation("Start", []interface{}{ccid, argsCopy, envCopy, filesToUpload, builder})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub(ccid, args, env, filesToUpload, builder)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.startReturns.result1
}

func (fake *VM) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *VM) StartArgsForCall(i int) (ccintf.CCID, []string, []string, map[string][]byte, container_test.Builder) {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return fake.startArgsForCall[i].ccid, fake.startArgsForCall[i].args, fake.startArgsForCall[i].env, fake.startArgsForCall[i].filesToUpload, fake.startArgsForCall[i].builder
}

func (fake *VM) StartReturns(result1 error) {
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *VM) StartReturnsOnCall(i int, result1 error) {
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

func (fake *VM) Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {
	fake.stopMutex.Lock()
	ret, specificReturn := fake.stopReturnsOnCall[len(fake.stopArgsForCall)]
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
		ccid       ccintf.CCID
		timeout    uint
		dontkill   bool
		dontremove bool
	}{ccid, timeout, dontkill, dontremove})
	fake.recordInvocation("Stop", []interface{}{ccid, timeout, dontkill, dontremove})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub(ccid, timeout, dontkill, dontremove)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.stopReturns.result1
}

func (fake *VM) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *VM) StopArgsForCall(i int) (ccintf.CCID, uint, bool, bool) {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return fake.stopArgsForCall[i].ccid, fake.stopArgsForCall[i].timeout, fake.stopArgsForCall[i].dontkill, fake.stopArgsForCall[i].dontremove
}

func (fake *VM) StopReturns(result1 error) {
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

func (fake *VM) StopReturnsOnCall(i int, result1 error) {
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

func (fake *VM) Invocations() map[string][][]interface{} {
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
