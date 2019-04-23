
package mock

import (
	"archive/tar"
	"sync"
)

type Platform struct {
	GenerateDockerBuildStub        func(string, []byte, *tar.Writer) error
	generateDockerBuildMutex       sync.RWMutex
	generateDockerBuildArgsForCall []struct {
		arg1 string
		arg2 []byte
		arg3 *tar.Writer
	}
	generateDockerBuildReturns struct {
		result1 error
	}
	generateDockerBuildReturnsOnCall map[int]struct {
		result1 error
	}
	GenerateDockerfileStub        func() (string, error)
	generateDockerfileMutex       sync.RWMutex
	generateDockerfileArgsForCall []struct {
	}
	generateDockerfileReturns struct {
		result1 string
		result2 error
	}
	generateDockerfileReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetDeploymentPayloadStub        func(string) ([]byte, error)
	getDeploymentPayloadMutex       sync.RWMutex
	getDeploymentPayloadArgsForCall []struct {
		arg1 string
	}
	getDeploymentPayloadReturns struct {
		result1 []byte
		result2 error
	}
	getDeploymentPayloadReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetMetadataAsTarEntriesStub        func([]byte) ([]byte, error)
	getMetadataAsTarEntriesMutex       sync.RWMutex
	getMetadataAsTarEntriesArgsForCall []struct {
		arg1 []byte
	}
	getMetadataAsTarEntriesReturns struct {
		result1 []byte
		result2 error
	}
	getMetadataAsTarEntriesReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	ValidateCodePackageStub        func([]byte) error
	validateCodePackageMutex       sync.RWMutex
	validateCodePackageArgsForCall []struct {
		arg1 []byte
	}
	validateCodePackageReturns struct {
		result1 error
	}
	validateCodePackageReturnsOnCall map[int]struct {
		result1 error
	}
	ValidatePathStub        func(string) error
	validatePathMutex       sync.RWMutex
	validatePathArgsForCall []struct {
		arg1 string
	}
	validatePathReturns struct {
		result1 error
	}
	validatePathReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Platform) GenerateDockerBuild(arg1 string, arg2 []byte, arg3 *tar.Writer) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.generateDockerBuildMutex.Lock()
	ret, specificReturn := fake.generateDockerBuildReturnsOnCall[len(fake.generateDockerBuildArgsForCall)]
	fake.generateDockerBuildArgsForCall = append(fake.generateDockerBuildArgsForCall, struct {
		arg1 string
		arg2 []byte
		arg3 *tar.Writer
	}{arg1, arg2Copy, arg3})
	fake.recordInvocation("GenerateDockerBuild", []interface{}{arg1, arg2Copy, arg3})
	fake.generateDockerBuildMutex.Unlock()
	if fake.GenerateDockerBuildStub != nil {
		return fake.GenerateDockerBuildStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.generateDockerBuildReturns
	return fakeReturns.result1
}

func (fake *Platform) GenerateDockerBuildCallCount() int {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	return len(fake.generateDockerBuildArgsForCall)
}

func (fake *Platform) GenerateDockerBuildCalls(stub func(string, []byte, *tar.Writer) error) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = stub
}

