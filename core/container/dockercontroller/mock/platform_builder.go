
package mock

import (
	"io"
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
)

type PlatformBuilder struct {
	GenerateDockerBuildStub        func(*ccprovider.ChaincodeContainerInfo, []byte) (io.Reader, error)
	generateDockerBuildMutex       sync.RWMutex
	generateDockerBuildArgsForCall []struct {
		arg1 *ccprovider.ChaincodeContainerInfo
		arg2 []byte
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

func (fake *PlatformBuilder) GenerateDockerBuild(arg1 *ccprovider.ChaincodeContainerInfo, arg2 []byte) (io.Reader, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.generateDockerBuildMutex.Lock()
	ret, specificReturn := fake.generateDockerBuildReturnsOnCall[len(fake.generateDockerBuildArgsForCall)]
	fake.generateDockerBuildArgsForCall = append(fake.generateDockerBuildArgsForCall, struct {
		arg1 *ccprovider.ChaincodeContainerInfo
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("GenerateDockerBuild", []interface{}{arg1, arg2Copy})
	fake.generateDockerBuildMutex.Unlock()
	if fake.GenerateDockerBuildStub != nil {
		return fake.GenerateDockerBuildStub(arg1, arg2)
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

func (fake *PlatformBuilder) GenerateDockerBuildCalls(stub func(*ccprovider.ChaincodeContainerInfo, []byte) (io.Reader, error)) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = stub
}

func (fake *PlatformBuilder) GenerateDockerBuildArgsForCall(i int) (*ccprovider.ChaincodeContainerInfo, []byte) {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	argsForCall := fake.generateDockerBuildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
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
