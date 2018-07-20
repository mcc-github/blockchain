
package mock

import (
	"archive/tar"
	"sync"

	platforms_test "github.com/mcc-github/blockchain/core/chaincode/platforms"
)

type Platform struct {
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct{}
	nameReturns     struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	ValidatePathStub        func(path string) error
	validatePathMutex       sync.RWMutex
	validatePathArgsForCall []struct {
		path string
	}
	validatePathReturns struct {
		result1 error
	}
	validatePathReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateCodePackageStub        func(code []byte) error
	validateCodePackageMutex       sync.RWMutex
	validateCodePackageArgsForCall []struct {
		code []byte
	}
	validateCodePackageReturns struct {
		result1 error
	}
	validateCodePackageReturnsOnCall map[int]struct {
		result1 error
	}
	GetDeploymentPayloadStub        func(path string) ([]byte, error)
	getDeploymentPayloadMutex       sync.RWMutex
	getDeploymentPayloadArgsForCall []struct {
		path string
	}
	getDeploymentPayloadReturns struct {
		result1 []byte
		result2 error
	}
	getDeploymentPayloadReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GenerateDockerfileStub        func() (string, error)
	generateDockerfileMutex       sync.RWMutex
	generateDockerfileArgsForCall []struct{}
	generateDockerfileReturns     struct {
		result1 string
		result2 error
	}
	generateDockerfileReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GenerateDockerBuildStub        func(path string, code []byte, tw *tar.Writer) error
	generateDockerBuildMutex       sync.RWMutex
	generateDockerBuildArgsForCall []struct {
		path string
		code []byte
		tw   *tar.Writer
	}
	generateDockerBuildReturns struct {
		result1 error
	}
	generateDockerBuildReturnsOnCall map[int]struct {
		result1 error
	}
	GetMetadataProviderStub        func(code []byte) platforms_test.MetadataProvider
	getMetadataProviderMutex       sync.RWMutex
	getMetadataProviderArgsForCall []struct {
		code []byte
	}
	getMetadataProviderReturns struct {
		result1 platforms_test.MetadataProvider
	}
	getMetadataProviderReturnsOnCall map[int]struct {
		result1 platforms_test.MetadataProvider
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Platform) Name() string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct{}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.nameReturns.result1
}