func (fake *Platform) GenerateDockerBuildArgsForCall(i int) (string, []byte, *tar.Writer) {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	argsForCall := fake.generateDockerBuildArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Platform) GenerateDockerBuildReturns(result1 error) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = nil
	fake.generateDockerBuildReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) GenerateDockerBuildReturnsOnCall(i int, result1 error) {
	fake.generateDockerBuildMutex.Lock()
	defer fake.generateDockerBuildMutex.Unlock()
	fake.GenerateDockerBuildStub = nil
	if fake.generateDockerBuildReturnsOnCall == nil {
		fake.generateDockerBuildReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.generateDockerBuildReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Platform) GenerateDockerfile() (string, error) {
	fake.generateDockerfileMutex.Lock()
	ret, specificReturn := fake.generateDockerfileReturnsOnCall[len(fake.generateDockerfileArgsForCall)]
	fake.generateDockerfileArgsForCall = append(fake.generateDockerfileArgsForCall, struct {
	}{})
	fake.recordInvocation("GenerateDockerfile", []interface{}{})
	fake.generateDockerfileMutex.Unlock()
	if fake.GenerateDockerfileStub != nil {
		return fake.GenerateDockerfileStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.generateDockerfileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Platform) GenerateDockerfileCallCount() int {
	fake.generateDockerfileMutex.RLock()
	defer fake.generateDockerfileMutex.RUnlock()
	return len(fake.generateDockerfileArgsForCall)
}

func (fake *Platform) GenerateDockerfileCalls(stub func() (string, error)) {
	fake.generateDockerfileMutex.Lock()
	defer fake.generateDockerfileMutex.Unlock()
	fake.GenerateDockerfileStub = stub
}

func (fake *Platform) GenerateDockerfileReturns(result1 string, result2 error) {
	fake.generateDockerfileMutex.Lock()
	defer fake.generateDockerfileMutex.Unlock()
	fake.GenerateDockerfileStub = nil
	fake.generateDockerfileReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *Platform) GenerateDockerfileReturnsOnCall(i int, result1 string, result2 error) {
	fake.generateDockerfileMutex.Lock()
	defer fake.generateDockerfileMutex.Unlock()
	fake.GenerateDockerfileStub = nil
	if fake.generateDockerfileReturnsOnCall == nil {
		fake.generateDockerfileReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.generateDockerfileReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *Platform) GetDeploymentPayload(arg1 string) ([]byte, error) {
	fake.getDeploymentPayloadMutex.Lock()
	ret, specificReturn := fake.getDeploymentPayloadReturnsOnCall[len(fake.getDeploymentPayloadArgsForCall)]
	fake.getDeploymentPayloadArgsForCall = append(fake.getDeploymentPayloadArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetDeploymentPayload", []interface{}{arg1})
	fake.getDeploymentPayloadMutex.Unlock()
	if fake.GetDeploymentPayloadStub != nil {
		return fake.GetDeploymentPayloadStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getDeploymentPayloadReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Platform) GetDeploymentPayloadCallCount() int {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	return len(fake.getDeploymentPayloadArgsForCall)
}

func (fake *Platform) GetDeploymentPayloadCalls(stub func(string) ([]byte, error)) {
	fake.getDeploymentPayloadMutex.Lock()
	defer fake.getDeploymentPayloadMutex.Unlock()
	fake.GetDeploymentPayloadStub = stub
}

func (fake *Platform) GetDeploymentPayloadArgsForCall(i int) string {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	argsForCall := fake.getDeploymentPayloadArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Platform) GetDeploymentPayloadReturns(result1 []byte, result2 error) {
	fake.getDeploymentPayloadMutex.Lock()
	defer fake.getDeploymentPayloadMutex.Unlock()
	fake.GetDeploymentPayloadStub = nil
	fake.getDeploymentPayloadReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Platform) GetDeploymentPayloadReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Platform) GetMetadataAsTarEntries(arg1 []byte) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.getMetadataAsTarEntriesMutex.Lock()
	ret, specificReturn := fake.getMetadataAsTarEntriesReturnsOnCall[len(fake.getMetadataAsTarEntriesArgsForCall)]
	fake.getMetadataAsTarEntriesArgsForCall = append(fake.getMetadataAsTarEntriesArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("GetMetadataAsTarEntries", []interface{}{arg1Copy})
	fake.getMetadataAsTarEntriesMutex.Unlock()
	if fake.GetMetadataAsTarEntriesStub != nil {
		return fake.GetMetadataAsTarEntriesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getMetadataAsTarEntriesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Platform) GetMetadataAsTarEntriesCallCount() int {
	fake.getMetadataAsTarEntriesMutex.RLock()
	defer fake.getMetadataAsTarEntriesMutex.RUnlock()
	return len(fake.getMetadataAsTarEntriesArgsForCall)
}

func (fake *Platform) GetMetadataAsTarEntriesCalls(stub func([]byte) ([]byte, error)) {
	fake.getMetadataAsTarEntriesMutex.Lock()
	defer fake.getMetadataAsTarEntriesMutex.Unlock()
	fake.GetMetadataAsTarEntriesStub = stub
}

func (fake *Platform) GetMetadataAsTarEntriesArgsForCall(i int) []byte {
	fake.getMetadataAsTarEntriesMutex.RLock()
	defer fake.getMetadataAsTarEntriesMutex.RUnlock()
	argsForCall := fake.getMetadataAsTarEntriesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Platform) GetMetadataAsTarEntriesReturns(result1 []byte, result2 error) {
	fake.getMetadataAsTarEntriesMutex.Lock()
	defer fake.getMetadataAsTarEntriesMutex.Unlock()
	fake.GetMetadataAsTarEntriesStub = nil
	fake.getMetadataAsTarEntriesReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Platform) GetMetadataAsTarEntriesReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getMetadataAsTarEntriesMutex.Lock()
	defer fake.getMetadataAsTarEntriesMutex.Unlock()
	fake.GetMetadataAsTarEntriesStub = nil
	if fake.getMetadataAsTarEntriesReturnsOnCall == nil {
		fake.getMetadataAsTarEntriesReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getMetadataAsTarEntriesReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Platform) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *Platform) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *Platform) NameCalls(stub func() string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *Platform) NameReturns(result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *Platform) NameReturnsOnCall(i int, result1 string) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Platform) ValidateCodePackage(arg1 []byte) error {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.validateCodePackageMutex.Lock()
	ret, specificReturn := fake.validateCodePackageReturnsOnCall[len(fake.validateCodePackageArgsForCall)]
	fake.validateCodePackageArgsForCall = append(fake.validateCodePackageArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("ValidateCodePackage", []interface{}{arg1Copy})
	fake.validateCodePackageMutex.Unlock()
	if fake.ValidateCodePackageStub != nil {
		return fake.ValidateCodePackageStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateCodePackageReturns
	return fakeReturns.result1
}

func (fake *Platform) ValidateCodePackageCallCount() int {
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	return len(fake.validateCodePackageArgsForCall)
}

func (fake *Platform) ValidateCodePackageCalls(stub func([]byte) error) {
	fake.validateCodePackageMutex.Lock()
	defer fake.validateCodePackageMutex.Unlock()
	fake.ValidateCodePackageStub = stub
}

func (fake *Platform) ValidateCodePackageArgsForCall(i int) []byte {
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	argsForCall := fake.validateCodePackageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Platform) ValidateCodePackageReturns(result1 error) {
	fake.validateCodePackageMutex.Lock()
	defer fake.validateCodePackageMutex.Unlock()
	fake.ValidateCodePackageStub = nil
	fake.validateCodePackageReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) ValidateCodePackageReturnsOnCall(i int, result1 error) {
	fake.validateCodePackageMutex.Lock()
	defer fake.validateCodePackageMutex.Unlock()
	fake.ValidateCodePackageStub = nil
	if fake.validateCodePackageReturnsOnCall == nil {
		fake.validateCodePackageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateCodePackageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Platform) ValidatePath(arg1 string) error {
	fake.validatePathMutex.Lock()
	ret, specificReturn := fake.validatePathReturnsOnCall[len(fake.validatePathArgsForCall)]
	fake.validatePathArgsForCall = append(fake.validatePathArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ValidatePath", []interface{}{arg1})
	fake.validatePathMutex.Unlock()
	if fake.ValidatePathStub != nil {
		return fake.ValidatePathStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validatePathReturns
	return fakeReturns.result1
}

func (fake *Platform) ValidatePathCallCount() int {
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	return len(fake.validatePathArgsForCall)
}

func (fake *Platform) ValidatePathCalls(stub func(string) error) {
	fake.validatePathMutex.Lock()
	defer fake.validatePathMutex.Unlock()
	fake.ValidatePathStub = stub
}

func (fake *Platform) ValidatePathArgsForCall(i int) string {
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	argsForCall := fake.validatePathArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Platform) ValidatePathReturns(result1 error) {
	fake.validatePathMutex.Lock()
	defer fake.validatePathMutex.Unlock()
	fake.ValidatePathStub = nil
	fake.validatePathReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) ValidatePathReturnsOnCall(i int, result1 error) {
	fake.validatePathMutex.Lock()
	defer fake.validatePathMutex.Unlock()
	fake.ValidatePathStub = nil
	if fake.validatePathReturnsOnCall == nil {
		fake.validatePathReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validatePathReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Platform) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	fake.generateDockerfileMutex.RLock()
	defer fake.generateDockerfileMutex.RUnlock()
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	fake.getMetadataAsTarEntriesMutex.RLock()
	defer fake.getMetadataAsTarEntriesMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Platform) recordInvocation(key string, args []interface{}) {
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
