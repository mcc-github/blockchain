
package mock

import (
	"io"
	"sync"
)

type PlatformBuilder struct {
	GenerateDockerBuildStub        func(string, string, string, string, []byte) (io.Reader, error)
	generateDockerBuildMutex       sync.RWMutex
	generateDockerBuildArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 []byte
	}
	generateDockerBuildReturns struct {
		result1 io.Reader
		result2 error
	}
	generateDockerBuildReturnsOnCall map[int]struct {
		result1 io.Reader
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PlatformBuilder) GenerateDockerBuild(arg1 string, arg2 string, arg3 string, arg4 string, arg5 []byte) (io.Reader, error) {
	var arg5Copy []byte
	if arg5 != nil {
		arg5Copy = make([]byte, len(arg5))
		copy(arg5Copy, arg5)
	}
	fake.generateDockerBuildMutex.Lock()
	ret, specificReturn := fake.generateDockerBuildReturnsOnCall[len(fake.generateDockerBuildArgsForCall)]
	fake.generateDockerBuildArgsForCall = append(fake.generateDockerBuildArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 []byte
	}{arg1, arg2, arg3, arg4, arg5Copy})
	fake.recordInvocation("GenerateDockerBuild", []interface{}{arg1, arg2, arg3, arg4, arg5Copy})
	fake.generateDockerBuildMutex.Unlock()
	if fake.GenerateDockerBuildStub != nil {
		return fake.GenerateDockerBuildStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.generateDockerBuildReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PlatformBuilder) GenerateDockerBuildCallCount() int {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	return len(fake.generateDockerBuildArgsForCall)
}

func (fake *PlatformBuilder) GenerateDockerBuildCalls(stub func(string, string, string, string, []byte) (io.Reader, error)) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = stub
}

func (fake *PlatformBuilder) GenerateDockerBuildArgsForCall(i int) (string, string, string, string, []byte) {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	argsForCall := fake.generateDockerBuildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *PlatformBuilder) GenerateDockerBuildReturns(result1 io.Reader, result2 error) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = nil
	fake.generateDockerBuildReturns = struct {
		result1 io.Reader
		result2 error
	}{result1, result2}
}

func (fake *PlatformBuilder) GenerateDockerBuildReturnsOnCall(i int, result1 io.Reader, result2 error) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = nil
	if fake.generateDockerBuildReturnsOnCall == nil {
		fake.generateDockerBuildReturnsOnCall = make(map[int]struct {
			result1 io.Reader
			result2 error
		})
	}
	fake.generateDockerBuildReturnsOnCall[i] = struct {
		result1 io.Reader
		result2 error
	}{result1, result2}
}

func (fake *PlatformBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PlatformBuilder) recordInvocation(key string, args []interface{}) {
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