func (fake *Platform) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *Platform) NameReturns(result1 string) {
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *Platform) NameReturnsOnCall(i int, result1 string) {
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

func (fake *Platform) ValidatePath(path string) error {
	fake.validatePathMutex.Lock()
	ret, specificReturn := fake.validatePathReturnsOnCall[len(fake.validatePathArgsForCall)]
	fake.validatePathArgsForCall = append(fake.validatePathArgsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("ValidatePath", []interface{}{path})
	fake.validatePathMutex.Unlock()
	if fake.ValidatePathStub != nil {
		return fake.ValidatePathStub(path)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validatePathReturns.result1
}

func (fake *Platform) ValidatePathCallCount() int {
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	return len(fake.validatePathArgsForCall)
}

func (fake *Platform) ValidatePathArgsForCall(i int) string {
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	return fake.validatePathArgsForCall[i].path
}

func (fake *Platform) ValidatePathReturns(result1 error) {
	fake.ValidatePathStub = nil
	fake.validatePathReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) ValidatePathReturnsOnCall(i int, result1 error) {
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

func (fake *Platform) ValidateCodePackage(code []byte) error {
	var codeCopy []byte
	if code != nil {
		codeCopy = make([]byte, len(code))
		copy(codeCopy, code)
	}
	fake.validateCodePackageMutex.Lock()
	ret, specificReturn := fake.validateCodePackageReturnsOnCall[len(fake.validateCodePackageArgsForCall)]
	fake.validateCodePackageArgsForCall = append(fake.validateCodePackageArgsForCall, struct {
		code []byte
	}{codeCopy})
	fake.recordInvocation("ValidateCodePackage", []interface{}{codeCopy})
	fake.validateCodePackageMutex.Unlock()
	if fake.ValidateCodePackageStub != nil {
		return fake.ValidateCodePackageStub(code)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validateCodePackageReturns.result1
}

func (fake *Platform) ValidateCodePackageCallCount() int {
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	return len(fake.validateCodePackageArgsForCall)
}

func (fake *Platform) ValidateCodePackageArgsForCall(i int) []byte {
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	return fake.validateCodePackageArgsForCall[i].code
}

func (fake *Platform) ValidateCodePackageReturns(result1 error) {
	fake.ValidateCodePackageStub = nil
	fake.validateCodePackageReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) ValidateCodePackageReturnsOnCall(i int, result1 error) {
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

func (fake *Platform) GetDeploymentPayload(path string) ([]byte, error) {
	fake.getDeploymentPayloadMutex.Lock()
	ret, specificReturn := fake.getDeploymentPayloadReturnsOnCall[len(fake.getDeploymentPayloadArgsForCall)]
	fake.getDeploymentPayloadArgsForCall = append(fake.getDeploymentPayloadArgsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("GetDeploymentPayload", []interface{}{path})
	fake.getDeploymentPayloadMutex.Unlock()
	if fake.GetDeploymentPayloadStub != nil {
		return fake.GetDeploymentPayloadStub(path)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getDeploymentPayloadReturns.result1, fake.getDeploymentPayloadReturns.result2
}

func (fake *Platform) GetDeploymentPayloadCallCount() int {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	return len(fake.getDeploymentPayloadArgsForCall)
}

func (fake *Platform) GetDeploymentPayloadArgsForCall(i int) string {
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	return fake.getDeploymentPayloadArgsForCall[i].path
}

func (fake *Platform) GetDeploymentPayloadReturns(result1 []byte, result2 error) {
	fake.GetDeploymentPayloadStub = nil
	fake.getDeploymentPayloadReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Platform) GetDeploymentPayloadReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Platform) GenerateDockerfile() (string, error) {
	fake.generateDockerfileMutex.Lock()
	ret, specificReturn := fake.generateDockerfileReturnsOnCall[len(fake.generateDockerfileArgsForCall)]
	fake.generateDockerfileArgsForCall = append(fake.generateDockerfileArgsForCall, struct{}{})
	fake.recordInvocation("GenerateDockerfile", []interface{}{})
	fake.generateDockerfileMutex.Unlock()
	if fake.GenerateDockerfileStub != nil {
		return fake.GenerateDockerfileStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.generateDockerfileReturns.result1, fake.generateDockerfileReturns.result2
}

func (fake *Platform) GenerateDockerfileCallCount() int {
	fake.generateDockerfileMutex.RLock()
	defer fake.generateDockerfileMutex.RUnlock()
	return len(fake.generateDockerfileArgsForCall)
}

func (fake *Platform) GenerateDockerfileReturns(result1 string, result2 error) {
	fake.GenerateDockerfileStub = nil
	fake.generateDockerfileReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *Platform) GenerateDockerfileReturnsOnCall(i int, result1 string, result2 error) {
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

func (fake *Platform) GenerateDockerBuild(path string, code []byte, tw *tar.Writer) error {
	var codeCopy []byte
	if code != nil {
		codeCopy = make([]byte, len(code))
		copy(codeCopy, code)
	}
	fake.generateDockerBuildMutex.Lock()
	ret, specificReturn := fake.generateDockerBuildReturnsOnCall[len(fake.generateDockerBuildArgsForCall)]
	fake.generateDockerBuildArgsForCall = append(fake.generateDockerBuildArgsForCall, struct {
		path string
		code []byte
		tw   *tar.Writer
	}{path, codeCopy, tw})
	fake.recordInvocation("GenerateDockerBuild", []interface{}{path, codeCopy, tw})
	fake.generateDockerBuildMutex.Unlock()
	if fake.GenerateDockerBuildStub != nil {
		return fake.GenerateDockerBuildStub(path, code, tw)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.generateDockerBuildReturns.result1
}

func (fake *Platform) GenerateDockerBuildCallCount() int {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	return len(fake.generateDockerBuildArgsForCall)
}

func (fake *Platform) GenerateDockerBuildArgsForCall(i int) (string, []byte, *tar.Writer) {
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	return fake.generateDockerBuildArgsForCall[i].path, fake.generateDockerBuildArgsForCall[i].code, fake.generateDockerBuildArgsForCall[i].tw
}

func (fake *Platform) GenerateDockerBuildReturns(result1 error) {
	fake.GenerateDockerBuildStub = nil
	fake.generateDockerBuildReturns = struct {
		result1 error
	}{result1}
}

func (fake *Platform) GenerateDockerBuildReturnsOnCall(i int, result1 error) {
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

func (fake *Platform) GetMetadataProvider(code []byte) platforms_test.MetadataProvider {
	var codeCopy []byte
	if code != nil {
		codeCopy = make([]byte, len(code))
		copy(codeCopy, code)
	}
	fake.getMetadataProviderMutex.Lock()
	ret, specificReturn := fake.getMetadataProviderReturnsOnCall[len(fake.getMetadataProviderArgsForCall)]
	fake.getMetadataProviderArgsForCall = append(fake.getMetadataProviderArgsForCall, struct {
		code []byte
	}{codeCopy})
	fake.recordInvocation("GetMetadataProvider", []interface{}{codeCopy})
	fake.getMetadataProviderMutex.Unlock()
	if fake.GetMetadataProviderStub != nil {
		return fake.GetMetadataProviderStub(code)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getMetadataProviderReturns.result1
}

func (fake *Platform) GetMetadataProviderCallCount() int {
	fake.getMetadataProviderMutex.RLock()
	defer fake.getMetadataProviderMutex.RUnlock()
	return len(fake.getMetadataProviderArgsForCall)
}

func (fake *Platform) GetMetadataProviderArgsForCall(i int) []byte {
	fake.getMetadataProviderMutex.RLock()
	defer fake.getMetadataProviderMutex.RUnlock()
	return fake.getMetadataProviderArgsForCall[i].code
}

func (fake *Platform) GetMetadataProviderReturns(result1 platforms_test.MetadataProvider) {
	fake.GetMetadataProviderStub = nil
	fake.getMetadataProviderReturns = struct {
		result1 platforms_test.MetadataProvider
	}{result1}
}

func (fake *Platform) GetMetadataProviderReturnsOnCall(i int, result1 platforms_test.MetadataProvider) {
	fake.GetMetadataProviderStub = nil
	if fake.getMetadataProviderReturnsOnCall == nil {
		fake.getMetadataProviderReturnsOnCall = make(map[int]struct {
			result1 platforms_test.MetadataProvider
		})
	}
	fake.getMetadataProviderReturnsOnCall[i] = struct {
		result1 platforms_test.MetadataProvider
	}{result1}
}

func (fake *Platform) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.validatePathMutex.RLock()
	defer fake.validatePathMutex.RUnlock()
	fake.validateCodePackageMutex.RLock()
	defer fake.validateCodePackageMutex.RUnlock()
	fake.getDeploymentPayloadMutex.RLock()
	defer fake.getDeploymentPayloadMutex.RUnlock()
	fake.generateDockerfileMutex.RLock()
	defer fake.generateDockerfileMutex.RUnlock()
	fake.generateDockerBuildMutex.RLock()
	defer fake.generateDockerBuildMutex.RUnlock()
	fake.getMetadataProviderMutex.RLock()
	defer fake.getMetadataProviderMutex.RUnlock()
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
