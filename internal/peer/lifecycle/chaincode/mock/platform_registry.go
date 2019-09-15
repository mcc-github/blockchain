
package mock

import (
	"sync"
)

type PlatformRegistry struct {
	GetDeploymentPayloadStub        func(string, string) ([]byte, error)
	getDeploymentPayloadMutex       sync.RWMutex
	getDeploymentPayloadArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getDeploymentPayloadReturns struct {
		result1 []byte
		result2 error
	}
	getDeploymentPayloadReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	NormalizePathStub        func(string, string) (string, error)
	normalizePathMutex       sync.RWMutex
	normalizePathArgsForCall []struct {
		arg1 string
		arg2 string
	}
	normalizePathReturns struct {
		result1 string
		result2 error
	}
	normalizePathReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PlatformRegistry) GetDeploymentPayload(arg1 string, arg2 string) ([]byte, error) {
	fake.getDeploymentPayloadMutex.Lock()
	ret, specificReturn := fake.getDeploymentPayloadReturnsOnCall[len(fake.getDeploymentPayloadArgsForCall)]
	fake.getDeploymentPayloadArgsForCall = append(fake.getDeploymentPayloadArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetDeploymentPayload", []interface{}{arg1, arg2})
	fake.getDeploymentPayloadMutex.Unlock()
	if fake.GetDeploymentPayloadStub != nil {
		return fake.GetDeploymentPayloadStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getDeploymentPayloadReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PlatformRegistry) GetDeploymentPayloadCallCount() int {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	return len(fake.getDeploymentPayloadArgsForCall)
}

func (fake *PlatformRegistry) GetDeploymentPayloadCalls(stub func(string, string) ([]byte, error)) {
	fake.getDeploymentPayloadMutex.Lock()
	defer fake.getDeploymentPayloadMutex.Unlock()
	fake.GetDeploymentPayloadStub = stub
}

func (fake *PlatformRegistry) GetDeploymentPayloadArgsForCall(i int) (string, string) {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	argsForCall := fake.getDeploymentPayloadArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PlatformRegistry) GetDeploymentPayloadReturns(result1 []byte, result2 error) {
	fake.getDeploymentPayloadMutex.Lock()
	defer fake.getDeploymentPayloadMutex.Unlock()
	fake.GetDeploymentPayloadStub = nil
	fake.getDeploymentPayloadReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *PlatformRegistry) GetDeploymentPayloadReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getDeploymentPayloadMutex.Lock()
	defer fake.getDeploymentPayloadMutex.Unlock()
	fake.GetDeploymentPayloadStub = nil
	if fake.getDeploymentPayloadReturnsOnCall == nil {
		fake.getDeploymentPayloadReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getDeploymentPayloadReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *PlatformRegistry) NormalizePath(arg1 string, arg2 string) (string, error) {
	fake.normalizePathMutex.Lock()
	ret, specificReturn := fake.normalizePathReturnsOnCall[len(fake.normalizePathArgsForCall)]
	fake.normalizePathArgsForCall = append(fake.normalizePathArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("NormalizePath", []interface{}{arg1, arg2})
	fake.normalizePathMutex.Unlock()
	if fake.NormalizePathStub != nil {
		return fake.NormalizePathStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.normalizePathReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PlatformRegistry) NormalizePathCallCount() int {
	fake.normalizePathMutex.RLock()
	defer fake.normalizePathMutex.RUnlock()
	return len(fake.normalizePathArgsForCall)
}

func (fake *PlatformRegistry) NormalizePathCalls(stub func(string, string) (string, error)) {
	fake.normalizePathMutex.Lock()
	defer fake.normalizePathMutex.Unlock()
	fake.NormalizePathStub = stub
}

func (fake *PlatformRegistry) NormalizePathArgsForCall(i int) (string, string) {
	fake.normalizePathMutex.RLock()
	defer fake.normalizePathMutex.RUnlock()
	argsForCall := fake.normalizePathArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *PlatformRegistry) NormalizePathReturns(result1 string, result2 error) {
	fake.normalizePathMutex.Lock()
	defer fake.normalizePathMutex.Unlock()
	fake.NormalizePathStub = nil
	fake.normalizePathReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *PlatformRegistry) NormalizePathReturnsOnCall(i int, result1 string, result2 error) {
	fake.normalizePathMutex.Lock()
	defer fake.normalizePathMutex.Unlock()
	fake.NormalizePathStub = nil
	if fake.normalizePathReturnsOnCall == nil {
		fake.normalizePathReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.normalizePathReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *PlatformRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	fake.normalizePathMutex.RLock()
	defer fake.normalizePathMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PlatformRegistry) recordInvocation(key string, args []interface{}) {
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
